package ws

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"goflylivechat/common"
	"goflylivechat/models"
	"log"
	"time"
)

func NewVisitorServer(c *gin.Context) {
	visitorInfo := models.FindVisitorByVistorId(c.Query("visitor_id"))
	if visitorInfo.VisitorId == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "visitor not found",
		})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	prepareConn(conn)
	defer conn.Close()

	user := &User{
		Conn:       conn,
		Name:       visitorInfo.Name,
		Avator:     visitorInfo.Avator,
		Id:         visitorInfo.VisitorId,
		To_id:      visitorInfo.ToId,
		UpdateTime: time.Now(),
	}
	go models.UpdateVisitorStatus(visitorInfo.VisitorId, 1)

	AddVisitorToList(user)

	for {
		messageType, receive, err := conn.ReadMessage()
		if err != nil {
			visitor := removeVisitorConn(conn)
			if visitor != nil {
				log.Println("remove visitor websocket:", visitor.Id)
				VisitorOffline(visitor.To_id, visitor.Id, visitor.Name)
			}
			log.Println("read visitor websocket failed:", err)
			return
		}
		touchConnReadDeadline(conn)

		message <- &Message{
			conn:        conn,
			content:     receive,
			context:     c,
			messageType: messageType,
		}
	}
}

func AddVisitorToList(user *User) {
	oldUser := setVisitor(user)
	if oldUser != nil && oldUser.Conn != nil && oldUser.Conn != user.Conn {
		msg := TypeMessage{
			Type: "close",
			Data: user.Id,
		}
		str, _ := json.Marshal(msg)
		oldUser.Mux.Lock()
		if err := oldUser.Conn.WriteMessage(websocket.TextMessage, str); err != nil {
			log.Println("close old visitor websocket failed:", err)
		}
		oldUser.Mux.Unlock()
		oldUser.Conn.Close()
	}

	lastMessage := models.FindLastMessageByVisitorId(user.Id)
	userInfo := make(map[string]string)
	userInfo["uid"] = user.Id
	userInfo["username"] = user.Name
	userInfo["avator"] = user.Avator
	userInfo["last_message"] = lastMessage.Content
	if userInfo["last_message"] == "" {
		userInfo["last_message"] = "新訪客"
	}
	msg := TypeMessage{
		Type: "userOnline",
		Data: userInfo,
	}
	str, _ := json.Marshal(msg)
	OneKefuMessage(user.To_id, str)
}

func VisitorOnline(kefuId string, visitor models.Visitor) {
	lastMessage := models.FindLastMessageByVisitorId(visitor.VisitorId)
	userInfo := make(map[string]string)
	userInfo["uid"] = visitor.VisitorId
	userInfo["username"] = visitor.Name
	userInfo["avator"] = visitor.Avator
	userInfo["last_message"] = lastMessage.Content
	if userInfo["last_message"] == "" {
		userInfo["last_message"] = "新訪客"
	}
	msg := TypeMessage{
		Type: "userOnline",
		Data: userInfo,
	}
	str, _ := json.Marshal(msg)
	OneKefuMessage(kefuId, str)
}

func VisitorOffline(kefuId string, visitorId string, visitorName string) {
	models.UpdateVisitorStatus(visitorId, 0)
	userInfo := make(map[string]string)
	userInfo["uid"] = visitorId
	userInfo["name"] = visitorName
	msg := TypeMessage{
		Type: "userOffline",
		Data: userInfo,
	}
	str, _ := json.Marshal(msg)
	OneKefuMessage(kefuId, str)
}

func VisitorNotice(visitorId string, notice string) {
	msg := TypeMessage{
		Type: "notice",
		Data: notice,
	}
	str, _ := json.Marshal(msg)
	writeVisitorMessage(visitorId, str)
}

func VisitorMessage(visitorId, content string, kefuInfo models.User) {
	msg := TypeMessage{
		Type: "message",
		Data: ClientMessage{
			Name:    common.PublicKefuName,
			Avator:  kefuInfo.Avator,
			Id:      kefuInfo.Name,
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			ToId:    visitorId,
			Content: content,
			IsKefu:  "no",
		},
	}
	str, _ := json.Marshal(msg)
	writeVisitorMessage(visitorId, str)
}

func writeVisitorMessage(visitorId string, str []byte) bool {
	visitor, ok := getVisitor(visitorId)
	if !ok || visitor.Conn == nil {
		return false
	}

	visitor.Mux.Lock()
	visitor.Conn.SetWriteDeadline(time.Now().Add(websocketWriteWait))
	err := visitor.Conn.WriteMessage(websocket.TextMessage, str)
	visitor.Mux.Unlock()
	if err != nil {
		log.Println("send websocket message to visitor failed:", visitorId, err)
		visitor.Conn.Close()
		if removeVisitorIfCurrent(visitor) {
			VisitorOffline(visitor.To_id, visitor.Id, visitor.Name)
		}
		return false
	}
	return true
}

func BroadcastVisitors(msg TypeMessage) {
	str, _ := json.Marshal(msg)
	for visitorId := range VisitorSnapshot() {
		writeVisitorMessage(visitorId, str)
	}
}

func CloseVisitor(visitorId string, msg TypeMessage) bool {
	visitor, ok := getVisitor(visitorId)
	if !ok || visitor.Conn == nil {
		return false
	}
	str, _ := json.Marshal(msg)
	writeVisitorMessage(visitorId, str)
	visitor.Conn.Close()
	removeVisitorIfCurrent(visitor)
	return true
}

func VisitorAutoReply(visitorInfo models.Visitor, kefuInfo models.User, content string) {
	reply := models.FindReplyItemByUserIdTitle(kefuInfo.Name, content)
	if reply.Content != "" {
		time.Sleep(1 * time.Second)
		VisitorMessage(visitorInfo.VisitorId, reply.Content, kefuInfo)
		KefuMessage(visitorInfo.VisitorId, reply.Content, kefuInfo)
		models.CreateMessage(kefuInfo.Name, visitorInfo.VisitorId, reply.Content, "kefu")
	}
	if !IsKefuOnline(kefuInfo.Name) {
		time.Sleep(1 * time.Second)
		config := models.FindConfigByUserId(kefuInfo.Name, "OfflineMessage")
		if config.ConfValue == "" || reply.Content != "" {
			return
		}
		VisitorMessage(visitorInfo.VisitorId, config.ConfValue, kefuInfo)
		models.CreateMessage(kefuInfo.Name, visitorInfo.VisitorId, config.ConfValue, "kefu")
	}
}

func CleanVisitorExpire() {
	go func() {
		log.Println("cleanVisitorExpire start...")
		for {
			for _, user := range VisitorSnapshot() {
				diff := time.Now().Sub(user.UpdateTime).Seconds()
				if diff < common.VisitorExpire {
					continue
				}
				msg := TypeMessage{
					Type: "auto_close",
					Data: user.Id,
				}
				str, _ := json.Marshal(msg)
				writeVisitorMessage(user.Id, str)
				log.Println(user.Name + ":cleanVisitorExpire finshed")
			}
			t := time.NewTimer(time.Second * 5)
			<-t.C
		}
	}()
}

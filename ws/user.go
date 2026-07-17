package ws

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"goflylivechat/common"
	"goflylivechat/models"
	"goflylivechat/tools"
	"log"
	"time"
)

func NewKefuServer(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")
	kefuInfo := models.FindUser(kefuName.(string))
	if kefuInfo.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "user not found",
		})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	prepareConn(conn)

	kefu := &User{
		Conn:   conn,
		Name:   kefuInfo.Nickname,
		Id:     kefuInfo.Name,
		Avator: kefuInfo.Avator,
	}
	AddKefuToList(kefu)
	defer func() {
		removeKefuConn(kefu.Id, conn)
		conn.Close()
	}()

	for {
		messageType, receive, err := conn.ReadMessage()
		if err != nil {
			log.Println("read kefu websocket failed:", err)
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

func AddKefuToList(kefu *User) {
	oldUser := setKefu(kefu)
	if oldUser == nil || oldUser.Conn == nil || oldUser.Conn == kefu.Conn {
		return
	}

	msg := TypeMessage{
		Type: "close",
		Data: kefu.Id,
	}
	str, _ := json.Marshal(msg)
	oldUser.Mux.Lock()
	if err := oldUser.Conn.WriteMessage(websocket.TextMessage, str); err != nil {
		log.Println("close old kefu websocket failed:", err)
	}
	oldUser.Mux.Unlock()
	oldUser.Conn.Close()
}

func OneKefuMessage(toId string, str []byte) bool {
	return writeKefuMessage(toId, str, true)
}

func writeKefuMessage(toId string, str []byte, logSend bool) bool {
	kefu, ok := getKefu(toId)
	if !ok || kefu.Conn == nil {
		return false
	}

	kefu.Mux.Lock()
	kefu.Conn.SetWriteDeadline(time.Now().Add(websocketWriteWait))
	err := kefu.Conn.WriteMessage(websocket.TextMessage, str)
	kefu.Mux.Unlock()
	if logSend {
		tools.Logger().Println("send_kefu_message", err, string(str))
	}
	if err != nil {
		log.Println("send websocket message to kefu failed:", toId, err)
		kefu.Conn.Close()
		removeKefuConn(toId, kefu.Conn)
		return false
	}
	return true
}

func writeKefuHeartbeat(toId string, str []byte) bool {
	kefu, ok := getKefu(toId)
	if !ok || kefu.Conn == nil {
		return false
	}

	kefu.Mux.Lock()
	kefu.Conn.SetWriteDeadline(time.Now().Add(websocketWriteWait))
	err := kefu.Conn.WriteMessage(websocket.TextMessage, str)
	if err == nil {
		kefu.MissedPing = 0
		kefu.Mux.Unlock()
		return true
	}
	kefu.MissedPing++
	missed := kefu.MissedPing
	conn := kefu.Conn
	kefu.Mux.Unlock()

	log.Println("send websocket heartbeat to kefu failed:", toId, err, "missed:", missed)
	if missed >= kefuHeartbeatMissLimit {
		conn.Close()
		removeKefuConn(toId, conn)
	}
	return false
}

func KefuMessage(visitorId, content string, kefuInfo models.User) {
	msg := TypeMessage{
		Type: "message",
		Data: ClientMessage{
			Name:    kefuInfo.Nickname,
			Avator:  common.KefuAvatar,
			Id:      visitorId,
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			ToId:    visitorId,
			Content: content,
			IsKefu:  "yes",
		},
	}
	str, _ := json.Marshal(msg)
	OneKefuMessage(kefuInfo.Name, str)
}

func SendPingToKefuClient() {
	msg := TypeMessage{
		Type: "many pong",
	}
	str, _ := json.Marshal(msg)
	for _, kefuId := range kefuIDs() {
		writeKefuHeartbeat(kefuId, str)
	}
}

package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"goflylivechat/models"
	"goflylivechat/tools"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type User struct {
	Conn       *websocket.Conn
	Name       string
	Id         string
	Avator     string
	To_id      string
	Role_id    string
	Mux        sync.Mutex
	UpdateTime time.Time
}
type Message struct {
	conn        *websocket.Conn
	context     *gin.Context
	content     []byte
	messageType int
	Mux         sync.Mutex
}
type TypeMessage struct {
	Type interface{} `json:"type"`
	Data interface{} `json:"data"`
}
type ClientMessage struct {
	Name      string `json:"name"`
	Avator    string `json:"avator"`
	Id        string `json:"id"`
	VisitorId string `json:"visitor_id"`
	Group     string `json:"group"`
	Time      string `json:"time"`
	ToId      string `json:"to_id"`
	Content   string `json:"content"`
	City      string `json:"city"`
	ClientIp  string `json:"client_ip"`
	Refer     string `json:"refer"`
	IsKefu    string `json:"is_kefu"`
}

var ClientList = make(map[string]*User)
var KefuList = make(map[string]*User)
var message = make(chan *Message, 10)
var upgrader = websocket.Upgrader{}
var Mux sync.RWMutex

func getVisitor(visitorId string) (*User, bool) {
	Mux.RLock()
	defer Mux.RUnlock()
	user, ok := ClientList[visitorId]
	return user, ok && user != nil
}

func getKefu(kefuId string) (*User, bool) {
	Mux.RLock()
	defer Mux.RUnlock()
	user, ok := KefuList[kefuId]
	return user, ok && user != nil
}

func VisitorExists(visitorId string) bool {
	_, ok := getVisitor(visitorId)
	return ok
}

func IsKefuOnline(kefuId string) bool {
	_, ok := getKefu(kefuId)
	return ok
}

func VisitorCount() int {
	Mux.RLock()
	defer Mux.RUnlock()
	return len(ClientList)
}

func TouchVisitor(visitorId string) {
	user, ok := getVisitor(visitorId)
	if ok {
		user.UpdateTime = time.Now()
	}
}

func VisitorSnapshot() map[string]*User {
	Mux.RLock()
	defer Mux.RUnlock()
	users := make(map[string]*User, len(ClientList))
	for id, user := range ClientList {
		users[id] = user
	}
	return users
}

func kefuIDs() []string {
	Mux.RLock()
	defer Mux.RUnlock()
	ids := make([]string, 0, len(KefuList))
	for id := range KefuList {
		ids = append(ids, id)
	}
	return ids
}

func setVisitor(user *User) *User {
	Mux.Lock()
	defer Mux.Unlock()
	oldUser := ClientList[user.Id]
	ClientList[user.Id] = user
	return oldUser
}

func setKefu(user *User) *User {
	Mux.Lock()
	defer Mux.Unlock()
	oldUser := KefuList[user.Id]
	KefuList[user.Id] = user
	return oldUser
}

func removeVisitorConn(conn *websocket.Conn) *User {
	Mux.Lock()
	defer Mux.Unlock()
	for id, user := range ClientList {
		if user != nil && user.Conn == conn {
			delete(ClientList, id)
			return user
		}
	}
	return nil
}

func removeVisitorIfCurrent(user *User) bool {
	if user == nil {
		return false
	}
	Mux.Lock()
	defer Mux.Unlock()
	current := ClientList[user.Id]
	if current != user {
		return false
	}
	delete(ClientList, user.Id)
	return true
}

func removeKefuConn(kefuId string, conn *websocket.Conn) bool {
	Mux.Lock()
	defer Mux.Unlock()
	current := KefuList[kefuId]
	if current == nil || current.Conn != conn {
		return false
	}
	delete(KefuList, kefuId)
	return true
}

func userByConn(conn *websocket.Conn) *User {
	Mux.RLock()
	defer Mux.RUnlock()
	for _, user := range ClientList {
		if user != nil && user.Conn == conn {
			return user
		}
	}
	for _, user := range KefuList {
		if user != nil && user.Conn == conn {
			return user
		}
	}
	return nil
}

func writeConnMessage(conn *websocket.Conn, str []byte) error {
	user := userByConn(conn)
	if user != nil {
		user.Mux.Lock()
		defer user.Mux.Unlock()
	}
	return conn.WriteMessage(websocket.TextMessage, str)
}

func init() {
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	go UpdateVisitorStatusCron()
}
func SendServerJiang(title string, content string, domain string) string {
	noticeServerJiang, err := strconv.ParseBool(models.FindConfig("NoticeServerJiang"))
	serverJiangAPI := models.FindConfig("ServerJiangAPI")
	if err != nil || !noticeServerJiang || serverJiangAPI == "" {
		log.Println("do not notice serverjiang:", serverJiangAPI, noticeServerJiang)
		return ""
	}
	sendStr := fmt.Sprintf("%s%s", title, content)
	desp := title + ":" + content + "[登录](http://" + domain + "/main)"
	url := serverJiangAPI + "?text=" + sendStr + "&desp=" + desp
	//log.Println(url)
	res := tools.Get(url)
	return res
}
func SendFlyServerJiang(title string, content string, domain string) string {
	return ""
}

// 定时给更新数据库状态
func UpdateVisitorStatusCron() {
	for {
		visitors := models.FindVisitorsOnline()
		for _, visitor := range visitors {
			if visitor.VisitorId == "" {
				continue
			}
			if !VisitorExists(visitor.VisitorId) {
				models.UpdateVisitorStatus(visitor.VisitorId, 0)
			}
		}
		SendPingToKefuClient()
		time.Sleep(60 * time.Second)
	}
}

// 后端广播发送消息
func WsServerBackend() {
	for {
		message := <-message
		var typeMsg TypeMessage
		json.Unmarshal(message.content, &typeMsg)
		conn := message.conn
		if typeMsg.Type == nil || typeMsg.Data == nil {
			continue
		}
		msgType := typeMsg.Type.(string)
		log.Println("客户端:", string(message.content))

		switch msgType {
		//心跳
		case "ping":
			msg := TypeMessage{
				Type: "pong",
			}
			str, _ := json.Marshal(msg)
			if err := writeConnMessage(conn, str); err != nil {
				log.Println("write websocket pong failed:", err)
			}
		case "inputing":
			data := typeMsg.Data.(map[string]interface{})
			from := data["from"].(string)
			to := data["to"].(string)
			//限流
			if tools.LimitFreqSingle("inputing:"+from, 1, 2) {
				OneKefuMessage(to, message.content)
			}
		}

	}
}
func UpdateVisitorUser(visitorId string, toId string) {
	guest, ok := getVisitor(visitorId)
	if ok {
		guest.To_id = toId
	}
}

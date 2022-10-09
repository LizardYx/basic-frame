package middleware

import (
	"basic-frame/util/consts"
	"basic-frame/util/i18n/Localizer"
	"encoding/json"
	"github.com/gorilla/websocket"
)

// Client websocket客户端
type Client struct {
	ID     uint64
	Name   string
	Socket *websocket.Conn
	Send   chan []byte
}

// ClientManager websocket客户端管理
type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

/*
	websocket消息
	ID: 消息事件ID
	Content: 消息内容
	Desc: 消息描述
	Status: 状态(-1:失败 1:成功)
*/
type Message struct {
	ID      int         `json:"id"`
	Content interface{} `json:"content,omitempty"`
	Desc    string      `json:"desc,omitempty"`
	Status  int         `json:"status,omitempty"`
}

// WSDataCenterEventStruct 数据中心的websocket事件
type WSDataCenterEventStruct struct {
	EventID int         `json:"event_id"`
	Params  interface{} `json:"params"`
}

var Manager = ClientManager{
	Clients:    make(map[*Client]bool),
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func (a *ClientManager) Start() {
	for {
		select {
		case conn := <-a.Register:
			// 注册客户端
			a.Clients[conn] = true
			//content := "Websocket连接成功"
			content := Localizer.I18n.Translate("Websocket connect Successfully")
			jsonMessage, _ := json.Marshal(&Message{ID: consts.WSEventIDConnectSuccess, Content: content})
			var clientsID []uint64
			clientsID = append(clientsID, conn.ID)
			a.SendMessage(jsonMessage, clientsID)
		case conn := <-a.Unregister:
			// 注销客户端
			if _, ok := a.Clients[conn]; ok {
				close(conn.Send)
				delete(a.Clients, conn)
				//content := "Websocket连接关闭"
				content := Localizer.I18n.Translate("Websocket connect Closed")
				jsonMessage, _ := json.Marshal(&Message{ID: consts.WSEventIDConnectClosed, Content: content})
				var clientsID []uint64
				clientsID = append(clientsID, conn.ID)
				a.SendMessage(jsonMessage, clientsID)
			}
		case message := <-a.Broadcast:
			// 客户端广播消息
			for conn := range a.Clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(a.Clients, conn)
				}
			}
		}
	}
}

// SendMessage 发送消息给指定客户端组
func (a *ClientManager) SendMessage(message []byte, clientsID []uint64) {
	for _, clientID := range clientsID {
		for conn := range a.Clients {
			if conn.ID == clientID {
				conn.Send <- message
			}
		}
	}
}

// SendBroadcast 发送广播消息
func (a *ClientManager) SendBroadcast(message []byte) {
	Manager.Broadcast <- message
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()
	for {
		var WSDataCenterEvent WSDataCenterEventStruct
		// 读取客户端发送的消息
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		if string(message[:]) == "close" {
			Manager.Unregister <- c
		}
		if err := json.Unmarshal(message, &WSDataCenterEvent); err == nil {

		}
	}
}

func (c *Client) Write() {
	defer func() {
		// 关闭Socket连接
		c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

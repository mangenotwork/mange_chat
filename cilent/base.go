package cilent

import (
	"bytes"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mangenotwork/mange_chat/obj"
)

var (
	newline = []byte{'\n'}

	space = []byte{' '}

	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Message struct {
	Name         string `json:"name"`
	Time         string `json:"time"`
	RoomManCount int    `json:"count"`     // 匿名房间人数
	Data         string `json:"data"`      //输入内容
	MsgState     string `json:"msg_state"` //消息读取状态 已读/未读
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// -----------------------------------------------------------
//						匿名聊天室
// -----------------------------------------------------------

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func AnonymityWritePump(c *obj.AnonymityClient) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func AnonymityReadPump(c *obj.AnonymityClient) {
	defer obj.OutAnonymityRoom(c)
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Println("Read message = ", string(message), c.Name)

		// 消息内容
		m := &Message{
			Name:         c.Name,
			Time:         time.Now().Format("2006-01-02 15:04:05"),
			RoomManCount: obj.GetAnonymityRoomCount(),
			Data:         string(message),
		}
		log.Println("write message : ", string(message), m)

		//序列化
		data, err := json.Marshal(&m)
		if err != nil {
			log.Println("序列化失败,error=", err)
		}

		//广播到每个client
		for client, _ := range obj.AnonymityRoom {
			select {
			case client.Send <- data:
			default:
				obj.OutAnonymityRoom(client)
			}
		}
	}
}

// -----------------------------------------------------------
//						指定聊天室
// -----------------------------------------------------------

// 指定房间
func RoomWritePump(c *obj.UserC) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:

			if c.Conn == nil {
				continue
			}

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func RoomReadPump(c *obj.UserC) {
	defer func() {
		obj.UserOutRoom(c)
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Println("Read message = ", string(message), c.Name)

		//获取所在房间
		log.Println("RoomName = ", c.RoomName)
		room := obj.GetRoom(c.RoomName)
		log.Println("获取所在房间 = ", room)

		// 消息内容
		m := &Message{
			Name:         c.Name,
			Time:         time.Now().Format("2006-01-02 15:04:05"),
			RoomManCount: room.GetManCount(),
			Data:         string(message),
		}
		log.Println("write message : ", string(message), m)

		//序列化
		data, err := json.Marshal(&m)
		if err != nil {
			log.Println("序列化失败,error=", err)
		}

		//广播到每个client
		for client, _ := range room.AllUser {
			select {
			case client.Send <- data:
			default:
				obj.UserOutRoom(c)
			}
		}
	}
}

// -----------------------------------------------------------
//						1对1聊天
// -----------------------------------------------------------

func OnebyoneRoomWritePump(c *obj.UserC) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:

			if c.Conn == nil {
				continue
			}

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func OnebyoneRoomReadPump(c *obj.UserC) {
	defer func() {
		obj.UserOutOnebyoneRoom(c)
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Println("Read message = ", string(message), c.Name)

		//获取所在房间
		log.Println("RoomName = ", c.RoomName)
		room := obj.GetOnebyoneRoom(c.RoomName)
		log.Println("获取所在房间 = ", room)

		mesgState := "未读"
		roomMan := len(room.AllUser)
		if roomMan >= 2 {
			mesgState = "已读"
		}

		// 消息内容
		m := &Message{
			Name:         c.Name,
			Time:         time.Now().Format("2006-01-02 15:04:05"),
			RoomManCount: room.GetManCount(),
			Data:         string(message),
			MsgState:     mesgState,
		}
		log.Println("write message : ", string(message), m)

		//序列化
		data, err := json.Marshal(&m)
		if err != nil {
			log.Println("序列化失败,error=", err)
		}

		//广播到每个client
		for client, _ := range room.AllUser {
			select {
			case client.Send <- data:
			default:
				obj.UserOutOnebyoneRoom(c)
			}
		}
	}
}

// -----------------------------------------------------------
//						大厅
// -----------------------------------------------------------

var WLock sync.Mutex

func LobbyWritePump(c *obj.User) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if c.Conn != nil {
			c.Conn.Close()
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:
			//log.Println("conn = ", c.Conn)
			//WLock.Lock()
			//defer WLock.Unlock()
			if c.Conn == nil {
				continue
			}

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			err = w.Close()
			if err != nil {
				return
			}
		case <-ticker.C:
			if c.Conn == nil {
				return
			}
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func LobbyReadPump(c *obj.User) {
	defer obj.OutLobby(c)

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {

		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Println("Read message = ", string(message), c.Name)

		// 消息内容
		m := &Message{
			Name:         c.Name,
			Time:         time.Now().Format("2006-01-02 15:04:05"),
			RoomManCount: obj.GetAnonymityRoomCount(),
			Data:         string(message),
		}
		log.Println("write message : ", string(message), m)

		//序列化
		data, err := json.Marshal(&m)
		if err != nil {
			log.Println("序列化失败,error=", err)
		}

		log.Println(obj.Lobby)

		//广播到每个client
		for client, _ := range obj.Lobby {

			select {
			case client.Send <- data:
			default:
				obj.OutLobby(c)
			}
		}
	}
}

// 输出的用户列表
type UserInfo struct {
	Name   string `json:"user_name"`
	Online bool   `json:"online"`
}

func LobbyReadPump2(c *obj.User) {
	for {
		select {
		case message := <-c.Cmd:

			userList := make([]*UserInfo, 0)
			roomList := make([]*RoomInfo, 0)

			//获取当前用户列表
			for k, v := range obj.AllUser {
				online := false
				if v.Conn != nil {
					online = true
				}
				userList = append(userList, &UserInfo{
					Name:   k,
					Online: online,
				})
			}

			userlist, err := json.Marshal(&userList)
			log.Println("userlist = ", string(userlist))

			//获取当前房间列表
			for k, _ := range obj.AllRoom {
				roomList = append(roomList, &RoomInfo{
					Name: k,
				})
			}

			// 消息内容
			m := &LobbyMsg{
				Name:     c.Name,
				UserList: userList,
				RoomList: roomList,
			}
			log.Println("write message : ", string(message), m)

			//序列化
			data, err := json.Marshal(&m)
			if err != nil {
				log.Println("序列化失败,error=", err)
			}

			log.Println(obj.Lobby)

			//广播到每个client
			for client, _ := range obj.Lobby {

				if client.Send == nil {
					continue
				}

				select {
				case client.Send <- data:
				default:
					obj.OutLobby(c)
				}
			}

		}
	}

}

// 输出的房间列表
type RoomInfo struct {
	Name string `json:"room_name"`
}

// 大厅数据结构
type LobbyMsg struct {
	Name     string      `json:"my_name"`
	UserList []*UserInfo `json:"user_list"` // 当前用户列表
	RoomList []*RoomInfo `json:"room_list"` // 当前房间列表
}

//
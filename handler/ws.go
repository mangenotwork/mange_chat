package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/mange_chat/cilent"
	"github.com/mangenotwork/mange_chat/common/utils"
	"github.com/mangenotwork/mange_chat/dao"
	"github.com/mangenotwork/mange_chat/obj"
)

// 首页的socket
func WSIndex(c *gin.Context) {

	name := c.Query("name")
	log.Println("name = ", name)

	conn, err := cilent.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("连接错误 : ", err)
		return
	}
	log.Println("conn = ", conn)

	//用户连接
	u := obj.GetUser(name)
	log.Println("u1 = ", u)
	if u == nil {
		u = &obj.User{
			Name: name,
		}
	}
	u.Conn = conn
	u.Send = make(chan []byte, 256)
	u.Cmd = make(chan []byte, 256)
	obj.UserLogin(u)
	log.Println("u = ", u)

	//进入大厅
	obj.IntoLobby(u)

	u.Cmd <- []byte("上线")

	go cilent.LobbyWritePump(u)
	go cilent.LobbyReadPump(u)
	go cilent.LobbyReadPump2(u)

}

// 匿名聊天室, websocket 连接
func WSAnonymity(c *gin.Context) {
	// 建立websocket连接
	//header := nil
	conn, err := cilent.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("连接错误 : ", err)
		return
	}
	log.Println("conn = ", conn)

	name := c.Query("name")
	log.Println("name = ", name)

	if name == "" {
		name = utils.RandChar(12)
	}

	userName, _ := c.Cookie("user")

	user := &obj.AnonymityClient{
		Name:    name,
		User:    userName,
		Conn:    conn,
		Send:    make(chan []byte, 256),
		SendImg: make(chan []byte, 256),
	}
	obj.AddAnonymityRoom(user)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go cilent.AnonymityWritePump(user)
	go cilent.AnonymityReadPump(user)
	go cilent.AnonymityReadPump2(user)
}

//指定聊天室  websocket 连接
func WSRoom(c *gin.Context) {
	roomName := c.Query("room_name")
	log.Println("room_name = ", roomName)

	userName := c.Query("user_name")
	log.Println("user_name = ", userName)

	if roomName == "" || userName == "" {
		log.Println("房间或用户为空")
		return
	}

	//用户连接
	u := &obj.UserC{
		Token: utils.RandChar(10),
		Name:  userName,
	}

	// 建立websocket连接
	//header := nil
	conn, err := cilent.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("连接错误 : ", err)
		return
	}
	log.Println("conn = ", conn)

	u.Conn = conn
	u.Send = make(chan []byte, 256)
	u.SendIMg = make(chan []byte, 256)
	u.RoomName = roomName

	// 用户进入聊天室
	r := obj.GetRoom(roomName)

	if r == nil {
		log.Println("创建房间 ")
		// 创建房间
		r = &obj.Room{
			Name:    roomName,
			AllUser: make(map[*obj.UserC]bool, 0),
		}
		// 登记房间
		obj.AddRoom(r)
	}
	log.Println("r = ", r)
	r.Into(u)

	//用户上线
	new(dao.DaoMsg).UserToOnline(u.Name)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go cilent.RoomWritePump(u)
	go cilent.RoomReadPump(u)
	go cilent.RoomReadPump2(u)
}

//一对一聊天
func WSOnebyone(c *gin.Context) {
	roomName := c.Query("room")
	log.Println("room = ", roomName)

	userName := c.Query("myname")
	log.Println("myname = ", userName)

	youName := c.Query("youname")
	log.Println("youName = ", youName)

	if roomName == "" || userName == "" || youName == "" {
		log.Println("房间或用户为空")
		return
	}

	//用户连接
	u := &obj.UserC{
		Token: utils.RandChar(10),
		Name:  userName,
		You:   youName,
	}

	// 建立websocket连接
	//header := nil
	conn, err := cilent.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("连接错误 : ", err)
		return
	}
	log.Println("conn = ", conn)

	u.Conn = conn
	u.Send = make(chan []byte, 256)
	u.SendIMg = make(chan []byte, 256)
	u.RoomName = roomName

	// 用户进入聊天室
	r := obj.GetOnebyoneRoom(roomName)

	flg := true
	for k, _ := range r.AllUser {
		if k.Name == u.Name {
			k.Conn = u.Conn
			flg = false
			break
		}
	}
	if flg {
		r.Into(u)
	}
	log.Println("room obj = ", r)

	go cilent.OnebyoneRoomWritePump(u)
	go cilent.OnebyoneRoomReadPump(u)
	go cilent.OnebyoneRoomReadPump2(u)
}

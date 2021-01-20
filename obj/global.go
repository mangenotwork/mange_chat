package obj

import (
	"sync"

	"github.com/gorilla/websocket"
)

// -----------------------------------------------------------
//						匿名聊天室
// -----------------------------------------------------------

var (
	// 匿名房间
	AnonymityRoom map[*AnonymityClient]bool = make(map[*AnonymityClient]bool, 0)
	// 匿名房间锁
	AnonymityRoomLock sync.Mutex
)

// 匿名用户
type AnonymityClient struct {
	Name string

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte
}

// 加入匿名房间
func AddAnonymityRoom(c *AnonymityClient) {
	AnonymityRoomLock.Lock()
	defer AnonymityRoomLock.Unlock()
	AnonymityRoom[c] = true
}

// 获取匿名房间人数
func GetAnonymityRoomCount() int {
	return len(AnonymityRoom)
}

// 退出匿名房间
func OutAnonymityRoom(c *AnonymityClient) {
	AnonymityRoomLock.Lock()
	defer AnonymityRoomLock.Unlock()
	delete(AnonymityRoom, c)
	close(c.Send)
}

// 匿名的名称是否重复
func AnonymityName(name string) bool {
	for k, _ := range AnonymityRoom {
		AnonymityRoomLock.Lock()
		i := k.Name
		AnonymityRoomLock.Unlock()
		if i == name {
			return true
		}
	}
	return false
}

// -----------------------------------------------------------
//						用户
// -----------------------------------------------------------

// 用户
type User struct {
	Name string

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	Cmd chan []byte

	//当前所在房间 目前是但房间聊天
	RoomName string

	//是否在线
	Online bool
}

// 当前登录用户
var (
	AllUser     map[string]*User = make(map[string]*User, 0)
	AllUserLock sync.Mutex
)

// 用户登录
func UserLogin(u *User) {
	AllUserLock.Lock()
	defer AllUserLock.Unlock()
	AllUser[u.Name] = u
}

// 获取用户并连接
func GetUser(userName string) *User {
	AllUserLock.Lock()
	defer AllUserLock.Unlock()
	return AllUser[userName]
}

// 匿名的名称是否重复
func AllUserName(name string) bool {
	for k, _ := range AllUser {
		if k == name {
			return true
		}
	}
	return false
}

// 当前在线用户
var (
	Online     map[string]*User = make(map[string]*User, 0)
	OnlineLock sync.Mutex
)

// -----------------------------------------------------------
//						指定聊天室
// -----------------------------------------------------------

// 聊天室用户连接对象
type UserC struct {
	Token string

	Name string

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	//当前所在房间 目前是但房间聊天
	RoomName string

	//一对一聊天 对象
	You string
}

// 所有房间列表
var (
	AllRoom     map[string]*Room = make(map[string]*Room, 0)
	AllRoomLock sync.Mutex
)

// 房间
type Room struct {
	Id      int
	Name    string
	Lock    sync.Mutex
	AllUser map[*UserC]bool
}

// 登记房间
func AddRoom(room *Room) {
	AllRoomLock.Lock()
	defer AllRoomLock.Unlock()
	AllRoom[room.Name] = room
}

// 获取所有房间
func GetRoom(roomName string) *Room {
	AllRoomLock.Lock()
	defer AllRoomLock.Unlock()
	return AllRoom[roomName]
}

// 用户进入房间
func (r *Room) Into(u *UserC) {
	r.Lock.Lock()
	r.AllUser[u] = true
	r.Lock.Unlock()
}

// 用户退出房间
func (r *Room) Out(u *UserC) {
	delete(r.AllUser, u)
}

//获取当前房间人数
func (r *Room) GetManCount() int {
	return len(r.AllUser)
}

// 用户退出房间
func UserOutRoom(u *UserC) {
	r := GetRoom(u.RoomName)
	r.Out(u)
}

// -----------------------------------------------------------
//						大厅
// -----------------------------------------------------------

var (
	// 匿名房间
	Lobby map[*User]bool = make(map[*User]bool, 0)
	// 匿名房间锁
	LobbyLock sync.Mutex
)

// 进入大厅
func IntoLobby(c *User) {
	LobbyLock.Lock()
	defer LobbyLock.Unlock()
	Lobby[c] = true
}

// 退出大厅
func OutLobby(c *User) {
	LobbyLock.Lock()
	defer LobbyLock.Unlock()
	//delete(Lobby, c)
	c.Cmd <- []byte("下线")
	c.Conn = nil

	//close(c.Send)
	//close(c.Cmd)
}

// -----------------------------------------------------------
//						一对一聊天室
// -----------------------------------------------------------

// 所有房间列表
var (
	AllOnebyoneRoom     map[string]*Room = make(map[string]*Room, 0)
	AllOnebyoneRoomLock sync.Mutex
)

// 创建房间
func AddOnebyoneRoom(name string, room *Room) {
	AllOnebyoneRoomLock.Lock()
	defer AllOnebyoneRoomLock.Unlock()
	AllOnebyoneRoom[name] = room
}

// 获取房间
func GetOnebyoneRoom(roomName string) *Room {
	AllOnebyoneRoomLock.Lock()
	defer AllOnebyoneRoomLock.Unlock()
	return AllOnebyoneRoom[roomName]
}

// 用户退出房间
func UserOutOnebyoneRoom(u *UserC) {
	r := GetOnebyoneRoom(u.RoomName)
	r.Out(u)
}

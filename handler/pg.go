package handler

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/mange_chat/obj"
)

// 输出的用户列表
type UserList struct {
	Name   string
	Online bool
}

// 首页
func PGIndex(c *gin.Context) {
	user, err := c.Cookie("user")
	if err != nil || user == "" {
		log.Println("err = ", err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	//获取一次用户列表
	userList := make([]*UserList, 0)
	for k, v := range obj.AllUser {
		online := false
		if v.Conn != nil {
			online = true
		}
		userList = append(userList, &UserList{
			Name:   k,
			Online: online,
		})
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"user":      user,
		"user_list": userList,
	})
}

// 登录页面
func PGLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// 登录
func Login(c *gin.Context) {
	cookie, _ := c.Cookie("user")
	log.Println("cookie = ", cookie)
	if cookie != "" && obj.AllUserName(cookie) {
		c.Redirect(http.StatusFound, "/")
	}

	name := c.Param("name")
	if obj.AllUserName(name) {
		c.HTML(http.StatusOK, "error.html", gin.H{
			"err": "名称已存在",
		})
		return
	}

	user := &obj.User{
		Name: name,
	}
	obj.UserLogin(user)
	c.SetCookie("user", name, 60*60*24*7, "/", "0.0.0.0/24", false, true)
	c.Redirect(http.StatusFound, "/")
}

// 匿名聊天室
func PGAnonymity(c *gin.Context) {

	name, _ := c.Cookie("AnonymityUser")
	if name == "" {
		//当前聊天室人数
		number := len(obj.AnonymityRoom)
		log.Println("number = ", number)
		name = fmt.Sprintf("匿名者-%d", number)
	}

	c.SetCookie("AnonymityUser", name, 60*60*24*7, "/", "0.0.0.0/24", false, true)

	c.HTML(http.StatusOK, "anonymity.html", gin.H{
		"name": name,
	})
}

// 登录到匿名房间
func LoginAnonymity(c *gin.Context) {

	manCount := obj.GetAnonymityRoomCount()
	c.HTML(http.StatusOK, "login_anonymity.html", gin.H{
		"man_count": manCount,
	})
}

// 指定房间聊天
func LoginRoom(c *gin.Context) {
	//获取所有房间
	allRoom := obj.AllRoom

	//请求类型
	getType := c.Query("type")
	log.Println("getType = ", getType)

	user_name := c.Query("user_name")
	log.Println("user_name = ", user_name)
	room_name := c.Query("room_name")
	log.Println("room_name = ", room_name)

	switch getType {

	//创建用户
	case "login":
		// 创建用户并登录
		user := &obj.User{
			Name: user_name,
		}
		obj.UserLogin(user)

	//创建房间
	case "createroom":
		// 创建房间
		room := &obj.Room{
			Name:    room_name,
			AllUser: make(map[*obj.UserC]bool, 0),
		}
		// 登记房间
		obj.AddRoom(room)

	default:
		log.Println("default")
	}

	c.HTML(http.StatusOK, "login_room.html", gin.H{
		"user_name": user_name,
		"room":      room_name,
		"allroom":   allRoom,
	})

}

// 指定聊天室
func PGRoom(c *gin.Context) {

	// 获取user
	user, err := c.Cookie("user")
	if err != nil || user == "" {
		log.Println("err = ", err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	room_name := c.Query("room_name")
	log.Println("room_name = ", room_name)

	c.HTML(http.StatusOK, "room.html", gin.H{
		"user_name": user,
		"room_name": room_name,
	})
}

// 一对一聊天
func LoginOnebyone(c *gin.Context) {

	c.HTML(http.StatusOK, "login_onebyone.html", gin.H{})
}

// 一对一聊天
func PGOnebyone(c *gin.Context) {

	my_name, err := c.Cookie("user")
	if err != nil || my_name == "" {
		log.Println("err = ", err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	you_name := c.Query("you_name")
	log.Println("you_name = ", you_name)

	//名称进行排序生成key
	a := sort.StringSlice{my_name, you_name}
	sort.Sort(a)
	log.Println("名称进行排序生成key = ", a)
	room_name := a[0] + ":" + a[1]

	//寻找聊天室，没有则创建
	roomObj := obj.AllOnebyoneRoom[room_name]
	log.Println("寻找聊天室 = ", roomObj)
	if roomObj == nil {
		room := &obj.Room{
			Name:    room_name,
			AllUser: make(map[*obj.UserC]bool, 0),
		}
		obj.AddOnebyoneRoom(room_name, room)
	}

	log.Println("room_name = ", room_name)

	c.HTML(http.StatusOK, "onebyone.html", gin.H{
		"my_name":   my_name,
		"you_name":  you_name,
		"room_name": room_name,
	})
}

// 创建群聊
func CreateRoom(c *gin.Context) {

	// 获取user
	user, err := c.Cookie("user")
	if err != nil || user == "" {
		log.Println("err = ", err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	//获取user 对应的客户端
	u := obj.GetUser(user)

	room_name := c.Query("room_name")
	log.Println("room_name = ", room_name)
	// 创建房间
	room := &obj.Room{
		Name:    room_name,
		AllUser: make(map[*obj.UserC]bool, 0),
	}
	// 登记房间
	obj.AddRoom(room)
	c.JSON(http.StatusOK, "ok")

	u.Cmd <- []byte("创建")
}

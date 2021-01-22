package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/mange_chat/common/utils"
	"github.com/mangenotwork/mange_chat/dao"
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

	//获取用户信息
	myInfo := new(dao.DaoMsg).GetUserInfo(user)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"user":      user,
		"user_list": userList,
		"myInfo":    myInfo,
	})
}

// 登录页面
func PGLogin(c *gin.Context) {

	//头像
	headImgs := []string{
		"https://up.enterdesk.com/edpic/c3/e0/9e/c3e09ec5f8e67a4ef95260a34ec40ab7.jpg",
		"https://up.enterdesk.com/edpic/34/67/4a/34674ab55a0f9659a069aef1582e275d.jpg",
		"https://up.enterdesk.com/edpic/ff/0b/6d/ff0b6d250839a43eadab5e759dc4338c.jpg",
		"https://up.enterdesk.com/edpic/7f/0b/e0/7f0be0eac43000111eed8d72e5d16322.jpg",
		"https://up.enterdesk.com/edpic/6e/84/e8/6e84e8db61a3480d4fcc82d198186ee7.jpg",
		"https://up.enterdesk.com/edpic/6c/0a/8f/6c0a8fbb714d065d451976a57e20cb4a.jpg",
		"https://up.enterdesk.com/edpic/46/68/31/4668311d1c89c54acfe1b8bc9ab71e97.jpg",
		"https://up.enterdesk.com/edpic/5d/62/eb/5d62eb6cb26ded70c70c7a536ce38bb2.jpg",
		"https://up.enterdesk.com/edpic/42/06/a2/4206a2819113a2f67ec2da4e8d5d3d10.jpg",
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"imgs": headImgs,
	})
}

// 登录
func Login(c *gin.Context) {
	name := c.Param("name")

	// 如果不存在就新建
	if !new(dao.DaoMsg).UserNameRepeat(name) {
		img := c.Query("img")
		log.Println("头像 = ", img)
		if img == "" {
			img = "https://up.enterdesk.com/edpic/c3/e0/9e/c3e09ec5f8e67a4ef95260a34ec40ab7.jpg"
		}
		// 用户信息
		userInfo := map[string]string{
			"name": name,
			"img":  img,
			"date": time.Now().Format("2006-01-02 15:04:05"),
		}
		// 存储用户信息到redis
		new(dao.DaoMsg).SetUser(name, userInfo)
	}

	//用户进入大厅
	user := &obj.User{
		Name: name,
	}
	obj.UserLogin(user)

	c.SetCookie("user", name, 60*60*24*7, "/", "0.0.0.0/24", false, true)
	c.Redirect(http.StatusFound, "/")
}

// 登出
func LoginOut(c *gin.Context) {
	c.SetCookie("user", "", 60*60*24*7, "/", "0.0.0.0/24", false, true)
	c.Redirect(http.StatusFound, "/")
}

// 匿名聊天室
func PGAnonymity(c *gin.Context) {
	//当前聊天室人数
	number := len(obj.AnonymityRoom)

	name, _ := c.Cookie("AnonymityUser")
	if name == "" {
		name = utils.RandChar(12)
	}

	c.SetCookie("AnonymityUser", name, 60*60*24*7, "/", "0.0.0.0/24", false, true)

	c.HTML(http.StatusOK, "anonymity.html", gin.H{
		"name":   name,
		"number": number,
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

	// TODO  未读到已读
	// 获取所有消息

	// TODO  未读消息红点数清除

	//获取历史聊天记录
	historyMsg := new(dao.DaoMsg).Get(room_name)
	historyMsgJson, err := json.Marshal(&historyMsg)
	if err != nil {
		log.Printf("序列号错误 err=%v\n", err)
	}

	c.HTML(http.StatusOK, "onebyone.html", gin.H{
		"my_name":     my_name,
		"you_name":    you_name,
		"room_name":   room_name,
		"history_msg": string(historyMsgJson),
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

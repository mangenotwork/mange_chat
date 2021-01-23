package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/mange_chat/common/utils"
	"github.com/mangenotwork/mange_chat/dao"
	"github.com/mangenotwork/mange_chat/obj"
)

// 输出的用户列表
type UserList struct {
	Name   string
	Img    string
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

	// 在redis获取用户列表
	userList := make([]*UserList, 0)
	allUser := new(dao.DaoMsg).GetAllUser()
	for k, v := range allUser {
		online := false
		if utils.Str2Int(v) > 0 {
			online = true
		}
		info := new(dao.DaoMsg).GetUserInfo(k)
		userList = append(userList, &UserList{
			Name:   k,
			Img:    info["ing"],
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

	//上线
	new(dao.DaoMsg).UserToOnline(name)

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

	roomName := c.Query("room_name")
	log.Println("room_name = ", roomName)

	// 群聊当前在线人数
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
	count := len(r.AllUser)

	historyMsg := new(dao.DaoMsg).GetRoomMsg(roomName)
	historyMsgJson, err := json.Marshal(&historyMsg)
	if err != nil {
		log.Printf("序列号错误 err=%v\n", err)
	}

	c.HTML(http.StatusOK, "room.html", gin.H{
		"user_name":   user,
		"room_name":   roomName,
		"count":       count,
		"history_msg": string(historyMsgJson),
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

	//未读条数清理
	new(dao.DaoMsg).EmptyUnreadMsg(my_name, you_name)

	//获取历史聊天记录
	historyMsg := new(dao.DaoMsg).Get(room_name)

	// //处理对方发的消息 未读到已读
	// for k, v := range historyMsg {
	// 	if strings.Index(v, fmt.Sprintf(`"name":"%s"`, my_name)) == -1 {
	// 		v = strings.Replace(v, "未读", "已读", -1)
	// 		log.Println(k+1, v)
	// 	}
	// }

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

	//获取所有room
	allroom := new(dao.DaoMsg).GetRoomList()
	for k, _ := range allroom {
		if room_name == k {
			c.JSON(http.StatusOK, "群聊已经存在")
			return
		}
	}

	// 创建房间
	room := &obj.Room{
		Name:    room_name,
		AllUser: make(map[*obj.UserC]bool, 0),
	}
	// 登记房间
	obj.AddRoom(room)

	new(dao.DaoMsg).NewRoom(room_name)

	//广播到每个在大厅的用户
	u.Cmd <- []byte("创建")

	c.JSON(http.StatusOK, "创建成功")

}

//	上传图片
func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	log.Println(file.Filename, err)
	//log.Println(file, err)

	myname := c.PostForm("myname")
	log.Println("myname = ", myname)

	youname := c.PostForm("youname")
	log.Println("youname = ", youname)

	roomname := c.PostForm("roomname")
	log.Println("roomname = ", roomname)

	// 聊天类型
	sendType := c.PostForm("type")

	ext := strings.Split(file.Filename, ".")
	extStr := ext[len(ext)-1]
	newFileName := fmt.Sprintf("%d.%s", time.Now().UnixNano(), extStr)

	// 上传文件到指定的路径
	err = c.SaveUploadedFile(file, "./img/"+newFileName)
	log.Println(err)
	if err != nil {
		c.JSON(http.StatusOK, "")
	}

	newFileNameUrl := "/img/" + newFileName
	imgShow := fmt.Sprintf(`<img src="%s" alt="" style="width: 250px;box-shadow:2px 2px 3px 3px rgba(0,0,0,.5);margin: 10px;">`, newFileNameUrl)

	switch sendType {

	// 一对一聊天
	case "obo":
		//获取房间
		//名称进行排序生成key
		a := sort.StringSlice{myname, youname}
		sort.Sort(a)
		log.Println("名称进行排序生成key = ", a)
		roomName := a[0] + ":" + a[1]
		r := obj.GetOnebyoneRoom(roomName)
		log.Println("r = ", r)
		for k, _ := range r.AllUser {
			if k.Name == myname {
				log.Println("Send obo")
				k.SendIMg <- []byte(imgShow)
				break
			}
		}

	//群聊
	case "room":
		r := obj.GetRoom(roomname)
		for k, _ := range r.AllUser {
			if k.Name == myname {
				log.Println("Send room : ", imgShow)
				k.SendIMg <- []byte(imgShow)
				break
			}
		}

	//匿名聊天
	case "anonymity":
		c := obj.GetAnonymity(myname)
		log.Println("Send anonymity : ", imgShow)
		c.SendImg <- []byte(imgShow)

	}

	c.JSON(http.StatusOK, newFileNameUrl)
}

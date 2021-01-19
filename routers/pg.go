// 页面
package routers

import (
	"github.com/mangenotwork/mange_chat/handler"
)

func PGRouter() {
	Router.GET("/", handler.PGIndex)               // 首页
	Router.GET("/login", handler.PGLogin)          // 登录页面
	Router.GET("/login/:name", handler.Login)      // 登录
	Router.GET("/anonymity", handler.PGAnonymity)  // 匿名聊天室
	Router.GET("/create/room", handler.CreateRoom) // 创建群聊
	//Router.GET("/login/anonymity", handler.LoginAnonymity) // 进入匿名聊天室
	//Router.GET("/login/room", handler.LoginRoom)           // 指定房间聊天
	Router.GET("/room", handler.PGRoom) // 指定聊天室
	//Router.GET("/login/onebyone", handler.LoginOnebyone)   // 一对一聊天
	Router.GET("/onebyone", handler.PGOnebyone) //
}

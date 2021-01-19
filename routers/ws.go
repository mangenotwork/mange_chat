package routers

import (
	"github.com/mangenotwork/mange_chat/handler"
)

func WSRouter() {
	WS := Router.Group("/ws")
	{
		WS.Any("/index", handler.WSIndex)         // 首页的socket
		WS.Any("/anonymity", handler.WSAnonymity) // 匿名聊天室
		WS.Any("/room", handler.WSRoom)           // 指定聊天室
		WS.Any("/onebyone", handler.WSOnebyone)   //一对一聊天
	}
}

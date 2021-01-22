// 页面
package routers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/mange_chat/handler"
)

func PGRouter() {
	Router.Use(IsLogin())
	Router.GET("/", handler.PGIndex)               // 首页
	Router.GET("/login", handler.PGLogin)          // 登录页面
	Router.GET("/login/:name", handler.Login)      // 登录
	Router.GET("/loginout", handler.LoginOut)      //登出
	Router.GET("/anonymity", handler.PGAnonymity)  // 匿名聊天室
	Router.GET("/create/room", handler.CreateRoom) // 创建群聊
	Router.GET("/room", handler.PGRoom)            // 指定聊天室
	Router.GET("/onebyone", handler.PGOnebyone)    //一对一聊天

}

func IsLogin() gin.HandlerFunc {

	return func(c *gin.Context) {
		log.Println(c.Request.URL.String())
		url := c.Request.URL.String()
		user, err := c.Cookie("user")
		if (err != nil || user == "") && strings.Index(url, "login") == -1 {
			log.Println("err = ", err)
			c.Redirect(http.StatusFound, "/login")
			return
		}

		c.Next()
		return
	}
}

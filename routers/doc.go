// 包含了所有路由与中间件

package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
}

func Routers() *gin.Engine {

	//静态目录配置
	// Router.Static("/static", "static")
	// Router.Static("/install/static", "static")
	Router.StaticFS("/static", http.Dir("./static"))

	//模板
	Router.LoadHTMLGlob("views/*")

	PGRouter()
	WSRouter()

	return Router
}

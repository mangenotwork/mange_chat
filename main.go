package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/mange_chat/routers"
)

func main() {
	gin.SetMode(gin.DebugMode)
	s := routers.Routers()
	port := "8555"
	s.Run(":" + port)
}

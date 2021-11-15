package router

import (
	"github.com/Alan-Lxc/crypto_contest/dcssweb/handlers"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/middleware"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/systeminit"
	"github.com/gin-gonic/gin"
)

var controll systeminit.Controll

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	v1 := r.Group("/api")
	{
		v1.GET("/ping", handlers.Ping)

	}
	r.POST("/api/newsecret", handlers.NewSecret)
	r.GET("/api/getsecretlist", handlers.Getsecretlist)
	return r
}

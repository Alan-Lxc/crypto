package router

import (
	"github.com/Alan-Lxc/crypto_contest/dcssweb/controller"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/middleware"
	"github.com/gin-gonic/gin"
)

//var controll systeminit.Controll

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)

	r.POST("/api/secret/newsecret", controller.NewSecret)
	r.GET("/api/secret/updatesecretcounter", controller.UpdateSecretCounter)
	r.GET("/api/secret/deletesecret", controller.DeleteSecret)
	r.GET("/api/secret/getsecret", controller.GetSecret)
	r.GET("/api/secret/reconstructsecret", controller.ReconstructSecret)
	r.GET("/api/secret/getsecretlist", controller.GetSecretList)
	r.GET("/api/secret/handoffsecret", controller.HandoffSecret)

	r.GET("/api/unit/getunitlist", controller.GetUnitList)
	r.GET("/api/unit/getunitlog", controller.GetUnitLog)
	return r
}

//func NewRouter() *gin.Engine {
//	r := gin.Default()
//	r.Use(middleware.Cors())
//	v1 := r.Group("/api")
//	{
//		v1.GET("/ping", controller.Ping)
//
//	}
//	r.POST("/api/newsecret", controller.NewSecret)
//	r.GET("/api/getsecretlist", controller.Getsecretlist)
//	return r
//}

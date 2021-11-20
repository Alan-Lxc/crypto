package main

import (
	"github.com/Alan-Lxc/crypto_contest/dcssweb/common"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/middleware"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/router"
	"github.com/Alan-Lxc/crypto_contest/src/controller"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	common.InitConfig()
	common.InitDB()
	controller.Controller = controller.New()
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r = router.CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run(":8889"))
	//r.Run() // listen and serve on 0.0.0.0:8080
	//fmt.Println("helloworld")
}

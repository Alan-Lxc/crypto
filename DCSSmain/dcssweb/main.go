package main

import (
	"github.com/Alan-Lxc/crypto_contest/dcssweb/common"
	_ "github.com/Alan-Lxc/crypto_contest/dcssweb/controller"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/router"
	_ "github.com/Alan-Lxc/crypto_contest/dcssweb/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	common.InitConfig()
	common.InitDB()

	r := gin.Default()
	r = router.CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}

	panic(r.Run(":8889"))
	//r.Run() // listen and serve on 0.0.0.0:8080
	//fmt.Println("helloworld")
}

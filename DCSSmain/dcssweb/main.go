package main

import (
	"github.com/Alan-Lxc/crypto_contest/dcssweb/handlers"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/router"
)

func main() {
	go handlers.Init_control()
	r := router.NewRouter()
	r.Run("localhost:8080")
}

package main

import "github.com/Alan-Lxc/crypto_contest/dcssweb/router"

func main() {
	r := router.NewRouter()
	r.Run("localhost:8080")
}

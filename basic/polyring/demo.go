package main

import (
	"fmt"
	"github.com/Alan-Lxc/crypto_contest/basic/point"
	"github.com/ncw/gmp"
)

func main() {
	a := gmp.NewInt(64)
	b := gmp.NewInt(87)
	c := point.NewPoint(a, b, nil)
	fmt.Println(c)
}

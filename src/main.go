package main

import (
	"fmt"
	//. "github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/Alan-Lxc/crypto_contest/src/nodes"
	//"github.com/ncw/gmp"
	//"log"
)

func main() {
	fmt.Println("hello world")
	nodes.Demo_test()

}

//func main()  {
//	op1 := FromVec(1, 1, 1, 1, 1, 1)
//	result := NewEmpty()
//
//	err := result.Multiply(op1, op1)
//	if err != nil {
//		fmt.Println("error")
//	}
//	fmt.Println(result.GetAllCoeff())
//	expected := FromVec(1, 2, 3, 4, 5, 6, 5, 4, 3, 2, 1)
//	fmt.Println(expected.GetAllCoeff())
//}

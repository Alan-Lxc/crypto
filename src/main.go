package main

import (
	"fmt"
	. "github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/ncw/gmp"
	"log"
)

//func Demo_test() {
//	var node1 = Node{
//		label:   1,
//		counter: 3,
//	}
//	var node2 = Node{
//		label:   2,
//		counter: 3,
//	}
//	var node3 = Node{
//		label:   3,
//		counter: 3,
//	}
//	var nodes = [3]Node{node1, node2, node3}
//	for i := 0; i < 3; i++ {
//		nodes[i].ClientSharePhase3(nodes)
//	}
//
//}

func main() {
	fmt.Println("hello world")
	//nodes.Demo_test()
	x := make([]*gmp.Int, 2)
	x[0] = gmp.NewInt(1)
	x[1] = gmp.NewInt(2)
	tmp, _ := NewPoly(1)
	product := NewConstant(1)

	// tmp(x) = x - x[i]
	//置一次项系数 为1
	tmp.SetCoeffWithInt(1, 1)
	// note only the first degree points are used
	//得到全排列，product就是l(x),每一位代表一个多项式系数
	for i := 0; i <= 1; i++ {
		tmp.GetPtrtoConstant().Neg(x[i])
		for i := 0; i < 2; i++ {
			num, _ := tmp.GetCoeff(i)
			log.Println(num.Int64())
		}
		//product反复乘以tmp
		product.MulSelf(tmp)
	}
	for i := 0; i < 3; i++ {
		num, _ := product.GetCoeff(i)
		log.Println(num.Int64())
	}

}

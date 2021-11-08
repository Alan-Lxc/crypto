package main

import (
	"fmt"
	"github.com/Alan-Lxc/crypto_contest/src/basic/commitment"
	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/Nik-U/pbc"
	"github.com/ncw/gmp"
	"math/rand"
	"time"
)

func testtt(degree int, devide int, fixedRandState *rand.Rand, p *gmp.Int, dpc commitment.DLPolyCommit) {
	polyy, _ := poly.NewRand(degree, fixedRandState, p)
	c := dpc.NewG1()
	dpc.Commit(c, polyy)
	dpc.Commit(c, polyy)
	if devide != 1 {
		testtt(devide, devide/2, fixedRandState, p, dpc)
		testtt(devide, devide/2, fixedRandState, p, dpc)
	}
}
func testtt2(degree int, devide int, c *pbc.Element, p *gmp.Int, dpc commitment.DLPolyCommit) {
	dpc.VerifyEval(c, p, p, c)
	if devide != 1 {
		testtt2(devide, devide/2, c, p, dpc)
	}
}
func testtt3(degree int, devide int, fixedRandState *rand.Rand, p *gmp.Int) {
	poly2, _ := poly.NewRand(degree/2, fixedRandState, p)
	ploy3, _ := poly.NewPoly(0)
	//poly1, _ := poly.NewRand(degree, fixedRandState, p)
	//ploy4, _ :=poly.NewPoly(0)
	//poly.DivMod(poly1,poly2,p,&ploy3,&ploy4)
	ploy3.Multiply(poly2, poly2)
	ploy3.Multiply(poly2, poly2)
	ploy3.Multiply(poly2, poly2)
	ploy3.Multiply(poly2, poly2)
	if devide != 1 {
		testtt3(devide, devide/2, fixedRandState, p)
		testtt3(devide, devide/2, fixedRandState, p)
	}
}
func testAMT(degree int) {
	fixedRandState := rand.New(rand.NewSource(int64(3)))
	p := gmp.NewInt(0)
	p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)
	dpc := commitment.DLPolyCommit{}
	dpc.SetupFix(degree*2 + 1)

	polyy, _ := poly.NewRand(degree, fixedRandState, p)
	c := dpc.NewG1()
	dpc.Commit(c, polyy)
	s1 := time.Now()
	s2 := time.Now()
	e1 := time.Now()
	e2 := time.Now()
	testtt(degree, degree*2+1, fixedRandState, p, dpc)
	testtt3(degree, degree*2+1, fixedRandState, p)
	for i := 0; i < degree*2+1; i++ {
		testtt2(degree, degree*2+1, c, p, dpc)
	}
	s2 = time.Now()
	e1 = time.Now()
	for i := 0; i < degree*2+1; i++ {
		dpc.CreateWitness(c, polyy, gmp.NewInt(int64(i)))
		dpc.VerifyEval(c, p, p, c)
	}

	e2 = time.Now()
	//fmt.Println(s1, s2, '\n', e1, e2)
	fmt.Println(degree, e2.Sub(s2).Nanoseconds(), e1.Sub(s1).Nanoseconds())
}

func main() {

	testAMT(300)

}

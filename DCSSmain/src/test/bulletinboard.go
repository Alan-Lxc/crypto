package main

import (
	"flag"
	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/Alan-Lxc/crypto_contest/src/bulletboard"
	"github.com/ncw/gmp"
	"math/rand"
)

func main() {
	cnt := flag.Int("c", 2, "Enter number of nodes")
	degree := flag.Int("d", 1, "Enter the polynomial degree")
	metadataPath := flag.String("path", "/mpss/metadata", "Enter the metadata path")
	aws := flag.Bool("aws", false, "if test on real aws")
	flag.Parse()

	fixedRandState := rand.New(rand.NewSource(int64(3)))
	p := gmp.NewInt(0)
	p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)
	tmp := gmp.NewInt(0)
	tmp.SetString("1234567899876543210", 10)
	polyy, _ := poly.NewRand(*degree, fixedRandState, p)
	polyy.SetCoeffWithGmp(0, tmp)
	polyyy := make([]poly.Poly, *cnt)
	for i := 0; i < *cnt; i++ {
		y := gmp.NewInt(0)

		polyyy[i], _ = poly.NewRand(*degree, fixedRandState, p)
		polyyy[i].SetCoeffWithGmp(0, y)
	}
	bb, _ := bulletboard.New(*degree, *cnt, *metadataPath, polyyy)
	bb.Serve(*aws)
}

package mainclient

import (
	"../bulletboard"
	"../clock"
	"flag"
	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/Alan-Lxc/crypto_contest/src/nodes"
	"github.com/ncw/gmp"
	"math/rand"
	"time"
)

func Initial(s0 string, degree int, counter int, metadataPath string) {
	//bulletboard
	aws := flag.Bool("aws", false, "if test on real aws")
	flag.Parse()
	bb, _ := bulletboard.New(degree, counter, metadataPath)
	bb.Serve(*aws)

	GeneratePoly(s0, degree, counter, metadataPath)
}

func Handoff(metadataPath string) {

	clock, _ := clock.New(metadataPath)
	clock.Connect()
	clock.ClientStartEpoch()
}

func GetAnswer(metadataPath string) string {
	clock, _ := clock.New(metadataPath)
	clock.Connect()
	ans := clock.ClientGetAnswer()
	return ans
}

func GeneratePoly(s0 string, degree int, counter int, metadataPath string) {

	//randState := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	fixedRandState := rand.New(rand.NewSource(int64(3)))
	p := gmp.NewInt(0)
	p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)
	tmp := gmp.NewInt(0)
	tmp.SetString(s0, 10)

	poly, _ := poly.NewRand(degree, fixedRandState, p)
	poly.SetCoeffWithGmp(0, tmp)

	//nn := make([]nodes.Node, counter)
	//????
	for i := 0; i < counter; i++ {
		x := int32(i)
		y := gmp.NewInt(0)
		poly.EvalMod(gmp.NewInt(int64(x)), p, y)
		n, _ := nodes.New(degree, i, counter, metadataPath, p, y)
		n.Service()
	}

}

func Gets0() string {

}

func main() {
}

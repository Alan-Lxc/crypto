package nodes

import (
	"errors"
	"github.com/Alan-Lxc/crypto_contest/src/basic/point"
	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/ncw/gmp"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Node struct {
	//Label of Node
	Label int
	//Total number of Nodes
	Total int
	//Degree of polynomial
	Degree int
	//the polynomial was set on Z_p
	P *gmp.Int
	// Rand source
	Randstate *rand.Rand
	//To store the point(shares) sent from other node
	RecPoint []*point.Point
	//To recode the share that have already received
	RecCounter int
	//The poly reconstructed with the shares
	RecPoly poly.Poly
}

func (node *Node) GetLabel() int {
	if node != nil {
		return node.Label
	} else {
		return 0
	}
}

func New(degree, label, counter int, logPath string) (Node, error) {
	if label < 0 {
		return Node{}, errors.New("Label must be a non-negative number!")
	}
	file, _ := os.Create(logPath + "/log" + strconv.Itoa(label))
	defer file.Close()
	if counter < 0 {
		return Node{}, errors.New("Counter must be a non-negative number!")
	}
	randState := rand.New(rand.NewSource(time.Now().Local().UnixNano()))
	p := gmp.NewInt(0)
	p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)

	return Node{
		Label:      label,
		Total:      counter,
		Degree:     degree,
		P:          p,
		Randstate:  randState,
		RecPoint:   nil,
		RecCounter: 0,
		RecPoly:    poly.Poly{},
	}, nil
}

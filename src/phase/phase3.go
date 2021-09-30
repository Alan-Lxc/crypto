package phase

import (
	"github.com/Alan-Lxc/crypto_contest/src/basic/point"
	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/ncw/gmp"
	"log"
	"math/rand"
	"sync"
)

//written by kzl

type message struct {
	Index int32
	X     int32
	Y     []byte
}

type Node struct {
	//Label of Node
	label int
	//Total number of Nodes
	total int
	//Degree of polynomial
	degree int
	//the polynomial was set on Z_p
	p *gmp.Int
	// Rand source
	randState *rand.Rand

	//Polynomial State
	secretShares []*point.Point

	//To store the point(shares) sent from other node
	recPoint []*point.Point
	//To recode the share that have already received
	recCounter int
	//The poly reconstructed with the shares
	recPoly poly.Poly
	//Mutex to control
	mutex sync.Mutex
	//Secret shares of node p(a0,y)
	secretShare []*point.Point
	// Proactivization Polynomial
	proPoly *poly.Poly
	//New Polynomials after phase3
	newPoly poly.Poly
	//Number of Nodes
	counter int

	// Counter for New Secret Shares
	shareCnt *int
}

func (node *Node) Phase3() {
	node.newPoly.Add(*node.recPoly, *node.proPoly)

	for i := 1; i < node.counter; i++ {
		eval := gmp.NewInt(0)
		node.newPoly.EvalMod(gmp.NewInt(int64(i+1)), node.p, eval)

		if i != node.label-1 {
			log.Printf("node %d send point message to node %d in phase 3", node.label, i+1)
			msg := &message{
				Index: int32(node.label),
				X:     int32(i + 1),
				Y:     eval.Bytes(),
			}
			go func(i int, msg *message) {}(i, msg)
		} else {
			node.secretShares[i].Y.Set(eval)
			node.mutex.Lock()

			*node.shareCnt = *node.shareCnt + 1
			flag := (*node.shareCnt == node.counter)
			node.mutex.Unlock()
			if flag {
				*node.shareCnt = 0

			}
		}

	}
}

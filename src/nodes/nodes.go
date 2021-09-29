package nodes

import (
	"github.com/Alan-Lxc/crypto_contest/src/basic/point"
	"github.com/ncw/gmp"
	"log"
	"math/rand"
)

type Node struct {
	//label of Node
	label int
	//Total number of Nodes
	total int
	//degree of polynomial
	degree int
	//the polynomial was set on Z_p
	p *gmp.Int
	// Rand source
	randstate *rand.Rand
	//To store the point(shares) sent from other node
	recPoint []*point.Point
	//To recode the share that have already received
	recCounter int
}

func (node Node) GetMsgFromNode(pointmsg point.Pointmsg) {
	index := pointmsg.GetIndex()
	log.Println("Phase 1 :[Node %d] receive point message from [Node %d]", node.label, index)
	p := pointmsg.GetPoint()
	//Receive the point and store
	node.recPoint[node.recCounter] = p
	node.recCounter += 1

	if node.recCounter == node.total {
		node.recCounter = 0
		node.Phase1()
	}
}
func (node Node) Phase1() {
	log.Printf("[Node %d] now start phase1", node.label)

}

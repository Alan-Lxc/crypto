package phase

//
//import (
//	"github.com/Alan-Lxc/crypto_contest/src/basic/point"
//	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
//	"github.com/ncw/gmp"
//	"log"
//	"math/rand"
//	"sync"
//)
//
////written by kzl
//
//type Node struct {
//	//Label of Node
//	label int
//	//Total number of Nodes
//	total int
//	//Degree of polynomial
//	degree int
//	//the polynomial was set on Z_p
//	p *gmp.Int
//	// Rand source
//	randState *rand.Rand
//
//	//Polynomial State
//	secretShares []*point.Point
//
//	//To store the point(shares) sent from other node
//	recPoint []*point.Point
//	//To recode the share that have already received
//	recCounter int
//	//The poly reconstructed with the shares
//	recPoly poly.Poly
//	//Mutex to control
//	mutex sync.Mutex
//	//Secret shares of node p(a0,y)
//	secretShare []*point.Point
//	// Proactivization Polynomial
//	proPoly *poly.Poly
//	//New Polynomials after phase3
//	newPoly poly.Poly
//	//Number of Nodes
//	counter int
//
//	// Counter for New Secret Shares
//	shareCnt *int
//
//	//the other nodes in the committee
//	Client []Node
//}
//
//type message struct {
//	Index int32
//	X     int32
//	Y     []byte
//}
//
//func (msg *message) getY() (Y []byte) {
//	return msg.Y
//}
//
//func (msg *message) getIndex() (Index int32) {
//	return msg.Index
//}
//
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
//
////第三部分的分享， 是全份额重建
////通过新的多项式B（x，i），对包括自己的每一个node发送B（i，i）
////最后重建多项式
//
////重建secretShare
//func (node *Node) SharePhase3(msg message) error {
//	index := msg.getIndex()
//	Y := msg.getY()
//
//	node.secretShares[index-1].Y.SetBytes(Y)
//	*node.shareCnt = *node.shareCnt + 1
//	flag := (*node.shareCnt == node.counter)
//	if flag {
//		*node.shareCnt = 0
//
//	}
//	return nil
//
//}
//func (node *Node) ClientSharePhase3(nodes [3]Node) {
//	node.newPoly.Add(node.recPoly, node.proPoly)
//
//	for i := 1; i < node.counter; i++ {
//		value := gmp.NewInt(0)
//		node.newPoly.EvalMod(gmp.NewInt(int64(i+1)), node.p, value)
//
//		if i != node.label-1 {
//			log.Printf("node %d send point message to node %d in phase 3", node.label, i+1)
//			msg := &message{
//				Index: int32(node.label),
//				X:     int32(i + 1),
//				Y:     value.Bytes(),
//			}
//			//把消息发送给不同的节点
//			nodes[i].SharePhase3(*msg)
//
//		} else {
//			node.secretShares[i].Y.Set(value)
//		}
//
//	}
//}

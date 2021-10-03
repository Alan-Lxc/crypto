package phase

//
//import (
//	"context"
//	//"github.com/Alan-Lxc/crypto_contest/src/basic/point"
//	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
//	//"github.com/Nik-U/pbc"
//	"github.com/ncw/gmp"
//	"log"
//	//"math/big"
//	"math/rand"
//	"sync"
//	//"time"
//)
//
////type Node struct {
////	//Label of Node
////	label int
////	//Total number of Nodes
////	total int
////	//Degree of polynomial
////	degree int
////	//the polynomial was set on Z_p
////	p *gmp.Int
////	// Rand source
////	randState *rand.Rand
////	//To store the point(shares) sent from other node
////	recPoint []*point.Point
////	//To recode the share that have already received
////	recCounter int
////	//The poly reconstructed with the shares
////	recPoly poly.Poly
////	//Mutex to control
////	mutex sync.Mutex
////	//Secret shares of node p(a0,y)
////	secretShare []*point.Point
////
////	// Lagrange Coefficient
////	lambda []*gmp.Int
////	//	0Shares
////	_0Shares     []*gmp.Int
////	_0ShareSum   *gmp.Int
////	_0ShareCount *int
////
////	//poly Q
////	proPoly poly.Poly
////}
//
//func (node *Node) ClientSharePhase2() {
//	// Generate Random Numbers
//	for i := 0; i < node.total-1; i++ {
//		node._0Shares[i].Rand(node.randState, gmp.NewInt(10))
//		node._0Shares[node.total-1].Sub(node._0Shares[node.total-1], gmp.NewInt(0).Mul(node._0Shares[i], node.lambda[i]))
//	}
//	//0-share means the S
//	node._0Shares[node.total-1].Mod(node._0Shares[node.total-1], node.p)
//	node._0Shares[node.total-1].Mul(node._0Shares[node.total-1], gmp.NewInt(0).ModInverse(node.lambda[node.total-1], node.p))
//	node._0Shares[node.total-1].Mod(node._0Shares[node.total-1], node.p)
//
//	//to get sum for \sum_counter
//	node.mutex.Lock()
//	node._0ShareSum.Add(node._0ShareSum, node._0Shares[node.label-1])
//	*node._0ShareCount = *node._0ShareCount + 1
//	_0shareSumFinish := *node._0ShareCount == node.total
//	node.mutex.Unlock()
//
//	if _0shareSumFinish {
//		*node._0ShareCount = 0
//		node._0ShareSum.Mod(node._0ShareSum, node.p)
//
//		//get a rand poly_tmp with 0-share
//		//rand a poly_tmp polynomial
//		polyTmp, _ := poly.NewRand(node.degree, node.randState, node.p)
//		polyTmp.SetCoeffWithGmp(0, node._0ShareSum)
//		node.proPoly.ResetTo(polyTmp)
//
//		//node.ClientWritePhase2()
//	}
//
//	// share 0-share
//	var wg sync.WaitGroup
//	for i := 0; i < node.total; i++ {
//		if i != node.label-1 {
//			log.Printf("[node %d] send message to [node %d] in phase 2", node.label, i+1)
//			msg := &pb.ZeroMsg{
//				Index: int32(node.label),
//				Share: node._0Shares[i].Bytes(),
//			}
//			wg.Add(1)
//			go func(i int, msg *pb.ZeroMsg) {
//				defer wg.Done()
//				ctx, cancel := context.WithCancel(context.Background())
//				defer cancel()
//				node.nClient[i].SharePhase2(ctx, msg)
//			}(i, msg)
//		}
//	}
//	wg.Wait()
//}
//func (node *Node) SharePhase2(ctx context.Context, msg *pb.ZeroMsg) (*pb.AckMsg, error) {
//	*node.totMsgSize = *node.totMsgSize + proto.Size(msg)
//	index := msg.GetIndex()
//	log.Printf("[node %d] receive zero message from [node %d] in phase 2", node.total, index)
//	inter := gmp.NewInt(0)
//	inter.SetBytes(msg.GetShare())
//
//	//to get sum for \sum_counter
//	node.mutex.Lock()
//	node._0ShareSum.Add(node._0ShareSum, node._0Shares[node.label-1])
//	*node._0ShareCount = *node._0ShareCount + 1
//	_0shareSumFinish := *node._0ShareCount == node.total
//	node.mutex.Unlock()
//
//	if _0shareSumFinish {
//		*node._0ShareCount = 0
//		node._0ShareSum.Mod(node._0ShareSum, node.p)
//
//		//get a rand polyTmp with 0-share
//		//rand a polyTm
//		polyTmp, _ := poly.NewRand(node.degree, node.randState, node.p)
//		polyTmp.SetCoeffWithGmp(0, node._0ShareSum)
//
//		node.proPoly.ResetTo(polyTmp)
//		//node.ClientWritePhase2()
//	}
//	return &pb.AckMsg{}, nil
//}

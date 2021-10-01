package phase

import (
	"context"
	"github.com/Nik-U/pbc"
	"github.com/ncw/gmp"
	"log"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

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
//}
//parameter should have t,
//and optimistic or not

//consider that for the optimistic path,
//just add Q(x),in other words,we can just add the point to each poly point
//which cost is n^2
//for pessimistic path ,t -> 2t, than add Q(x) ,than 2t->t
//in this way,t->2t & 2t->t cost n^3

//1st step ,generate Q(x)
//observe std method
//means that each committee contain 2t poly point
//not really generate Q,but use the seed for each committee,just calc the add-value for 2t point
//so this step is work during step 2, don't need to accomplish separately

//2nd step ,read ploy from phase1,add Q(x) using \Polynomial or\ poly point(uncertain)
//then write back, using database ,tcp ,etc.
//more
//  1st judge optimistic or not
//  2nd get a committee,get the poly point t or 2t
//  3rd add add-value to each point,return
//  interpolation is not need in this phase when calc the s0,then using interpolation
//  more for 3rd
//  need to judge t or 2t,then according to the label for node, then using the seed can calc the add-value rapidly
//  ep. B(x,2) has 3 poly point 1,2,3
//  Q use seed for committee No.2 can get

// Phase2 written by Gary

func (node *Node) ClientSharePhase2() {
	// Generate Random Numbers
	for i := 0; i < node.counter-1; i++ {
		node.zeroShares[i].Rand(node.randState, gmp.NewInt(10))
		inter := gmp.NewInt(0)
		inter.Mul(node.zeroShares[i], node.lambda[i])
		node.zeroShares[node.counter-1].Sub(node.zeroShares[node.counter-1], inter)
	}
	//0-share means the S
	node.zeroShares[node.counter-1].Mod(node.zeroShares[node.counter-1], node.p)
	inter := gmp.NewInt(0)
	inter.ModInverse(node.lambda[node.counter-1], node.p)
	node.zeroShares[node.counter-1].Mul(node.zeroShares[node.counter-1], inter)
	node.zeroShares[node.counter-1].Mod(node.zeroShares[node.counter-1], node.p)

	//to get sum for \sum_counter
	node.mutex.Lock()
	node.zeroShare.Add(node.zeroShare, node.zeroShares[node.label-1])
	*node.zeroCnt = *node.zeroCnt + 1
	flag := *node.zeroCnt == node.counter
	node.mutex.Unlock()

	if flag {
		*node.zeroCnt = 0
		node.zeroShare.Mod(node.zeroShare, node.p)
		//get a rand poly with 0-share
		//rand a poly polynomial
		poly, _ := polyring.NewRand(node.degree, node.randState, node.p)
		poly.SetCoefficient(0, 0)
		poly.SetCoefficientBig(0, node.zeroShare)

		node.proPoly.ResetTo(poly.DeepCopy())
		//node.ClientWritePhase2()
	}

	// share 0-share
	var wg sync.WaitGroup
	for i := 0; i < node.counter; i++ {
		if i != node.label-1 {
			log.Printf("[node %d] send message to [node %d] in phase 2", node.label, i+1)
			msg := &pb.ZeroMsg{
				Index: int32(node.label),
				Share: node.zeroShares[i].Bytes(),
			}
			wg.Add(1)
			go func(i int, msg *pb.ZeroMsg) {
				defer wg.Done()
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				node.nClient[i].SharePhase2(ctx, msg)
			}(i, msg)
		}
	}
	wg.Wait()
}

// Share Phase 2
// The server function which takes the sent message of zero shares and sum them up to get the final share and generate the proactivization polynomial according to the zero share. It then calls ClientWritePhase2 to write the commitment of zeroshare, zeropolynomial and the witness at zero on the bulletinboard.
func (node *Node) SharePhase2(ctx context.Context, msg *pb.ZeroMsg) (*pb.AckMsg, error) {
	*node.totMsgSize = *node.totMsgSize + proto.Size(msg)
	index := msg.GetIndex()
	log.Printf("[node %d] receive zero message from [node %d] in phase 2", node.label, index)
	inter := gmp.NewInt(0)
	inter.SetBytes(msg.GetShare())

	node.mutex.Lock()
	node.zeroShare.Add(node.zeroShare, inter)
	*node.zeroCnt = *node.zeroCnt + 1
	flag := *node.zeroCnt == node.counter
	node.mutex.Unlock()

	if flag {
		*node.zeroCnt = 0
		node.zeroShare.Mod(node.zeroShare, node.p)

		//get a rand poly with 0-share
		poly, _ := polyring.NewRand(node.degree, node.randState, node.p)
		poly.SetCoefficient(0, 0)
		poly.SetCoefficientBig(0, node.zeroShare)

		node.proPoly.ResetTo(poly.DeepCopy())
		//node.ClientWritePhase2()
	}
	return &pb.AckMsg{}, nil
}

//
//func (node *Node) ClientWritePhase2() {
//  log.Printf("[node %d] write bulletinboard in phase 2", node.label)
//  ctx, cancel := context.WithCancel(context.Background())
//  defer cancel()
//  msg := &pb.Cmt2Msg{
//    Index:       int32(node.label),
//    Sharecmt:    node.zeroShareCmt.CompressedBytes(),
//    Polycmt:     node.zeroPolyCmt.CompressedBytes(),
//    Zerowitness: node.zeroPolyWit.CompressedBytes(),
//  }
//  node.bClient.WritePhase2(ctx, msg)
//}
//
//// After the bulletinboard has received the writing of all nodes, it will start a client call to this function telling the nodes to read the commitment on it.
//func (node *Node) StartVerifPhase2(ctx context.Context, in *pb.EmptyMsg) (*pb.AckMsg, error) {
//  log.Printf("[node %d] start verification in phase 2", node.label)
//  node.ClientReadPhase2()
//  return &pb.AckMsg{}, nil
//}
//
//// Read from bulletinboard and does the verification in phase 2.
//func (node *Node) ClientReadPhase2() {
//  log.Printf("[node %d] read bulletinboard in phase 2", node.label)
//  ctx, cancel := context.WithCancel(context.Background())
//  defer cancel()
//  stream, err := node.bClient.ReadPhase2(ctx, &pb.EmptyMsg{})
//  if err != nil {
//    log.Fatalf("client failed to read phase2: %v", err)
//  }
//  for i := 0; i < node.counter; i++ {
//    msg, err := stream.Recv()
//    *node.totMsgSize = *node.totMsgSize + proto.Size(msg)
//    if err != nil {
//      log.Fatalf("client failed to receive in read phase1: %v", err)
//    }
//    index := msg.GetIndex()
//    sharecmt := msg.GetSharecmt()
//    polycmt := msg.GetPolycmt()
//    zerowitness := msg.GetZerowitness()
//    node.zerosumShareCmt[index-1].SetCompressedBytes(sharecmt)
//    inter := node.dpc.NewG1()
//    inter.SetString(node.zerosumShareCmt[index-1].String(), 10)
//    node.zerosumPolyCmt[index-1].SetCompressedBytes(polycmt)
//    node.midPolyCmt[index-1].Mul(inter, node.zerosumPolyCmt[index-1])
//    node.zerosumPolyWit[index-1].SetCompressedBytes(zerowitness)
//  }
//  exponentSum := node.dc.NewG1()
//  exponentSum.Set1()
//  for i := 0; i < node.counter; i++ {
//    lambda := big.NewInt(0)
//    lambda.SetString(node.lambda[i].String(), 10)
//    tmp := node.dc.NewG1()
//    tmp.PowBig(node.zerosumShareCmt[i], lambda)
//    // log.Printf("label: %d #share %d\nlambda %s\nzeroshareCmt %s\ntmp %s", node.label, i+1, lambda.String(), node.zerosumShareCmt[i].String(), tmp.String())
//    exponentSum.Mul(exponentSum, tmp)
//  }
//  // log.Printf("%d exponentSum: %s", node.label, exponentSum.String())
//  if !exponentSum.Is1() {
//    panic("Proactivization Verification 1 failed")
//  }
//  flag := true
//  for i := 0; i < node.counter; i++ {
//    if !node.dpc.VerifyEval(node.zerosumPolyCmt[i], gmp.NewInt(0), gmp.NewInt(0), node.zerosumPolyWit[i]) {
//      flag = false
//    }
//  }
//  if !flag {
//    panic("Proactivization Verification 2 failed")
//  }
//  *node.e2 = time.Now()
//  *node.s3 = time.Now()
//  node.ClientSharePhase3()
//}

package nodes

import (
	"context"
	"errors"
	"fmt"
	. "github.com/Alan-Lxc/crypto_contest/src/basic/getprime"
	"github.com/Alan-Lxc/crypto_contest/src/basic/interpolation"
	"github.com/Alan-Lxc/crypto_contest/src/basic/point"
	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/ncw/gmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type Node struct {
	//Label of Node
	label int
	//Total number of Nodes
	counter int
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
	recPoly *poly.Poly
	//Mutex to control
	mutex sync.Mutex

	// Proactivization Polynomial
	proPoly *poly.Poly
	//New Polynomials after phase3
	newPoly poly.Poly

	// Lagrange Coefficient
	lambda []*gmp.Int
	//	0Shares
	_0Shares     []*gmp.Int
	_0ShareSum   *gmp.Int
	_0ShareCount *int

	// Counter for New Secret Shares
	shareCnt *int

	//the other nodes in the committee
	Client []*Node

	//Secret shares of node p(a0,y)
	secretShare []*point.Point
	//IP_address of node
	ipAddress []string
	//board IP address
	ipOfBoard string
	//clientconn
	clinetConn []*grpc.ClientConn
	//nodeService
	//nodeService []pb.nodeservice
	//boardconn
	boardConn *grpc.ClientConn
	//boardService
	//boardService pb.boardservice
}

func (node *Node) GetLabel() int {
	if node != nil {
		return node.label
	} else {
		return 0
	}

}

func (node *Node) connect(ptrs []*Node) {
	for i := 0; i < node.counter; i++ {

		if i != node.label-1 {
			node.Client[i] = ptrs[i]
		}
	}
}

//phase1
func (node *Node) GetMsgFromNode(pointmsg point.Pointmsg) {
	index := pointmsg.GetIndex()
	log.Printf("Phase 1 :[Node %d] receive point message from [Node %d]", node.label, index)
	p := pointmsg.GetPoint()
	//Receive the point and store
	node.mutex.Lock()
	node.recPoint[node.recCounter] = p
	fmt.Println(p)
	node.recCounter += 1
	flag := (node.recCounter == node.counter)
	node.mutex.Unlock()
	if flag {
		node.recCounter = 0
		node.Phase1()
	}
}
func (node *Node) SendMsgToNode(nodes [3]Node) {
	p := point.Point{
		X: node.secretShares[node.label-1].X,
		Y: node.secretShares[node.label-1].Y,
		//PolyWit: node.secretShare[node.label-1].PolyWit,
	}
	node.mutex.Lock()
	node.recPoint[node.recCounter] = &p
	node.recCounter += 1
	flag := (node.recCounter == node.counter)
	node.mutex.Unlock()
	if flag {
		node.recCounter = 0
		node.Phase1()
	}
	for i := 0; i < node.counter; i++ {
		if i != node.label-1 {
			log.Printf("[Node %d] send point message to [Node %d]", node.label, i+1)
			msg := point.Pointmsg{}
			msg.SetIndex(node.label)
			msg.SetPoint(node.secretShares[i])
			(*node.Client[i]).GetMsgFromNode(msg)

			//node.sent(msg)
		}
	}

}
func (node *Node) Phase1() {
	log.Printf("[Node %d] now start phase1", node.label)
	x_point := make([]*gmp.Int, node.degree+1)
	y_point := make([]*gmp.Int, node.degree+1)
	for i := 0; i <= node.degree; i++ {
		point := node.recPoint[i]
		//x_point = append(x_point, gmp.NewInt(int64(point.X)))
		x_point[i] = gmp.NewInt(int64(point.X))
		y_point[i] = point.Y
		//y_point = append(y_point, point.Y)
	}
	p, err := interpolation.LagrangeInterpolate(node.degree, x_point, y_point, node.p)
	if err != nil {
		for i := 0; i < 3; i++ {
			num, _ := p.GetCoeff(i)
			log.Println(num.Int64())
		}
		//for i := 0; i <= node.degree; i++ {
		//	log.Print(x_point[i])
		//	log.Print(y_point[i])
		//}
		log.Print(err)
		panic("Interpolation failed")
	}
	node.recPoly = &p
	fmt.Printf("Interpolation finished")
	//node.Phase2()

}

////phase1
//func GetMsgFromNode(nodeReceive Node, nodeSend Node, pointmsg point.Pointmsg) {
//	index := pointmsg.GetIndex()
//	log.Printf("Phase 1 :[Node %d] receive point message from [Node %d]\n", nodeReceive.GetLabel, index)
//	p := pointmsg.GetPoint()
//	//Receive the point and store
//	nodeReceive.recPoint[nodeReceive.recCounter] = p
//	nodeReceive.recCounter += 1
//
//	if nodeReceive.recCounter == nodeReceive.counter {
//		nodeReceive.recCounter = 0
//		Phase1(nodeReceive)
//	}
//}
//func Phase1(node Node) {
//	log.Printf("[Node %d] now start phase1", node.label)
//	x_point := make([]*gmp.Int, node.degree+1)
//	y_point := make([]*gmp.Int, node.degree+1)
//	for i := 0; i <= node.degree; i++ {
//		x_point[i] = node.recPoint[i].X
//		y_point[i] = node.recPoint[i].Y
//	}
//	p, err := interpolation.LagrangeInterpolate(node.Degree, x_point, y_point, node.P)
//	if err != nil {
//		for i := 0; i <= node.Degree; i++ {
//			log.Print(x_point[i])
//			log.Print(y_point[i])
//		}
//		log.Print(err)
//		panic("Interpolation failed")
//	}
//	node.RecPoly = p
//	log.Printf("Interpolation finished")
//	//Phase2(node)
//}
type ZeroMsg struct {
	Index int32
	Share []byte
}

func (msg *ZeroMsg) GetIndex() int32 {
	return msg.Index
}
func (msg *ZeroMsg) GetShare() []byte {
	return msg.Share
}

//phase2
func (node *Node) ClientSharePhase2() {
	// Generate Random Numbers
	for i := 0; i < node.counter-1; i++ {
		node._0Shares[i].Rand(node.randState, gmp.NewInt(10))
		node._0Shares[node.counter-1].Sub(node._0Shares[node.counter-1], gmp.NewInt(0).Mul(node._0Shares[i], node.lambda[i]))
	}
	//0-share means the S
	node._0Shares[node.counter-1].Mod(node._0Shares[node.counter-1], node.p)
	node._0Shares[node.counter-1].Mul(node._0Shares[node.counter-1], gmp.NewInt(0).ModInverse(node.lambda[node.counter-1], node.p))
	node._0Shares[node.counter-1].Mod(node._0Shares[node.counter-1], node.p)

	//to get sum for \sum_counter
	node.mutex.Lock()
	node._0ShareSum.Add(node._0ShareSum, node._0Shares[node.label-1])
	*node._0ShareCount = *node._0ShareCount + 1
	_0shareSumFinish := *node._0ShareCount == node.counter
	node.mutex.Unlock()

	if _0shareSumFinish {
		*node._0ShareCount = 0
		node._0ShareSum.Mod(node._0ShareSum, node.p)

		//get a rand poly_tmp with 0-share
		//rand a poly_tmp polynomial
		polyTmp, _ := poly.NewRand(node.degree, node.randState, node.p)
		polyTmp.SetCoeffWithGmp(0, node._0ShareSum)
		node.proPoly.ResetTo(polyTmp.DeepCopy())

		//node.ClientWritePhase2()
	}

	// share 0-share
	var wg sync.WaitGroup
	for i := 0; i < node.counter; i++ {
		if i != node.label-1 {
			log.Printf("[node %d] send message to [node %d] in phase 2", node.label, i+1)
			msg := ZeroMsg{
				Index: int32(node.label),
				Share: node._0Shares[i].Bytes(),
			}
			wg.Add(1)
			go func(i int, msg ZeroMsg) {
				defer wg.Done()
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				node.Client[i].SharePhase2(ctx, msg)
			}(i, msg)
		}
	}
	wg.Wait()
}

func (node *Node) SharePhase2(ctx context.Context, msg ZeroMsg) error {
	//*node.totMsgSize = *node.totMsgSize + proto.Size(msg)
	index := msg.GetIndex()
	log.Printf("[node %d] receive zero message from [node %d] in phase 2", node.counter, index)
	inter := gmp.NewInt(0)
	inter.SetBytes(msg.GetShare())

	//to get sum for \sum_counter
	node.mutex.Lock()
	node._0ShareSum.Add(node._0ShareSum, node._0Shares[node.label-1])
	*node._0ShareCount = *node._0ShareCount + 1
	_0shareSumFinish := *node._0ShareCount == node.counter
	node.mutex.Unlock()

	if _0shareSumFinish {
		*node._0ShareCount = 0
		node._0ShareSum.Mod(node._0ShareSum, node.p)

		//get a rand polyTmp with 0-share
		//rand a polyTm
		polyTmp, _ := poly.NewRand(node.degree, node.randState, node.p)
		polyTmp.SetCoeffWithGmp(0, node._0ShareSum)

		node.proPoly.ResetTo(polyTmp)
		//node.ClientWritePhase2()
	}
	return nil
}

//phase3
type message struct {
	Index int32
	X     int32
	Y     []byte
}

func (msg *message) getY() (Y []byte) {
	return msg.Y
}

func (msg *message) getIndex() (Index int32) {
	return msg.Index
}

func Demo_test() {
	var nodes [3]Node
	var modp *gmp.Int
	modp = GetPrime(256)
	for i := 0; i < 3; i++ {
		nodes[i], _ = New(1, i+1, 3, "/home/kzl/Desktop", modp)
	}

	for i := 0; i < 3; i++ {

		nodes[i].connect([]*Node{&nodes[0], &nodes[1], &nodes[2]})
	}
	for i := 0; i < 3; i++ {
		nodes[i].SendMsgToNode(nodes)
	}

	for i := 0; i < 3; i++ {
		//nodes[i].ClientSharePhase2()
	}

	for i := 0; i < 3; i++ {
		//nodes[i].ClientSharePhase3()
	}

}

func (node *Node) NodeConnect() {
	boradConn, err := grpc.Dial(node.ipOfBoard, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Fail to connect board:%v", err)
	}
	node.boardConn = boradConn
	//node.boardService = pb
	for i := 0; i < node.counter; i++ {
		if i != node.label-1 {
			clientconn, err := grpc.Dial(node.ipAddress[i], grpc.WithInsecure())
			if err != nil {
				log.Fatalf("[Node %d] Fail to connect to other node:%v", node.label, err)
			}
			node.clinetConn[i] = clientconn
			//node.nodeService[i] = pb
		}
	}
}
func (node *Node) Service() {
	port := node.ipAddress[node.label-1]
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("[Node %d] fail to listen:%v", node.label, err)
	}
	server := grpc.NewServer()
	//pb.re	(s,node)
	reflection.Register(server)
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("[Node %d] fail to provide service", node.label)
	}
	log.Printf("[Node %d] now serve on %s", node.label, node.ipAddress[node.label-1])
}

func New(degree, label, counter int, logPath string, modp *gmp.Int) (Node, error) {
	if label < 0 {
		return Node{}, errors.New("Label must be a non-negative number!")
	}
	file, _ := os.Create(logPath + "/log" + strconv.Itoa(label))
	defer file.Close()
	if counter < 0 {
		return Node{}, errors.New("Counter must be a non-negative number!")
	}
	randState := rand.New(rand.NewSource(time.Now().Local().UnixNano()))
	//p := gmp.NewInt(0)
	//Maybe We can generate a big prime?
	//p.SetString(rand.)
	//p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)

	fixedRandState := rand.New(rand.NewSource(int64(3)))

	secretShares := make([]*point.Point, counter)
	tmpPoly, err := poly.NewRand(degree, fixedRandState, modp)
	for i := 0; i < counter; i++ {
		if err != nil {
			panic("Error initializing random tmpPoly")
		}
		x := int32(label)
		y := gmp.NewInt(0)
		//w := dpc.NewG1()
		tmpPoly.EvalMod(gmp.NewInt(int64(x)), modp, y)
		//dpc.CreateWitness(w, tmpPoly, gmp.NewInt(int64(x)))

		secretShares[i] = point.NewPoint(x, y)
	}

	proPoly, _ := poly.NewPoly(degree)
	lambda := make([]*gmp.Int, counter)
	// Calculate Lagrange Interpolation
	denominator := poly.NewConstant(1)
	tmp, _ := poly.NewPoly(1)
	tmp.SetCoeffWithInt(1, 1)
	for i := 0; i < counter; i++ {
		tmp.GetPtrtoConstant().Neg(gmp.NewInt(int64(i + 1)))
		denominator.MulSelf(tmp)
	}
	for i := 0; i < counter; i++ {
		lambda[i] = gmp.NewInt(0)
		deno, _ := poly.NewPoly(0)
		tmp.GetPtrtoConstant().Neg(gmp.NewInt(int64(i + 1)))

		deno.Divide(denominator, tmp)
		deno.EvalMod(gmp.NewInt(0), modp, lambda[i])
		inter := gmp.NewInt(0)
		deno.EvalMod(gmp.NewInt(int64(i+1)), modp, inter)
		interInv := gmp.NewInt(0)
		interInv.ModInverse(inter, modp)
		lambda[i].Mul(lambda[i], interInv)
		lambda[i].Mod(lambda[i], modp)
	}

	_0Shares := make([]*gmp.Int, counter)
	for i := 0; i < counter; i++ {
		_0Shares[i] = gmp.NewInt(0)
	}
	_0ShareCount := 0
	_0ShareSum := gmp.NewInt(0)

	recPoint := make([]*point.Point, counter)
	recCounter := 0
	Client := make([]*Node, 3)
	return Node{
		label:        label,
		counter:      counter,
		degree:       degree,
		p:            modp,
		randState:    randState,
		recPoint:     recPoint,
		recCounter:   recCounter,
		recPoly:      &poly.Poly{},
		Client:       Client,
		secretShares: secretShares,
		proPoly:      &proPoly,
		_0Shares:     _0Shares,
		_0ShareCount: &_0ShareCount,
		_0ShareSum:   _0ShareSum,
		lambda:       lambda,
	}, nil

}

//第三部分的分享， 是全份额重建
//通过新的多项式B（x，i），对包括自己的每一个node发送B（i，i）
//最后重建多项式

//重建secretShare
func (node *Node) SharePhase3(msg message) error {
	index := msg.getIndex()
	Y := msg.getY()

	node.secretShares[index-1].Y.SetBytes(Y)
	*node.shareCnt = *node.shareCnt + 1
	flag := (*node.shareCnt == node.counter)
	if flag {
		*node.shareCnt = 0

	}
	return nil

}
func (node *Node) ClientSharePhase3() {
	node.newPoly.Add(*node.recPoly, *node.proPoly)
	fmt.Println(*node.recPoly)
	for i := 1; i < node.counter; i++ {
		value := gmp.NewInt(0)
		node.newPoly.EvalMod(gmp.NewInt(int64(i+1)), node.p, value)

		if i != node.label-1 {
			log.Printf("node %d send point message to node %d in phase 3", node.label, i+1)
			msg := &message{
				Index: int32(node.label),
				X:     int32(i + 1),
				Y:     value.Bytes(),
			}
			//把消息发送给不同的节点
			node.Client[i].SharePhase3(*msg)

		} else {
			node.secretShares[i].Y.Set(value)
		}

	}
}

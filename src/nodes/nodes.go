package nodes

import (
	"errors"
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
	total int
	//Degree of polynomial
	degree int
	//the polynomial was set on Z_p
	p *gmp.Int
	// Rand source
	randState *rand.Rand
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
	//IP_address of node
	ipAddress []string
	//board IP address
	ipOfBoard string
	//clientconn
	clinetConn []*grpc.ClientConn
	//nodeService
	nodeService []pb.nodeservice
	//boardconn
	boardConn *grpc.ClientConn
	//boardService
	boardService pb.boardservice
}

func (node *Node) GetMsgFromNode(pointmsg point.Pointmsg) {
	index := pointmsg.GetIndex()
	log.Println("Phase 1 :[Node %d] receive point message from [Node %d]", node.label, index)
	p := pointmsg.GetPoint()
	//Receive the point and store
	node.mutex.Lock()
	node.recPoint[node.recCounter] = p
	node.recCounter += 1
	flag := (node.recCounter == node.total)
	node.mutex.Unlock()
	if flag {
		node.recCounter = 0
		node.Phase1()
	}
}
func (node *Node) SendMsgToNode() {
	p := point.Point{
		X:       node.secretShare[node.label-1].X,
		Y:       node.secretShare[node.label-1].Y,
		PolyWit: node.secretShare[node.label-1].PolyWit,
	}
	node.mutex.Lock()
	node.recPoint[node.recCounter] = &p
	node.recCounter = node.recCounter + 1
	flag := (node.recCounter == node.total)
	node.mutex.Unlock()
	if flag {
		node.recCounter = 0
		node.Phase1()
	}
	for i := 0; i < node.total; i++ {
		if i != node.label-1 {
			log.Printf("[Node %d] send point message to [Node %d]", node.label, i+1)
			msg := point.Pointmsg{}
			msg.SetIndex(node.label)
			msg.SetPoint(node.secretShare[i])
			//node.sent(msg)
		}
	}

}
func (node *Node) Phase1() {
	log.Printf("[Node %d] now start phase1", node.label)
	x_point := make([]*gmp.Int, node.degree+1)
	y_point := make([]*gmp.Int, node.degree+1)
	for i := 0; i <= node.degree; i++ {
		x_point[i] = node.recPoint[i].X
		y_point[i] = node.recPoint[i].Y
	}
	p, err := interpolation.LagrangeInterpolate(node.degree, x_point, y_point, node.p)
	if err != nil {
		for i := 0; i <= node.degree; i++ {
			log.Print(x_point[i])
			log.Print(y_point[i])
		}
		log.Print(err)
		panic("Interpolation failed")
	}
	node.recPoly = p
	log.Printf("Interpolation finished")
	//node.Phase2()
}

func (node *Node) GetLabel() int {
	if node != nil {
		return node.label
	} else {
		return 0
	}
}

func (node *Node) NodeConnect() {
	boradConn, err := grpc.Dial(node.ipOfBoard, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Fail to connect board:%v", err)
	}
	node.boardConn = boradConn
	//node.boardService = pb
	for i := 0; i < node.total; i++ {
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
	//var modp:= getprime.Getprime(256)
	return Node{
		label:      label,
		total:      counter,
		degree:     degree,
		p:          modp,
		randState:  randState,
		recPoint:   nil,
		recCounter: 0,
		recPoly:    poly.Poly{},
	}, nil
}

package main

import (
	"errors"
	"fmt"
	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/Alan-Lxc/crypto_contest/src/bulletboard"
	"github.com/Alan-Lxc/crypto_contest/src/control"
	"github.com/Alan-Lxc/crypto_contest/src/nodes"
	pb "github.com/Alan-Lxc/crypto_contest/src/service"
	"github.com/ncw/gmp"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
)

type client struct {
	ip string
	//IP addr of node in secret share
	ipList []string
	//IP addr of board
	ipBorad string
	//nodeConn
	nodeConn []*grpc.ClientConn
	//nodeservice
	nodeService []pb.NodeServiceClient
	//boardconn
	boardConn *grpc.ClientConn
	//boardService
	boardService pb.BulletinBoardServiceClient
	//metadatapath
	metadataPath string

	//degree of poly
	degree int
	//counter
	counter int
	//controller
	control control.Control
}

func ReadIpList(metadataPath string) []string {
	ipData, err := ioutil.ReadFile(metadataPath + "/ip_list")
	if err != nil {
		log.Fatalf("Failed to read iplist %v\n", err)
	}
	return strings.Split(string(ipData), "\n")
}
func newClient(degree, counter int, metadataPath string, ip string) (client, error) {

	ipRaw := ReadIpList(metadataPath)[0 : counter+1]
	bip := ipRaw[0]
	ipList := ipRaw[1 : counter+1]
	if degree < 0 {
		return client{}, errors.New(fmt.Sprintf("Can't generate a poly that degree < 0 "))

	}
	if counter < 0 {
		return client{}, errors.New(fmt.Sprintf("Can't generate a commitee smaller than 0"))
	}
	nConn := make([]*grpc.ClientConn, counter)
	nClient := make([]pb.NodeServiceClient, counter)
	return client{
		ip:           ip,
		ipBorad:      bip,
		ipList:       ipList,
		nodeService:  nClient,
		nodeConn:     nConn,
		degree:       degree,
		counter:      counter,
		metadataPath: metadataPath,
	}, nil
}

func (c *client) InitandConnect(s0 string) {
	log.Printf("Start to connecttion")
	degree := c.degree
	counter := c.counter
	fixedRandState := rand.New(rand.NewSource(int64(3)))
	p := gmp.NewInt(0)
	p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)
	tmp := gmp.NewInt(0)
	tmp.SetString(s0, 10)
	polyy, _ := poly.NewRand(degree, fixedRandState, p)
	polyy.SetCoeffWithGmp(0, tmp)
	polyyy := make([]poly.Poly, counter)
	for i := 0; i < counter; i++ {
		x := int32(i + 1)
		y := gmp.NewInt(0)
		polyy.EvalMod(gmp.NewInt(int64(x)), p, y)

		polyyy[i], _ = poly.NewRand(degree, fixedRandState, p)
		polyyy[i].SetCoeffWithGmp(0, y)
	}
	//nn := make([]nodes.Node, counter)
	//board	Init
	go newBoard(c.degree, c.counter, c.metadataPath, polyyy)
	time.Sleep(6)
	bconn, err := grpc.Dial(c.ipBorad, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Client could not connect to bulletboard")
	}
	c.boardConn = bconn
	c.boardService = pb.NewBulletinBoardServiceClient(bconn)
	//node
	//nn := make([]nodes.Node, counter)
	for i := 0; i < counter; i++ {
		//x := int32(i + 1)
		//y := gmp.NewInt(0)
		//polyy.EvalMod(gmp.NewInt(int64(x)), p, y)
		coeff := polyyy[i].GetAllCoeff()
		node, _ := nodes.New(degree, i+1, counter, c.metadataPath, coeff)
		go node.Service()
		time.Sleep(1)
		nConn, err := grpc.Dial(c.ipList[i], grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Client could not connect to node %d", i+1)
		}
		c.nodeConn[i] = nConn
		c.nodeService[i] = pb.NewNodeServiceClient(nConn)
	}
	time.Sleep(6)
	log.Printf("client has connected to committee and board")

	controll, err := control.New(c.metadataPath)
	if err != nil {
		log.Fatalf("Fail to connect controller")
	}
	log.Printf("Has connected to controller")
	c.control = controll
	c.control.Connect()
	//time.Sleep(6)
	c.control.StartHandoff()
}
func newBoard(degree int, ccounter int, metadataPath string, polyyy []poly.Poly) {
	bb, _ := bulletboard.New(degree, ccounter, metadataPath, polyyy)
	bb.Serve(false)
}
func main() {
	//cnt := flag.Int("c", 2, "Enter number of nodes")
	//degree := flag.Int("d", 1, "Enter the polynomial degree")
	//metadataPath := flag.String("path", "/mpss/metadata", "Enter the metadata path")
	//s0 := flag.String("secret","1234567899876543210","Enter the secret")
	////aws := flag.Bool("aws", false, "if test on real aws")
	//flag.Parse()
	client1, err := newClient(2, 5, "/home/gary/GolandProjects/crypto_contest/src/metadata", "192.168.0.1")
	if err != nil {
		log.Fatalf("Can't create a new client:%v", err)
	}
	client1.InitandConnect("0")
	//log.Printf("Done")
	//client1.control.StartHandoff()
}

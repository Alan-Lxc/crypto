package controller

import (
	"context"
	"fmt"
	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/Alan-Lxc/crypto_contest/src/bulletboard"
	"github.com/Alan-Lxc/crypto_contest/src/nodes"
	pb "github.com/Alan-Lxc/crypto_contest/src/service"
	"github.com/ncw/gmp"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"time"
)

type Controll struct {
	//nodeConn
	nodeConn []*grpc.ClientConn
	//nodeservice
	nodeService []pb.NodeServiceClient
	//boardconn
	boardConn []*grpc.ClientConn
	//boardService
	boardService []pb.BulletinBoardServiceClient
	//metadatapath
	ipList          []string
	bulletboardList []string
	//bb num
	bbNum   int
	nodeNum int
}

var Controller *Controll

//
func Initsystem() *Controll {
	//metadatapath := "./src/metadata"
	metadatapath := "/home/kzl/Desktop/test/crypto_contest/DCSSmain/src/metadata"
	var nodeConnnect []*nodes.Node
	nConn := make([]*grpc.ClientConn, 100) //get from sql and new
	nodeService := make([]pb.NodeServiceClient, 100)
	ipList := nodes.ReadIpList(metadatapath + "/ip_list")
	for i := 0; i < 100; i++ {
		node, err := nodes.New_for_web(i+1, metadatapath)
		nodeConnnect = append(nodeConnnect, node)
		if err != nil {
			println(err)
		}
		go node.Serve_for_web()
		Conn, err := grpc.Dial(ipList[i], grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Fail to connect with %s:%v", ipList[i], err)
		}
		nConn[i] = Conn
		nodeService[i] = pb.NewNodeServiceClient(Conn)
	}
	boradList := nodes.ReadIpList(metadatapath + "/bulletboard_list")
	controll := new(Controll)
	controll.ipList = ipList
	controll.nodeConn = nConn
	controll.nodeService = nodeService
	controll.bulletboardList = boradList
	controll.boardConn = make([]*grpc.ClientConn, 10)
	controll.boardService = make([]pb.BulletinBoardServiceClient, 10)
	controll.bbNum = 10
	controll.nodeNum = 100
	return controll
}
func (controll *Controll) NewSecret(secretid int, degree int, counter int, s0 string) {
	fmt.Println(controll.bbNum)
	metadatapath := "./src/metadata"
	fixedRandState := rand.New(rand.NewSource(int64(3)))
	p := gmp.NewInt(0)
	p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)
	tmp := gmp.NewInt(0)
	//tmp.SetString(s0, 10)
	tmp.SetString(s0, 10)
	polyy, _ := poly.NewRand(degree, fixedRandState, p)
	polyy.SetCoeffWithGmp(0, tmp)
	polyyy := make([]poly.Poly, counter)
	for i := 0; i < counter; i++ {
		x := int32(i + 1)
		y := gmp.NewInt(0)
		polyy.EvalMod(gmp.NewInt(int64(x)), p, y)

		polyyy[i], _ = poly.NewRand(degree*2, fixedRandState, p)
		polyyy[i].SetCoeffWithGmp(0, y)
	}
	bb, _ := bulletboard.New_bulletboard_for_web(degree, counter, metadatapath, secretid, polyyy)
	go bb.Serve(false)
	time.Sleep(2)
	bconn, err := grpc.Dial(bb.Getbip(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("System could not connect to bulletboard %d", secretid)
	}
	if secretid > controll.bbNum {
		tmp1 := make([]*grpc.ClientConn, secretid-controll.bbNum)
		controll.boardConn = append(controll.boardConn, tmp1...)
		tmp2 := make([]pb.BulletinBoardServiceClient, secretid-controll.bbNum)
		controll.boardService = append(controll.boardService, tmp2...)
		controll.bbNum = secretid
	}
	controll.boardConn[secretid-1] = bconn
	controll.boardService[secretid-1] = pb.NewBulletinBoardServiceClient(bconn)

	if counter > controll.nodeNum {
		tmp1 := make([]*grpc.ClientConn, counter-controll.nodeNum)
		controll.nodeConn = append(controll.nodeConn, tmp1...)
		tmp2 := make([]pb.NodeServiceClient, counter-controll.nodeNum)
		controll.nodeService = append(controll.nodeService, tmp2...)
		for i := controll.nodeNum; i < counter; i++ {
			Conn, err := grpc.Dial(controll.ipList[i], grpc.WithInsecure())
			if err != nil {
				log.Fatalf("Fail to connect with %s:%v", controll.ipList[i], err)
			}
			controll.nodeConn[i] = Conn
			controll.nodeService[i] = pb.NewNodeServiceClient(Conn)
		}
		controll.nodeNum = counter
	}
	for i := 0; i < counter; i++ {
		coeff := polyyy[i].GetAllCoeff()
		Coeff := make([][]byte, len(coeff))
		for i := 0; i < len(coeff); i++ {
			tmp := make([]byte, len(coeff[i].Bytes()))
			tmp = coeff[i].Bytes()
			Coeff[i] = tmp
		}
		msg := pb.InitMsg{
			Degree:   int32(degree),
			Counter:  int32(counter),
			Secretid: int32(secretid),
			Coeff:    Coeff,
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		controll.nodeService[i].Initsecret(ctx, &msg)
	}
}

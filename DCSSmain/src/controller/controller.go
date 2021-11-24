package controller

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/common"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/model"
	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/Alan-Lxc/crypto_contest/src/bulletinboard"
	"github.com/Alan-Lxc/crypto_contest/src/nodes"
	pb "github.com/Alan-Lxc/crypto_contest/src/service"
	"github.com/ncw/gmp"
	"google.golang.org/grpc"
	"log"
	//"math"
	"math/rand"
	"sync"
	"time"
)

type Controll struct {
	//nodeConn
	nodeConn []*grpc.ClientConn
	//nodeservice
	nodeService []pb.NodeServiceClient
	//boardconn
	boardConn *grpc.ClientConn
	//boardService
	boardService pb.BulletinBoardServiceClient
	//metadatapath
	node              []*nodes.Node
	Bulletinboard     *Bulletinboard.BulletinBoard
	ipList            []string
	BulletinboardList []string
	//bb num
	bbNum   int
	nodeNum int
}

// FORGOT TO KILL THR THREAD OF CONTROLL !!
var Controller *Controll

//这里写了metadatapath，后面就不需要写了
var metadatapath = "//home/gary/GolandProjects/crypto_contest4/DCSSmain/src/metadata"

func New() *Controll {
	return new(Controll)
}
func (controll *Controll) Release(counter int) {

	for i := 0; i < counter; i++ {
		controll.node[i].DeleteServe()
	}
	controll.Bulletinboard.DeleteServe()
	controll.node = nil
	controll.Bulletinboard = nil
}

func (controll *Controll) Connect(counter int) {
	for i := 0; i < counter; i++ {
		controll.node[i].NodeConnect()
	}
	controll.Bulletinboard.Connect()
}
func (controll *Controll) SetSecret(counter int) {
	for i := 0; i < counter; i++ {
		controll.node[i].SetSecret()
	}
	controll.Bulletinboard.SetSecret()
}
func (controll *Controll) Initsystem(degree, counter int, metadatapath string, secretid int) {
	db := common.GetDB()
	if db != nil {

	}
	//var nodeConnnect []*nodes.Node
	nConn := make([]*grpc.ClientConn, counter) //get from sql and new
	nodeService := make([]pb.NodeServiceClient, counter)
	ipRaw := nodes.ReadIpList(metadatapath)[0 : counter+1]
	ipList := ipRaw[1 : counter+1]
	nn := make([]*nodes.Node, counter)
	for i := 0; i < counter; i++ {
		node, err := nodes.NewForWeb(degree, i+1, counter, metadatapath, secretid)
		go node.ServeForWeb()
		// here need to change merge NODE
		newunit := model.Unit{
			UnitId: node.GetLabel(),
			UnitIp: ipList[node.GetLabel()-1],
			//Secretnum: 0,
		}
		db.Create(&newunit)
		nn[i] = &node
		//nodeConnnect = append(nodeConnnect, &node)
		if err != nil {
			println(err)
		}
		Conn, err := grpc.Dial(ipList[i], grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Fail to connect with %s:%v", ipList[i], err)
		}
		nConn[i] = Conn
		nodeService[i] = pb.NewNodeServiceClient(Conn)
	}
	bb, _ := Bulletinboard.NewBulletboardForWeb(degree, counter, metadatapath, secretid)
	go bb.Serve(false)
	time.Sleep(3)
	bconn, _ := grpc.Dial(bb.Getbip(), grpc.WithInsecure())
	boardConn := bconn
	boardService := pb.NewBulletinBoardServiceClient(bconn)

	//controll := new(Controll)
	controll.ipList = ipList
	controll.nodeConn = nConn
	controll.nodeService = nodeService
	//controll.BulletinboardList = boradList
	controll.boardConn = boardConn
	controll.boardService = boardService
	controll.node = nn
	controll.Bulletinboard = &bb
	//controll.bbNum = 1
	//controll.nodeNum = counter
	//return controll
}
func (controll *Controll) NewSecret(secretid int, degree int, counter int, s0 string) {

	//fmt.Println(controll.bbNum)
	fixedRandState := rand.New(rand.NewSource(int64(3)))
	p := gmp.NewInt(0)
	p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)
	tmp := gmp.NewInt(0)
	//tmp.SetString(s0, 10)
	tmp.SetString(s0, 16)
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
	controll.Initsystem(degree, counter, metadatapath, secretid)

	var wg sync.WaitGroup
	for i := 0; i < counter; i++ {
		coeff := polyyy[i].GetAllCoeff()
		tmpLength := len(coeff)
		Coeff := make([][]byte, tmpLength)
		for j := 0; j < tmpLength; j++ {
			Coeff[j] = coeff[j].Bytes()
		}
		msg := &pb.InitMsg{
			Degree:   int32(degree),
			Counter:  int32(counter),
			Secretid: int32(secretid),
			Coeff:    Coeff,
		}
		wg.Add(1)
		go func(i int, msg *pb.InitMsg) {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			controll.nodeService[i].Initsecret(ctx, msg)
			defer cancel()
		}(i, msg)
	}
	wg.Wait()
	controll.Release(counter)
}

func (controll *Controll) Handoff(secretid int, degree int, counter int) {

	log.Printf("Start to Handoff")
	//polyyy := controll.Getmessage(secretid, degree, counter)
	//get degree, counter, metadatapath, secretid, polyyy
	controll.Initsystem(degree, counter, metadatapath, secretid)
	controll.Connect(counter)
	controll.SetSecret(counter)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := controll.boardService.StartEpoch(ctx, &pb.RequestMsg{})
	if err != nil {
		log.Fatalf("Start Handoff Fail:%v", err)
	}
	//time.Sleep(5*time.Second)
	for {
		i := 0
		for i = 0; i < counter; i++ {
			if controll.node[i].Finish() == false {
				time.Sleep(500 * time.Millisecond)
				break
			}
		}

		if i == counter {
			break
		}
	}
	controll.Release(counter)
}

func (controll *Controll) Reconstruct(secretid int, degree int, counter int) string {
	log.Printf("Start to Reconstruction")

	//get degree, counter, metadatapath, secretid, polyyy
	//polyyy := controll.Getmessage(secretid, degree, counter)
	controll.Initsystem(degree, counter, metadatapath, secretid)
	controll.Connect(counter)
	controll.SetSecret(counter)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := controll.boardService.StartReconstruct(ctx, &pb.RequestMsg{})
	if err != nil {
		log.Fatalf("Start Reconstruction Fail:%v", err)
	}
	// Set to string
	Secret := hex.EncodeToString(controll.Bulletinboard.GetReconstructSecret().Bytes())
	fmt.Println(Secret)
	//Secret
	controll.Release(counter)
	return Secret
	//return "123"
}

func (controll *Controll) ModifyCommittee(secretid int, degree int, oldn int, newn int) {
	log.Printf("Start to ModifyCommittee")

	//get degree, counter, metadatapath, secretid, polyyy
	counter := 0
	if oldn > newn {
		counter = oldn
	} else {
		counter = newn
	}
	//polyyy := controll.Getmessage(secretid, degree, counter)
	controll.Initsystem(degree, counter, metadatapath, secretid)
	controll.Connect(counter)
	var wg sync.WaitGroup
	if oldn > newn {
		for i := newn; i < counter; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				ctx, cancel := context.WithCancel(context.Background())
				controll.nodeService[i].DeleteSecret(ctx, &pb.RequestMsg{})
				defer cancel()
			}(i)
		}
	} else {
		for i := oldn; i < counter; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				ctx, cancel := context.WithCancel(context.Background())
				controll.nodeService[i].AddSecret(ctx, &pb.RequestMsg{})
				defer cancel()
			}(i)
		}
	}
	wg.Wait()
	controll.Release(counter)
}

//package controller
//
//import (
//	"context"
//	"fmt"
//	"github.com/Alan-Lxc/crypto_contest/dcssweb/common"
//	"github.com/Alan-Lxc/crypto_contest/dcssweb/model"
//	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
//	"github.com/Alan-Lxc/crypto_contest/src/Bulletinboard"
//	"github.com/Alan-Lxc/crypto_contest/src/nodes"
//	pb "github.com/Alan-Lxc/crypto_contest/src/service"
//	"github.com/ncw/gmp"
//	"google.golang.org/grpc"
//	"log"
//	"math/rand"
//	"time"
//)
//
//type Controll struct {
//	//nodeConn
//	nodeConn []*grpc.ClientConn
//	//nodeservice
//	nodeService []pb.NodeServiceClient
//	//boardconn
//	boardConn []*grpc.ClientConn
//	//boardService
//	boardService []pb.BulletinBoardServiceClient
//	//metadatapath
//	ipList          []string
//	BulletinboardList []string
//	//bb num
//	bbNum   int
//	nodeNum int
//}
//
//var Controller *Controll
//
////这里写了metadatapath，后面就不需要写了
//var metadatapath = "/home/kzl/Desktop/test/crypto_contest/DCSSmain/src/metadata"
//
//func Initsystem() *Controll {
//	db := common.GetDB()
//	if db != nil {
//
//	}
//	var nodeConnnect []*nodes.Node
//	nConn := make([]*grpc.ClientConn, 100) //get from sql and new
//	nodeService := make([]pb.NodeServiceClient, 100)
//	ipList := nodes.ReadIpList(metadatapath + "/ip_list")
//	for i := 0; i < 100; i++ {
//		node, err := nodes.New_for_web(i+1, metadatapath)
//		newunit := model.Unit{
//			UnitId:  node.GetLabel(),
//			UnitIp:  node.IpAddress[node.GetLabel()],
//			//Secretnum: 0,
//		}
//		db.Create(&newunit)
//		nodeConnnect = append(nodeConnnect, node)
//		if err != nil {
//			println(err)
//		}
//		go node.Serve_for_web()
//		Conn, err := grpc.Dial(ipList[i], grpc.WithInsecure())
//		if err != nil {
//			log.Fatalf("Fail to connect with %s:%v", ipList[i], err)
//		}
//		nConn[i] = Conn
//		nodeService[i] = pb.NewNodeServiceClient(Conn)
//	}
//	boradList := nodes.ReadIpList(metadatapath + "/Bulletinboard_list")
//	controll := new(Controll)
//	controll.ipList = ipList
//	controll.nodeConn = nConn
//	controll.nodeService = nodeService
//	controll.BulletinboardList = boradList
//	controll.boardConn = make([]*grpc.ClientConn, 10)
//	controll.boardService = make([]pb.BulletinBoardServiceClient, 10)
//	controll.bbNum = 10
//	controll.nodeNum = 100
//	return controll
//}
//func (controll *Controll) NewSecret(secretid int, degree int, counter int, s0 string) {
//	fmt.Println(controll.bbNum)
//	fixedRandState := rand.New(rand.NewSource(int64(3)))
//	p := gmp.NewInt(0)
//	p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)
//	tmp := gmp.NewInt(0)
//	//tmp.SetString(s0, 10)
//	tmp.SetString(s0, 10)
//	polyy, _ := poly.NewRand(degree, fixedRandState, p)
//	polyy.SetCoeffWithGmp(0, tmp)
//	polyyy := make([]poly.Poly, counter)
//	for i := 0; i < counter; i++ {
//		x := int32(i + 1)
//		y := gmp.NewInt(0)
//		polyy.EvalMod(gmp.NewInt(int64(x)), p, y)
//
//		polyyy[i], _ = poly.NewRand(degree*2, fixedRandState, p)
//		polyyy[i].SetCoeffWithGmp(0, y)
//	}
//	bb, _ := Bulletinboard.New_Bulletinboard_for_web(degree, counter, metadatapath, secretid, polyyy)
//	go bb.Serve(false)
//	time.Sleep(2)
//	bconn, err := grpc.Dial(bb.Getbip(), grpc.WithInsecure())
//	if err != nil {
//		log.Fatalf("System could not connect to Bulletinboard %d", secretid)
//	}
//	if secretid > controll.bbNum {
//		tmp1 := make([]*grpc.ClientConn, secretid-controll.bbNum)
//		controll.boardConn = append(controll.boardConn, tmp1...)
//		tmp2 := make([]pb.BulletinBoardServiceClient, secretid-controll.bbNum)
//		controll.boardService = append(controll.boardService, tmp2...)
//		controll.bbNum = secretid
//	}
//	controll.boardConn[secretid-1] = bconn
//	controll.boardService[secretid-1] = pb.NewBulletinBoardServiceClient(bconn)
//
//	if counter > controll.nodeNum {
//		tmp1 := make([]*grpc.ClientConn, counter-controll.nodeNum)
//		controll.nodeConn = append(controll.nodeConn, tmp1...)
//		tmp2 := make([]pb.NodeServiceClient, counter-controll.nodeNum)
//		controll.nodeService = append(controll.nodeService, tmp2...)
//		for i := controll.nodeNum; i < counter; i++ {
//			Conn, err := grpc.Dial(controll.ipList[i], grpc.WithInsecure())
//			if err != nil {
//				log.Fatalf("Fail to connect with %s:%v", controll.ipList[i], err)
//			}
//			controll.nodeConn[i] = Conn
//			controll.nodeService[i] = pb.NewNodeServiceClient(Conn)
//		}
//		controll.nodeNum = counter
//	}
//	for i := 0; i < counter; i++ {
//		coeff := polyyy[i].GetAllCoeff()
//		Coeff := make([][]byte, len(coeff))
//		for i := 0; i < len(coeff); i++ {
//// 			tmp := make([]byte, len(coeff[i].Bytes()))
//// 			tmp = coeff[i].Bytes()
//			Coeff[i] = coeff[i].Bytes()
//		}
//		msg := pb.InitMsg{
//			Degree:   int32(degree),
//			Counter:  int32(counter),
//			Secretid: int32(secretid),
//			Coeff:    Coeff,
//		}
//		ctx, cancel := context.WithCancel(context.Background())
//		defer cancel()
//		controll.nodeService[i].Initsecret(ctx, &msg)
//	}
//	//controll.Handoff(secretid)
//}
//
//func (controll *Controll) Handoff(secretid int) {
//	log.Printf("Start to Handoff")
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	_, err := controll.boardService[secretid-1].StartEpoch(ctx, &pb.RequestMsg{})
//	if err != nil {
//		log.Fatalf("Start Handoff Fail:%v", err)
//	}
//}

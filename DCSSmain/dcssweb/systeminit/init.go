package systeminit

import (
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
	ip_list          []string
	bulletboard_list []string
	//bb num
	bb_num int
}

func Initsystem() *Controll {
	//metadatapath := "./src/metadata"
	metadatapath := "./src/metadata"
	var nodeConnnect []*nodes.Node
	nConn := make([]*grpc.ClientConn, 100)
	nodeService := make([]pb.NodeServiceClient, 100)
	ip_list := nodes.ReadIpList(metadatapath + "ip_list")
	for i := 0; i < 100; i++ {
		node, err := nodes.New_for_web(i, metadatapath)
		nodeConnnect = append(nodeConnnect, node)
		if err != nil {
			println(err)
		}
		go node.Serve_for_web()
		Conn, err := grpc.Dial(ip_list[i], grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Fail to connect with %s:%v", ip_list[i], err)
		}
		nConn[i] = Conn
		nodeService[i] = pb.NewNodeServiceClient(Conn)
	}
	borad_list := nodes.ReadIpList(metadatapath + "bulletboard_list")
	controll := new(Controll)
	controll.ip_list = ip_list
	controll.nodeConn = nConn
	controll.nodeService = nodeService
	controll.bulletboard_list = borad_list
	controll.boardConn = make([]*grpc.ClientConn, 10)
	controll.boardService = make([]pb.BulletinBoardServiceClient, 10)
	controll.bb_num = 10
	return controll
}
func NewSecret(secretid int, degree int, counter int, s0 string, controll Controll) *bulletboard.BulletinBoard {
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
	if secretid > controll.bb_num {
		tmp1 := make([]*grpc.ClientConn, secretid-controll.bb_num)
		controll.boardConn = append(controll.boardConn, tmp1...)
		tmp2 := make([]pb.BulletinBoardServiceClient, secretid-controll.bb_num)
		controll.boardService = append(controll.boardService, tmp2...)
		controll.bb_num = secretid
	}
	controll.boardConn[secretid-1] = bconn
	controll.boardService[secretid-1] = pb.NewBulletinBoardServiceClient(bconn)
	return &bb
}

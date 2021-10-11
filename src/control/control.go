package control

import (
	"context"
	"flag"
	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	"github.com/Alan-Lxc/crypto_contest/src/bulletboard"
	"github.com/Alan-Lxc/crypto_contest/src/nodes"
	pb "github.com/Alan-Lxc/crypto_contest/src/service"
	"github.com/ncw/gmp"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
)

type Control struct {
	metadataPath string
	boardIp      string

	boardConn    *grpc.ClientConn
	boardService pb.BulletinBoardServiceClient
}

func (c *Control) Connect() {
	boarConn, err := grpc.Dial(c.boardIp, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Fail to connect borad:%v", err)
	}
	c.boardConn = boarConn
	c.boardService = pb.NewBulletinBoardServiceClient(boarConn)
}
func (c *Control) Disconnect() {
	log.Println("Disconnect the board")
	c.boardConn.Close()
}
func (c *Control) StartHandoff() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Printf("Start to Handoff")
	_, err := c.boardService.StartEpoch(ctx, &pb.RequestMsg{})
	if err != nil {
		log.Fatalf("Start Handoff Fail:%v", err)
	}
}
func ReadIpList(metadataPath string) []string {
	ipData, err := ioutil.ReadFile(metadataPath + "/ip_list")
	if err != nil {
		log.Fatalf("clock failed to read iplist %v\n", err)
	}
	return strings.Split(string(ipData), "\n")
}

// New returns a network node structure
func New(degree int, counter int, metadataPath string, s0 string) (Control, error) {
	bip := ReadIpList(metadataPath)[0]

	//randState := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	fixedRandState := rand.New(rand.NewSource(int64(3)))
	p := gmp.NewInt(0)
	p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)
	tmp := gmp.NewInt(0)
	tmp.SetString(s0, 10)

	polyy, _ := poly.NewRand(degree, fixedRandState, p)
	polyy.SetCoeffWithGmp(0, tmp)

	polyyy := make([]poly.Poly, counter)
	nn := make([]nodes.Node, counter)
	for i := 0; i < counter; i++ {
		y := gmp.NewInt(0)
		polyy.EvalMod(gmp.NewInt(int64(x)), p, y)

		polyyy[i], _ = poly.NewRand(degree, fixedRandState, p)
		polyyy[i].SetCoeffWithGmp(0, y)
	}

	aws := flag.Bool("aws", false, "if test on real aws")
	flag.Parse()
	bb, _ := bulletboard.New(degree, counter, metadataPath, polyyy)
	bb.Serve(*aws)

	for i := 0; i < counter; i++ {
		x := int32(i)
		y := gmp.NewInt(0)
		polyy.EvalMod(gmp.NewInt(int64(x)), p, y)
		coeff := polyyy[i].GetAllCoeff()
		nn[i], _ = nodes.New(degree, i, counter, metadataPath, coeff)
		nn[i].Service()
	}

	return Control{
		metadataPath: metadataPath,
		boardIp:      bip,
	}, nil
}

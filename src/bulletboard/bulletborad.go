package bulletboard

import (
	"context"
	"errors"
	"fmt"
	"github.com/Alan-Lxc/crypto_contest/src/basic/commitment"
	"github.com/Alan-Lxc/crypto_contest/src/basic/interpolation"
	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	pb "github.com/Alan-Lxc/crypto_contest/src/service"
	"github.com/golang/protobuf/proto"
	"github.com/ncw/gmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
)

// BulletinBoard Simulator Structure
type BulletinBoard struct {
	// Metadata Directory Path
	metadataPath string
	// Degree
	degree int
	// Counter
	counter int
	//the polynomial was set on Z_p
	p *gmp.Int
	// BulletinBoard IP Address
	bip string
	// IP
	ipList []string
	// Rand
	randState *rand.Rand
	// Reconstruction BulletinBoard
	reconstructionContent []*pb.Cmt1Msg
	// Proactivization BulletinBoard
	proCnt                 *int
	proactivizationContent []*pb.CommitMsg
	// Share Distribution BulletinBoard
	shaCnt *int
	//Secretcnt
	secretCnt *int
	//RecontructSecret
	recontructSecret []*gmp.Int
	//the secret
	secret *gmp.Int
	// Mutexes
	mutex sync.Mutex

	nConn   []*grpc.ClientConn
	nClient []pb.NodeServiceClient

	// Metrics
	totMsgSize *int
	//poly
	randPoly *poly.Poly
	//logger file pointer
	log *log.Logger
}

func (bb *BulletinBoard) GetCoeffofNodeSecretShares2(ctx context.Context, msg *pb.RequestMsg) (*pb.CoeffMsg, error) {

	degree := bb.randPoly.GetDegree()
	coeff := make([][]byte, degree)
	for i := 0; i <= degree; i++ {
		tmp, err := bb.randPoly.GetCoeff(i)
		if err != nil {
			panic("error")
		}
		coeff[i] = tmp.Bytes()
	}
	return &pb.CoeffMsg{Coeff: coeff}, nil
}

func (bb *BulletinBoard) StartEpoch(ctx context.Context, in *pb.RequestMsg) (*pb.ResponseMsg, error) {
	bb.log.Print("[bulletinboard] start epoch")
	bb.ClientStartPhase1()
	return &pb.ResponseMsg{}, nil
}

func (bb *BulletinBoard) ReadPhase1(in *pb.RequestMsg, stream pb.BulletinBoardService_ReadPhase1Server) error {
	bb.log.Print("[bulletinboard] is being read in phase 1")
	for i := 0; i < bb.degree*2+1; i++ {
		if err := stream.Send(bb.reconstructionContent[i]); err != nil {
			bb.log.Fatalf("bulletinboard failed to read phase1: %v", err)
			return err
		}
	}
	return nil
}

func (bb *BulletinBoard) WritePhase1(ctx context.Context, msg *pb.Cmt1Msg) (*pb.ResponseMsg, error) {
	//fmt.Println("written index is xxxx")
	*bb.totMsgSize = *bb.totMsgSize + proto.Size(msg)
	bb.log.Print("[bulletinboard] is being writcommitmentten in phase 3")
	index := msg.GetIndex()
	bb.reconstructionContent[index-1] = msg
	bb.mutex.Lock()
	*bb.shaCnt = *bb.shaCnt + 1
	flag := (*bb.shaCnt == bb.degree*2+1)
	bb.mutex.Unlock()
	if flag {
		*bb.shaCnt = 0
		bb.ClientStartVerifPhase1()
	}
	return &pb.ResponseMsg{}, nil
}

func (bb *BulletinBoard) ClientStartVerifPhase1() {
	var wg sync.WaitGroup
	for i := 0; i < bb.counter; i++ {
		bb.log.Print("[bulletinboard] start verification in phase 3")
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			bb.nClient[i].Phase1Verify(ctx, &pb.RequestMsg{})
		}(i)
	}
	wg.Wait()
	*bb.totMsgSize = 0
}

func (bb *BulletinBoard) WritePhase2(ctx context.Context, msg *pb.CommitMsg) (*pb.ResponseMsg, error) {
	*bb.totMsgSize = *bb.totMsgSize + proto.Size(msg)
	bb.log.Print("[bulletinboard] is being written in phase 2")
	index := msg.GetIndex()
	bb.proactivizationContent[index-1] = msg
	bb.mutex.Lock()
	*bb.proCnt = *bb.proCnt + 1
	flag := (*bb.proCnt == bb.degree*2+1)
	bb.mutex.Unlock()
	if flag {
		*bb.proCnt = 0
		bb.ClientStartVerifPhase2()
	}
	return &pb.ResponseMsg{}, nil
}

func (bb *BulletinBoard) ReadPhase2(in *pb.RequestMsg, stream pb.BulletinBoardService_ReadPhase2Server) error {
	bb.log.Print("[bulletinboard] is being read in phase 2")
	for i := 0; i < bb.degree*2+1; i++ {
		if err := stream.Send(bb.proactivizationContent[i]); err != nil {
			bb.log.Fatalf("bulletinboard failed to read phase2: %v", err)
			return err
		}
	}
	return nil
}

func (bb *BulletinBoard) WritePhase3(ctx context.Context, msg *pb.Cmt1Msg) (*pb.ResponseMsg, error) {
	//fmt.Println("written index is xxxx")
	*bb.totMsgSize = *bb.totMsgSize + proto.Size(msg)
	bb.log.Print("[bulletinboard] is being written in phase 3")
	index := msg.GetIndex()
	bb.reconstructionContent[index-1] = msg
	bb.mutex.Lock()
	*bb.shaCnt = *bb.shaCnt + 1
	flag := (*bb.shaCnt == bb.degree*2+1)
	bb.mutex.Unlock()
	if flag {
		*bb.shaCnt = 0
		bb.ClientStartVerifPhase3()
	}
	return &pb.ResponseMsg{}, nil
}

func (bb *BulletinBoard) WritePhase32(ctx context.Context, msg *pb.Cmt1Msg) (*pb.ResponseMsg, error) {
	//fmt.Println("written index is xxxx")
	*bb.totMsgSize = *bb.totMsgSize + proto.Size(msg)
	bb.log.Print("[bulletinboard] is being written in phase 3 2")
	index := msg.GetIndex()
	bb.reconstructionContent[index-1] = msg
	bb.mutex.Lock()
	*bb.shaCnt = *bb.shaCnt + 1
	flag := (*bb.shaCnt == bb.degree*2+1)
	bb.mutex.Unlock()
	if flag {
		*bb.shaCnt = 0
	}
	return &pb.ResponseMsg{}, nil
}

func (bb *BulletinBoard) ReconstructSecret(ctx context.Context, msg *pb.PointMsg) (*pb.ResponseMsg, error) {
	*bb.totMsgSize = *bb.totMsgSize + proto.Size(msg)
	bb.log.Print("[bulletinboard] recontruct secret")
	index := msg.GetIndex()
	Y := gmp.NewInt(0).SetBytes(msg.GetY())
	//fmt.Println(len(bb.recontructSecret),bb.recontructSecret)
	bb.recontructSecret[index-1] = Y
	bb.mutex.Lock()
	*bb.secretCnt = *bb.secretCnt + 1
	flag := (*bb.secretCnt == bb.degree*2+1)
	bb.mutex.Unlock()
	if flag {
		//bb.log.Println(*bb.secretCnt,bb.recontructSecret)
		*bb.secretCnt = 0
		bb.SecretPrint()
	}
	return &pb.ResponseMsg{}, nil
}
func (bb *BulletinBoard) SecretPrint() {
	X := make([]*gmp.Int, bb.degree*2+1)
	for i := 0; i < bb.degree*2+1; i++ {
		X[i] = gmp.NewInt(int64(i + 1))
	}
	//bb.log.Println(bb.recontructSecret)
	flg := 1
	for flg == 1 {
		flg = 0
		for i := 0; i < bb.degree*2+1; i++ {
			if bb.recontructSecret[i] == nil {
				flg = 1
				break
			}
		}
	}
	polytmp, _ := interpolation.LagrangeInterpolate(bb.degree, X, bb.recontructSecret[:bb.degree*2+1], bb.p)
	polytmp.EvalMod(gmp.NewInt(0), bb.p, bb.secret)
	bb.log.Print("[bulletinboard] the secret is ", bb.secret)

	*bb.proCnt = 0
	*bb.shaCnt = 0
	//f, _ := os.OpenFile(bb.metadataPath+"/log0", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//defer f.Close()
	bb.log.Printf("totMsgSize,%d\n", *bb.totMsgSize)
	*bb.totMsgSize = 0
}
func (bb *BulletinBoard) ReadPhase3(in *pb.RequestMsg, stream pb.BulletinBoardService_ReadPhase3Server) error {
	bb.log.Print("[bulletinboard] is being read in phase 3")
	for i := 0; i < bb.degree*2+1; i++ {
		if err := stream.Send(bb.reconstructionContent[i]); err != nil {
			bb.log.Fatalf("bulletinboard failed to read phase2: %v", err)
			return err
		}
	}
	return nil
}

func (bb *BulletinBoard) Connect() {
	for i := 0; i < bb.counter; i++ {
		nConn, err := grpc.Dial(bb.ipList[i], grpc.WithInsecure())
		if err != nil {
			bb.log.Fatalf("bulletinboard did not connect: %v", err)
		}
		bb.nConn[i] = nConn
		bb.nClient[i] = pb.NewNodeServiceClient(nConn)
	}
}

func (bb *BulletinBoard) Disconnect() {
	for i := 0; i < bb.counter; i++ {
		bb.nConn[i].Close()
	}
}

func (bb *BulletinBoard) Serve(aws bool) {
	port := bb.bip
	if aws {
		port = "0.0.0.0:12001"
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		bb.log.Fatalf("bulletinboard failed to listen %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBulletinBoardServiceServer(s, bb)
	reflection.Register(s)
	bb.log.Printf("bulletinboard serve on " + bb.bip)
	//bb.log.Printf("hello\n")
	//if err := s.Serve(lis); err != nil {
	//	bb.log.Fatalf("bulletinboard failed to serve %v", err)
	//}
	//bb.log.Printf("hello")
	err = s.Serve(lis)
	bb.log.Printf("[bulletboard]hello")
	if err != nil {
		bb.log.Fatalf("bulletinboard failed to serve %v", err)
	}
}

func (bb *BulletinBoard) ClientStartPhase1() {
	if bb.nConn[0] == nil {
		bb.Connect()
	}
	var wg sync.WaitGroup
	for i := 0; i < bb.counter; i++ {
		bb.log.Print("[bulletinboard] start phase 1")
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			bb.nClient[i].Phase1GetStart(ctx, &pb.RequestMsg{})
		}(i)
	}
	wg.Wait()
}

func (bb *BulletinBoard) ClientStartVerifPhase2() {
	var wg sync.WaitGroup
	for i := 0; i < bb.counter; i++ {
		bb.log.Print("[bulletinboard] start verification in phase 2")
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			bb.nClient[i].Phase2Verify(ctx, &pb.RequestMsg{})
		}(i)
	}
	wg.Wait()
}
func GetCoeffFromMsg(msg *pb.CoeffMsg) []*gmp.Int {
	tmp := msg.GetCoeff()
	l := len(tmp)
	res := make([]*gmp.Int, l)
	for i := 0; i < l; i++ {
		n := gmp.NewInt(0)
		n.SetBytes(tmp[i])
		res[i] = n
	}
	return res
}
func (bb *BulletinBoard) ClientStartVerifPhase3() {
	var wg sync.WaitGroup
	for i := 0; i < bb.counter; i++ {
		bb.log.Print("[bulletinboard] start verification in phase 3")
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			bb.nClient[i].Phase3Verify(ctx, &pb.RequestMsg{})
		}(i)
	}
	wg.Wait()
}

func ReadIpList(metadataPath string) []string {
	ipData, err := ioutil.ReadFile(metadataPath + "/ip_list")
	if err != nil {
		log.Fatalf("bulletinboard failed to read iplist: %v", err)
	}
	return strings.Split(string(ipData), "\n")
}

// New returns a network node structure
func New(degree int, counter int, metadataPath string, Polyyy []poly.Poly) (BulletinBoard, error) {
	//f, _ := os.Create(metadataPath + "/log0")
	//defer f.Close()

	fileName := metadataPath + "/bulletboard.log"
	tmplogger, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		tmplogger, err = os.Create(fileName)
	}
	//os.Truncate(fileName, 0)
	log := log.New(tmplogger, "", log.LstdFlags)

	if counter < 0 {
		return BulletinBoard{}, errors.New(fmt.Sprintf("counter must be non-negative, got %d", counter))
	}

	//fixedRandState := rand.New(rand.NewSource(int64(3)))
	p := gmp.NewInt(0)
	p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)
	dpc := commitment.DLPolyCommit{}
	dpc.SetupFix(counter)

	ipRaw := ReadIpList(metadataPath)[0 : counter+1]
	bip := ipRaw[0]
	ipList := ipRaw[1 : counter+1]

	proCnt := 0
	shaCnt := 0
	secretCnt := 0
	reconstructionContent := make([]*pb.Cmt1Msg, counter)
	reconstructSecret := make([]*gmp.Int, counter)
	secret := gmp.NewInt(0)
	//polyp, err := poly.NewRand(degree, fixedRandState, p)
	//if err != nil {
	//	bb.log.Fatal("Error initializing random poly")
	//}
	for i := 0; i < counter; i++ {
		c := dpc.NewG1()
		dpc.Commit(c, Polyyy[i])
		//fmt.Println(Polyyy[i].GetDegree())
		cBytes := c.CompressedBytes()
		msg := &pb.Cmt1Msg{
			Index:   int32(i + 1),
			Polycmt: cBytes,
		}
		reconstructionContent[i] = msg
	}
	proactivizationContent := make([]*pb.CommitMsg, counter)

	nConn := make([]*grpc.ClientConn, counter)
	nClient := make([]pb.NodeServiceClient, counter)

	totMsgSize := 0

	return BulletinBoard{
		degree:                 degree,
		p:                      p,
		metadataPath:           metadataPath,
		recontructSecret:       reconstructSecret,
		counter:                counter,
		bip:                    bip,
		ipList:                 ipList,
		proCnt:                 &proCnt,
		shaCnt:                 &shaCnt,
		secretCnt:              &secretCnt,
		secret:                 secret,
		reconstructionContent:  reconstructionContent,
		proactivizationContent: proactivizationContent,
		nConn:                  nConn,
		nClient:                nClient,
		totMsgSize:             &totMsgSize,
		log:                    log,
	}, nil
}

package Bulletinboard

import (
	"context"
	"errors"
	"fmt"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/common"
	"github.com/Alan-Lxc/crypto_contest/dcssweb/model"
	"github.com/Alan-Lxc/crypto_contest/src/basic/commitment"
	"github.com/Alan-Lxc/crypto_contest/src/basic/interpolation"
	"github.com/Alan-Lxc/crypto_contest/src/basic/poly"
	model1 "github.com/Alan-Lxc/crypto_contest/src/model"
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
	"strconv"
	"strings"
	"sync"
)

// BulletinBoard Simulator Structure
type BulletinBoard struct {
	id int
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

	// Reconstruction BulletinBoard2
	reconstructionContent2 []*pb.Cmt1Msg

	// Reconstruction BulletinBoard3
	reconstructionContent3 []*pb.Cmt1Msg

	// Reconstruction BulletinBoard4
	reconstructionContent4 []*pb.Cmt1Msg
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
	//Server
	server  *grpc.Server
	nConn   []*grpc.ClientConn
	nClient []pb.NodeServiceClient

	// Metrics
	totMsgSize *int
	//poly
	randPoly *poly.Poly
	//logger file pointer
	log    *log.Logger
	dc     *commitment.DLCommit
	dpc    *commitment.DLPolyCommit
	finish bool
}

func (bb *BulletinBoard) Getfinish() bool {
	return bb.finish
}
func (bb *BulletinBoard) Getbip() string {
	return bb.bip
}
func (bb *BulletinBoard) GetReconstructSecret() *gmp.Int {
	return bb.secret
}

//func (bb *BulletinBoard) GetCoeffofNodeSecretShares2(ctx context.Context, msg *pb.RequestMsg) (*pb.CoeffMsg, error) {
//
//	degree := bb.randPoly.GetDegree()
//	coeff := make([][]byte, degree)
//	for i := 0; i <= degree; i++ {
//		tmp, err := bb.randPoly.GetCoeff(i)
//		if err != nil {
//			panic("error")
//		}
//		coeff[i] = tmp.Bytes()
//	}
//
//}

func (bb *BulletinBoard) StartEpoch(ctx context.Context, in *pb.RequestMsg) (*pb.ResponseMsg, error) {
	bb.log.Print("[Bulletinboard] start epoch")
	for i := 0; i < bb.counter; i++ {
		bb.reconstructionContent[i] = bb.reconstructionContent4[i]
	}
	//secretid:=in.GetSecretid()
	bb.ClientStartPhase1(bb.id)
	return &pb.ResponseMsg{}, nil
}

func (bb *BulletinBoard) ReadPhase1(in *pb.RequestMsg, stream pb.BulletinBoardService_ReadPhase1Server) error {
	bb.log.Print("[Bulletinboard] is being read in phase 1")
	for i := 0; i < bb.degree*2+1; i++ {
		if err := stream.Send(bb.reconstructionContent[i]); err != nil {
			bb.log.Fatalf("[Bulletinboard] failed to read phase1: %v", err)
			return err
		}
	}
	return nil
}

func (bb *BulletinBoard) WritePhase1(ctx context.Context, msg *pb.Cmt1Msg) (*pb.ResponseMsg, error) {
	//fmt.Println("written index is xxxx")
	*bb.totMsgSize = *bb.totMsgSize + proto.Size(msg)
	bb.log.Print("[Bulletinboard] is being writcommitmentten in phase 3")
	index := msg.GetIndex()
	bb.reconstructionContent2[index-1] = msg
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
		bb.log.Print("[Bulletinboard] start verification in phase 3")
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			bb.nClient[i].Phase1Verify(ctx, &pb.RequestMsg{})
		}(i)
	}
	wg.Wait()
	//*bb.totMsgSize = 0
	return
}

func (bb *BulletinBoard) ReadPhase12(in *pb.RequestMsg, stream pb.BulletinBoardService_ReadPhase12Server) error {
	bb.log.Print("[Bulletinboard] is being read in phase 1 2")
	for i := 0; i < bb.degree*2+1; i++ {
		if err := stream.Send(bb.reconstructionContent2[i]); err != nil {
			bb.log.Fatalf("[Bulletinboard] failed to read phase1: %v", err)
			return err
		}
	}
	return nil
}
func (bb *BulletinBoard) WritePhase2(ctx context.Context, msg *pb.CommitMsg) (*pb.ResponseMsg, error) {
	*bb.totMsgSize = *bb.totMsgSize + proto.Size(msg)
	bb.log.Print("[Bulletinboard] is being written in phase 2")
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
	bb.log.Print("[Bulletinboard] is being read in phase 2")
	for i := 0; i < bb.degree*2+1; i++ {
		if err := stream.Send(bb.proactivizationContent[i]); err != nil {
			bb.log.Fatalf("[Bulletinboard] failed to read phase2: %v", err)
			return err
		}
	}
	return nil
}

func (bb *BulletinBoard) WritePhase3(ctx context.Context, msg *pb.Cmt1Msg) (*pb.ResponseMsg, error) {
	//fmt.Println("written index is xxxx")
	*bb.totMsgSize = *bb.totMsgSize + proto.Size(msg)
	index := msg.GetIndex()
	bb.log.Print("[Bulletinboard] is being written in phase 3", index)
	bb.reconstructionContent3[index-1] = msg
	bb.mutex.Lock()
	*bb.shaCnt = *bb.shaCnt + 1
	flag := (*bb.shaCnt == bb.degree*2+1)
	bb.mutex.Unlock()
	if flag {
		*bb.shaCnt = 0
		bb.ClientStartVerifPhase3()
	}
	bb.finish = true
	return &pb.ResponseMsg{}, nil
}

func (bb *BulletinBoard) WritePhase32(ctx context.Context, msg *pb.Cmt1Msg) (*pb.ResponseMsg, error) {
	//fmt.Println("written index is xxxx")
	*bb.totMsgSize = *bb.totMsgSize + proto.Size(msg)
	bb.log.Print("[Bulletinboard] is being written in phase 3 2")
	index := msg.GetIndex()
	bb.reconstructionContent4[index-1] = msg
	bb.mutex.Lock()
	*bb.shaCnt = *bb.shaCnt + 1
	flag := (*bb.shaCnt == bb.counter)
	bb.mutex.Unlock()
	if flag {
		*bb.shaCnt = 0
	}

	*bb.proCnt = 0
	*bb.shaCnt = 0
	//f, _ := os.OpenFile(bb.metadataPath+"/log0", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//defer f.Close()
	bb.log.Printf("totMsgSize,%d\n", *bb.totMsgSize)
	fmt.Println("[Bulletinboard] totmsgsize is %d", *bb.totMsgSize)
	*bb.totMsgSize = 0
	return &pb.ResponseMsg{}, nil
}

func (bb *BulletinBoard) ReconstructSecret(ctx context.Context, msg *pb.PointMsg) (*pb.ResponseMsg, error) {
	*bb.totMsgSize = *bb.totMsgSize + proto.Size(msg)
	bb.log.Print("[Bulletinboard] reconstruct secret")
	index := msg.GetIndex()
	Y := gmp.NewInt(0).SetBytes(msg.GetY())
	//fmt.Println(len(bb.recontructSecret),bb.recontructSecret)
	bb.recontructSecret[index-1] = Y
	bb.mutex.Lock()
	*bb.secretCnt = *bb.secretCnt + 1
	flag := (*bb.secretCnt == bb.counter)
	//fmt.Println(bb.counter, *bb.secretCnt)
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
	//flg := 1
	//for flg == 1 {
	//	flg = 0
	//	for i := 0; i < bb.degree*2+1; i++ {
	//		if bb.recontructSecret[i] == nil {
	//			flg = 1
	//			break
	//		}
	//	}
	//}
	polytmp, _ := interpolation.LagrangeInterpolate(bb.degree, X, bb.recontructSecret[:bb.degree*2+1], bb.p)
	polytmp.EvalMod(gmp.NewInt(0), bb.p, bb.secret)
	bb.log.Print("[Bulletinboard] the secret is ", bb.secret)
	return
}
func (bb *BulletinBoard) ReadPhase3(in *pb.RequestMsg, stream pb.BulletinBoardService_ReadPhase3Server) error {
	bb.log.Print("[Bulletinboard] is being read in phase 3")
	for i := 0; i < bb.degree*2+1; i++ {
		if err := stream.Send(bb.reconstructionContent3[i]); err != nil {
			bb.log.Fatalf("[Bulletinboard] failed to read phase2: %v", err)
			return err
		}
	}
	return nil
}

func (bb *BulletinBoard) Connect() {
	for i := 0; i < bb.counter; i++ {
		nConn, err := grpc.Dial(bb.ipList[i], grpc.WithInsecure())
		if err != nil {
			bb.log.Fatalf("[Bulletinboard] did not connect: %v", err)
		}
		bb.nConn[i] = nConn
		bb.nClient[i] = pb.NewNodeServiceClient(nConn)
	}
	return
}

func (bb *BulletinBoard) Disconnect() {
	for i := 0; i < bb.counter; i++ {
		bb.nConn[i].Close()
	}
	return
}

func (bb *BulletinBoard) DeleteServe() {
	bb.log.Printf("[Bulletinboard] delete serve on " + bb.bip)
	bb.server.Stop()

	return
}
func (bb *BulletinBoard) Serve(aws bool) {
	port := bb.bip
	if aws {
		port = "0.0.0.0:12001"
	}
	bb.log.Printf("[Bulletinboard] serve on " + bb.bip)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		bb.log.Fatalf("[Bulletinboard] failed to listen %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBulletinBoardServiceServer(s, bb)
	reflection.Register(s)
	bb.server = s
	//bb.log.Printf("hello\n")
	//if err := s.Serve(lis); err != nil {
	//	bb.log.Fatalf("[Bulletinboard] failed to serve %v", err)
	//}
	//bb.log.Printf("hello")
	err = s.Serve(lis)
	bb.log.Printf("[Bulletinboard] hello")
	if err != nil {
		bb.log.Fatalf("[Bulletinboard] failed to serve %v", err)
	}
	return
}
func (bb *BulletinBoard) ClientStartPhase1(secretid int) {
	var wg sync.WaitGroup
	bb.log.Print("[Bulletinboard] start phase 1")
	for i := 0; i < bb.counter; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			msg := pb.StartMsg{Secretid: int32(secretid)}
			_, err := bb.nClient[i].Phase1GetStart(ctx, &msg)
			if err != nil {
				return
			}
		}(i)
	}
	wg.Wait()
	return
}

func (bb *BulletinBoard) ClientStartVerifPhase2() {
	var wg sync.WaitGroup
	bb.log.Print("[Bulletinboard] start verification in phase 2")
	for i := 0; i < bb.counter; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			bb.nClient[i].Phase2Verify(ctx, &pb.RequestMsg{})
		}(i)
	}
	wg.Wait()
	return
}

func (bb *BulletinBoard) ClientStartVerifPhase3() {
	var wg sync.WaitGroup
	bb.log.Print("[Bulletinboard] start verification in phase 3")
	for i := 0; i < bb.counter; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			bb.nClient[i].Phase3Verify(ctx, &pb.RequestMsg{})
		}(i)
	}
	wg.Wait()
	return
}

func ReadIpList(metadataPath string) []string {
	ipData, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		log.Fatalf("[Bulletinboard] failed to read iplist: %v", err)
	}
	return strings.Split(string(ipData), "\n")
}

// New returns a network node structure
func New(degree int, counter int, metadataPath string, Polyyy []poly.Poly) (BulletinBoard, error) {
	//f, _ := os.Create(metadataPath + "/log0")
	//defer f.Close()

	fileName := metadataPath + "/Bulletinboard.logger"
	tmplogger, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		tmplogger, err = os.Create(fileName)
	}
	//os.Truncate(fileName, 0)
	logger := log.New(tmplogger, "", log.LstdFlags)

	if counter < 0 {
		return BulletinBoard{}, errors.New(fmt.Sprintf("counter must be non-negative, got %d", counter))
	}

	//fixedRandState := rand.New(rand.NewSource(int64(3)))
	p := gmp.NewInt(0)
	p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)
	dpc := commitment.DLPolyCommit{}
	dpc.SetupFix(counter)

	ipRaw := ReadIpList(metadataPath + "/ip_list")[0 : counter+1]
	bip := ipRaw[0]
	ipList := ipRaw[1 : counter+1]

	proCnt := 0
	shaCnt := 0
	secretCnt := 0
	reconstructionContent := make([]*pb.Cmt1Msg, counter)
	reconstructionContent2 := make([]*pb.Cmt1Msg, counter)
	reconstructionContent3 := make([]*pb.Cmt1Msg, counter)
	reconstructionContent4 := make([]*pb.Cmt1Msg, counter)
	reconstructSecret := make([]*gmp.Int, counter)
	secret := gmp.NewInt(0)
	//polyp, err := poly.NewRand(degree, fixedRandState, p)
	//if err != nil {
	//	bb.logger.Fatal("Error initializing random poly")
	//}
	for i := 0; i < counter; i++ {
		c := dpc.NewG1()
		dpc.Commit(c, Polyyy[i])
		//fmt.Println(i+1,"witness is ",c)
		//fmt.Println(Polyyy[i].GetDegree())
		cBytes := c.CompressedBytes()
		msg := &pb.Cmt1Msg{
			Index:   int32(i + 1),
			Polycmt: cBytes,
		}
		reconstructionContent4[i] = msg
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
		reconstructionContent2: reconstructionContent2,
		reconstructionContent3: reconstructionContent3,
		reconstructionContent4: reconstructionContent4,
		proactivizationContent: proactivizationContent,
		nConn:                  nConn,
		nClient:                nClient,
		totMsgSize:             &totMsgSize,
		log:                    logger,
	}, nil
}
func (bb *BulletinBoard) GetMessageOfNode(secretid, label int) poly.Poly {

	db := common.GetDB()
	var secretshares []model1.Secretshare
	result := db.Where("secret_id =?", secretid).Where("unit_id", label).Find(&secretshares)
	rowNum := result.RowsAffected
	var newsecret model.Secret
	db.Where("id = ? ", secretid).First(&newsecret)
	degree := newsecret.Degree
	//counter := newsecretshare.Counter
	//secretid := int(newsecretshare.SecretId)
	coeff := make([]*gmp.Int, 2*degree+1)
	for i := 0; int64(i) < rowNum; i++ {
		var newsecretshare model1.Secretshare
		db.Where("secret_id = ? and unit_id = ? and row_num =? ", secretid, label, i).Find(&newsecretshare)
		//Data存放秘密份额,多项式
		Data := newsecretshare.Data
		//fmt.Println(Data)
		coeff[i] = gmp.NewInt(0)
		coeff[i].SetBytes(Data)
	}
	tmpPoly, _ := poly.NewPoly(len(coeff) - 1)
	tmpPoly.SetbyCoeff(coeff)
	return tmpPoly
	//	return a poly
}
func (bb *BulletinBoard) Getmessage(secretid int, degree int, counter int) []poly.Poly {

	polyyy := make([]poly.Poly, counter)
	for i := 0; i < counter; i++ {
		polyyy[i] = bb.GetMessageOfNode(secretid, i+1)
	}
	return polyyy
}
func (bb *BulletinBoard) SetSecret() {
	Polyyy := bb.Getmessage(bb.id, bb.degree, bb.counter)
	for i := 0; i < bb.counter; i++ {
		c := bb.dpc.NewG1()
		bb.dpc.Commit(c, Polyyy[i])
		//fmt.Println(i+1,"witness is ",c)
		//fmt.Println(Polyyy[i].GetDegree())
		cBytes := c.CompressedBytes()
		msg := &pb.Cmt1Msg{
			Index:   int32(i + 1),
			Polycmt: cBytes,
		}
		bb.reconstructionContent4[i] = msg
	}
	return

}
func NewBulletboardForWeb(degree, counter int, metadataPath string, secretid int) (BulletinBoard, error) {

	fileName := metadataPath + "/Secretid" + strconv.Itoa(secretid) + "Bulletinboard.log"
	tmplogger, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	tmplogger = os.Stdout
	if err != nil {
		tmplogger, err = os.Create(fileName)
	}
	//os.Truncate(fileName, 0)
	logger := log.New(tmplogger, "", log.LstdFlags)

	if counter < 0 {
		return BulletinBoard{}, errors.New(fmt.Sprintf("counter must be non-negative, got %d", counter))
	}

	//fixedRandState := rand.New(rand.NewSource(int64(3)))
	p := gmp.NewInt(0)
	p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)
	dpc := commitment.DLPolyCommit{}
	dpc.SetupFix(counter)

	ipRaw := ReadIpList(metadataPath + "/ip_list")[0 : counter+1]
	bip := ipRaw[0]
	ipList := ipRaw[1 : counter+1]

	proCnt := 0
	shaCnt := 0
	secretCnt := 0
	reconstructionContent := make([]*pb.Cmt1Msg, counter)
	reconstructionContent2 := make([]*pb.Cmt1Msg, counter)
	reconstructionContent3 := make([]*pb.Cmt1Msg, counter)
	reconstructionContent4 := make([]*pb.Cmt1Msg, counter)
	reconstructSecret := make([]*gmp.Int, counter)
	secret := gmp.NewInt(0)
	//polyp, err := poly.NewRand(degree, fixedRandState, p)
	//if err != nil {
	//	bb.logger.Fatal("Error initializing random poly")
	//}

	proactivizationContent := make([]*pb.CommitMsg, counter)

	nConn := make([]*grpc.ClientConn, counter)
	nClient := make([]pb.NodeServiceClient, counter)

	totMsgSize := 0
	finish := false

	return BulletinBoard{
		id:                     secretid,
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
		reconstructionContent2: reconstructionContent2,
		reconstructionContent3: reconstructionContent3,
		reconstructionContent4: reconstructionContent4,
		proactivizationContent: proactivizationContent,
		nConn:                  nConn,
		nClient:                nClient,
		totMsgSize:             &totMsgSize,
		log:                    logger,
		dpc:                    &dpc,
		finish:                 finish,
	}, nil
}

func (bb *BulletinBoard) StartReconstruct(ctx context.Context, in *pb.RequestMsg) (*pb.ResponseMsg, error) {

	bb.log.Print("[Bulletinboard] start reconstruction")
	var wg sync.WaitGroup
	for i := 0; i < bb.counter; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			bb.nClient[i].Reconstruct(ctx, &pb.RequestMsg{})
		}(i)
	}
	wg.Wait()
	return &pb.ResponseMsg{}, nil
}

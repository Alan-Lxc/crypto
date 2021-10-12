package control

import (
	"context"
	pb "github.com/Alan-Lxc/crypto_contest/src/service"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"strings"
)

type Control struct {
	metadataPath string
	boardIp      string

	boardConn    *grpc.ClientConn
	BoardService pb.BulletinBoardServiceClient
}

func (c *Control) Connect() {
	boarConn, err := grpc.Dial(c.boardIp, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Fail to connect borad:%v", err)
	}
	c.boardConn = boarConn
	c.BoardService = pb.NewBulletinBoardServiceClient(c.boardConn)
}
func (c *Control) Disconnect() {
	log.Println("Disconnect the board")
	c.boardConn.Close()
}
func (c *Control) StartHandoff() {
	log.Printf("Start to Handoff")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := c.BoardService.StartEpoch(ctx, &pb.RequestMsg{})
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
func New(metadataPath string) (Control, error) {
	bip := ReadIpList(metadataPath)[0]

	return Control{
		metadataPath: metadataPath,
		boardIp:      bip,
	}, nil
}

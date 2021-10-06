package client

import (
	pb "github.com/Alan-Lxc/crypto_contest/src/service"
	"google.golang.org/grpc"
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
	nodeService []*pb.NodeServiceServer
	//boardconn
	boardConn *grpc.ClientConn
	//boardService
	boardService *pb.BulletinBoardServiceServer
	//metadatapath
	metadataPath string
}

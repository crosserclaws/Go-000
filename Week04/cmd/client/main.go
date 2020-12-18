package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"

	pb "week04/api/video/v1"
)

const (
	address         = "localhost:50051"
	defaultID int64 = 0
)

func main() {
	conn := getGrpcConnection()
	defer conn.Close()
	c := pb.NewVideoInformerClient(conn)

	id := getQueryID()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r := getQueryReply(ctx, c, id)

	log.Printf("Video info: name=%v, count=%v", r.GetName(), r.GetCount())
}

func getQueryReply(ctx context.Context, c pb.VideoInformerClient, id int64) *pb.VideoInfoReply {
	r, err := c.GetVideoInfo(ctx, &pb.VideoInfoRequest{Id: id})
	if err != nil {
		log.Fatalf("Failed to get video info: %v", err)
	}
	return r
}

func getQueryID() int64 {
	var id int64 = defaultID
	if len(os.Args) > 1 {
		n, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
			log.Fatalln("Failed to parse int from args[1]:", os.Args[1])
		}
		id = n
	}
	return id
}

func getGrpcConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	return conn
}

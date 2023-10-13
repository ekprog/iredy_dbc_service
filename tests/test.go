package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	pb "microservice/pkg/pb/api"
	"time"
)

func main() {
	addr := "127.0.0.1:8073"
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	cl := pb.NewToDoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := cl.CreateProject(ctx, &pb.CreateProjectRequest{
		Name: "dsfs",
	})
	if err != nil {
		panic(err)
	}
	log.Printf("%v %v", res.Status, res.Id)
}

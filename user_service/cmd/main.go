package main

import (
	pb "github.com/argo-agorshechnikov/gRPC-microservices/user_service/proto"
)

type Server struct {
	pb.UnimplementedUserServiceServer
}

func main() {

}

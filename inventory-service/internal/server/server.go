package server

import (
	"log"
	"net"

	"github.com/shoe-store/inventory-service/internal/service"
	"github.com/shoe-store/inventory-service/proto"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	service    *service.InventoryService
}

func NewServer(service *service.InventoryService) *Server {
	return &Server{
		grpcServer: grpc.NewServer(),
		service:    service,
	}
}

func (s *Server) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	proto.RegisterInventoryServiceServer(s.grpcServer, s.service)

	log.Printf("Starting gRPC server on port %s", port)
	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}

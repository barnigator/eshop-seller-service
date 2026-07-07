package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	port       int
}

func New(port int) *Server {
	grpcServer := grpc.NewServer()
	return &Server{grpcServer, port}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	return s.grpcServer.Serve(listener)
}

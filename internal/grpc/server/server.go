package server

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/barnigator/eshop-seller-service/internal/grpc/interceptor"
	sellerv1 "github.com/barnigator/protos/gen/go/seller/v1"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	port       int
}

func New(port int, handler sellerv1.SellerServiceServer, log *slog.Logger) *Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.RequestID(),
			interceptor.Logging(log),
		),
	)

	sellerv1.RegisterSellerServiceServer(grpcServer, handler)

	return &Server{grpcServer, port}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	sellerv1 "github.com/barnigator/protos/gen/go/seller/v1"
)

func main() {
	cc, err := grpc.NewClient(
		"localhost:44045",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fatal("failed to connect: %v\n", err)
	}
	defer cc.Close()

	client := sellerv1.NewSellerServiceClient(cc)

	req := &sellerv1.GetSellerStatusRequest{
		SellerId: "550e8400-e29b-41d4-a716-446655440000",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetSellerStatus(ctx, req)
	if err != nil {
		fatal("failed to get seller status: %v\n", err)
	}

	fmt.Println("seller status:", resp.Status)
}

func fatal(format string, err error) {
	fmt.Fprintf(os.Stderr, format, err)
	os.Exit(1)
}

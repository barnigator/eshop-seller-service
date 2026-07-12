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

	// тест 1. валидный uuid, продавец есть
	checkSellerStatus(client, "550e8400-e29b-41d4-a716-446655440000")

	// тест 2. валидный uuid, продавца нет
	checkSellerStatus(client, "550e8400-e29b-41d4-a716-446655441111")

	//	тест 3. невалидный uuid
	checkSellerStatus(client, "invalid uuid")

}

func fatal(format string, err error) {
	fmt.Fprintf(os.Stderr, format, err)
	os.Exit(1)
}

func checkSellerStatus(client sellerv1.SellerServiceClient, sellerID string) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &sellerv1.GetSellerStatusRequest{
		SellerId: sellerID,
	}

	resp, err := client.GetSellerStatus(ctx, req)
	if err != nil {
		fmt.Printf("failed to get seller status: %v\n", err)
		return
	}

	fmt.Println("seller status:", resp.Status)

}

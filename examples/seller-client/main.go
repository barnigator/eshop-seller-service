package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

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
	checkGetSellerStatus(client, "550e8400-e29b-41d4-a716-446655440000")

	// тест 2. валидный uuid, продавца нет
	checkGetSellerStatus(client, "550e8400-e29b-41d4-a716-446655441111")

	//	тест 3. невалидный uuid
	checkGetSellerStatus(client, "invalid uuid")

	//	тест 4. создаем продавца
	checkCreateSeller(client, "22211111-1111-1111-1111-111111111111", "Nike", "  cool  ")

	//	тест 5. создаем продавца с тем же id, но с другим брендом
	checkCreateSeller(client, "22211111-1111-1111-1111-111111111111", "  Adidas  ", "nice    ")

	//	тест 6. создаем продавца с тем же id, c тем же брендом
	checkCreateSeller(client, "22211111-1111-1111-1111-111111111111", "adidas", "very nice")

	// тест 7. получаем продавца, который есть
	checkGetSeller(client, "550e8400-e29b-41d4-a716-446655440000")

	//	тест 8. получаем продавца, которого нет
	checkGetSeller(client, "000e8400-e29b-41d4-a716-446655440000")

	//	тест 9. получаем продавца с помощью невалидного uuid
	checkGetSeller(client, "invalid uuid")

	// тест 10. получаем список продавцов по user_id
	checkListSellersByUserID(client, "33311111-1111-1111-1111-111111111111")

	// тест 11. получаем пустой список продавцов по user_id
	checkListSellersByUserID(client, "44411111-1111-1111-1111-111111111111")

	checkUpdateSeller(client, "550e8400-e29b-41d4-a716-446655440000", "New Brand", "",
		[]string{"brand_name"},
	)

	checkUpdateSeller(
		client,
		"550e8400-e29b-41d4-a716-446655440000",
		"",
		"New description",
		[]string{"description"},
	)

	checkUpdateSeller(
		client,
		"550e8400-e29b-41d4-a716-446655440000",
		"Another Brand",
		"Another description",
		[]string{"brand_name", "description"},
	)

	checkUpdateSeller(
		client,
		"550e8400-e29b-41d4-a716-446655440000",
		"",
		"",
		[]string{"description"},
	)

	checkUpdateSeller(
		client,
		"550e8400-e29b-41d4-a716-446655440000",
		"",
		"",
		[]string{},
	)

	checkUpdateSeller(
		client,
		"550e8400-e29b-41d4-a716-446655440000",
		"",
		"",
		[]string{"status"},
	)

	checkUpdateSeller(
		client,
		uuid.NewString(),
		"New Brand",
		"",
		[]string{"brand_name"},
	)

	sellerId := "550e8400-e29b-41d4-a716-446655440000"

	checkArchiveSeller(client, sellerId)

	checkGetSellerStatus(client, sellerId)

	checkArchiveSeller(client, sellerId)

	checkGetSellerStatus(client, sellerId)

	checkArchiveSeller(client, uuid.NewString())

	checkArchiveSeller(client, "")
}

func fatal(format string, err error) {
	fmt.Fprintf(os.Stderr, format, err)
	os.Exit(1)
}

func checkGetSellerStatus(client sellerv1.SellerServiceClient, sellerID string) {

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

func checkCreateSeller(client sellerv1.SellerServiceClient, userID, brandName, description string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &sellerv1.CreateSellerRequest{
		UserId:      userID,
		BrandName:   brandName,
		Description: description,
	}

	resp, err := client.CreateSeller(ctx, req)
	if err != nil {
		fmt.Printf("failed to create seller: %v\n", err)
		return
	}

	fmt.Println("seller created:", resp.Seller)
}

func checkGetSeller(client sellerv1.SellerServiceClient, sellerID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &sellerv1.GetSellerRequest{
		SellerId: sellerID,
	}

	resp, err := client.GetSeller(ctx, req)
	if err != nil {
		fmt.Printf("failed to get seller: %v\n", err)
		return
	}

	fmt.Println("seller:", resp.Seller)
}

func checkListSellersByUserID(client sellerv1.SellerServiceClient, userID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &sellerv1.ListSellersByUserIDRequest{
		UserId: userID,
	}

	resp, err := client.ListSellersByUserID(ctx, req)
	if err != nil {
		fmt.Printf("failed to list sellers: %v\n", err)
		return
	}

	fmt.Println("sellers:", resp.Sellers)
}

func checkUpdateSeller(client sellerv1.SellerServiceClient, sellerID string, brandName string, description string, paths []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &sellerv1.UpdateSellerRequest{
		SellerId:    sellerID,
		BrandName:   brandName,
		Description: description,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: paths,
		},
	}

	resp, err := client.UpdateSeller(ctx, req)
	if err != nil {
		fmt.Printf("failed to update seller: %v\n", err)
		return
	}

	fmt.Println("seller:", resp.Seller)
}

func checkArchiveSeller(client sellerv1.SellerServiceClient, sellerID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &sellerv1.ArchiveSellerRequest{
		SellerId: sellerID,
	}

	empty, err := client.ArchiveSeller(ctx, req)
	if err != nil {
		fmt.Printf("failed to archive seller: %v\n", err)
		return
	}

	if empty == nil {
		fmt.Println("unexpected nil response")
		return
	}

	fmt.Println("seller archived successfully")
}

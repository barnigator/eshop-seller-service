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

const (
	serverAddress  = "localhost:44045"
	requestTimeout = 5 * time.Second
	invalidID      = "invalid uuid"
)

func main() {
	connection, err := grpc.NewClient(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fatal("failed to create gRPC client: %v\n", err)
	}
	defer connection.Close()

	client := sellerv1.NewSellerServiceClient(connection)

	runGetSellerStatusScenarios(client)
	runCreateSellerScenarios(client)
	runGetSellerScenarios(client)
	runListSellersByUserIDScenarios(client)
	runUpdateSellerScenarios(client)
	runArchiveSellerScenarios(client)
	runDeleteSellerScenarios(client)
}

func runGetSellerStatusScenarios(client sellerv1.SellerServiceClient) {
	printSection("GetSellerStatus")

	sellerID, _, err := createTestSeller(client)
	if err != nil {
		fatal("failed to create test seller: %v\n", err)
	}

	// Получение статуса существующего продавца.
	checkGetSellerStatus(client, sellerID)

	// Получение статуса несуществующего продавца.
	checkGetSellerStatus(client, uuid.NewString())

	// Получение статуса по некорректному UUID.
	checkGetSellerStatus(client, invalidID)
}

func runCreateSellerScenarios(client sellerv1.SellerServiceClient) {
	printSection("CreateSeller")

	userID := uuid.NewString()
	firstBrand := uniqueValue("Nike")
	secondBrand := uniqueValue("Adidas")

	// Создание первого продавца пользователя.
	checkCreateSeller(
		client,
		userID,
		firstBrand,
		"  cool  ",
	)

	// Создание второго продавца того же пользователя.
	checkCreateSeller(
		client,
		userID,
		"  "+secondBrand+"  ",
		"nice    ",
	)

	// Попытка создать продавца с уже существующим названием бренда.
	checkCreateSeller(
		client,
		userID,
		secondBrand,
		"very nice",
	)
}

func runGetSellerScenarios(client sellerv1.SellerServiceClient) {
	printSection("GetSeller")

	sellerID, _, err := createTestSeller(client)
	if err != nil {
		fatal("failed to create test seller: %v\n", err)
	}

	// Получение существующего продавца.
	checkGetSeller(client, sellerID)

	// Получение несуществующего продавца.
	checkGetSeller(
		client,
		uuid.NewString(),
	)

	// Получение продавца по некорректному UUID.
	checkGetSeller(client, invalidID)
}

func runListSellersByUserIDScenarios(client sellerv1.SellerServiceClient) {
	printSection("ListSellersByUserID")

	_, userID, err := createTestSeller(client)
	if err != nil {
		fatal("failed to create test seller: %v\n", err)
	}

	// Получение непустого списка продавцов пользователя.
	checkListSellersByUserID(client, userID)

	// Получение пустого списка продавцов пользователя.
	checkListSellersByUserID(client, uuid.NewString())
}

func runUpdateSellerScenarios(client sellerv1.SellerServiceClient) {
	printSection("UpdateSeller")

	sellerID, _, err := createTestSeller(client)
	if err != nil {
		fatal("failed to create test seller: %v\n", err)
	}

	newBrand := uniqueValue("New-Brand")
	anotherBrand := uniqueValue("Another-Brand")

	// Обновление названия бренда.
	checkUpdateSeller(
		client,
		sellerID,
		newBrand,
		"",
		[]string{"brand_name"},
	)

	// Обновление описания.
	checkUpdateSeller(
		client,
		sellerID,
		"",
		"New description",
		[]string{"description"},
	)

	// Одновременное обновление названия бренда и описания.
	checkUpdateSeller(
		client,
		sellerID,
		anotherBrand,
		"Another description",
		[]string{"brand_name", "description"},
	)

	// Очистка описания.
	checkUpdateSeller(
		client,
		sellerID,
		"",
		"",
		[]string{"description"},
	)

	// Попытка обновления с пустым FieldMask.
	checkUpdateSeller(
		client,
		sellerID,
		"",
		"",
		[]string{},
	)

	// Попытка обновления неподдерживаемого поля.
	checkUpdateSeller(
		client,
		sellerID,
		"",
		"",
		[]string{"status"},
	)

	// Попытка обновления несуществующего продавца.
	checkUpdateSeller(
		client,
		uuid.NewString(),
		uniqueValue("Missing-brand"),
		"",
		[]string{"brand_name"},
	)
}

func runArchiveSellerScenarios(client sellerv1.SellerServiceClient) {
	printSection("ArchiveSeller")

	sellerID, _, err := createTestSeller(client)
	if err != nil {
		fatal("failed to create test seller: %v\n", err)
	}
	// Архивирование существующего продавца.
	checkArchiveSeller(client, sellerID)

	// Проверка статуса после архивирования.
	checkGetSellerStatus(client, sellerID)

	// Повторное архивирование уже архивированного продавца.
	checkArchiveSeller(client, sellerID)

	// Проверка статуса после повторного архивирования.
	checkGetSellerStatus(client, sellerID)

	// Попытка архивировать несуществующего продавца.
	checkArchiveSeller(client, uuid.NewString())

	// Попытка архивировать продавца без seller_id.
	checkArchiveSeller(client, "")
}

func runDeleteSellerScenarios(client sellerv1.SellerServiceClient) {
	printSection("DeleteSeller")

	sellerID, _, err := createTestSeller(client)
	if err != nil {
		fatal("failed to create test seller: %v\n", err)
	}
	// Удаление существующего продавца.
	checkDeleteSeller(client, sellerID)

	// Повторное удаление уже удалённого продавца.
	checkDeleteSeller(client, sellerID)

	// Проверка недоступности продавца после удаления.
	checkGetSeller(client, sellerID)

	// Попытка удалить несуществующего продавца.
	checkDeleteSeller(client, uuid.NewString())

	// Попытка удалить продавца без seller_id.
	checkDeleteSeller(client, "")
}

func checkGetSellerStatus(
	client sellerv1.SellerServiceClient,
	sellerID string,
) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
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

func checkCreateSeller(
	client sellerv1.SellerServiceClient,
	userID string,
	brandName string,
	description string,
) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
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

func checkGetSeller(
	client sellerv1.SellerServiceClient,
	sellerID string,
) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
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

func checkListSellersByUserID(
	client sellerv1.SellerServiceClient,
	userID string,
) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
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

func checkUpdateSeller(
	client sellerv1.SellerServiceClient,
	sellerID string,
	brandName string,
	description string,
	paths []string,
) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
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

	fmt.Println("seller updated:", resp.Seller)
}

func checkArchiveSeller(
	client sellerv1.SellerServiceClient,
	sellerID string,
) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req := &sellerv1.ArchiveSellerRequest{
		SellerId: sellerID,
	}

	resp, err := client.ArchiveSeller(ctx, req)
	if err != nil {
		fmt.Printf("failed to archive seller: %v\n", err)
		return
	}

	if resp == nil {
		fmt.Println("failed to archive seller: unexpected nil response")
		return
	}

	fmt.Println("seller archived successfully")
}

func checkDeleteSeller(
	client sellerv1.SellerServiceClient,
	sellerID string,
) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req := &sellerv1.DeleteSellerRequest{
		SellerId: sellerID,
	}

	resp, err := client.DeleteSeller(ctx, req)
	if err != nil {
		fmt.Printf("failed to delete seller: %v\n", err)
		return
	}

	if resp == nil {
		fmt.Println("failed to delete seller: unexpected nil response")
		return
	}

	fmt.Println("seller deleted successfully")
}

func createTestSeller(client sellerv1.SellerServiceClient) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req := &sellerv1.CreateSellerRequest{
		UserId:      uuid.NewString(),
		BrandName:   uniqueValue("Brand"),
		Description: uniqueValue("Description"),
	}

	resp, err := client.CreateSeller(ctx, req)
	if err != nil {
		return "", "", err
	}

	if resp.Seller == nil {
		return "", "", fmt.Errorf("create seller returned nil seller")
	}

	return resp.Seller.Id, resp.Seller.UserId, nil
}

func uniqueValue(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, uuid.NewString()[:8])
}

func printSection(name string) {
	fmt.Printf("\n========== %s ==========\n", name)
}

func fatal(format string, err error) {
	fmt.Fprintf(os.Stderr, format, err)
	os.Exit(1)
}

//go:build integration

package postgres

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestSellerRepository_CreateSeller(t *testing.T) {
	repo := newRepo(t)

	seller := domain.Seller{
		UserID:      uuid.New(),
		BrandName:   "Integration Test Brand",
		Description: "integration test seller",
		Status:      domain.SellerStatusPending,
	}

	sellerResult, err := repo.CreateSeller(context.Background(), seller)
	if err != nil {
		t.Fatalf("create seller: %v", err)
	}

	if sellerResult.ID == uuid.Nil {
		t.Fatalf("seller id must not be nil")
	}

	if sellerResult.UserID != seller.UserID {
		t.Fatalf("unexpected user ID: got %v, want %v", sellerResult.UserID, seller.UserID)
	}

	if sellerResult.BrandName != seller.BrandName {
		t.Fatalf("unexpected brand name: got %q, want %q", sellerResult.BrandName, seller.BrandName)
	}

	if sellerResult.Description != seller.Description {
		t.Fatalf("unexpected description: got %q, want %q", sellerResult.Description, seller.Description)
	}

	if sellerResult.Status != seller.Status {
		t.Fatalf("unexpected status: got %v, want %v", sellerResult.Status, seller.Status)
	}

	storedSeller, err := repo.GetSellerByID(context.Background(), sellerResult.ID)
	if err != nil {
		t.Fatalf("get created seller: %v", err)
	}

	if storedSeller != sellerResult {
		t.Fatalf("unexpected stored seller: got %+v, want %+v", storedSeller, sellerResult)
	}
}

func TestSellerRepository_CreateSeller_BrandAlreadyExists(t *testing.T) {
	repo := newRepo(t)

	userUUID := uuid.New()

	seller1 := domain.Seller{
		UserID:      userUUID,
		BrandName:   "Integration Test Brand",
		Description: "integration test seller1",
		Status:      domain.SellerStatusPending,
	}

	seller2 := domain.Seller{
		UserID:      userUUID,
		BrandName:   "Integration Test Brand",
		Description: "integration test seller2",
		Status:      domain.SellerStatusPending,
	}

	_, err := repo.CreateSeller(context.Background(), seller1)
	if err != nil {
		t.Fatalf("create seller: %v", err)
	}

	_, err = repo.CreateSeller(context.Background(), seller2)
	if !errors.Is(err, domain.ErrBrandAlreadyExists) {
		t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrBrandAlreadyExists)
	}
}

func TestSellerRepository_GetSellerByID(t *testing.T) {
	repo := newRepo(t)

	seller := domain.Seller{
		UserID:      uuid.New(),
		BrandName:   "Integration Test Brand",
		Description: "integration test seller",
		Status:      domain.SellerStatusPending,
	}

	sellerResult, err := repo.CreateSeller(context.Background(), seller)
	if err != nil {
		t.Fatalf("create seller: %v", err)
	}

	storedSeller, err := repo.GetSellerByID(context.Background(), sellerResult.ID)
	if err != nil {
		t.Fatalf("get seller by ID: %v", err)
	}

	if storedSeller != sellerResult {
		t.Fatalf("unexpected seller: got %+v, want %+v", storedSeller, sellerResult)
	}
}

func TestSellerRepository_GetSellerByID_NotFound(t *testing.T) {
	repo := newRepo(t)

	_, err := repo.GetSellerByID(context.Background(), uuid.New())
	if !errors.Is(err, domain.ErrSellerNotFound) {
		t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrSellerNotFound)
	}
}

func TestSellerRepository_ListSellersByUserID(t *testing.T) {
	repo := newRepo(t)

	userUUID := uuid.New()

	sellers := []domain.Seller{
		{
			UserID:      userUUID,
			BrandName:   "Adidas",
			Description: "integration test seller1",
			Status:      domain.SellerStatusPending,
		},
		{
			UserID:      userUUID,
			BrandName:   "Puma",
			Description: "integration test seller2",
			Status:      domain.SellerStatusPending,
		},
		{
			UserID:      userUUID,
			BrandName:   "Nike",
			Description: "integration test seller3",
			Status:      domain.SellerStatusPending,
		},
		{
			UserID:      userUUID,
			BrandName:   "Demix",
			Description: "integration test seller4",
			Status:      domain.SellerStatusPending,
		},
	}

	createdSellers := make(map[uuid.UUID]domain.Seller, len(sellers))

	for _, seller := range sellers {
		createdSeller, err := repo.CreateSeller(context.Background(), seller)
		if err != nil {
			t.Fatalf("create seller %q: %v", seller.BrandName, err)
		}

		createdSellers[createdSeller.ID] = createdSeller
	}

	_, err := repo.CreateSeller(context.Background(), domain.Seller{
		UserID:      uuid.New(),
		BrandName:   "Adidas",
		Description: "must not be returned",
		Status:      domain.SellerStatusPending,
	})
	if err != nil {
		t.Fatalf("create seller for another user: %v", err)
	}

	sellersResult, err := repo.ListSellersByUserID(context.Background(), userUUID)
	if err != nil {
		t.Fatalf("list sellers by user ID: %v", err)
	}

	if len(sellersResult) != len(sellers) {
		t.Fatalf("unexpected sellers amount: got %d, want %d", len(sellersResult), len(sellers))
	}

	for _, got := range sellersResult {
		want, ok := createdSellers[got.ID]
		if !ok {
			t.Fatalf("unexpected seller returned: %+v", got)
		}

		if got != want {
			t.Fatalf("unexpected seller: got %+v, want %+v",
				got,
				want)
		}
	}
}

func TestSellerRepository_ListSellersByUserID_Empty(t *testing.T) {
	repo := newRepo(t)

	sellers, err := repo.ListSellersByUserID(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("list sellers by user id: %v", err)
	}

	if len(sellers) != 0 {
		t.Fatalf("unexpected sellers amount: got %d, want 0", len(sellers))
	}
}

func TestSellerRepository_UpdateSeller(t *testing.T) {
	repo := newRepo(t)

	seller := domain.Seller{
		UserID:      uuid.New(),
		BrandName:   "Integration Test Brand",
		Description: "integration test seller",
		Status:      domain.SellerStatusPending,
	}

	sellerCreated, err := repo.CreateSeller(context.Background(), seller)
	if err != nil {
		t.Fatalf("create seller: %v", err)
	}

	brandName := "Adidas"
	description := "New description"

	sellerUpdated, err := repo.UpdateSeller(context.Background(), sellerCreated.ID, &brandName, &description)
	if err != nil {
		t.Fatalf("update seller: %v", err)
	}

	if sellerUpdated.BrandName != brandName {
		t.Fatalf("unexpected brand name: got %v, want %v", sellerUpdated.BrandName, brandName)
	}

	if sellerUpdated.Description != description {
		t.Fatalf("unexpected description: got %v, want %v", sellerUpdated.Description, description)
	}

	sellerGot, err := repo.GetSellerByID(context.Background(), sellerCreated.ID)
	if err != nil {
		t.Fatalf("get seller: %v", err)
	}

	if sellerGot.BrandName != brandName {
		t.Fatalf("unexpected brand name: got %v, want %v", sellerGot.BrandName, brandName)
	}

	if sellerGot.Description != description {
		t.Fatalf("unexpected description: got %v, want %v", sellerGot.Description, description)
	}
}

func TestSellerRepository_UpdateSeller_BrandName(t *testing.T) {
	repo := newRepo(t)

	seller := domain.Seller{
		UserID:      uuid.New(),
		BrandName:   "Integration Test Brand",
		Description: "integration test seller",
		Status:      domain.SellerStatusPending,
	}

	sellerCreated, err := repo.CreateSeller(context.Background(), seller)
	if err != nil {
		t.Fatalf("create seller: %v", err)
	}

	brandName := "New name"
	var description *string

	sellerUpdated, err := repo.UpdateSeller(context.Background(), sellerCreated.ID, &brandName, description)
	if err != nil {
		t.Fatalf("update seller: %v", err)
	}

	if sellerUpdated.BrandName != brandName {
		t.Fatalf("unexpected brand name: got %v, want %v", sellerUpdated.BrandName, brandName)
	}

	if sellerUpdated.Description != sellerCreated.Description {
		t.Fatalf("unexpected description: got %v, want %v", sellerUpdated.Description, sellerCreated.Description)
	}
}

func TestSellerRepository_UpdateSeller_Description(t *testing.T) {
	repo := newRepo(t)

	seller := domain.Seller{
		UserID:      uuid.New(),
		BrandName:   "Integration Test Brand",
		Description: "integration test seller",
		Status:      domain.SellerStatusPending,
	}

	sellerCreated, err := repo.CreateSeller(context.Background(), seller)
	if err != nil {
		t.Fatalf("create seller: %v", err)
	}

	var brandName *string
	description := "New description"

	sellerUpdated, err := repo.UpdateSeller(context.Background(), sellerCreated.ID, brandName, &description)
	if err != nil {
		t.Fatalf("update seller: %v", err)
	}

	if sellerUpdated.BrandName != sellerCreated.BrandName {
		t.Fatalf("unexpected brand name: got %v, want %v", sellerUpdated.BrandName, sellerCreated.BrandName)
	}

	if sellerUpdated.Description != description {
		t.Fatalf("unexpected description: got %v, want %v", sellerUpdated.Description, description)
	}
}

func TestSellerRepository_UpdateSeller_ClearDescription(t *testing.T) {
	repo := newRepo(t)

	seller := domain.Seller{
		UserID:      uuid.New(),
		BrandName:   "Integration Test Brand",
		Description: "integration test seller",
		Status:      domain.SellerStatusPending,
	}

	sellerCreated, err := repo.CreateSeller(context.Background(), seller)
	if err != nil {
		t.Fatalf("create seller: %v", err)
	}

	var brandName *string
	description := ""

	sellerUpdated, err := repo.UpdateSeller(context.Background(), sellerCreated.ID, brandName, &description)
	if err != nil {
		t.Fatalf("update seller: %v", err)
	}

	if sellerUpdated.BrandName != sellerCreated.BrandName {
		t.Fatalf("unexpected brand name: got %v, want %v", sellerUpdated.BrandName, sellerCreated.BrandName)
	}

	if sellerUpdated.Description != description {
		t.Fatalf("unexpected description: got %v, want %v", sellerUpdated.Description, description)
	}
}

func TestSellerRepository_UpdateSeller_NotFound(t *testing.T) {
	repo := newRepo(t)

	brandName := "test brand"
	description := "test description"

	_, err := repo.UpdateSeller(context.Background(), uuid.New(), &brandName, &description)
	if !errors.Is(err, domain.ErrSellerNotFound) {
		t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrSellerNotFound)
	}
}

func TestSellerRepository_UpdateSeller_BrandAlreadyExists(t *testing.T) {
	repo := newRepo(t)

	userUUID := uuid.New()

	seller1 := domain.Seller{
		UserID:      userUUID,
		BrandName:   "Adidas",
		Description: "integration test seller1",
		Status:      domain.SellerStatusPending,
	}

	seller2 := domain.Seller{
		UserID:      userUUID,
		BrandName:   "Puma",
		Description: "integration test seller2",
		Status:      domain.SellerStatusPending,
	}

	seller1Created, err := repo.CreateSeller(context.Background(), seller1)
	if err != nil {
		t.Fatalf("create seller: %v", err)
	}

	seller2Created, err := repo.CreateSeller(context.Background(), seller2)
	if err != nil {
		t.Fatalf("create seller: %v", err)
	}

	_, err = repo.UpdateSeller(context.Background(), seller2Created.ID, &seller1Created.BrandName, nil)
	if !errors.Is(err, domain.ErrBrandAlreadyExists) {
		t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrBrandAlreadyExists)
	}
}

func TestSellerRepository_ArchiveSeller(t *testing.T) {
	repo := newRepo(t)

	seller := domain.Seller{
		UserID:      uuid.New(),
		BrandName:   "Adidas",
		Description: "cool brand",
		Status:      domain.SellerStatusPending,
	}

	sellerCreated, err := repo.CreateSeller(context.Background(), seller)
	if err != nil {
		t.Fatalf("create seller: %v", err)
	}

	err = repo.ArchiveSeller(context.Background(), sellerCreated.ID)
	if err != nil {
		t.Fatalf("archive seller: %v", err)
	}

	archivedSeller, err := repo.GetSellerByID(context.Background(), sellerCreated.ID)
	if err != nil {
		t.Fatalf("get seller by id: %v", err)
	}

	if archivedSeller.Status != domain.SellerStatusArchived {
		t.Fatalf("unexpected seller status: got %v, want %v", archivedSeller.Status, domain.SellerStatusArchived)
	}
}

func TestSellerRepository_ArchiveSeller_Double(t *testing.T) {
	repo := newRepo(t)

	seller := domain.Seller{
		UserID:      uuid.New(),
		BrandName:   "Adidas",
		Description: "cool brand",
		Status:      domain.SellerStatusPending,
	}

	sellerCreated, err := repo.CreateSeller(context.Background(), seller)
	if err != nil {
		t.Fatalf("create seller: %v", err)
	}

	err = repo.ArchiveSeller(context.Background(), sellerCreated.ID)
	if err != nil {
		t.Fatalf("first attempt archive seller: %v", err)
	}

	err = repo.ArchiveSeller(context.Background(), sellerCreated.ID)
	if err != nil {
		t.Fatalf("second attempt archive seller: %v", err)
	}

	archivedSeller, err := repo.GetSellerByID(context.Background(), sellerCreated.ID)
	if err != nil {
		t.Fatalf("get seller by id: %v", err)
	}

	if archivedSeller.Status != domain.SellerStatusArchived {
		t.Fatalf("unexpected seller status: got %v, want %v", archivedSeller.Status, domain.SellerStatusArchived)
	}
}

func TestSellerRepository_ArchiveSeller_NotFound(t *testing.T) {
	repo := newRepo(t)

	err := repo.ArchiveSeller(context.Background(), uuid.New())
	if !errors.Is(err, domain.ErrSellerNotFound) {
		t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrSellerNotFound)
	}
}

func TestSellerRepository_ArchiveSeller_DeletedSeller(t *testing.T) {
	repo := newRepo(t)

	seller := domain.Seller{
		UserID:      uuid.New(),
		BrandName:   "Adidas",
		Description: "cool brand",
		Status:      domain.SellerStatusPending,
	}

	createdSeller, err := repo.CreateSeller(context.Background(), seller)
	if err != nil {
		t.Fatalf("create seller: %v", err)
	}

	err = repo.DeleteSeller(context.Background(), createdSeller.ID)
	if err != nil {
		t.Fatalf("delete seller: %v", err)
	}

	err = repo.ArchiveSeller(context.Background(), createdSeller.ID)
	if !errors.Is(err, domain.ErrSellerNotFound) {
		t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrSellerNotFound)
	}
}

func TestSellerRepository_DeleteSeller(t *testing.T) {
	repo := newRepo(t)

	seller := domain.Seller{
		UserID:      uuid.New(),
		BrandName:   "Adidas",
		Description: "cool brand",
		Status:      domain.SellerStatusPending,
	}

	createdSeller, err := repo.CreateSeller(context.Background(), seller)
	if err != nil {
		t.Fatalf("create seller: %v", err)
	}

	err = repo.DeleteSeller(context.Background(), createdSeller.ID)
	if err != nil {
		t.Fatalf("delete seller: %v", err)
	}

	_, err = repo.GetSellerByID(context.Background(), createdSeller.ID)
	if !errors.Is(err, domain.ErrSellerNotFound) {
		t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrSellerNotFound)
	}
}

func TestSellerRepository_DeleteSeller_DoubleDelete(t *testing.T) {
	repo := newRepo(t)

	seller := domain.Seller{
		UserID:      uuid.New(),
		BrandName:   "Adidas",
		Description: "cool brand",
		Status:      domain.SellerStatusPending,
	}

	createdSeller, err := repo.CreateSeller(context.Background(), seller)
	if err != nil {
		t.Fatalf("create seller: %v", err)
	}

	err = repo.DeleteSeller(context.Background(), createdSeller.ID)
	if err != nil {
		t.Fatalf("delete seller: %v", err)
	}

	err = repo.DeleteSeller(context.Background(), createdSeller.ID)
	if !errors.Is(err, domain.ErrSellerNotFound) {
		t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrSellerNotFound)
	}
}

func TestSellerRepository_DeleteSeller_NotFound(t *testing.T) {
	repo := newRepo(t)

	err := repo.DeleteSeller(context.Background(), uuid.New())
	if !errors.Is(err, domain.ErrSellerNotFound) {
		t.Fatalf("unexpected error: got %v, want %v", err, domain.ErrSellerNotFound)
	}
}

func newRepo(t *testing.T) *SellerRepository {
	t.Helper()

	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_TEST_DSN is not set")
	}

	pool, err := NewPool(context.Background(), dsn)
	if err != nil {
		t.Fatalf("connect test database: %v", err)
	}
	t.Cleanup(pool.Close)

	truncateSellers(t, pool)

	return New(pool)
}

func truncateSellers(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(context.Background(), "TRUNCATE TABLE sellers")
	if err != nil {
		t.Fatalf("truncate sellers table: %v", err)
	}
}

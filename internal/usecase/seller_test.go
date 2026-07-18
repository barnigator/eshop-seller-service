package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	"github.com/google/uuid"
)

type fakeSellerRepository struct {
	sellers []domain.Seller
	seller  domain.Seller
	called  bool
	err     error

	receivedSellerID    uuid.UUID
	receivedSeller      domain.Seller
	receivedBrandName   *string
	receivedDescription *string
}

func (f *fakeSellerRepository) GetSellerByID(_ context.Context, sellerID uuid.UUID) (domain.Seller, error) {
	f.called = true
	f.receivedSellerID = sellerID

	return f.seller, f.err
}

func (f *fakeSellerRepository) CreateSeller(_ context.Context, seller domain.Seller) (domain.Seller, error) {
	f.receivedSeller = seller
	f.called = true

	return f.seller, f.err
}

func (f *fakeSellerRepository) ListSellersByUserID(_ context.Context, sellerID uuid.UUID) ([]domain.Seller, error) {
	f.receivedSellerID = sellerID
	f.called = true

	return f.sellers, f.err
}

func (f *fakeSellerRepository) UpdateSeller(_ context.Context, sellerID uuid.UUID, brandName *string, description *string) (domain.Seller, error) {
	f.receivedSellerID = sellerID
	f.receivedBrandName = brandName
	f.receivedDescription = description
	f.called = true

	return f.seller, f.err
}

func (f *fakeSellerRepository) ArchiveSeller(_ context.Context, sellerID uuid.UUID) error {
	f.receivedSellerID = sellerID
	f.called = true

	return f.err
}

func (f *fakeSellerRepository) DeleteSeller(_ context.Context, sellerID uuid.UUID) error {
	f.receivedSellerID = sellerID
	f.called = true

	return f.err
}

func TestSellerUseCase_GetSeller(t *testing.T) {
	tests := []struct {
		name               string
		sellerID           string
		repositoryData     domain.Seller
		repositoryErr      error
		expectedSeller     domain.Seller
		expectedErr        error
		expectedRepoCalled bool
	}{
		{
			name:     "success",
			sellerID: "550e8400-e29b-41d4-a716-446655440000",
			repositoryData: domain.Seller{
				ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				BrandName: "Adidas",
				Status:    domain.SellerStatusActive,
			},
			repositoryErr: nil,
			expectedSeller: domain.Seller{
				ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				BrandName: "Adidas",
				Status:    domain.SellerStatusActive,
			},
			expectedErr:        nil,
			expectedRepoCalled: true,
		},
		{
			name:               "empty seller_id",
			sellerID:           "",
			expectedErr:        domain.ErrSellerIDRequired,
			expectedRepoCalled: false,
		},
		{
			name:               "invalid seller_id",
			sellerID:           "invalid uuid",
			expectedErr:        domain.ErrInvalidSellerID,
			expectedRepoCalled: false,
		},
		{
			name:               "repository error",
			sellerID:           "550e8400-e29b-41d4-a716-446655440000",
			repositoryErr:      domain.ErrSellerNotFound,
			expectedErr:        domain.ErrSellerNotFound,
			expectedRepoCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeSellerRepository{
				seller: tt.repositoryData,
				err:    tt.repositoryErr,
			}

			uc := New(repo)

			seller, err := uc.GetSeller(context.Background(), tt.sellerID)

			if repo.called != tt.expectedRepoCalled {
				t.Fatalf("unexpected repository call state: got %v, want %v", repo.called, tt.expectedRepoCalled)
			}

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("unexpected error: got %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if seller.ID != tt.expectedSeller.ID {
				t.Fatalf("unexpected seller ID: got %v, want %v", seller.ID, tt.expectedSeller.ID)
			}

			if seller.BrandName != tt.expectedSeller.BrandName {
				t.Fatalf("unexpected seller BrandName: got %v, want %v", seller.BrandName, tt.expectedSeller.BrandName)
			}

			if seller.Status != tt.expectedSeller.Status {
				t.Fatalf("unexpected seller Status: got %v, want %v", seller.Status, tt.expectedSeller.Status)
			}

		})
	}
}

func TestSellerUseCase_CreateSeller(t *testing.T) {
	tests := []struct {
		name                    string
		userID                  string
		brandName               string
		description             string
		createResult            domain.Seller
		createErr               error
		expectedSeller          domain.Seller
		expectedErr             error
		expectedRepositoryInput domain.Seller
		expectedRepoCalled      bool
	}{
		{
			name:        "success",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			brandName:   "Adidas",
			description: "cool brand",
			createResult: domain.Seller{
				ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				BrandName:   "Adidas",
				Description: "cool brand",
				Status:      domain.SellerStatusPending,
			},
			createErr: nil,
			expectedSeller: domain.Seller{
				ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				BrandName:   "Adidas",
				Description: "cool brand",
				Status:      domain.SellerStatusPending,
			},
			expectedErr: nil,
			expectedRepositoryInput: domain.Seller{
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				BrandName:   "Adidas",
				Description: "cool brand",
				Status:      domain.SellerStatusPending,
			},
			expectedRepoCalled: true,
		},
		{
			name:               "empty user_id",
			userID:             "",
			expectedErr:        domain.ErrUserIDRequired,
			expectedRepoCalled: false,
		},
		{
			name:               "invalid user_id",
			userID:             "invalid user_id",
			expectedErr:        domain.ErrInvalidUserID,
			expectedRepoCalled: false,
		},
		{
			name:               "empty brand_name",
			userID:             "550e8400-e29b-41d4-a716-446655440000",
			brandName:          "",
			expectedErr:        domain.ErrBrandNameRequired,
			expectedRepoCalled: false,
		},
		{
			name:               "space brand_name",
			userID:             "550e8400-e29b-41d4-a716-446655440000",
			brandName:          "       	  ",
			expectedErr:        domain.ErrBrandNameRequired,
			expectedRepoCalled: false,
		},
		{
			name:               "too long brand_name",
			userID:             "550e8400-e29b-41d4-a716-446655440000",
			brandName:          strings.Repeat("a", 121),
			expectedErr:        domain.ErrBrandNameTooLong,
			expectedRepoCalled: false,
		},
		{
			name:        "brand_name_120_symbols_long",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			brandName:   strings.Repeat("a", 120),
			description: "cool brand",
			createResult: domain.Seller{
				ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				BrandName:   strings.Repeat("a", 120),
				Description: "cool brand",
				Status:      domain.SellerStatusPending,
			},
			createErr: nil,
			expectedSeller: domain.Seller{
				ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				BrandName:   strings.Repeat("a", 120),
				Description: "cool brand",
				Status:      domain.SellerStatusPending,
			},
			expectedErr: nil,
			expectedRepositoryInput: domain.Seller{
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				BrandName:   strings.Repeat("a", 120),
				Description: "cool brand",
				Status:      domain.SellerStatusPending,
			},
			expectedRepoCalled: true,
		},
		{
			name:        "repository returns brand already exists",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			brandName:   "Adidas",
			description: "cool brand",
			createErr:   domain.ErrBrandAlreadyExists,
			expectedErr: domain.ErrBrandAlreadyExists,
			expectedRepositoryInput: domain.Seller{
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				BrandName:   "Adidas",
				Description: "cool brand",
				Status:      domain.SellerStatusPending,
			},
			expectedRepoCalled: true,
		},
		{
			name:        "normalize check",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			brandName:   " Adidas	  ",
			description: "  cool brand   ",
			createResult: domain.Seller{
				ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				BrandName:   "Adidas",
				Description: "cool brand",
				Status:      domain.SellerStatusPending,
			},
			createErr: nil,
			expectedSeller: domain.Seller{
				ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				BrandName:   "Adidas",
				Description: "cool brand",
				Status:      domain.SellerStatusPending,
			},
			expectedErr: nil,
			expectedRepositoryInput: domain.Seller{
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				BrandName:   "Adidas",
				Description: "cool brand",
				Status:      domain.SellerStatusPending,
			},
			expectedRepoCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeSellerRepository{
				seller: tt.createResult,
				err:    tt.createErr,
			}

			uc := New(repo)

			seller, err := uc.CreateSeller(context.Background(), tt.userID, tt.brandName, tt.description)

			if repo.called != tt.expectedRepoCalled {
				t.Fatalf("unexpected repository call state: got %v, want %v", repo.called, tt.expectedRepoCalled)
			}

			if tt.expectedRepoCalled && repo.receivedSeller != tt.expectedRepositoryInput {
				t.Fatalf("unexpected seller RepositoryInput: got %v, want %v", repo.receivedSeller, tt.expectedRepositoryInput)
			}

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("unexpected error: got %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if seller.ID != tt.expectedSeller.ID {
				t.Fatalf("unexpected seller ID: got %v, want %v", seller.ID, tt.expectedSeller.ID)
			}

			if seller.UserID != tt.expectedSeller.UserID {
				t.Fatalf(
					"unexpected seller UserID: got %v, want %v",
					seller.UserID,
					tt.expectedSeller.UserID,
				)
			}

			if seller.BrandName != tt.expectedSeller.BrandName {
				t.Fatalf("unexpected seller BrandName: got %v, want %v", seller.BrandName, tt.expectedSeller.BrandName)
			}

			if seller.Description != tt.expectedSeller.Description {
				t.Fatalf("unexpected seller Description: got %v, want %v", seller.Description, tt.expectedSeller.Description)
			}

			if seller.Status != tt.expectedSeller.Status {
				t.Fatalf("unexpected seller Status: got %v, want %v", seller.Status, tt.expectedSeller.Status)
			}
		})
	}
}

func TestSellerUseCase_GetSellerStatus(t *testing.T) {
	var randomErr = errors.New("database unavailable")
	tests := []struct {
		name               string
		sellerID           string
		repositorySeller   domain.Seller
		repositoryErr      error
		expectedStatus     domain.SellerStatus
		expectedErr        error
		expectedRepoCalled bool
	}{
		{
			name:     "success",
			sellerID: "550e8400-e29b-41d4-a716-446655440000",
			repositorySeller: domain.Seller{
				ID:     uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				Status: domain.SellerStatusActive,
			},
			repositoryErr:      nil,
			expectedStatus:     domain.SellerStatusActive,
			expectedErr:        nil,
			expectedRepoCalled: true,
		},
		{
			name:               "empty seller_id",
			sellerID:           "",
			expectedErr:        domain.ErrSellerIDRequired,
			expectedRepoCalled: false,
		},
		{
			name:               "invalid seller_id",
			sellerID:           "invalid uuid",
			expectedErr:        domain.ErrInvalidSellerID,
			expectedRepoCalled: false,
		},
		{
			name:               "seller not found",
			sellerID:           "550e8400-e29b-41d4-a716-446655440000",
			repositoryErr:      domain.ErrSellerNotFound,
			expectedErr:        domain.ErrSellerNotFound,
			expectedRepoCalled: true,
		},
		{
			name:               "random error",
			sellerID:           "550e8400-e29b-41d4-a716-446655440000",
			repositoryErr:      randomErr,
			expectedErr:        randomErr,
			expectedRepoCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeSellerRepository{
				seller: tt.repositorySeller,
				err:    tt.repositoryErr,
			}

			uc := New(repo)

			status, err := uc.GetSellerStatus(context.Background(), tt.sellerID)

			if repo.called != tt.expectedRepoCalled {
				t.Fatalf("unexpected repository call state: got %v, want %v", repo.called, tt.expectedRepoCalled)
			}

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("unexpected error: got %v, want %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if status != tt.expectedStatus {
				t.Fatalf("unexpected seller status: got %v, want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestSellerUseCase_ListSellersByUserID(t *testing.T) {
	var randomErr = errors.New("database unavailable")

	tests := []struct {
		name               string
		userID             string
		repositorySellers  []domain.Seller
		repositoryErr      error
		expectedSellers    []domain.Seller
		expectedErr        error
		expectedRepoCalled bool
	}{
		{
			name:   "success with few brands",
			userID: "33311111-1111-1111-1111-111111111111",
			repositorySellers: []domain.Seller{
				{
					ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					UserID:      uuid.MustParse("33311111-1111-1111-1111-111111111111"),
					BrandName:   "Adidas",
					Description: "cool brand",
					Status:      domain.SellerStatusActive,
				},
				{
					ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
					UserID:      uuid.MustParse("33311111-1111-1111-1111-111111111111"),
					BrandName:   "Nike",
					Description: "cool brand",
					Status:      domain.SellerStatusActive,
				},
				{
					ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
					UserID:      uuid.MustParse("33311111-1111-1111-1111-111111111111"),
					BrandName:   "Puma",
					Description: "cool brand",
					Status:      domain.SellerStatusActive,
				},
			},
			repositoryErr: nil,
			expectedSellers: []domain.Seller{
				{
					ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					UserID:      uuid.MustParse("33311111-1111-1111-1111-111111111111"),
					BrandName:   "Adidas",
					Description: "cool brand",
					Status:      domain.SellerStatusActive,
				},
				{
					ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
					UserID:      uuid.MustParse("33311111-1111-1111-1111-111111111111"),
					BrandName:   "Nike",
					Description: "cool brand",
					Status:      domain.SellerStatusActive,
				},
				{
					ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
					UserID:      uuid.MustParse("33311111-1111-1111-1111-111111111111"),
					BrandName:   "Puma",
					Description: "cool brand",
					Status:      domain.SellerStatusActive,
				},
			},
			expectedErr:        nil,
			expectedRepoCalled: true,
		},
		{
			name:               "success with empty list",
			userID:             "44411111-1111-1111-1111-111111111111",
			repositorySellers:  []domain.Seller{},
			repositoryErr:      nil,
			expectedSellers:    []domain.Seller{},
			expectedErr:        nil,
			expectedRepoCalled: true,
		},
		{
			name:               "empty user_id",
			userID:             "",
			expectedErr:        domain.ErrUserIDRequired,
			expectedRepoCalled: false,
		},
		{
			name:               "invalid user_id",
			userID:             "invalid user_id",
			expectedErr:        domain.ErrInvalidUserID,
			expectedRepoCalled: false,
		},
		{
			name:               "random repository err",
			userID:             "55511111-1111-1111-1111-111111111111",
			repositoryErr:      randomErr,
			expectedErr:        randomErr,
			expectedRepoCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeSellerRepository{
				sellers: tt.repositorySellers,
				err:     tt.repositoryErr,
			}

			uc := New(repo)

			sellers, err := uc.ListSellersByUserID(context.Background(), tt.userID)

			if repo.called != tt.expectedRepoCalled {
				t.Fatalf("unexpected repository call state: got %v, want %v", repo.called, tt.expectedRepoCalled)
			}

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("unexpected error: got %v, want %v", err, tt.expectedErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(sellers) != len(tt.expectedSellers) {
				t.Fatalf("unexpected sellers amount: got %v, want %v", len(sellers), len(tt.expectedSellers))
			}

			for i := range tt.expectedSellers {
				got := sellers[i]
				want := tt.expectedSellers[i]

				if got != want {
					t.Fatalf(
						"unexpected seller at index %d: got %+v, want %+v",
						i,
						got,
						want,
					)
				}
			}
		})
	}
}

func TestSellerUseCase_UpdateSeller(t *testing.T) {
	tests := []struct {
		name                string
		sellerID            string
		brandName           *string
		description         *string
		expectedBrandName   *string
		expectedDescription *string
		expectedErr         error
		expectedRepoCalled  bool
	}{
		{
			name:               "success brand",
			sellerID:           "550e8400-e29b-41d4-a716-446655440002",
			brandName:          strPtr("Adidas"),
			expectedBrandName:  strPtr("Adidas"),
			expectedRepoCalled: true,
		},
		{
			name:                "success description",
			sellerID:            "550e8400-e29b-41d4-a716-446655440002",
			description:         strPtr("New description"),
			expectedDescription: strPtr("New description"),
			expectedRepoCalled:  true,
		},
		{
			name:                "success both",
			sellerID:            "550e8400-e29b-41d4-a716-446655440002",
			brandName:           strPtr("Adidas"),
			description:         strPtr("New description"),
			expectedBrandName:   strPtr("Adidas"),
			expectedDescription: strPtr("New description"),
			expectedRepoCalled:  true,
		},
		{
			name:               "empty seller id",
			sellerID:           "",
			expectedErr:        domain.ErrSellerIDRequired,
			expectedRepoCalled: false,
		},
		{
			name:               "invalid seller id",
			sellerID:           "invalid uuid",
			expectedErr:        domain.ErrInvalidSellerID,
			expectedRepoCalled: false,
		},
		{
			name:               "nil fields",
			sellerID:           "550e8400-e29b-41d4-a716-446655440002",
			expectedErr:        domain.ErrNoFieldsToUpdate,
			expectedRepoCalled: false,
		},
		{
			name:               "brand empty",
			sellerID:           "550e8400-e29b-41d4-a716-446655440002",
			brandName:          strPtr(""),
			expectedErr:        domain.ErrBrandNameRequired,
			expectedRepoCalled: false,
		},
		{
			name:               "brand too long",
			sellerID:           "550e8400-e29b-41d4-a716-446655440002",
			brandName:          strPtr(strings.Repeat("a", 121)),
			expectedErr:        domain.ErrBrandNameTooLong,
			expectedRepoCalled: false,
		},
		{
			name:               "brand space",
			sellerID:           "550e8400-e29b-41d4-a716-446655440002",
			brandName:          strPtr("  	"),
			expectedErr:        domain.ErrBrandNameRequired,
			expectedRepoCalled: false,
		},
		{
			name:               "normalize brand",
			sellerID:           "550e8400-e29b-41d4-a716-446655440002",
			brandName:          strPtr("  Adidas	 	"),
			expectedBrandName:  strPtr("Adidas"),
			expectedRepoCalled: true,
		},
		{
			name:                "normalize description",
			sellerID:            "550e8400-e29b-41d4-a716-446655440002",
			description:         strPtr(" 	 New description	  "),
			expectedDescription: strPtr("New description"),
			expectedRepoCalled:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeSellerRepository{}

			uc := New(repo)

			_, err := uc.UpdateSeller(context.Background(), tt.sellerID, tt.brandName, tt.description)

			if repo.called != tt.expectedRepoCalled {
				t.Fatalf("unexpected repository call state: got %v, want %v", repo.called, tt.expectedRepoCalled)
			}

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("unexpected error: got %v, want %v", err, tt.expectedErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectedBrandName != nil && *repo.receivedBrandName != *tt.expectedBrandName {
				t.Fatalf("unexpected brand name: got %v, want %v", repo.receivedBrandName, tt.expectedBrandName)
			}

			if tt.expectedDescription != nil && *repo.receivedDescription != *tt.expectedDescription {
				t.Fatalf("unexpected description: got %v, want %v", repo.receivedDescription, tt.expectedDescription)
			}

		})
	}
}

func TestSellerUseCase_ArchiveSeller(t *testing.T) {
	sellerUUID := uuid.New()
	randomErr := errors.New("database error")
	tests := []struct {
		name               string
		sellerID           string
		repositoryErr      error
		expectedUUID       uuid.UUID
		expectedErr        error
		expectedRepoCalled bool
	}{
		{
			name:               "success",
			sellerID:           sellerUUID.String(),
			expectedUUID:       sellerUUID,
			expectedRepoCalled: true,
		},
		{
			name:               "empty seller_id",
			sellerID:           "",
			expectedErr:        domain.ErrSellerIDRequired,
			expectedRepoCalled: false,
		},
		{
			name:               "invalid uuid",
			sellerID:           "invalid uuid",
			expectedErr:        domain.ErrInvalidSellerID,
			expectedRepoCalled: false,
		},
		{
			name:               "seller not found",
			sellerID:           sellerUUID.String(),
			repositoryErr:      domain.ErrSellerNotFound,
			expectedErr:        domain.ErrSellerNotFound,
			expectedRepoCalled: true,
		},
		{
			name:               "random error",
			sellerID:           sellerUUID.String(),
			repositoryErr:      randomErr,
			expectedErr:        randomErr,
			expectedRepoCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeSellerRepository{
				err: tt.repositoryErr,
			}

			uc := New(repo)

			err := uc.ArchiveSeller(context.Background(), tt.sellerID)

			if repo.called != tt.expectedRepoCalled {
				t.Fatalf("unexpected sellerRepo call state: got %v, want %v", repo.called, tt.expectedRepoCalled)
			}

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("unexpected error: got %v, want %v", err, tt.expectedErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)

			}

			if repo.receivedSellerID != tt.expectedUUID {
				t.Fatalf("unexpected seller id: got %v, want %v", repo.receivedSellerID, tt.expectedUUID)
			}
		})
	}
}

func TestSellerUseCase_DeleteSeller(t *testing.T) {
	sellerUUID := uuid.New()
	randomErr := errors.New("database error")

	tests := []struct {
		name               string
		sellerID           string
		repositoryErr      error
		expectedSellerID   uuid.UUID
		expectedErr        error
		expectedRepoCalled bool
	}{
		{
			name:               "success",
			sellerID:           sellerUUID.String(),
			expectedSellerID:   sellerUUID,
			expectedRepoCalled: true,
		},
		{
			name:               "empty seller id",
			sellerID:           "",
			expectedErr:        domain.ErrSellerIDRequired,
			expectedRepoCalled: false,
		},
		{
			name:               "invalid seller id",
			sellerID:           "invalid seller id",
			expectedErr:        domain.ErrInvalidSellerID,
			expectedRepoCalled: false,
		},
		{
			name:               "seller not found",
			sellerID:           sellerUUID.String(),
			repositoryErr:      domain.ErrSellerNotFound,
			expectedErr:        domain.ErrSellerNotFound,
			expectedRepoCalled: true,
		},
		{
			name:               "random error",
			sellerID:           sellerUUID.String(),
			repositoryErr:      randomErr,
			expectedErr:        randomErr,
			expectedRepoCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeSellerRepository{
				err: tt.repositoryErr,
			}

			uc := New(repo)

			err := uc.DeleteSeller(context.Background(), tt.sellerID)

			if repo.called != tt.expectedRepoCalled {
				t.Fatalf("unexpected sellerRepo call state: got %v, want %v", repo.called, tt.expectedRepoCalled)
			}

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("unexpected error: got %v, want %v", err, tt.expectedErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if repo.receivedSellerID != tt.expectedSellerID {
				t.Fatalf("unexpected seller id: got %v, want %v", repo.receivedSellerID, tt.expectedSellerID)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}

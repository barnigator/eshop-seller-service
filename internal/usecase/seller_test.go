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
	seller domain.Seller
	err    error

	receivedSeller domain.Seller
	createResult   domain.Seller
	createCalled   bool
	createErr      error
}

func (f *fakeSellerRepository) GetSellerByID(_ context.Context, _ uuid.UUID) (domain.Seller, error) {
	return f.seller, f.err
}

func (f *fakeSellerRepository) CreateSeller(_ context.Context, seller domain.Seller) (domain.Seller, error) {
	f.receivedSeller = seller
	f.createCalled = true
	return f.createResult, f.createErr
}

func TestSellerUseCase_GetSeller(t *testing.T) {
	tests := []struct {
		name           string
		sellerID       string
		repositoryData domain.Seller
		repositoryErr  error
		expectedSeller domain.Seller
		expectedErr    error
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
			expectedErr: nil,
		},
		{
			name:           "empty seller_id",
			sellerID:       "",
			repositoryData: domain.Seller{},
			repositoryErr:  nil,
			expectedSeller: domain.Seller{},
			expectedErr:    domain.ErrSellerIDRequired,
		},
		{
			name:           "invalid seller_id",
			sellerID:       "invalid uuid",
			repositoryData: domain.Seller{},
			repositoryErr:  nil,
			expectedSeller: domain.Seller{},
			expectedErr:    domain.ErrInvalidSellerID,
		},
		{
			name:           "repository error",
			sellerID:       "550e8400-e29b-41d4-a716-446655440000",
			repositoryData: domain.Seller{},
			repositoryErr:  domain.ErrSellerNotFound,
			expectedSeller: domain.Seller{},
			expectedErr:    domain.ErrSellerNotFound,
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
		expectedCreateCalled    bool
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
			expectedCreateCalled: true,
		},
		{
			name:                 "empty user_id",
			userID:               "",
			expectedErr:          domain.ErrUserIDRequired,
			expectedCreateCalled: false,
		},
		{
			name:                 "invalid user_id",
			userID:               "invalid user_id",
			expectedErr:          domain.ErrInvalidUserID,
			expectedCreateCalled: false,
		},
		{
			name:                 "empty brand_name",
			userID:               "550e8400-e29b-41d4-a716-446655440000",
			brandName:            "",
			expectedErr:          domain.ErrBrandNameRequired,
			expectedCreateCalled: false,
		},
		{
			name:                 "space brand_name",
			userID:               "550e8400-e29b-41d4-a716-446655440000",
			brandName:            "       	  ",
			expectedErr:          domain.ErrBrandNameRequired,
			expectedCreateCalled: false,
		},
		{
			name:                 "too long brand_name",
			userID:               "550e8400-e29b-41d4-a716-446655440000",
			brandName:            strings.Repeat("a", 121),
			expectedErr:          domain.ErrBrandNameTooLong,
			expectedCreateCalled: false,
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
			expectedCreateCalled: true,
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
			expectedCreateCalled: true,
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
			expectedCreateCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeSellerRepository{
				createResult: tt.createResult,
				createErr:    tt.createErr,
			}

			uc := New(repo)

			seller, err := uc.CreateSeller(context.Background(), tt.userID, tt.brandName, tt.description)

			if repo.createCalled != tt.expectedCreateCalled {
				t.Fatalf("unexpected repository call state: got %v, want %v", repo.createCalled, tt.expectedCreateCalled)
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

			if tt.expectedCreateCalled && repo.receivedSeller != tt.expectedRepositoryInput {
				t.Fatalf("unexpected seller RepositoryInput: got %v, want %v", repo.receivedSeller, tt.expectedRepositoryInput)
			}
		})
	}
}

func TestSellerUseCase_GetSellerStatus(t *testing.T) {
	tests := []struct {

	}
}

package usecase

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	"github.com/google/uuid"
)

type SellerRepository interface {
	GetSellerByID(ctx context.Context, sellerID uuid.UUID) (domain.Seller, error)
	CreateSeller(ctx context.Context, seller domain.Seller) (domain.Seller, error)
}

type SellerUseCase struct {
	repo SellerRepository
}

func New(repo SellerRepository) *SellerUseCase {
	return &SellerUseCase{repo: repo}
}

func (uc *SellerUseCase) GetSellerStatus(ctx context.Context, sellerID string) (domain.SellerStatus, error) {
	if sellerID == "" {
		return domain.SellerStatusUnspecified, domain.ErrSellerIDRequired
	}

	sellerUUID, err := uuid.Parse(sellerID)
	if err != nil {
		return domain.SellerStatusUnspecified, domain.ErrInvalidSellerID
	}

	seller, err := uc.repo.GetSellerByID(ctx, sellerUUID)
	if err != nil {
		return domain.SellerStatusUnspecified, err
	}

	return seller.Status, nil
}

func (uc *SellerUseCase) CreateSeller(ctx context.Context, userID string, brandName string, description string) (domain.Seller, error) {
	if userID == "" {
		return domain.Seller{}, domain.ErrUserIDRequired
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return domain.Seller{}, domain.ErrInvalidUserID
	}

	cleanBrandName := strings.TrimSpace(brandName)

	if cleanBrandName == "" {
		return domain.Seller{}, domain.ErrBrandNameRequired
	}

	if utf8.RuneCountInString(cleanBrandName) > 120 {
		return domain.Seller{}, domain.ErrBrandNameTooLong
	}

	cleanDescription := strings.TrimSpace(description)

	seller := domain.Seller{
		UserID:      userUUID,
		BrandName:   cleanBrandName,
		Description: cleanDescription,
		Status:      domain.SellerStatusPending,
	}

	createdSeller, err := uc.repo.CreateSeller(ctx, seller)
	if err != nil {
		return domain.Seller{}, err
	}

	return createdSeller, nil
}

func (uc *SellerUseCase) GetSeller(ctx context.Context, sellerID string) (domain.Seller, error) {
	if sellerID == "" {
		return domain.Seller{}, domain.ErrSellerIDRequired
	}

	sellerUUID, err := uuid.Parse(sellerID)
	if err != nil {
		return domain.Seller{}, domain.ErrInvalidSellerID
	}

	seller, err := uc.repo.GetSellerByID(ctx, sellerUUID)
	if err != nil {
		return domain.Seller{}, err
	}

	return seller, nil
}

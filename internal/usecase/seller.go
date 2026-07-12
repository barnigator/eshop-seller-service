package usecase

import (
	"context"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	"github.com/google/uuid"
)

type SellerRepository interface {
	GetSellerByID(ctx context.Context, sellerID uuid.UUID) (domain.Seller, error)
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
	id, err := uuid.Parse(sellerID)
	if err != nil {
		return domain.SellerStatusUnspecified, domain.ErrInvalidSellerID
	}
	seller, err := uc.repo.GetSellerByID(ctx, id)
	if err != nil {
		return domain.SellerStatusUnspecified, err
	}

	return seller.Status, nil
}

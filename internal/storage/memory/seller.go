package memory

import (
	"context"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	"github.com/google/uuid"
)

type SellerRepository struct {
	sellers map[uuid.UUID]domain.Seller
}

func New() *SellerRepository {
	sellers := make(map[uuid.UUID]domain.Seller)
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	sellers[id] = domain.Seller{
		ID:     id,
		Status: domain.SellerStatusActive,
	}

	return &SellerRepository{sellers: sellers}
}

func (s *SellerRepository) GetSellerByID(_ context.Context, sellerID uuid.UUID) (domain.Seller, error) {
	if seller, exists := s.sellers[sellerID]; exists {
		return seller, nil
	}

	return domain.Seller{}, domain.ErrSellerNotFound
}

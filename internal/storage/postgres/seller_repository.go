package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const getSellerByID = `
		SELECT
			id, 
			user_id, 
			brand_name, 
			description, 
			status 
		FROM sellers 
		WHERE id = $1  
		  AND deleted_at IS NULL;
`

type SellerRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *SellerRepository {
	return &SellerRepository{pool: pool}
}

func (s *SellerRepository) GetSellerByID(ctx context.Context, sellerID uuid.UUID) (domain.Seller, error) {
	var seller domain.Seller
	var status string

	err := s.pool.QueryRow(
		ctx,
		getSellerByID,
		sellerID,
	).Scan(
		&seller.ID,
		&seller.UserID,
		&seller.BrandName,
		&seller.Description,
		&status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Seller{}, domain.ErrSellerNotFound
		}
		return domain.Seller{}, fmt.Errorf("get seller by id: %w", err)
	}
	seller.Status, err = convertSellerStatus(status)
	if err != nil {
		return domain.Seller{}, fmt.Errorf("convert seller status: %w", err)
	}

	return seller, nil
}

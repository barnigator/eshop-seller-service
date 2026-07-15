package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	getSellerByIDQuery = `
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
	createSellerQuery = `
		INSERT INTO sellers (
			user_id,
			brand_name,
			description,
			status
		)
		VALUES ($1, $2, $3, $4)
		RETURNING 
			id, 
			user_id, 
			brand_name, 
			description, 
			status;
`
	listSellersByUserIDQuery = `
		SELECT 
			id,
			user_id,
			brand_name,
			description, 
			status
		FROM sellers
		WHERE user_id = $1
			AND deleted_at IS NULL
		ORDER BY created_at, id;
`
)

type SellerRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *SellerRepository {
	return &SellerRepository{pool: pool}
}

func (r *SellerRepository) GetSellerByID(ctx context.Context, sellerID uuid.UUID) (domain.Seller, error) {
	var seller domain.Seller
	var status string

	err := r.pool.QueryRow(
		ctx,
		getSellerByIDQuery,
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

	seller.Status, err = convertStringToSellerStatus(status)
	if err != nil {
		return domain.Seller{}, fmt.Errorf("convert string to seller status: %w", err)
	}

	return seller, nil
}

func (r *SellerRepository) CreateSeller(ctx context.Context, seller domain.Seller) (domain.Seller, error) {
	var createdSeller domain.Seller
	var status string

	status, err := convertSellerStatusToString(seller.Status)
	if err != nil {
		return domain.Seller{}, fmt.Errorf("convert seller status to string: %w", err)
	}

	err = r.pool.QueryRow(
		ctx,
		createSellerQuery,
		seller.UserID,
		seller.BrandName,
		seller.Description,
		status,
	).Scan(
		&createdSeller.ID,
		&createdSeller.UserID,
		&createdSeller.BrandName,
		&createdSeller.Description,
		&status,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "ux_sellers_user_brand_name_active" {
			return domain.Seller{}, domain.ErrBrandAlreadyExists
		}

		return domain.Seller{}, fmt.Errorf("create seller: %w", err)
	}

	createdSeller.Status, err = convertStringToSellerStatus(status)
	if err != nil {
		return domain.Seller{}, fmt.Errorf("convert string to seller status: %w", err)
	}

	return createdSeller, nil
}

func (r *SellerRepository) ListSellersByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Seller, error) {
	rows, err := r.pool.Query(
		ctx,
		listSellersByUserIDQuery,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list sellers by user id: %w", err)
	}
	defer rows.Close()

	sellers := make([]domain.Seller, 0)

	for rows.Next() {
		var seller domain.Seller
		var status string

		err = rows.Scan(
			&seller.ID,
			&seller.UserID,
			&seller.BrandName,
			&seller.Description,
			&status,
		)
		if err != nil {
			return nil, fmt.Errorf("scan seller row: %w", err)
		}

		seller.Status, err = convertStringToSellerStatus(status)
		if err != nil {
			return nil, fmt.Errorf("convert string to seller status: %w", err)
		}

		sellers = append(sellers, seller)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("iterate seller rows: %w", err)
	}

	return sellers, nil
}

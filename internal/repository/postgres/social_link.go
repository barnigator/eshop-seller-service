package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	"github.com/google/uuid"
)

const (
	addSocialLinkQuery = `
		INSERT INTO social_links (
			seller_id,
			type,
			url
		)
		SELECT
			$1,
			$2,
			$3
		FROM sellers s
		WHERE s.id = $1
		  AND s.deleted_at IS NULL
		RETURNING
			id,
			seller_id,
			type,
			url;
`
	listSocialLinksQuery = `
		SELECT
		    id,
		    seller_id,
		    type,
		    url
FROM social_links
		WHERE seller_id = $1
		ORDER BY (created_at, id);
`
)

func (r *Repository) AddSocialLink(ctx context.Context, link domain.SocialLink) (domain.SocialLink, error) {
	var addedLink domain.SocialLink

	err := r.pool.QueryRow(
		ctx,
		addSocialLinkQuery,
		link.SellerID,
		link.Type,
		link.URL,
	).Scan(
		&addedLink.ID,
		&addedLink.SellerID,
		&addedLink.Type,
		&addedLink.URL,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.SocialLink{}, domain.ErrSellerNotFound
		}

		return domain.SocialLink{}, fmt.Errorf("add social link: %w", err)
	}

	return addedLink, nil
}

func (r *Repository) ListSocialLinks(ctx context.Context, sellerID uuid.UUID) ([]domain.SocialLink, error) {
	rows, err := r.pool.Query(
		ctx,
		listSocialLinksQuery,
		sellerID,
	)
	if err != nil {
		return nil, fmt.Errorf("list social links: %w", err)
	}
	defer rows.Close()

	links := make([]domain.SocialLink, 0)

	for rows.Next() {
		var link domain.SocialLink

		err = rows.Scan(
			&link.ID,
			&link.SellerID,
			&link.Type,
			&link.URL,
		)
		if err != nil {
			return nil, fmt.Errorf("scan link row: %w", err)
		}

		links = append(links, link)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("iterate link rows: %w", err)
	}

	return links, nil
}

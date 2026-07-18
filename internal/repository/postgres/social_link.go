package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/barnigator/eshop-seller-service/internal/domain"
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

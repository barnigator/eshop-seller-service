package handler

import (
	"context"

	"github.com/barnigator/eshop-seller-service/internal/domain"
)

type SocialLinkUseCase interface {
	AddSocialLink(ctx context.Context, sellerID string, linkType domain.SocialLinkType, url string) (domain.SocialLink, error)
}

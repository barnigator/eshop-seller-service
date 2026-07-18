package usecase

import (
	"context"
	"strings"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	"github.com/google/uuid"
)

type SocialLinkRepository interface {
	AddSocialLink(ctx context.Context, link domain.SocialLink) (domain.SocialLink, error)
}

func (uc *UseCase) AddSocialLink(ctx context.Context, sellerID string, linkType domain.SocialLinkType, url string) (domain.SocialLink, error) {
	if sellerID == "" {
		return domain.SocialLink{}, domain.ErrSellerIDRequired
	}

	sellerUUID, err := uuid.Parse(sellerID)
	if err != nil {
		return domain.SocialLink{}, domain.ErrInvalidSellerID
	}

	if !linkType.IsValid() {
		return domain.SocialLink{}, domain.ErrInvalidSocialLinkType
	}

	cleanURL := strings.TrimSpace(url)

	if cleanURL == "" {
		return domain.SocialLink{}, domain.ErrURLRequired
	}

	link := domain.SocialLink{
		SellerID: sellerUUID,
		Type:     linkType,
		URL:      cleanURL,
	}

	return uc.socialLinkRepo.AddSocialLink(ctx, link)
}

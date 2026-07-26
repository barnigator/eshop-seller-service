package usecase

import (
	"context"
	"strings"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	"github.com/google/uuid"
)

type SocialLinkRepository interface {
	AddSocialLink(ctx context.Context, link domain.SocialLink) (domain.SocialLink, error)
	ListSocialLinks(ctx context.Context, sellerID uuid.UUID) ([]domain.SocialLink, error)
	DeleteSocialLink(ctx context.Context, linkID uuid.UUID) error
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

func (uc *UseCase) ListSocialLinks(ctx context.Context, sellerID string) ([]domain.SocialLink, error) {
	if sellerID == "" {
		return nil, domain.ErrSellerIDRequired
	}

	sellerUUID, err := uuid.Parse(sellerID)
	if err != nil {
		return nil, domain.ErrInvalidSellerID
	}

	return uc.socialLinkRepo.ListSocialLinks(ctx, sellerUUID)
}

func (uc *UseCase) DeleteSocialLink(ctx context.Context, linkID string) error {
	if linkID == "" {
		return domain.ErrSocialLinkIDRequired
	}

	linkUUID, err := uuid.Parse(linkID)
	if err != nil {
		return domain.ErrInvalidSocialLinkID
	}

	return uc.socialLinkRepo.DeleteSocialLink(ctx, linkUUID)
}

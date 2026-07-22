package handler

import (
	"context"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	sellerv1 "github.com/barnigator/protos/gen/go/seller/v1"
)

type SocialLinkUseCase interface {
	AddSocialLink(ctx context.Context, sellerID string, linkType domain.SocialLinkType, url string) (domain.SocialLink, error)
	ListSocialLinks(ctx context.Context, sellerID string) ([]domain.SocialLink, error)
}

func (h *Handler) AddSocialLink(ctx context.Context, req *sellerv1.AddSocialLinkRequest) (*sellerv1.SocialLinkResponse, error) {
	linkType := convertLinkTypeToDomainLinkType(req.Type)

	link, err := h.socialLinkUC.AddSocialLink(ctx, req.SellerId, linkType, req.Url)
	if err != nil {
		return nil, convertError(err)
	}

	return &sellerv1.SocialLinkResponse{
		Link: convertLink(link),
	}, nil
}

func (h *Handler) ListSocialLinks(ctx context.Context, req *sellerv1.ListSocialLinksRequest) (*sellerv1.ListSocialLinksResponse, error) {
	links, err := h.socialLinkUC.ListSocialLinks(ctx, req.SellerId)
	if err != nil {
		return nil, convertError(err)
	}

	return &sellerv1.ListSocialLinksResponse{
		Link: convertLinks(links),
	}, nil
}

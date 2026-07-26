package handler

import (
	"context"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	sellerv1 "github.com/barnigator/protos/gen/go/seller/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SocialLinkUseCase interface {
	AddSocialLink(ctx context.Context, sellerID string, linkType domain.SocialLinkType, url string) (domain.SocialLink, error)
	ListSocialLinks(ctx context.Context, sellerID string) ([]domain.SocialLink, error)
	DeleteSocialLink(ctx context.Context, linkID string) error
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
		Links: convertLinks(links),
	}, nil
}

func (h *Handler) DeleteSocialLink(ctx context.Context, req *sellerv1.DeleteSocialLinkRequest) (*emptypb.Empty, error) {
	err := h.socialLinkUC.DeleteSocialLink(ctx, req.LinkId)
	if err != nil {
		return nil, convertError(err)
	}

	return &emptypb.Empty{}, nil
}

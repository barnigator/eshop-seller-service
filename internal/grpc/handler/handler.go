package handler

import (
	"context"
	"fmt"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	sellerv1 "github.com/barnigator/protos/gen/go/seller/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SellerUseCase interface {
	GetSellerStatus(ctx context.Context, sellerID string) (domain.SellerStatus, error)
	CreateSeller(ctx context.Context, userID string, brandName string, description string) (domain.Seller, error)
	GetSeller(ctx context.Context, sellerID string) (domain.Seller, error)
	ListSellersByUserID(ctx context.Context, userID string) ([]domain.Seller, error)
	UpdateSeller(ctx context.Context, sellerID string, brandName *string, description *string) (domain.Seller, error)
	ArchiveSeller(ctx context.Context, sellerID string) error
	DeleteSeller(ctx context.Context, sellerID string) error
}

type Handler struct {
	uc SellerUseCase
	sellerv1.UnimplementedSellerServiceServer
}

func New(uc SellerUseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) GetSellerStatus(ctx context.Context, req *sellerv1.GetSellerStatusRequest) (*sellerv1.GetSellerStatusResponse, error) {
	st, err := h.uc.GetSellerStatus(ctx, req.SellerId)
	if err != nil {
		return nil, convertError(err)
	}
	return &sellerv1.GetSellerStatusResponse{
		Status: convertSellerStatus(st),
	}, nil
}

func (h *Handler) CreateSeller(ctx context.Context, req *sellerv1.CreateSellerRequest) (*sellerv1.SellerResponse, error) {
	seller, err := h.uc.CreateSeller(ctx, req.UserId, req.BrandName, req.Description)
	if err != nil {
		return nil, convertError(err)
	}
	return &sellerv1.SellerResponse{
		Seller: convertSeller(seller),
	}, nil
}

func (h *Handler) GetSeller(ctx context.Context, req *sellerv1.GetSellerRequest) (*sellerv1.SellerResponse, error) {
	seller, err := h.uc.GetSeller(ctx, req.SellerId)
	if err != nil {
		return nil, convertError(err)
	}

	return &sellerv1.SellerResponse{
		Seller: convertSeller(seller),
	}, nil
}

func (h *Handler) ListSellersByUserID(ctx context.Context, req *sellerv1.ListSellersByUserIDRequest) (*sellerv1.ListSellersResponse, error) {
	sellers, err := h.uc.ListSellersByUserID(ctx, req.UserId)
	if err != nil {
		return nil, convertError(err)
	}

	return &sellerv1.ListSellersResponse{
		Sellers: convertSellers(sellers),
	}, nil
}

func (h *Handler) UpdateSeller(ctx context.Context, req *sellerv1.UpdateSellerRequest) (*sellerv1.SellerResponse, error) {
	if req.UpdateMask == nil {
		return nil, invalidArgument("update_mask is required")
	}

	if len(req.UpdateMask.Paths) == 0 {
		return nil, invalidArgument("update_mask.paths must not be empty")
	}

	var brandName, description *string

	for _, path := range req.UpdateMask.Paths {
		switch path {
		case "brand_name":
			brandName = &req.BrandName
		case "description":
			description = &req.Description
		default:
			return nil, invalidArgument(
				fmt.Sprintf("unsupported update field: %s", path),
			)
		}
	}

	seller, err := h.uc.UpdateSeller(ctx, req.SellerId, brandName, description)
	if err != nil {
		return nil, convertError(err)
	}

	return &sellerv1.SellerResponse{
		Seller: convertSeller(seller),
	}, nil
}

func (h *Handler) ArchiveSeller(ctx context.Context, req *sellerv1.ArchiveSellerRequest) (*emptypb.Empty, error) {
	err := h.uc.ArchiveSeller(ctx, req.SellerId)
	if err != nil {
		return nil, convertError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) DeleteSeller(ctx context.Context, req *sellerv1.DeleteSellerRequest) (*emptypb.Empty, error) {
	err := h.uc.DeleteSeller(ctx, req.SellerId)
	if err != nil {
		return nil, convertError(err)
	}

	return &emptypb.Empty{}, nil
}

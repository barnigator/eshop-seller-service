package handler

import (
	"context"
	"errors"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	sellerv1 "github.com/barnigator/protos/gen/go/seller/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SellerUseCase interface {
	GetSellerStatus(ctx context.Context, sellerID string) (domain.SellerStatus, error)
	CreateSeller(ctx context.Context, userID string, brandName string, description string) (domain.Seller, error)
	GetSeller(ctx context.Context, sellerID string) (domain.Seller, error)
	ListSellersByUserID(ctx context.Context, userID string) ([]domain.Seller, error)
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
		switch {
		case errors.Is(err, domain.ErrSellerIDRequired):
			return nil, status.Error(codes.InvalidArgument, domain.ErrSellerIDRequired.Error())
		case errors.Is(err, domain.ErrInvalidSellerID):
			return nil, status.Error(codes.InvalidArgument, domain.ErrInvalidSellerID.Error())
		case errors.Is(err, domain.ErrSellerNotFound):
			return nil, status.Error(codes.NotFound, domain.ErrSellerNotFound.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &sellerv1.GetSellerStatusResponse{
		Status: convertSellerStatus(st),
	}, nil
}

func (h *Handler) CreateSeller(ctx context.Context, req *sellerv1.CreateSellerRequest) (*sellerv1.SellerResponse, error) {
	seller, err := h.uc.CreateSeller(ctx, req.UserId, req.BrandName, req.Description)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserIDRequired):
			return nil, status.Error(codes.InvalidArgument, domain.ErrUserIDRequired.Error())
		case errors.Is(err, domain.ErrInvalidUserID):
			return nil, status.Error(codes.InvalidArgument, domain.ErrInvalidUserID.Error())
		case errors.Is(err, domain.ErrBrandNameRequired):
			return nil, status.Error(codes.InvalidArgument, domain.ErrBrandNameRequired.Error())
		case errors.Is(err, domain.ErrBrandNameTooLong):
			return nil, status.Error(codes.InvalidArgument, domain.ErrBrandNameTooLong.Error())
		case errors.Is(err, domain.ErrBrandAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, domain.ErrBrandAlreadyExists.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &sellerv1.SellerResponse{
		Seller: convertSeller(seller),
	}, nil
}

func (h *Handler) GetSeller(ctx context.Context, req *sellerv1.GetSellerRequest) (*sellerv1.SellerResponse, error) {
	seller, err := h.uc.GetSeller(ctx, req.SellerId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSellerIDRequired):
			return nil, status.Error(codes.InvalidArgument, domain.ErrSellerIDRequired.Error())
		case errors.Is(err, domain.ErrInvalidSellerID):
			return nil, status.Error(codes.InvalidArgument, domain.ErrInvalidSellerID.Error())
		case errors.Is(err, domain.ErrSellerNotFound):
			return nil, status.Error(codes.NotFound, domain.ErrSellerNotFound.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &sellerv1.SellerResponse{
		Seller: convertSeller(seller),
	}, nil
}

func (h *Handler) ListSellersByUserID(ctx context.Context, req *sellerv1.ListSellersByUserIDRequest) (*sellerv1.ListSellersResponse, error) {
	sellers, err := h.uc.ListSellersByUserID(ctx, req.UserId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserIDRequired):
			return nil, status.Error(codes.InvalidArgument, domain.ErrUserIDRequired.Error())
		case errors.Is(err, domain.ErrInvalidUserID):
			return nil, status.Error(codes.InvalidArgument, domain.ErrInvalidUserID.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &sellerv1.ListSellersResponse{
		Sellers: convertSellers(sellers),
	}, nil
}

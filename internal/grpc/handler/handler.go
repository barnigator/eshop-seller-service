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

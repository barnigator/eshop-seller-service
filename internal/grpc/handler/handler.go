package handler

import (
	"context"

	sellerv1 "github.com/barnigator/protos/gen/go/seller/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	sellerv1.UnimplementedSellerServiceServer
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) GetSellerStatus(ctx context.Context, req *sellerv1.GetSellerStatusRequest) (*sellerv1.GetSellerStatusResponse, error) {
	if req.SellerId == "" {
		return nil, status.Error(codes.InvalidArgument, "seller_id is required")
	}

	return &sellerv1.GetSellerStatusResponse{
		Status: sellerv1.SellerStatus_SELLER_STATUS_ACTIVE,
	}, nil
}

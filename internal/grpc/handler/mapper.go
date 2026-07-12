package handler

import (
	"github.com/barnigator/eshop-seller-service/internal/domain"
	sellerv1 "github.com/barnigator/protos/gen/go/seller/v1"
)

func convertSellerStatus(status domain.SellerStatus) sellerv1.SellerStatus {
	switch status {
	case domain.SellerStatusActive:
		return sellerv1.SellerStatus_SELLER_STATUS_ACTIVE
	case domain.SellerStatusArchived:
		return sellerv1.SellerStatus_SELLER_STATUS_ARCHIVED
	case domain.SellerStatusBlocked:
		return sellerv1.SellerStatus_SELLER_STATUS_BLOCKED
	case domain.SellerStatusPending:
		return sellerv1.SellerStatus_SELLER_STATUS_PENDING
	case domain.SellerStatusRejected:
		return sellerv1.SellerStatus_SELLER_STATUS_REJECTED
	case domain.SellerStatusUnspecified:
		return sellerv1.SellerStatus_SELLER_STATUS_UNSPECIFIED
	}

	return sellerv1.SellerStatus_SELLER_STATUS_UNSPECIFIED
}

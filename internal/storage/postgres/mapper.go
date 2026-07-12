package postgres

import (
	"fmt"

	"github.com/barnigator/eshop-seller-service/internal/domain"
)

func convertSellerStatus(status string) (domain.SellerStatus, error) {
	switch status {
	case "active":
		return domain.SellerStatusActive, nil
	case "archived":
		return domain.SellerStatusArchived, nil
	case "blocked":
		return domain.SellerStatusBlocked, nil
	case "pending":
		return domain.SellerStatusPending, nil
	case "rejected":
		return domain.SellerStatusRejected, nil
	default:
		return domain.SellerStatusUnspecified, fmt.Errorf("%w: %s", ErrUnknownSellerStatus, status)
	}

}

package postgres

import (
	"fmt"

	"github.com/barnigator/eshop-seller-service/internal/domain"
)

func convertStringToSellerStatus(status string) (domain.SellerStatus, error) {
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
		return domain.SellerStatusUnspecified, fmt.Errorf("%w: %s", ErrUnknownDatabaseSellerStatus, status)
	}
}

func convertSellerStatusToString(status domain.SellerStatus) (string, error) {
	switch status {
	case domain.SellerStatusActive:
		return "active", nil
	case domain.SellerStatusArchived:
		return "archived", nil
	case domain.SellerStatusBlocked:
		return "blocked", nil
	case domain.SellerStatusPending:
		return "pending", nil
	case domain.SellerStatusRejected:
		return "rejected", nil
	default:
		return "", fmt.Errorf("%w: %v", ErrUnknownDomainSellerStatus, status)
	}
}

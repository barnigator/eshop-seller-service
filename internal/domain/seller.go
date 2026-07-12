package domain

import "github.com/google/uuid"

type SellerStatus int

const (
	SellerStatusUnspecified SellerStatus = iota
	SellerStatusPending
	SellerStatusActive
	SellerStatusRejected
	SellerStatusBlocked
	SellerStatusArchived
)

type Seller struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	BrandName   string
	Description string
	Status      SellerStatus
}

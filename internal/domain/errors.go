package domain

import "errors"

var (
	ErrSellerIDRequired = errors.New("seller_id is required")
	ErrSellerNotFound   = errors.New("seller not found")
	ErrInvalidSellerID  = errors.New("invalid seller_id")

	ErrUserIDRequired = errors.New("user_id is required")
	ErrInvalidUserID  = errors.New("invalid user_id")

	ErrBrandNameRequired  = errors.New("brand_name is required")
	ErrBrandNameTooLong   = errors.New("brand_name exceeds 120 characters")
	ErrBrandAlreadyExists = errors.New("brand already exists")
)

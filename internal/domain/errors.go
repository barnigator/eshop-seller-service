package domain

import "errors"

var (
	ErrSellerIDRequired = errors.New("seller_id is required")
	ErrSellerNotFound   = errors.New("seller does not exist")
	ErrInvalidSellerID  = errors.New("invalid seller id")
)

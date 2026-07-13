package postgres

import "errors"

var (
	ErrUnknownDatabaseSellerStatus = errors.New("invalid seller status from database")
	ErrUnknownDomainSellerStatus   = errors.New("invalid seller status from domain")
)

package postgres

import "errors"

var ErrUnknownSellerStatus = errors.New("invalid seller status from database")

package handler

import (
	"errors"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrUnknownLinkType = errors.New("unknown link type")

func invalidArgument(message string) error {
	return status.Error(codes.InvalidArgument, message)
}

func convertError(err error) error {
	switch {
	case errors.Is(err, domain.ErrSellerIDRequired):
		return status.Error(codes.InvalidArgument, domain.ErrSellerIDRequired.Error())
	case errors.Is(err, domain.ErrInvalidSellerID):
		return status.Error(codes.InvalidArgument, domain.ErrInvalidSellerID.Error())
	case errors.Is(err, domain.ErrUserIDRequired):
		return status.Error(codes.InvalidArgument, domain.ErrUserIDRequired.Error())
	case errors.Is(err, domain.ErrInvalidUserID):
		return status.Error(codes.InvalidArgument, domain.ErrInvalidUserID.Error())
	case errors.Is(err, domain.ErrBrandNameRequired):
		return status.Error(codes.InvalidArgument, domain.ErrBrandNameRequired.Error())
	case errors.Is(err, domain.ErrBrandNameTooLong):
		return status.Error(codes.InvalidArgument, domain.ErrBrandNameTooLong.Error())
	case errors.Is(err, domain.ErrBrandAlreadyExists):
		return status.Error(codes.AlreadyExists, domain.ErrBrandAlreadyExists.Error())
	case errors.Is(err, domain.ErrSellerNotFound):
		return status.Error(codes.NotFound, domain.ErrSellerNotFound.Error())
	case errors.Is(err, domain.ErrNoFieldsToUpdate):
		return status.Error(codes.InvalidArgument, domain.ErrNoFieldsToUpdate.Error())
	case errors.Is(err, domain.ErrSocialLinkIDRequired):
		return status.Error(codes.InvalidArgument, domain.ErrSocialLinkIDRequired.Error())
	case errors.Is(err, domain.ErrInvalidSocialLinkID):
		return status.Error(codes.InvalidArgument, domain.ErrInvalidSocialLinkID.Error())
	case errors.Is(err, domain.ErrSocialLinkNotFound):
		return status.Error(codes.NotFound, domain.ErrSocialLinkNotFound.Error())
	case errors.Is(err, domain.ErrInvalidSocialLinkType):
		return status.Error(codes.InvalidArgument, domain.ErrInvalidSocialLinkType.Error())
	case errors.Is(err, domain.ErrURLRequired):
		return status.Error(codes.InvalidArgument, domain.ErrURLRequired.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}

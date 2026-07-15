package handler

import (
	"errors"
	"fmt"
	"testing"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestConvertError(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		expectedErrCode codes.Code
		expectedMessage string
	}{
		{
			name:            "user_id required",
			err:             domain.ErrUserIDRequired,
			expectedErrCode: codes.InvalidArgument,
			expectedMessage: domain.ErrUserIDRequired.Error(),
		},
		{
			name:            "invalid seller_id",
			err:             domain.ErrInvalidSellerID,
			expectedErrCode: codes.InvalidArgument,
			expectedMessage: domain.ErrInvalidSellerID.Error(),
		},
		{
			name:            "seller not found",
			err:             domain.ErrSellerNotFound,
			expectedErrCode: codes.NotFound,
			expectedMessage: domain.ErrSellerNotFound.Error(),
		},
		{
			name:            "brand already exists",
			err:             domain.ErrBrandAlreadyExists,
			expectedErrCode: codes.AlreadyExists,
			expectedMessage: domain.ErrBrandAlreadyExists.Error(),
		},
		{
			name:            "random error",
			err:             errors.New("database error"),
			expectedErrCode: codes.Internal,
			expectedMessage: "internal error",
		},
		{
			name:            "wrapped error",
			err:             fmt.Errorf("repository: %w", domain.ErrSellerNotFound),
			expectedErrCode: codes.NotFound,
			expectedMessage: domain.ErrSellerNotFound.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := convertError(tt.err)

			st := status.Convert(err)

			if st.Code() != tt.expectedErrCode {
				t.Fatalf("unexpected err: got %v, want %v", st.Code(), tt.expectedErrCode)
			}

			if st.Message() != tt.expectedMessage {
				t.Fatalf("unexpected error message: got %v, want %v", st.Message(), tt.expectedMessage)
			}
		})
	}
}

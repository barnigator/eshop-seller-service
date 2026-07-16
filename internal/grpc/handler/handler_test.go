package handler

import (
	"context"
	"testing"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	sellerv1 "github.com/barnigator/protos/gen/go/seller/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type fakeSellerUsecase struct {
	seller  domain.Seller
	sellers []domain.Seller
	called  bool
	err     error

	receivedBrandName   *string
	receivedDescription *string
}

func (f *fakeSellerUsecase) GetSellerStatus(ctx context.Context, sellerID string) (domain.SellerStatus, error) {
	return f.seller.Status, f.err
}

func (f *fakeSellerUsecase) CreateSeller(_ context.Context, _ string, _ string, _ string) (domain.Seller, error) {
	return f.seller, f.err
}

func (f *fakeSellerUsecase) GetSeller(_ context.Context, _ string) (domain.Seller, error) {
	return f.seller, f.err
}

func (f *fakeSellerUsecase) ListSellersByUserID(_ context.Context, _ string) ([]domain.Seller, error) {
	return f.sellers, f.err
}

func (f *fakeSellerUsecase) UpdateSeller(_ context.Context, sellerID string, brandName *string, description *string) (domain.Seller, error) {
	f.called = true
	f.receivedBrandName = brandName
	f.receivedDescription = description

	return f.seller, f.err
}

func TestHandler_UpdateSeller(t *testing.T) {
	tests := []struct {
		name                string
		brandName           string
		description         string
		updateMask          *fieldmaskpb.FieldMask
		expectedBrandName   *string
		expectedDescription *string
		expectedCode        codes.Code
		expectedMessage     string
		expectedUCCalled    bool
	}{
		{
			name:             "nil update mask",
			updateMask:       nil,
			expectedCode:     codes.InvalidArgument,
			expectedMessage:  "update_mask is required",
			expectedUCCalled: false,
		},
		{
			name: "empty update mask",
			updateMask: &fieldmaskpb.FieldMask{
				Paths: []string{},
			},
			expectedCode:     codes.InvalidArgument,
			expectedMessage:  "update_mask.paths must not be empty",
			expectedUCCalled: false,
		},
		{
			name:        "both paths",
			brandName:   "Adidas",
			description: "cool brand",
			updateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"brand_name", "description"},
			},
			expectedBrandName:   ptrStr("Adidas"),
			expectedDescription: ptrStr("cool brand"),
			expectedUCCalled:    true,
		},
		{
			name:      "brand_name",
			brandName: "Adidas",
			updateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"brand_name"},
			},
			expectedBrandName: ptrStr("Adidas"),
			expectedUCCalled:  true,
		},
		{
			name:        "description",
			description: "cool brand",
			updateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"description"},
			},
			expectedDescription: ptrStr("cool brand"),
			expectedUCCalled:    true,
		},
		{
			name: "invalid path",
			updateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"invalid path"},
			},
			expectedCode:     codes.InvalidArgument,
			expectedMessage:  "unsupported update field: invalid path",
			expectedUCCalled: false,
		},
		{
			name:        "clear description",
			description: "",
			updateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"description"},
			},
			expectedDescription: ptrStr(""),
			expectedUCCalled:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &fakeSellerUsecase{}

			h := New(uc)

			req := &sellerv1.UpdateSellerRequest{
				BrandName:   tt.brandName,
				Description: tt.description,
				UpdateMask:  tt.updateMask,
			}

			_, err := h.UpdateSeller(context.Background(), req)

			if uc.called != tt.expectedUCCalled {
				t.Fatalf("unexpected usecase call state: got %v, want %v", uc.called, tt.expectedUCCalled)
			}

			st := status.Convert(err)

			if st.Code() != tt.expectedCode {
				t.Fatalf("unexpected err: got %v, want %v", st.Code(), tt.expectedCode)
			}

			if st.Message() != tt.expectedMessage {
				t.Fatalf("unexpected error message: got %v, want %v", st.Message(), tt.expectedMessage)
			}

			if tt.expectedBrandName == nil {
				if uc.receivedBrandName != nil {
					t.Fatalf("unexpected brand name: got %v, want nil", uc.receivedBrandName)
				}
			} else {
				if uc.receivedBrandName == nil {
					t.Fatalf("unexpected brand name: got nil, want %v", tt.expectedBrandName)
				}

				if *uc.receivedBrandName != *tt.expectedBrandName {
					t.Fatalf("unexpected brand name: got %v, want %v", *uc.receivedBrandName, *tt.expectedBrandName)
				}
			}

			if tt.expectedDescription == nil {
				if uc.receivedDescription != nil {
					t.Fatalf("unexpected description: got %v, want nil", uc.receivedDescription)
				}
			} else {
				if uc.receivedDescription == nil {
					t.Fatalf("unexpected description: got nil, want %v", tt.expectedDescription)
				}

				if *uc.receivedDescription != *tt.expectedDescription {
					t.Fatalf("unexpected description: got %v, want %v", *uc.receivedDescription, *tt.expectedDescription)
				}
			}
		})
	}
}

func ptrStr(s string) *string {
	return &s
}

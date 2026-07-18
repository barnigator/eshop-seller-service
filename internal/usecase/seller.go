package usecase

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/barnigator/eshop-seller-service/internal/domain"
	"github.com/google/uuid"
)

type SellerRepository interface {
	GetSellerByID(ctx context.Context, sellerID uuid.UUID) (domain.Seller, error)
	CreateSeller(ctx context.Context, seller domain.Seller) (domain.Seller, error)
	ListSellersByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Seller, error)
	UpdateSeller(ctx context.Context, sellerID uuid.UUID, brandName *string, description *string) (domain.Seller, error)
	ArchiveSeller(ctx context.Context, sellerID uuid.UUID) error
	DeleteSeller(ctx context.Context, sellerID uuid.UUID) error
}

type UseCase struct {
	sellerRepo     SellerRepository
	socialLinkRepo SocialLinkRepository
}

func New(sellerRepo SellerRepository, socialLinkRepo SocialLinkRepository) *UseCase {
	return &UseCase{
		sellerRepo:     sellerRepo,
		socialLinkRepo: socialLinkRepo,
	}
}

func (uc *UseCase) GetSellerStatus(ctx context.Context, sellerID string) (domain.SellerStatus, error) {
	if sellerID == "" {
		return domain.SellerStatusUnspecified, domain.ErrSellerIDRequired
	}

	sellerUUID, err := uuid.Parse(sellerID)
	if err != nil {
		return domain.SellerStatusUnspecified, domain.ErrInvalidSellerID
	}

	seller, err := uc.sellerRepo.GetSellerByID(ctx, sellerUUID)
	if err != nil {
		return domain.SellerStatusUnspecified, err
	}

	return seller.Status, nil
}

func (uc *UseCase) CreateSeller(ctx context.Context, userID string, brandName string, description string) (domain.Seller, error) {
	if userID == "" {
		return domain.Seller{}, domain.ErrUserIDRequired
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return domain.Seller{}, domain.ErrInvalidUserID
	}

	cleanBrandName := strings.TrimSpace(brandName)

	if cleanBrandName == "" {
		return domain.Seller{}, domain.ErrBrandNameRequired
	}

	if utf8.RuneCountInString(cleanBrandName) > 120 {
		return domain.Seller{}, domain.ErrBrandNameTooLong
	}

	cleanDescription := strings.TrimSpace(description)

	seller := domain.Seller{
		UserID:      userUUID,
		BrandName:   cleanBrandName,
		Description: cleanDescription,
		Status:      domain.SellerStatusPending,
	}

	createdSeller, err := uc.sellerRepo.CreateSeller(ctx, seller)
	if err != nil {
		return domain.Seller{}, err
	}

	return createdSeller, nil
}

func (uc *UseCase) GetSeller(ctx context.Context, sellerID string) (domain.Seller, error) {
	if sellerID == "" {
		return domain.Seller{}, domain.ErrSellerIDRequired
	}

	sellerUUID, err := uuid.Parse(sellerID)
	if err != nil {
		return domain.Seller{}, domain.ErrInvalidSellerID
	}

	seller, err := uc.sellerRepo.GetSellerByID(ctx, sellerUUID)
	if err != nil {
		return domain.Seller{}, err
	}

	return seller, nil
}

func (uc *UseCase) ListSellersByUserID(ctx context.Context, userID string) ([]domain.Seller, error) {
	if userID == "" {
		return nil, domain.ErrUserIDRequired
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, domain.ErrInvalidUserID
	}

	sellers, err := uc.sellerRepo.ListSellersByUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	return sellers, nil
}

func (uc *UseCase) UpdateSeller(ctx context.Context, sellerID string, brandName *string, description *string) (domain.Seller, error) {
	if sellerID == "" {
		return domain.Seller{}, domain.ErrSellerIDRequired
	}

	sellerUUID, err := uuid.Parse(sellerID)
	if err != nil {
		return domain.Seller{}, domain.ErrInvalidSellerID
	}

	if brandName == nil && description == nil {
		return domain.Seller{}, domain.ErrNoFieldsToUpdate
	}

	if brandName != nil {
		cleanBrandName := strings.TrimSpace(*brandName)
		if cleanBrandName == "" {
			return domain.Seller{}, domain.ErrBrandNameRequired
		}

		if utf8.RuneCountInString(cleanBrandName) > 120 {
			return domain.Seller{}, domain.ErrBrandNameTooLong
		}

		brandName = &cleanBrandName
	}

	if description != nil {
		cleanDescription := strings.TrimSpace(*description)

		description = &cleanDescription
	}

	seller, err := uc.sellerRepo.UpdateSeller(ctx, sellerUUID, brandName, description)
	if err != nil {
		return domain.Seller{}, err
	}

	return seller, nil
}

func (uc *UseCase) ArchiveSeller(ctx context.Context, sellerID string) error {
	if sellerID == "" {
		return domain.ErrSellerIDRequired
	}

	sellerUUID, err := uuid.Parse(sellerID)
	if err != nil {
		return domain.ErrInvalidSellerID
	}

	return uc.sellerRepo.ArchiveSeller(ctx, sellerUUID)
}

func (uc *UseCase) DeleteSeller(ctx context.Context, sellerID string) error {
	if sellerID == "" {
		return domain.ErrSellerIDRequired
	}

	sellerUUID, err := uuid.Parse(sellerID)
	if err != nil {
		return domain.ErrInvalidSellerID
	}

	return uc.sellerRepo.DeleteSeller(ctx, sellerUUID)
}

package handler

import (
	"github.com/barnigator/eshop-seller-service/internal/domain"
	sellerv1 "github.com/barnigator/protos/gen/go/seller/v1"
)

func convertSellerStatus(status domain.SellerStatus) sellerv1.SellerStatus {
	switch status {
	case domain.SellerStatusActive:
		return sellerv1.SellerStatus_SELLER_STATUS_ACTIVE
	case domain.SellerStatusArchived:
		return sellerv1.SellerStatus_SELLER_STATUS_ARCHIVED
	case domain.SellerStatusBlocked:
		return sellerv1.SellerStatus_SELLER_STATUS_BLOCKED
	case domain.SellerStatusPending:
		return sellerv1.SellerStatus_SELLER_STATUS_PENDING
	case domain.SellerStatusRejected:
		return sellerv1.SellerStatus_SELLER_STATUS_REJECTED
	case domain.SellerStatusUnspecified:
		return sellerv1.SellerStatus_SELLER_STATUS_UNSPECIFIED
	}

	return sellerv1.SellerStatus_SELLER_STATUS_UNSPECIFIED
}

func convertSeller(seller domain.Seller) *sellerv1.Seller {
	return &sellerv1.Seller{
		Id:          seller.ID.String(),
		UserId:      seller.UserID.String(),
		BrandName:   seller.BrandName,
		Description: seller.Description,
		Status:      convertSellerStatus(seller.Status),
	}
}

func convertSellers(sellers []domain.Seller) []*sellerv1.Seller {
	sellersResult := make([]*sellerv1.Seller, len(sellers))
	for i, seller := range sellers {
		sellersResult[i] = convertSeller(seller)
	}

	return sellersResult
}

func convertLinkTypeToDomainLinkType(linkType sellerv1.SocialLinkType) domain.SocialLinkType {
	switch linkType {
	case sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_AVITO:
		return domain.SocialLinkTypeAvito
	case sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_FACEBOOK:
		return domain.SocialLinkTypeFacebook
	case sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_INSTAGRAM:
		return domain.SocialLinkTypeInstagram
	case sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_OZON:
		return domain.SocialLinkTypeOzon
	case sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_TELEGRAM:
		return domain.SocialLinkTypeTelegram
	case sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_VK:
		return domain.SocialLinkTypeVK
	case sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_WEBSITE:
		return domain.SocialLinkTypeWebsite
	case sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_WHATSAPP:
		return domain.SocialLinkTypeWhatsApp
	case sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_WILDBERRIES:
		return domain.SocialLinkTypeWildberries
	case sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_YANDEX_MARKET:
		return domain.SocialLinkTypeYandexMarket
	case sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_YOUTUBE:
		return domain.SocialLinkTypeYoutube
	default:
		return domain.SocialLinkTypeUnspecified
	}
}

func convertDomainLinkTypeToLinkType(linkType domain.SocialLinkType) sellerv1.SocialLinkType {
	switch linkType {
	case domain.SocialLinkTypeYoutube:
		return sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_YOUTUBE
	case domain.SocialLinkTypeOzon:
		return sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_OZON
	case domain.SocialLinkTypeVK:
		return sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_VK
	case domain.SocialLinkTypeYandexMarket:
		return sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_YANDEX_MARKET
	case domain.SocialLinkTypeWebsite:
		return sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_WEBSITE
	case domain.SocialLinkTypeInstagram:
		return sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_INSTAGRAM
	case domain.SocialLinkTypeWildberries:
		return sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_WILDBERRIES
	case domain.SocialLinkTypeWhatsApp:
		return sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_WHATSAPP
	case domain.SocialLinkTypeTelegram:
		return sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_TELEGRAM
	case domain.SocialLinkTypeFacebook:
		return sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_FACEBOOK
	case domain.SocialLinkTypeAvito:
		return sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_AVITO
	default:
		return sellerv1.SocialLinkType_SOCIAL_LINK_TYPE_UNSPECIFIED
	}
}

func convertLink(link domain.SocialLink) *sellerv1.SocialLink {
	return &sellerv1.SocialLink{
		Id:       link.ID.String(),
		SellerId: link.SellerID.String(),
		Type:     convertDomainLinkTypeToLinkType(link.Type),
		Url:      link.URL,
	}
}

func convertLinks(links []domain.SocialLink) []*sellerv1.SocialLink {
	var linksResult []*sellerv1.SocialLink

	for _, link := range links {
		linksResult = append(linksResult, convertLink(link))
	}

	return linksResult
}

package domain

import "github.com/google/uuid"

type SocialLinkType string

const (
	SocialLinkTypeUnspecified  SocialLinkType = ""
	SocialLinkTypeTelegram     SocialLinkType = "telegram"
	SocialLinkTypeVK           SocialLinkType = "vk"
	SocialLinkTypeInstagram    SocialLinkType = "instagram"
	SocialLinkTypeWebsite      SocialLinkType = "website"
	SocialLinkTypeWhatsApp     SocialLinkType = "whatsapp"
	SocialLinkTypeYoutube      SocialLinkType = "youtube"
	SocialLinkTypeFacebook     SocialLinkType = "facebook"
	SocialLinkTypeOzon         SocialLinkType = "ozon"
	SocialLinkTypeWildberries  SocialLinkType = "wildberries"
	SocialLinkTypeYandexMarket SocialLinkType = "yandexmarket"
	SocialLinkTypeAvito        SocialLinkType = "avito"
)

func (t SocialLinkType) IsValid() bool {
	switch t {
	case SocialLinkTypeAvito,
		SocialLinkTypeFacebook,
		SocialLinkTypeInstagram,
		SocialLinkTypeOzon,
		SocialLinkTypeVK,
		SocialLinkTypeTelegram,
		SocialLinkTypeWebsite,
		SocialLinkTypeWhatsApp,
		SocialLinkTypeYoutube,
		SocialLinkTypeWildberries,
		SocialLinkTypeYandexMarket:
		return true
	default:
		return false
	}
}

type SocialLink struct {
	ID       uuid.UUID
	SellerID uuid.UUID
	Type     SocialLinkType
	URL      string
}

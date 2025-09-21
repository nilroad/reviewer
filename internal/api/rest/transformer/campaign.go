package transformer

import (
	"club/internal/constant"
	"club/internal/domain"
	"time"

	"git.oceantim.com/backend/packages/golang/essential/cdnmaker"
)

type CampaignResponse struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type DealResponse struct {
	ID              uint64    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	DiscountedPrice uint64    `json:"discounted_price"`
	Product         *Product  `json:"product,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CampaignTransformer struct {
	cdn                *cdnmaker.CDNMaker
	productTransformer *ProductTransformer
}

func NewCampaignTransformer(cdnHandler *cdnmaker.Handler, productTransformer *ProductTransformer) *CampaignTransformer {
	return &CampaignTransformer{
		cdn:                cdnHandler.Use(constant.ClubCDN),
		productTransformer: productTransformer,
	}
}

func (t *CampaignTransformer) Transform(campaign *domain.Campaign) CampaignResponse {
	response := CampaignResponse{
		ID:          campaign.ID,
		Title:       campaign.Title,
		Subtitle:    campaign.Subtitle,
		Description: campaign.Description,
		ImageURL:    t.cdn.GenerateFullURL(campaign.ImageFilePath),
		CreatedAt:   campaign.CreatedAt,
	}

	return response
}

func (t *CampaignTransformer) TransformMany(campaigns []*domain.Campaign) []CampaignResponse {
	result := make([]CampaignResponse, 0, len(campaigns))
	for _, campaign := range campaigns {
		result = append(result, t.Transform(campaign))
	}

	return result
}

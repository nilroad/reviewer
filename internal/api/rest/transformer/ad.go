package transformer

import (
	"club/internal/constant"
	"club/internal/domain"
	"time"

	"git.oceantim.com/backend/packages/golang/essential/cdnmaker"
	"git.oceantim.com/backend/packages/golang/go-hashids/v2"
)

type AdResponse struct {
	ID                string                     `json:"id"`
	RedirectURL       string                     `json:"redirect_url"`
	AspectRatio       constant.AspectRatio       `json:"aspect_ratio"`
	AspectRatioFactor constant.AspectRatioFactor `json:"aspect_ratio_factor"`
	CreatedAt         time.Time                  `json:"created_at"`
	ImageURL          string                     `json:"image_url"`
	Sort              uint64                     `json:"sort"`
}

type AdTransformer struct {
	cdn *cdnmaker.CDNMaker
}

func NewAdTransformer(cdnHandler *cdnmaker.Handler) *AdTransformer {
	return &AdTransformer{
		cdn: cdnHandler.Use(constant.ClubCDN),
	}
}

func (t *AdTransformer) Transform(ad *domain.Ad) AdResponse {
	return AdResponse{
		//nolint:all
		ID:                hashids.Encode([]int{int(ad.ID)}),
		RedirectURL:       ad.RedirectURL,
		AspectRatio:       ad.AspectRatio,
		AspectRatioFactor: ad.AspectRatioFactor,
		CreatedAt:         ad.CreatedAt,
		ImageURL:          t.cdn.GenerateFullURL(ad.Media.FilePath),
		Sort:              ad.Sort,
	}
}

func (t *AdTransformer) TransformMany(ads []*domain.Ad) []AdResponse {
	transformedAds := make([]AdResponse, 0, len(ads))
	for _, ad := range ads {
		transformedAds = append(transformedAds, t.Transform(ad))
	}

	return transformedAds
}

type AdminAdResponse struct {
	ID                uint64                     `json:"id"`
	RedirectURL       string                     `json:"redirect_url"`
	AspectRatio       constant.AspectRatio       `json:"aspect_ratio"`
	AspectRatioFactor constant.AspectRatioFactor `json:"aspect_ratio_factor"`
	CreatedAt         time.Time                  `json:"created_at"`
	ImageURL          string                     `json:"image_url"`
	Sort              uint64                     `json:"sort"`
	IsActive          bool                       `json:"is_active"`
}

type AdminAdTransformer struct {
	cdn *cdnmaker.CDNMaker
}

func NewAdminAdTransformer(cdnHandler *cdnmaker.Handler) *AdminAdTransformer {
	return &AdminAdTransformer{
		cdn: cdnHandler.Use(constant.ClubCDN),
	}
}

func (t *AdminAdTransformer) Transform(ad *domain.Ad) AdminAdResponse {
	return AdminAdResponse{
		ID:                ad.ID,
		RedirectURL:       ad.RedirectURL,
		AspectRatio:       ad.AspectRatio,
		AspectRatioFactor: ad.AspectRatioFactor,
		CreatedAt:         ad.CreatedAt,
		ImageURL:          t.cdn.GenerateFullURL(ad.Media.FilePath),
		Sort:              ad.Sort,
		IsActive:          ad.IsActive,
	}
}

func (t *AdminAdTransformer) TransformMany(ads []*domain.Ad) []AdminAdResponse {
	transformedAds := make([]AdminAdResponse, 0, len(ads))
	for _, ad := range ads {
		transformedAds = append(transformedAds, t.Transform(ad))
	}

	return transformedAds
}

package transformer

import (
	"club/internal/constant"
	"club/internal/domain"

	"git.oceantim.com/backend/packages/golang/essential/cdnmaker"
)

// CategoryResponse represents a product category in the API response
type CategoryResponse struct {
	ID          uint64  `json:"id"`
	Title       string  `json:"title"`
	Image       *string `json:"image"`
	HasDiscount bool    `json:"has_discount"`
}

// CategoryTransformer handles transforming product categories to API responses
type CategoryTransformer struct {
	cdn *cdnmaker.CDNMaker
}

// NewCategoryTransformer creates a new product category transformer
func NewCategoryTransformer(cdnMakerHandler *cdnmaker.Handler) *CategoryTransformer {
	return &CategoryTransformer{
		cdn: cdnMakerHandler.Use(constant.ClubCDN),
	}
}

// Transform converts a domain product category with products to an API response
func (t *CategoryTransformer) Transform(category *domain.Category) *CategoryResponse {
	if category == nil {
		return nil
	}

	var imageURL *string
	if category.ImageFilePath != nil {
		url := t.cdn.GenerateFullURL(*category.ImageFilePath)
		imageURL = &url
	}

	return &CategoryResponse{
		ID:          category.ID,
		Title:       category.Title,
		Image:       imageURL,
		HasDiscount: category.HasDiscount,
	}
}

// TransformMany converts multiple domain product categories to API responses
func (t *CategoryTransformer) TransformMany(categories []*domain.Category) []*CategoryResponse {
	if categories == nil {
		return nil
	}

	result := make([]*CategoryResponse, 0, len(categories))
	for _, category := range categories {
		result = append(result, t.Transform(category))
	}

	return result
}

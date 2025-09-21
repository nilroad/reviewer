package transformer

import (
	"club/internal/constant"
	"club/internal/domain"
	"math"

	"git.oceantim.com/backend/packages/golang/essential/cdnmaker"
	"git.oceantim.com/backend/packages/golang/go-hashids/v2"
)

type CategoryPreviewTransformer struct {
	cdn *cdnmaker.CDNMaker
}

func NewCategoryPreviewTransformer(cdnMakerHandler *cdnmaker.Handler) *CategoryPreviewTransformer {
	return &CategoryPreviewTransformer{
		cdn: cdnMakerHandler.Use(constant.ClubCDN),
	}
}

type CategoryPreviewResponse struct {
	ID           uint64             `json:"id"`
	Name         string             `json:"name"`
	ProductCount int64              `json:"product_count"`
	Products     []*ProductResponse `json:"products"`
}

// ProductResponse represents a product in the API response
type ProductResponse struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Price           uint64  `json:"price"`
	DiscountedPrice *uint64 `json:"discounted_price"`
	DiscountPercent *uint64 `json:"discount_percent"`
	ImageURL        *string `json:"image_url"`
	Stock           *uint64 `json:"stock"`
}

// Transform converts a domain product category with products to an API response
func (t *CategoryPreviewTransformer) Transform(category *domain.CategoryWithProducts) *CategoryPreviewResponse {
	if category == nil {
		return nil
	}

	productResponses := make([]*ProductResponse, 0, len(category.Products))
	for _, product := range category.Products {
		var imageURL *string
		if product.ImageFilePath != nil {
			//todo: temporary comment out cdn make
			//nolint:all
			// url := p.cdn.GenerateFullURL(*product.ImageFilePath)
			imageURL = product.ImageFilePath
		}

		var discountPercentRounded uint64
		if product.DiscountPercent != nil {
			discountPercentRounded = uint64(math.Round(*product.DiscountPercent))
		}

		productResponses = append(productResponses, &ProductResponse{
			//nolint:all
			ID:              hashids.Encode([]int{int(product.ID)}),
			Title:           product.Title,
			Description:     product.Description,
			Price:           product.Price,
			DiscountedPrice: product.DiscountedPrice,
			DiscountPercent: &discountPercentRounded,
			ImageURL:        imageURL,
			Stock:           product.Stock,
		})
	}

	return &CategoryPreviewResponse{
		ID:           category.ID,
		Name:         category.Name,
		ProductCount: category.ProductCount,
		Products:     productResponses,
	}
}

// TransformMany converts multiple domain product categories to API responses
func (t *CategoryPreviewTransformer) TransformMany(categories []*domain.CategoryWithProducts) []*CategoryPreviewResponse {
	if categories == nil {
		return nil
	}

	result := make([]*CategoryPreviewResponse, 0, len(categories))
	for _, category := range categories {
		result = append(result, t.Transform(category))
	}

	return result
}

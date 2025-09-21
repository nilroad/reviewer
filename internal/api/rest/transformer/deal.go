package transformer

import (
	"club/internal/constant"
	"club/internal/domain"
	"time"

	"git.oceantim.com/backend/packages/golang/essential/cdnmaker"
)

type DealListResponse struct {
	ID        uint64    `json:"id"`
	Product   *Product  `json:"product,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type DealTransformer struct {
	cdn                *cdnmaker.CDNMaker
	productTransformer *ProductTransformer
}

func NewDealTransformer(cdnHandler *cdnmaker.Handler, productTransformer *ProductTransformer) *DealTransformer {
	return &DealTransformer{
		cdn:                cdnHandler.Use(constant.ClubCDN),
		productTransformer: productTransformer,
	}
}

func (t *DealTransformer) Transform(deal *domain.Deal) DealListResponse {
	response := DealListResponse{
		ID:        deal.ID,
		CreatedAt: deal.CreatedAt,
	}

	response.Product = t.productTransformer.Transform(deal.Product)

	return response
}

func (t *DealTransformer) TransformMany(deals []*domain.Deal) []DealListResponse {
	result := make([]DealListResponse, 0, len(deals))
	for _, deal := range deals {
		result = append(result, t.Transform(deal))
	}

	return result
}

package transformer

import (
	"club/internal/domain"
)

type AdminProduct struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}

type AdminProductTransformer struct {
}

func NewAdminProductTransformer() *AdminProductTransformer {
	return &AdminProductTransformer{}
}

func (p *AdminProductTransformer) Transform(product *domain.Product) *AdminProduct {

	return &AdminProduct{
		ID:    product.ID,
		Title: product.Title,
	}
}

func (p *AdminProductTransformer) TransformMany(products []*domain.Product) []AdminProduct {
	result := make([]AdminProduct, 0, len(products))
	for _, list := range products {
		result = append(result, *p.Transform(list))
	}

	return result
}

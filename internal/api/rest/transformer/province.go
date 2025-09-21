package transformer

import (
	"club/internal/domain"
)

type Province struct {
	ID    uint64  `json:"id"`
	Name  string  `json:"name"`
	Title string  `json:"title"`
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
}

type ProvinceTransformer struct{}

func NewProvinceTransformer() *ProvinceTransformer {
	return &ProvinceTransformer{}
}

func (t *ProvinceTransformer) Transform(province *domain.Province) *Province {
	if province == nil {
		return nil
	}

	return &Province{
		ID:    province.ID,
		Title: province.Title,
	}
}

func (t *ProvinceTransformer) TransformMany(provinces []*domain.Province) []*Province {
	result := make([]*Province, 0, len(provinces))

	for _, province := range provinces {
		result = append(result, t.Transform(province))
	}

	return result
}

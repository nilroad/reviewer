package transformer

import (
	"club/internal/domain"
)

type City struct {
	ID         uint64 `json:"id"`
	Title      string `json:"title"`
	ProvinceID uint64 `json:"province_id"`
}

type CityTransformer struct{}

func NewCityTransformer() *CityTransformer {
	return &CityTransformer{}
}

func (t *CityTransformer) Transform(city *domain.City) *City {
	if city == nil {
		return nil
	}

	return &City{
		ID:         city.ID,
		Title:      city.Title,
		ProvinceID: city.ProvinceID,
	}
}

func (t *CityTransformer) TransformMany(cities []*domain.City) []*City {
	result := make([]*City, 0, len(cities))

	for _, city := range cities {
		result = append(result, t.Transform(city))
	}

	return result
}

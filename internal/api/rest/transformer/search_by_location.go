package transformer

import "club/internal/domain"

type SearchTransformer struct{}

func NewSearchTransformer() *SearchTransformer {
	return &SearchTransformer{}
}

type LocationResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CityTapsiResponse struct {
	Name   string `json:"name"`
	NameFa string `json:"name_fa"`
}

type ProvinceTapsiResponse struct {
	Name   string `json:"name"`
	NameFa string `json:"name_fa"`
}
type Response struct {
	Title    string                `json:"title"`
	Subtitle string                `json:"subtitle"`
	Location LocationResponse      `json:"location"`
	City     CityTapsiResponse     `json:"city"`
	Province ProvinceTapsiResponse `json:"province"`
}

func (t *SearchTransformer) Transform(s *domain.SearchLocation) *Response {
	if s == nil {
		return nil
	}

	return &Response{
		Title:    s.Title,
		Subtitle: s.Subtitle,
		Location: LocationResponse{
			Latitude:  s.Location.Latitude,
			Longitude: s.Location.Longitude,
		},
		City: CityTapsiResponse{
			Name:   s.City.Name,
			NameFa: s.City.NameFa,
		},
		Province: ProvinceTapsiResponse{
			Name:   s.Province.Name,
			NameFa: s.Province.NameFa,
		},
	}
}

func (t *SearchTransformer) TransformMany(items []*domain.SearchLocation) []*Response {
	result := make([]*Response, 0, len(items))
	for _, s := range items {
		result = append(result, t.Transform(s))
	}

	return result
}

package transformer

import (
	"club/internal/domain"
)

type Branch struct {
	ID        uint64  `json:"id"`
	Title     string  `json:"title"`
	Merchant  string  `json:"merchant"`
	City      string  `json:"city"`
	Province  string  `json:"province"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type BranchTransformer struct{}

func NewBranchTransformer() *BranchTransformer {
	return &BranchTransformer{}
}

func (t *BranchTransformer) Transform(branch *domain.Branch) *Branch {
	if branch == nil {
		return nil
	}

	result := &Branch{
		ID:        branch.ID,
		Title:     branch.Title,
		Address:   branch.Address,
		Latitude:  branch.Lat,
		Longitude: branch.Lon,
	}

	if branch.Merchant != nil {
		result.Merchant = branch.Merchant.Title
	}

	if branch.City != nil {
		result.City = branch.City.Title
		if branch.City.Province != nil {
			result.Province = branch.City.Province.Title
		}
	}

	return result
}

func (t *BranchTransformer) TransformMany(branches []*domain.Branch) []*Branch {
	result := make([]*Branch, 0, len(branches))

	for _, branch := range branches {
		result = append(result, t.Transform(branch))
	}

	return result
}

type SearchBranchTransformer struct{}

func NewSearchBranchTransformer() *SearchBranchTransformer {
	return &SearchBranchTransformer{}
}

type SearchBranchResponse struct {
	ID            uint64    `json:"id"`
	Title         string    `json:"title"`
	Merchant      string    `json:"merchant"`
	CityTitle     *string   `json:"city_title"`
	Province      *Province `json:"province"`
	ProvinceTitle *string   `json:"province_title"`
	Address       string    `json:"address"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
}

func (r *SearchBranchTransformer) Transform(branch *domain.Branch) SearchBranchResponse {
	var cityTitle *string
	var provinceTitle *string
	var province *Province

	if branch.City != nil {
		cityTitle = &branch.City.Title

		if branch.City.Province != nil {
			provinceTitle = &branch.City.Province.Title
			province = &Province{
				ID:    branch.City.Province.ID,
				Title: branch.City.Province.Title,
			}
		}
	}

	return SearchBranchResponse{
		ID:            branch.ID,
		Title:         branch.Title,
		Merchant:      branch.Merchant.Title,
		CityTitle:     cityTitle,
		Province:      province,
		ProvinceTitle: provinceTitle,
		Address:       branch.Address,
		Latitude:      branch.Lat,
		Longitude:     branch.Lon,
	}
}

func (r *SearchBranchTransformer) TransformMany(branches []*domain.Branch) []SearchBranchResponse {
	res := make([]SearchBranchResponse, 0, len(branches))
	for _, branch := range branches {
		res = append(res, r.Transform(branch))
	}

	return res
}

type AdminBranchResponse struct {
	ID            uint64  `json:"id"`
	Title         string  `json:"title"`
	CityID        uint64  `json:"city_id"`
	CityTitle     string  `json:"city_title"`
	ProvinceID    uint64  `json:"province_id"`
	ProvinceTitle string  `json:"province_title"`
	Address       string  `json:"address"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
}

type AdminBranchTransformer struct{}

func NewAdminBranchTransformer() *AdminBranchTransformer {
	return &AdminBranchTransformer{}
}

func (t *AdminBranchTransformer) Transform(branch *domain.Branch) *AdminBranchResponse {
	if branch == nil {
		return nil
	}

	result := &AdminBranchResponse{
		ID:        branch.ID,
		Title:     branch.Title,
		Address:   branch.Address,
		Latitude:  branch.Lat,
		Longitude: branch.Lon,
	}

	if branch.City != nil {
		result.CityID = branch.City.ID
		result.CityTitle = branch.City.Title
		if branch.City.Province != nil {
			result.ProvinceID = branch.City.Province.ID
			result.ProvinceTitle = branch.City.Province.Title
		}
	}

	return result
}

func (t *AdminBranchTransformer) TransformMany(branches []*domain.Branch) []*AdminBranchResponse {
	result := make([]*AdminBranchResponse, 0, len(branches))

	for _, branch := range branches {
		result = append(result, t.Transform(branch))
	}

	return result
}

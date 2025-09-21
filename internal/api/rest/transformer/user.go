package transformer

import (
	"club/internal/constant"
	"club/internal/domain"
	"time"
)

type UserResponse struct {
	ID         uint64                 `json:"id"`
	Name       *string                `json:"name"`
	LastName   *string                `json:"last_name"`
	Cellphone  string                 `json:"cellphone"`
	KYCStatus  constant.UserKYCStatus `json:"kyc_status"`
	Status     domain.UserStatus      `json:"status"`
	BirthDate  *time.Time             `json:"birth_date"`
	BranchID   *uint64                `json:"merchant_branch_id"`
	Email      *string                `json:"email"`
	PostalCode *string                `json:"postal_code"`
	Address    *string                `json:"address"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`

	City           *CityResponse `json:"city"`
	BranchResponse `json:"merchant_branch"`
}

type CityResponse struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	ProvinceID   uint64 `json:"province_id"`
	ProvinceName string `json:"province_name"`
}

type BranchResponse struct {
	ID         uint64 `json:"id"`
	Title      string `json:"title"`
	MerchantID uint64 `json:"merchant_id"`
}

type UserTransformer struct{}

func NewUserTransformer() *UserTransformer {
	return &UserTransformer{}
}

func (t *UserTransformer) Transform(user *domain.User) UserResponse {
	var city *CityResponse
	if user.City != nil {
		city = &CityResponse{
			ID:           user.City.ID,
			Name:         user.City.Title,
			ProvinceName: user.City.Province.Title,
			ProvinceID:   user.City.ProvinceID,
		}
	}

	response := UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		LastName:   user.LastName,
		Cellphone:  user.Cellphone,
		KYCStatus:  user.KYCStatus,
		Status:     user.Status,
		BirthDate:  user.BirthDate,
		BranchID:   user.BranchID,
		Email:      user.Email,
		PostalCode: user.PostalCode,
		Address:    user.Address,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		City:       city,
		BranchResponse: BranchResponse{
			ID:         user.Branch.ID,
			Title:      user.Branch.Title,
			MerchantID: user.Branch.MerchantID,
		},
	}

	return response
}

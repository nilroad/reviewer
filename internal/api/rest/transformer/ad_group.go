package transformer

import (
	"club/internal/constant"
	"club/internal/domain"
	"time"

	"github.com/dustin/go-humanize"
)

type AdGroupResponse struct {
	AdGroupType string       `json:"ad_group_type"`
	IsActive    bool         `json:"is_active"`
	Ads         []AdResponse `json:"ads"`
}

type AdGroupTransformer struct {
	adTransformer *AdTransformer
}

func NewAdGroupTransformer(adTransformer *AdTransformer) *AdGroupTransformer {
	return &AdGroupTransformer{
		adTransformer: adTransformer,
	}
}

func (t *AdGroupTransformer) Transform(group *domain.AdGroup) AdGroupResponse {
	return AdGroupResponse{
		AdGroupType: string(group.AdGroupType),
		IsActive:    group.IsActive,
		Ads:         t.adTransformer.TransformMany(group.Ads),
	}
}

func (t *AdGroupTransformer) TransformMany(groups []*domain.AdGroup) []AdGroupResponse {
	transformedGroups := make([]AdGroupResponse, 0, len(groups))
	for _, group := range groups {
		transformedGroups = append(transformedGroups, t.Transform(group))
	}

	return transformedGroups
}

type AdminAdGroupResponse struct {
	ID          uint64                 `json:"id"`
	Title       string                 `json:"title"`
	AdGroupType constant.AdGroupType   `json:"ad_group_type"`
	CreatedBy   *uint64                `json:"created_by"`
	CreatorUser *UserResponse          `json:"creator_user"`
	BranchCount uint64                 `json:"branch_count"`
	CreatedAt   time.Time              `json:"created_at"`
	IsActive    bool                   `json:"is_active"`
	StartAt     *time.Time             `json:"start_at"`
	EndAt       *time.Time             `json:"end_at"`
	Ads         []AdminAdResponse      `json:"ads"`
	Branches    []*AdminBranchResponse `json:"branches"`
}

type AdminAdGroupTransformer struct {
	userTransformer        *UserTransformer
	adminAdTransformer     *AdminAdTransformer
	adminBranchTransformer *AdminBranchTransformer
}

func NewAdminAdGroupTransformer(userTransformer *UserTransformer, adminAdTransformer *AdminAdTransformer, adminBranchTransformer *AdminBranchTransformer) *AdminAdGroupTransformer {
	return &AdminAdGroupTransformer{
		userTransformer:        userTransformer,
		adminAdTransformer:     adminAdTransformer,
		adminBranchTransformer: adminBranchTransformer,
	}
}

func (t *AdminAdGroupTransformer) Transform(adGroup *domain.AdGroup) AdminAdGroupResponse {
	var creatorUser *UserResponse
	if adGroup.CreatorUser != nil {
		user := t.userTransformer.Transform(adGroup.CreatorUser)
		creatorUser = &user
	}

	ads := t.adminAdTransformer.TransformMany(adGroup.Ads)
	branches := t.adminBranchTransformer.TransformMany(adGroup.Branches)

	return AdminAdGroupResponse{
		ID:          adGroup.ID,
		Title:       adGroup.Title,
		AdGroupType: adGroup.AdGroupType,
		CreatedBy:   adGroup.CreatedBy,
		BranchCount: adGroup.BranchCount,
		CreatedAt:   adGroup.CreatedAt,
		IsActive:    adGroup.IsActive,
		StartAt:     adGroup.StartAt,
		EndAt:       adGroup.EndAt,
		CreatorUser: creatorUser,
		Ads:         ads,
		Branches:    branches,
	}
}

func (t *AdminAdGroupTransformer) TransformMany(adGroups []*domain.AdGroup) []AdminAdGroupResponse {
	result := make([]AdminAdGroupResponse, 0, len(adGroups))
	for _, adGroup := range adGroups {
		result = append(result, t.Transform(adGroup))
	}

	return result
}

type AdminAdGroupConfigResponse struct {
	MaxFileSize          string               `json:"max_file_size"`
	QuadrupleAspectRatio constant.AspectRatio `json:"quadruple_aspect_ratio"`
	TripleAspectRatio    constant.AspectRatio `json:"triple_aspect_ratio"`
	SingleAspectRatio    constant.AspectRatio `json:"single_aspect_ratio"`
}

type AdminAdGroupConfigTransformer struct {
}

func NewAdminAdGroupConfigTransformer() *AdminAdGroupConfigTransformer {
	return &AdminAdGroupConfigTransformer{}
}

func (t *AdminAdGroupConfigTransformer) Transform() AdminAdGroupConfigResponse {
	return AdminAdGroupConfigResponse{
		MaxFileSize:          humanize.Bytes(constant.AdGroupMaxFileSize),
		QuadrupleAspectRatio: constant.AspectRatioQuadruple,
		TripleAspectRatio:    constant.AspectRatioTriple,
		SingleAspectRatio:    constant.AspectRatioSingle,
	}
}

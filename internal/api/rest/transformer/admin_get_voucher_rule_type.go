package transformer

import (
	"club/internal/constant"
	"club/internal/domain"
	"time"
)

type AdminGetVoucherRuleTypeResponse struct {
	VoucherTypeNumber            uint64                            `json:"voucher_type_number"`
	Title                        string                            `json:"title"`
	SubTitle                     string                            `json:"sub_title"`
	VoucherTitle                 string                            `json:"voucher_title"`
	MaximumDiscountAmount        uint64                            `json:"maximum_discount_amount"`
	MaximumExpirationDayNumber   uint64                            `json:"maximum_expiration_day_number"`
	DiscountPercent              float64                           `json:"discount_percent"`
	AssignedCount                uint64                            `json:"assigned_count"`
	MaximumAllowedCount          uint64                            `json:"maximum_allowed_count"`
	MinimumAcceptedInvoice       uint64                            `json:"minimum_accepted_invoice"`
	Description                  string                            `json:"description"`
	ImageURL                     *string                           `json:"image_url"`
	VoucherRuleTypeCategory      *constant.VoucherRuleTypeCategory `json:"voucher_rule_type_category"`
	VoucherRuleTypeCategoryTitle *string                           `json:"voucher_rule_type_category_title"`
	VoucherRuleTypeType          constant.VoucherRuleTypeType      `json:"voucher_rule_type_type"`
	Status                       constant.VoucherRuleTypeStatus    `json:"status"`
	Products                     []Product                         `json:"products"`
	CreatedAt                    time.Time                         `json:"created_at"`
	ExpirationDate               time.Time                         `json:"expiration_date"`
	IsUnmergeable                bool                              `json:"is_unmergeable"`
	VoucherCount                 uint64                            `json:"voucher_count"`
}

type AdminGetVoucherRuleTypeTransformer struct {
	productTransformer *ProductTransformer
}

func NewAdminGetVoucherRuleTypeTransformer(productTransformer *ProductTransformer) *AdminGetVoucherRuleTypeTransformer {
	return &AdminGetVoucherRuleTypeTransformer{
		productTransformer: productTransformer,
	}
}

func (t *AdminGetVoucherRuleTypeTransformer) Transform(voucherRuleType *domain.VoucherRuleType) *AdminGetVoucherRuleTypeResponse {
	var voucherRuleTypeCategory constant.VoucherRuleTypeCategory

	switch voucherRuleType.VoucherRuleTypeCategoryID {
	case constant.VoucherTypeFixDiscountVoucher:
		voucherRuleTypeCategory = constant.VoucherRuleTypeCategoryFixedDiscount
	case constant.VoucherTypePercentageDiscountVoucher:
		voucherRuleTypeCategory = constant.VoucherRuleTypeCategoryPercentageDiscount
	case constant.VoucherTypeFixDiscountQuestRewardVoucher:
		voucherRuleTypeCategory = constant.VoucherRuleFixDiscountQuestReward
	case constant.VoucherTypePercentageDiscountQuestRewardVoucher:
		voucherRuleTypeCategory = constant.VoucherRulePercentageDiscountQuestReward
	}

	var voucherRuleTypeType constant.VoucherRuleTypeType
	if voucherRuleType.IsProductVoucher {
		voucherRuleTypeType = constant.VoucherRuleTypeTypeProductVoucher
	} else {
		voucherRuleTypeType = constant.VoucherRuleTypeTypeInvoiceVoucher
	}

	return &AdminGetVoucherRuleTypeResponse{
		VoucherTypeNumber:            voucherRuleType.VoucherTypeNumber,
		MaximumDiscountAmount:        voucherRuleType.MaximumDiscountAmount,
		DiscountPercent:              voucherRuleType.DiscountPercent,
		AssignedCount:                voucherRuleType.AssignedCount,
		MaximumAllowedCount:          voucherRuleType.MaximumAllowedCount,
		MinimumAcceptedInvoice:       voucherRuleType.MinimumAcceptedInvoice,
		MaximumExpirationDayNumber:   voucherRuleType.MaximumExpirationDayNumber,
		Description:                  voucherRuleType.Description,
		ImageURL:                     voucherRuleType.ImageURL,
		Title:                        voucherRuleType.Title,
		SubTitle:                     voucherRuleType.SubTitle,
		VoucherTitle:                 voucherRuleType.VoucherTitle,
		VoucherRuleTypeCategory:      &voucherRuleTypeCategory,
		VoucherRuleTypeCategoryTitle: voucherRuleType.VoucherRuleTypeCategoryTitle,
		VoucherRuleTypeType:          voucherRuleTypeType,
		Status:                       voucherRuleType.Status,
		CreatedAt:                    voucherRuleType.CreatedAt,
		ExpirationDate:               voucherRuleType.ExpirationDate,
		Products:                     t.productTransformer.TransformMany(voucherRuleType.Products),
		IsUnmergeable:                voucherRuleType.IsUnmergeable,
		VoucherCount:                 voucherRuleType.VoucherCount,
	}
}

func (t *AdminGetVoucherRuleTypeTransformer) TransformMany(voucherRuleTypes []*domain.VoucherRuleType) []*AdminGetVoucherRuleTypeResponse {
	responses := make([]*AdminGetVoucherRuleTypeResponse, len(voucherRuleTypes))

	for i, voucherRuleType := range voucherRuleTypes {
		responses[i] = t.Transform(voucherRuleType)
	}

	return responses
}

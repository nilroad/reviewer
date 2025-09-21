package transformer

import (
	"club/internal/constant"
	"club/internal/domain"
	"time"
)

type VoucherResponse struct {
	VoucherNumber             uint64     `json:"voucher_number"`
	VoucherTypeNumber         uint64     `json:"voucher_type_number"`
	VoucherTypeTitle          string     `json:"voucher_type_title"`
	VoucherTitle              string     `json:"voucher_title"`
	ExpirationDate            time.Time  `json:"expiration_date"`
	VoucherStatusID           uint64     `json:"voucher_status_id"`
	VoucherStatusName         string     `json:"voucher_status_name"`
	VoucherStatusTitle        string     `json:"voucher_status_title"`
	CreateDate                time.Time  `json:"create_date"`
	SpentDate                 *time.Time `json:"spent_date,omitempty"`
	VoucherRuleTypeCategoryID uint64     `json:"voucher_rule_type_category_id"`
	VoucherPrice              uint64     `json:"voucher_price"`
	MinimumAcceptedInvoice    uint64     `json:"minimum_accepted_invoice"`
	DiscountPercent           float64    `json:"discount_percent"`
	MaximumDiscountAmount     uint64     `json:"maximum_discount_amount"`
	VoucherImageURL           *string    `json:"voucher_image_url"`
	IsProductVoucher          bool       `json:"is_product_voucher"`
}

type VoucherTransformer struct{}

func NewVoucherTransformer() *VoucherTransformer {
	return &VoucherTransformer{}
}

// Transform transforms a domain voucher to a response voucher
func (t *VoucherTransformer) Transform(voucher domain.Voucher) VoucherResponse {

	return VoucherResponse{
		VoucherNumber:             voucher.VoucherNumber,
		VoucherTypeNumber:         voucher.VoucherTypeNumber,
		VoucherTypeTitle:          voucher.VoucherTypeTitle,
		VoucherTitle:              voucher.VoucherTitle,
		ExpirationDate:            voucher.ExpirationDate,
		VoucherStatusID:           voucher.VoucherStatusID,
		VoucherStatusName:         voucher.VoucherStatusName,
		VoucherStatusTitle:        voucher.VoucherStatusTitle,
		CreateDate:                voucher.CreateDate,
		SpentDate:                 voucher.SpentDate,
		VoucherRuleTypeCategoryID: voucher.VoucherRuleTypeCategoryID,
		VoucherPrice:              voucher.VoucherPrice,
		MinimumAcceptedInvoice:    voucher.MinimumAcceptedInvoice,
		DiscountPercent:           voucher.DiscountPercent,
		MaximumDiscountAmount:     voucher.MaximumDiscountAmount,
		VoucherImageURL:           voucher.ImageURL,
		IsProductVoucher:          len(voucher.Items) > 0,
	}
}

// TransformMany transforms multiple domain vouchers to response vouchers
func (t *VoucherTransformer) TransformMany(vouchers []domain.Voucher) []VoucherResponse {
	result := make([]VoucherResponse, len(vouchers))
	for i, v := range vouchers {
		result[i] = t.Transform(v)
	}

	return result
}

type AdminVoucherResponse struct {
	VoucherNumber             uint64                           `json:"voucher_number"`
	VoucherTypeNumber         uint64                           `json:"voucher_type_number"`
	VoucherTypeTitle          string                           `json:"voucher_type_title"`
	VoucherTitle              string                           `json:"voucher_title"`
	ExpirationDate            time.Time                        `json:"expiration_date"`
	VoucherStatusID           uint64                           `json:"voucher_status_id"`
	VoucherStatusName         string                           `json:"voucher_status_name"`
	VoucherStatusTitle        string                           `json:"voucher_status_title"`
	VoucherStatus             constant.VoucherStatus           `json:"voucher_status"`
	CreateDate                time.Time                        `json:"create_date"`
	SpentDate                 *time.Time                       `json:"spent_date,omitempty"`
	VoucherRuleTypeCategoryID uint64                           `json:"voucher_rule_type_category_id"`
	VoucherRuleTypeCategory   constant.VoucherRuleTypeCategory `json:"voucher_rule_type_category"`
	VoucherPrice              uint64                           `json:"voucher_price"`
	MinimumAcceptedInvoice    uint64                           `json:"minimum_accepted_invoice"`
	DiscountPercent           float64                          `json:"discount_percent"`
	MaximumDiscountAmount     uint64                           `json:"maximum_discount_amount"`
	VoucherImageURL           *string                          `json:"voucher_image_url"`
	IsProductVoucher          bool                             `json:"is_product_voucher"`
	VoucherRuleTypeType       constant.VoucherRuleTypeType     `json:"voucher_rule_type_type"`
}

type AdminVoucherTransformer struct{}

func NewAdminVoucherTransformer() *AdminVoucherTransformer {
	return &AdminVoucherTransformer{}
}

// Transform transforms a domain voucher to a response voucher
func (t *AdminVoucherTransformer) Transform(voucher domain.Voucher) AdminVoucherResponse {

	var voucherRuleTypeCategory constant.VoucherRuleTypeCategory
	switch voucher.VoucherRuleTypeCategoryID {
	case constant.VoucherTypeFixDiscountVoucher, constant.VoucherTypeFixDiscountSelectableGiftVoucher:
		voucherRuleTypeCategory = constant.VoucherRuleTypeCategoryFixedDiscount
	case constant.VoucherTypePercentageDiscountVoucher, constant.VoucherTypePercentageDiscountSelectableGiftVoucher:
		voucherRuleTypeCategory = constant.VoucherRuleTypeCategoryPercentageDiscount
	}

	var voucherStatus constant.VoucherStatus
	switch voucher.VoucherStatusID {
	case constant.OzoneCardVoucherStatusActive:
		voucherStatus = constant.VoucherStatusActive
	case constant.OzoneCardVoucherStatusInactiveForPurchase:
		voucherStatus = constant.VoucherStatusInactiveForPurchase
	case constant.OzoneCardVoucherStatusExpired:
		voucherStatus = constant.VoucherStatusExpired
	case constant.OzoneCardVoucherStatusDeactive:
		voucherStatus = constant.VoucherStatusBlocked
	case constant.OzoneCardVoucherStatusSpent:
		voucherStatus = constant.VoucherStatusSpent
	}

	var voucherRuleTypeType constant.VoucherRuleTypeType
	if len(voucher.Items) > 0 {
		voucherRuleTypeType = constant.VoucherRuleTypeTypeProductVoucher
	} else {
		voucherRuleTypeType = constant.VoucherRuleTypeTypeInvoiceVoucher
	}

	return AdminVoucherResponse{
		VoucherNumber:             voucher.VoucherNumber,
		VoucherTypeNumber:         voucher.VoucherTypeNumber,
		VoucherTypeTitle:          voucher.VoucherTypeTitle,
		VoucherTitle:              voucher.VoucherTitle,
		ExpirationDate:            voucher.ExpirationDate,
		VoucherStatusID:           voucher.VoucherStatusID,
		VoucherStatusName:         voucher.VoucherStatusName,
		VoucherStatusTitle:        voucher.VoucherStatusTitle,
		VoucherStatus:             voucherStatus,
		CreateDate:                voucher.CreateDate,
		SpentDate:                 voucher.SpentDate,
		VoucherRuleTypeCategoryID: voucher.VoucherRuleTypeCategoryID,
		VoucherRuleTypeCategory:   voucherRuleTypeCategory,
		VoucherPrice:              voucher.VoucherPrice,
		MinimumAcceptedInvoice:    voucher.MinimumAcceptedInvoice,
		DiscountPercent:           voucher.DiscountPercent,
		MaximumDiscountAmount:     voucher.MaximumDiscountAmount,
		VoucherImageURL:           voucher.ImageURL,
		IsProductVoucher:          len(voucher.Items) > 0,
		VoucherRuleTypeType:       voucherRuleTypeType,
	}
}

// TransformMany transforms multiple domain vouchers to response vouchers
func (t *AdminVoucherTransformer) TransformMany(vouchers []domain.Voucher) []AdminVoucherResponse {
	result := make([]AdminVoucherResponse, len(vouchers))
	for i, v := range vouchers {
		result[i] = t.Transform(v)
	}

	return result
}

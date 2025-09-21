package transformer

import (
	"club/internal/constant"
	"club/internal/domain"
	"time"
)

type VoucherItemResponse struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}

type VoucherDetailResponse struct {
	VoucherNumber                uint64                `json:"voucher_number"`
	ExpirationDate               time.Time             `json:"expiration_date"`
	MinimumAcceptedInvoice       uint64                `json:"minimum_accepted_invoice"`
	MaximumDiscountAmount        uint64                `json:"maximum_discount_amount"`
	DiscountPercent              float64               `json:"discount_percent"`
	VoucherPrice                 *uint64               `json:"voucher_price"`
	VoucherStatusID              uint64                `json:"voucher_status_id"`
	VoucherStatusTitle           string                `json:"voucher_status_title"`
	ImageURL                     *string               `json:"image_url"`
	VoucherRuleTypeCategoryID    *uint64               `json:"voucher_rule_type_category_id"`
	VoucherRuleTypeTitle         string                `json:"voucher_rule_type_title"`
	VoucherRuleTypeDescription   string                `json:"voucher_rule_type_description"`
	VoucherRuleTypeCategoryTitle string                `json:"voucher_rule_type_category_title"`
	IsFixDiscount                bool                  `json:"is_fix_discount"`
	VoucherRuleTypeNumber        uint64                `json:"voucher_rule_type_number"`
	AcceptorMerchantIDs          []uint64              `json:"acceptor_merchant_ids,omitempty"`
	Items                        []VoucherItemResponse `json:"items,omitempty"`
	IsProductVoucher             bool                  `json:"is_product_voucher"`
}

type VoucherDetailTransformer struct{}

func NewVoucherDetailTransformer() *VoucherDetailTransformer {
	return &VoucherDetailTransformer{}
}

// Transform transforms a domain VoucherDetail to a response VoucherDetailResponse
func (t *VoucherDetailTransformer) Transform(voucherDetail domain.VoucherDetail) VoucherDetailResponse {
	// Transform items
	items := make([]VoucherItemResponse, 0, len(voucherDetail.Items))
	for _, item := range voucherDetail.Items {
		items = append(items, VoucherItemResponse{
			ID:    item.ID,
			Title: item.Title,
		})
	}

	return VoucherDetailResponse{
		VoucherNumber:                voucherDetail.VoucherNumber,
		ExpirationDate:               voucherDetail.ExpirationDate,
		MinimumAcceptedInvoice:       voucherDetail.MinimumAcceptedInvoice,
		MaximumDiscountAmount:        voucherDetail.MaximumDiscountAmount,
		DiscountPercent:              voucherDetail.DiscountPercent,
		VoucherPrice:                 voucherDetail.VoucherPrice,
		VoucherStatusID:              voucherDetail.VoucherStatusID,
		VoucherStatusTitle:           voucherDetail.VoucherStatusTitle,
		ImageURL:                     voucherDetail.ImageURL,
		VoucherRuleTypeCategoryID:    voucherDetail.VoucherRuleTypeCategoryID,
		VoucherRuleTypeTitle:         voucherDetail.VoucherRuleTypeTitle,
		VoucherRuleTypeDescription:   voucherDetail.VoucherRuleTypeDescription,
		VoucherRuleTypeCategoryTitle: voucherDetail.VoucherRuleTypeCategoryTitle,
		IsFixDiscount:                voucherDetail.IsFixDiscount,
		VoucherRuleTypeNumber:        voucherDetail.VoucherRuleTypeNumber,
		AcceptorMerchantIDs:          voucherDetail.AcceptorMerchantIDs,
		Items:                        items,
		IsProductVoucher:             voucherDetail.IsProductVoucher,
	}
}

type AdminVoucherItemResponse struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}

type AdminVoucherDetailResponse struct {
	VoucherNumber                 uint64                            `json:"voucher_number"`
	ExpirationDate                time.Time                         `json:"expiration_date"`
	MinimumAcceptedInvoice        uint64                            `json:"minimum_accepted_invoice"`
	MaximumDiscountAmount         uint64                            `json:"maximum_discount_amount"`
	DiscountPercent               float64                           `json:"discount_percent"`
	VoucherPrice                  *uint64                           `json:"voucher_price"`
	VoucherStatusID               uint64                            `json:"voucher_status_id"`
	VoucherStatusTitle            string                            `json:"voucher_status_title"`
	VoucherStatus                 constant.VoucherStatus            `json:"voucher_status"`
	ImageURL                      *string                           `json:"image_url"`
	VoucherRuleTypeCategoryID     *uint64                           `json:"voucher_rule_type_category_id"`
	VoucherRuleTypeCategory       *constant.VoucherRuleTypeCategory `json:"voucher_rule_type_category"`
	VoucherRuleTypeTitle          string                            `json:"voucher_rule_type_title"`
	VoucherRuleTypeDescription    string                            `json:"voucher_rule_type_description"`
	VoucherRuleTypeCategoryTitle  string                            `json:"voucher_rule_type_category_title"`
	IsFixDiscount                 bool                              `json:"is_fix_discount"`
	VoucherRuleTypeNumber         uint64                            `json:"voucher_rule_type_number"`
	AcceptorMerchantIDs           []uint64                          `json:"acceptor_merchant_ids,omitempty"`
	Items                         []AdminVoucherItemResponse        `json:"items,omitempty"`
	IsProductVoucher              bool                              `json:"is_product_voucher"`
	VoucherRuleTypeType           constant.VoucherRuleTypeType      `json:"voucher_rule_type_type"`
	MaximumExpirationDayNumber    uint64                            `json:"maximum_expiration_day_number"`
	IsUnmergeable                 bool                              `json:"is_unmergeable"`
	VoucherRuleTypeExpirationDate time.Time                         `json:"voucher_rule_type_expiration_date"`
}

type AdminVoucherDetailTransformer struct{}

func NewAdminVoucherDetailTransformer() *AdminVoucherDetailTransformer {
	return &AdminVoucherDetailTransformer{}
}

// Transform transforms a domain VoucherDetail to a response VoucherDetailResponse
func (t *AdminVoucherDetailTransformer) Transform(voucherDetail domain.VoucherDetail) AdminVoucherDetailResponse {
	// Transform items
	items := make([]AdminVoucherItemResponse, 0, len(voucherDetail.Items))
	for _, item := range voucherDetail.Items {
		items = append(items, AdminVoucherItemResponse{
			ID:    item.ID,
			Title: item.Title,
		})
	}

	var voucherRuleTypeCategory *constant.VoucherRuleTypeCategory
	if voucherDetail.VoucherRuleTypeCategoryID != nil {
		switch *voucherDetail.VoucherRuleTypeCategoryID {
		case constant.VoucherTypeFixDiscountVoucher, constant.VoucherTypeFixDiscountSelectableGiftVoucher:
			vrt := constant.VoucherRuleTypeCategoryFixedDiscount
			voucherRuleTypeCategory = &vrt
		case constant.VoucherTypePercentageDiscountVoucher, constant.VoucherTypePercentageDiscountSelectableGiftVoucher:
			vrt := constant.VoucherRuleTypeCategoryPercentageDiscount
			voucherRuleTypeCategory = &vrt
		}
	}

	var voucherStatus constant.VoucherStatus
	switch voucherDetail.VoucherStatusID {
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
	if len(voucherDetail.Items) > 0 {
		voucherRuleTypeType = constant.VoucherRuleTypeTypeProductVoucher
	} else {
		voucherRuleTypeType = constant.VoucherRuleTypeTypeInvoiceVoucher
	}

	return AdminVoucherDetailResponse{
		VoucherNumber:                 voucherDetail.VoucherNumber,
		ExpirationDate:                voucherDetail.ExpirationDate,
		MinimumAcceptedInvoice:        voucherDetail.MinimumAcceptedInvoice,
		MaximumDiscountAmount:         voucherDetail.MaximumDiscountAmount,
		DiscountPercent:               voucherDetail.DiscountPercent,
		VoucherPrice:                  voucherDetail.VoucherPrice,
		VoucherStatusID:               voucherDetail.VoucherStatusID,
		VoucherStatusTitle:            voucherDetail.VoucherStatusTitle,
		VoucherStatus:                 voucherStatus,
		ImageURL:                      voucherDetail.ImageURL,
		VoucherRuleTypeCategoryID:     voucherDetail.VoucherRuleTypeCategoryID,
		VoucherRuleTypeCategory:       voucherRuleTypeCategory,
		VoucherRuleTypeTitle:          voucherDetail.VoucherRuleTypeTitle,
		VoucherRuleTypeDescription:    voucherDetail.VoucherRuleTypeDescription,
		VoucherRuleTypeCategoryTitle:  voucherDetail.VoucherRuleTypeCategoryTitle,
		IsFixDiscount:                 voucherDetail.IsFixDiscount,
		VoucherRuleTypeNumber:         voucherDetail.VoucherRuleTypeNumber,
		AcceptorMerchantIDs:           voucherDetail.AcceptorMerchantIDs,
		Items:                         items,
		IsProductVoucher:              voucherDetail.IsProductVoucher,
		VoucherRuleTypeType:           voucherRuleTypeType,
		MaximumExpirationDayNumber:    voucherDetail.MaximumExpirationDayNumber,
		IsUnmergeable:                 voucherDetail.IsUnmergeable,
		VoucherRuleTypeExpirationDate: voucherDetail.VoucherRuleTypeExpirationDate,
	}
}

package transformer

import (
	"club/internal/constant"
	"club/internal/domain"
	"math"
	"time"

	"git.oceantim.com/backend/packages/golang/essential/cdnmaker"
	"git.oceantim.com/backend/packages/golang/go-hashids/v2"
)

type VoucherRuleTypeResponse struct {
	ExpirationDate               time.Time                        `json:"expiration_date"`
	MaximumDiscountAmount        uint64                           `json:"maximum_discount_amount"`
	VoucherTypeNumber            uint64                           `json:"voucher_type_number"`
	MaximumExpirationDayNumber   uint64                           `json:"maximum_expiration_day_number"`
	DiscountPercent              float64                          `json:"discount_percent"`
	MaximumAllowedCount          uint64                           `json:"maximum_allowed_count"`
	MinimumAcceptedInvoice       uint64                           `json:"minimum_accepted_invoice"`
	Title                        string                           `json:"title"`
	SubTitle                     string                           `json:"sub_title"`
	VoucherTitle                 string                           `json:"voucher_title"`
	Description                  string                           `json:"description"`
	VoucherPrice                 *uint64                          `json:"voucher_price"`
	ImageURL                     *string                          `json:"image_url"`
	VoucherRuleTypeCategoryID    uint64                           `json:"voucher_rule_type_category_id"`
	VoucherRuleTypeCategoryTitle *string                          `json:"voucher_rule_type_category_title"`
	Products                     []Product                        `json:"products"`
	IsProductVoucher             bool                             `json:"is_product_voucher"`
	AlreadyRedeemed              bool                             `json:"already_redeemed"`
	VoucherRuleTypeType          constant.VoucherRuleTypeType     `json:"voucher_rule_type_type"`
	VoucherRuleTypeCategory      constant.VoucherRuleTypeCategory `json:"voucher_rule_type_category"`
}

type VoucherRuleTypeTransformer struct {
	cdn *cdnmaker.CDNMaker
}

func NewVoucherRuleTypeTransformer(cdn *cdnmaker.CDNMaker) *VoucherRuleTypeTransformer {
	return &VoucherRuleTypeTransformer{
		cdn: cdn,
	}
}

// Transform transforms a domain VoucherRuleType to a response VoucherRuleType
func (t *VoucherRuleTypeTransformer) Transform(voucherRuleType domain.VoucherRuleType) VoucherRuleTypeResponse {

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

	products := make([]Product, 0, len(voucherRuleType.Products))
	for _, product := range voucherRuleType.Products {

		var discountPercentRounded uint64
		if product.DiscountPercent != nil {
			discountPercentRounded = uint64(math.Round(*product.DiscountPercent))
		}

		products = append(products, Product{
			ID:                     hashids.Encode([]int{int(product.ID)}), //nolint:all
			Title:                  product.Title,
			Description:            product.Description,
			Price:                  product.Price,
			DiscountedPrice:        product.DiscountedPrice,
			DiscountPercent:        &discountPercentRounded,
			VoucherDiscountedPrice: product.VoucherDiscountedPrice,
			VoucherDiscountPercent: product.VoucherDiscountPercent,
			ImageURL:               product.ImageFilePath,
		})
	}

	return VoucherRuleTypeResponse{
		ExpirationDate:               voucherRuleType.ExpirationDate,
		MaximumDiscountAmount:        voucherRuleType.MaximumDiscountAmount,
		VoucherTypeNumber:            voucherRuleType.VoucherTypeNumber,
		MaximumExpirationDayNumber:   voucherRuleType.MaximumExpirationDayNumber,
		DiscountPercent:              voucherRuleType.DiscountPercent,
		MaximumAllowedCount:          voucherRuleType.MaximumAllowedCount,
		MinimumAcceptedInvoice:       voucherRuleType.MinimumAcceptedInvoice,
		Title:                        voucherRuleType.Title,
		SubTitle:                     voucherRuleType.SubTitle,
		VoucherTitle:                 voucherRuleType.VoucherTitle,
		Description:                  voucherRuleType.Description,
		VoucherPrice:                 voucherRuleType.VoucherPrice,
		ImageURL:                     voucherRuleType.ImageURL,
		VoucherRuleTypeCategoryID:    voucherRuleType.VoucherRuleTypeCategoryID,
		VoucherRuleTypeCategoryTitle: voucherRuleType.VoucherRuleTypeCategoryTitle,
		Products:                     products,
		IsProductVoucher:             voucherRuleType.IsProductVoucher,
		AlreadyRedeemed:              voucherRuleType.AlreadyRedeemed,
		VoucherRuleTypeCategory:      voucherRuleTypeCategory,
		VoucherRuleTypeType:          voucherRuleTypeType,
	}
}

// TransformMany transforms multiple domain VoucherRuleType to response VoucherRuleType
func (t *VoucherRuleTypeTransformer) TransformMany(voucherRuleTypes []domain.VoucherRuleType) []VoucherRuleTypeResponse {
	result := make([]VoucherRuleTypeResponse, len(voucherRuleTypes))
	for i, v := range voucherRuleTypes {
		result[i] = t.Transform(v)
	}

	return result
}

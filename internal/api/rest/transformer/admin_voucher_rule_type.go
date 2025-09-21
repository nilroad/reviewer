package transformer

import "club/internal/domain"

type AdminCreateVoucherRuleTypeResponse struct {
	VoucherTypeNumber uint64 `json:"voucher_type_number"`
}

type AdminCreateVoucherRuleTypeTransformer struct{}

func NewAdminCreateVoucherRuleTypeTransformer() *AdminCreateVoucherRuleTypeTransformer {
	return &AdminCreateVoucherRuleTypeTransformer{}
}

func (t *AdminCreateVoucherRuleTypeTransformer) Transform(voucherRuleType *domain.VoucherRuleType) *AdminCreateVoucherRuleTypeResponse {
	return &AdminCreateVoucherRuleTypeResponse{
		VoucherTypeNumber: voucherRuleType.VoucherTypeNumber,
	}
}

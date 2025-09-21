package transformer

import (
	"club/internal/domain"
)

type BarcodeResponse struct {
	BarcodeData        string `json:"barcode_data"`
	ExpirationDatetime int64  `json:"expiration_datetime"`
}

type PaymentTransformer struct{}

func NewPaymentTransformer() *PaymentTransformer {
	return &PaymentTransformer{}
}

func (t *PaymentTransformer) Transform(data *domain.IPGInvoiceBarcode) *BarcodeResponse {
	return &BarcodeResponse{
		BarcodeData:        data.BarcodeData,
		ExpirationDatetime: data.ExpirationDatetime,
	}
}

package transformer

import (
	"club/internal/domain"
	"time"
)

type WalletTransformer struct{}

func NewWalletTransformer() *WalletTransformer {
	return &WalletTransformer{}
}

type WalletInfoResponse struct {
	Balance           uint64    `json:"balance"`
	WalletCode        string    `json:"wallet_code"`
	WalletStatusID    uint64    `json:"wallet_status_id"`
	WalletStatusName  string    `json:"wallet_status_name"`
	WalletStatusTitle string    `json:"wallet_status_title"`
	CreatedAt         time.Time `json:"created_at"`
}

func (t *WalletTransformer) Transform(w *domain.Wallet) *WalletInfoResponse {
	return &WalletInfoResponse{
		Balance:           w.Balance,
		WalletCode:        w.SegmentCode,
		WalletStatusID:    w.SegmentStatusID,
		WalletStatusName:  w.SegmentStatusName,
		WalletStatusTitle: w.SegmentStatusTitle,
		CreatedAt:         time.Now(), // TODO: change to created_at by ozone card proxy value
	}
}

type WalletTransactionResponse struct {
	RequestTypeID    int64     `json:"request_type_id"`
	RequestTypeTitle string    `json:"request_type_title"`
	CreateDate       time.Time `json:"create_date"`
	Amount           int64     `json:"amount"`
	FromClientID     int64     `json:"from_client_id"`
	ToClientID       int64     `json:"to_client_id"`
	FromClientCode   string    `json:"from_client_code"`
	ToClientCode     string    `json:"to_client_code"`
	FromClientName   string    `json:"from_client_name"`
	ToClientName     string    `json:"to_client_name"`
	Description      string    `json:"description"`
	IsReverse        bool      `json:"is_reverse"`
}

type WalletTransactionTransformer struct{}

func NewWalletTransactionTransformer() *WalletTransactionTransformer {
	return &WalletTransactionTransformer{}
}

func (t *WalletTransactionTransformer) Transform(wt *domain.WalletTransaction) *WalletTransactionResponse {
	return &WalletTransactionResponse{
		RequestTypeID:    wt.RequestTypeID,
		RequestTypeTitle: wt.RequestTypeTitle,
		CreateDate:       wt.CreateDate,
		Amount:           wt.Amount,
		FromClientID:     wt.FromClientID,
		ToClientID:       wt.ToClientID,
		FromClientCode:   wt.FromClientCode,
		ToClientCode:     wt.ToClientCode,
		FromClientName:   wt.FromClientName,
		ToClientName:     wt.ToClientName,
		Description:      wt.Description,
		IsReverse:        wt.IsReverse,
	}
}

func (t *WalletTransactionTransformer) TransformMany(wt []*domain.WalletTransaction) []WalletTransactionResponse {
	transactions := make([]WalletTransactionResponse, 0, len(wt))
	for _, transaction := range wt {
		transactions = append(transactions, *t.Transform(transaction))
	}

	return transactions
}

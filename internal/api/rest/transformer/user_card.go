package transformer

import (
	"club/internal/constant"
	"club/internal/domain"

	"git.oceantim.com/backend/packages/golang/essential/cdnmaker"
)

type UserCardResponse struct {
	ID         uint64    `json:"id"`
	CardNumber string    `json:"card_number"`
	Status     string    `json:"status"`
	BankInfo   *BankInfo `json:"bank_info"`
}

type UserCardTransformer struct {
	cdn *cdnmaker.CDNMaker
}

func NewUserCardTransformer(cdn *cdnmaker.Handler) *UserCardTransformer {
	return &UserCardTransformer{
		cdn: cdn.Use(constant.ClubCDN),
	}
}

func (t *UserCardTransformer) Transform(userCard *domain.UserCard) *UserCardResponse {
	if userCard == nil {
		return nil
	}

	response := &UserCardResponse{
		ID:         userCard.ID,
		CardNumber: userCard.CardNumber,
		Status:     string(userCard.Status),
		BankInfo: &BankInfo{
			BankTitle: userCard.BankInfo.BankTitle,
			Icon:      t.cdn.GenerateFullURL(userCard.BankInfo.Icon),
		},
	}

	return response
}

func (t *UserCardTransformer) TransformMany(userCards []*domain.UserCard) []*UserCardResponse {
	if userCards == nil {
		return nil
	}

	responses := make([]*UserCardResponse, 0, len(userCards))
	for _, userCard := range userCards {
		responses = append(responses, t.Transform(userCard))
	}

	return responses
}

type BankInfo struct {
	BankTitle string `json:"bank_title"`
	Icon      string `json:"icon"`
}

type CardInquiryTransformer struct {
	cdn *cdnmaker.CDNMaker
}

func NewCardInquiryTransformer(cdn *cdnmaker.Handler) *CardInquiryTransformer {
	return &CardInquiryTransformer{
		cdn: cdn.Use(constant.ClubCDN),
	}
}

func (t *CardInquiryTransformer) Transform(bankInfo *domain.BankInfo) *BankInfo {
	if bankInfo == nil {
		return nil
	}

	return &BankInfo{
		BankTitle: bankInfo.BankTitle,
		Icon:      t.cdn.GenerateFullURL(bankInfo.Icon),
	}
}

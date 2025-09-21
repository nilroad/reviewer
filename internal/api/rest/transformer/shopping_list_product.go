package transformer

import (
	"club/internal/constant"
	"club/internal/domain"
	"time"

	"git.oceantim.com/backend/packages/golang/essential/cdnmaker"
)

type ShoppingListResponse struct {
	ID           uint64    `json:"id"`
	Title        string    `json:"title"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ProductCount uint64    `json:"product_count"`
}

type ShoppingListTransformer struct {
	cdn *cdnmaker.CDNMaker
}

func NewShoppingListTransformer(cdn *cdnmaker.Handler) *ShoppingListTransformer {
	return &ShoppingListTransformer{
		cdn: cdn.Use(constant.ClubCDN),
	}
}

func (t *ShoppingListTransformer) Transform(userShoppingList *domain.UserShoppingList) ShoppingListResponse {
	return ShoppingListResponse{
		ID:           userShoppingList.ID,
		Title:        userShoppingList.Title,
		ProductCount: userShoppingList.ProductCount,
		CreatedAt:    userShoppingList.CreatedAt,
		UpdatedAt:    userShoppingList.UpdatedAt,
	}
}

func (t *ShoppingListTransformer) TransformMany(shoppingLists []*domain.UserShoppingList) []ShoppingListResponse {
	result := make([]ShoppingListResponse, 0, len(shoppingLists))
	for _, list := range shoppingLists {
		result = append(result, t.Transform(list))
	}

	return result
}

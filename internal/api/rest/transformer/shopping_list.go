package transformer

import (
	"club/internal/constant"
	"club/internal/domain"
	"time"

	"git.oceantim.com/backend/packages/golang/essential/cdnmaker"
	"git.oceantim.com/backend/packages/golang/go-hashids/v2"
)

type ShoppingListProductResponse struct {
	ID          string    `json:"id"`
	ProductID   *string   `json:"product_id,omitempty"`
	ProductName *string   `json:"product_name,omitempty"`
	Product     *Product  `json:"product,omitempty"`
	Count       uint64    `json:"count"`
	CreatedAt   time.Time `json:"created_at"`
}

type ShoppingListProductTransformer struct {
	cdn *cdnmaker.CDNMaker
}

func NewShoppingListProductTransformer(cdn *cdnmaker.Handler) *ShoppingListProductTransformer {
	return &ShoppingListProductTransformer{
		cdn: cdn.Use(constant.ClubCDN),
	}
}

// TransformProduct transforms a single UserShoppingListProduct to a ShoppingListProductResponse
func (t *ShoppingListProductTransformer) Transform(product *domain.UserShoppingListProduct) *ShoppingListProductResponse {
	if product == nil {
		return nil
	}

	response := &ShoppingListProductResponse{
		//nolint:all
		ID:          hashids.Encode([]int{int(product.ID)}),
		ProductName: product.ProductName,
		Count:       product.Count,
		CreatedAt:   product.CreatedAt,
	}

	// Handle optional ProductID
	if product.ProductID != nil {
		//nolint:all
		productID := hashids.Encode([]int{int(*product.ProductID)})
		response.ProductID = &productID
	}

	// Handle optional nested Product
	if product.Product != nil {
		productTransformer := ProductTransformer{
			cdn: t.cdn,
		}
		response.Product = productTransformer.Transform(product.Product)
	}

	return response
}

// TransformManyProducts transforms a slice of UserShoppingListProduct to ShoppingListProductResponse
func (t *ShoppingListProductTransformer) TransformMany(products []*domain.UserShoppingListProduct) []*ShoppingListProductResponse {
	if products == nil {
		return nil
	}

	result := make([]*ShoppingListProductResponse, 0, len(products))
	for _, product := range products {
		result = append(result, t.Transform(product))
	}

	return result
}

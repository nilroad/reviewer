package transformer

import (
	"club/internal/domain"
	"math"

	"git.oceantim.com/backend/packages/golang/essential/cdnmaker"
	"git.oceantim.com/backend/packages/golang/go-hashids/v2"
)

type Product struct {
	ID                     string     `json:"id"`
	Title                  string     `json:"title"`
	Description            string     `json:"description"`
	Price                  uint64     `json:"price"`
	DiscountedPrice        *uint64    `json:"discounted_price"`
	DiscountPercent        *uint64    `json:"discount_percent"`
	VoucherDiscountedPrice *uint64    `json:"voucher_discounted_price,omitempty"`
	VoucherDiscountPercent *float64   `json:"voucher_discount_percent,omitempty"`
	ImageURL               *string    `json:"image_url"`
	Stock                  *uint64    `json:"stock"`
	Categories             []Category `json:"categories"`
	BrandTitle             *string    `json:"brand_title"`
}

type Category struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type ProductTransformer struct {
	cdn *cdnmaker.CDNMaker
}

func NewProductTransformer(cdn *cdnmaker.CDNMaker) *ProductTransformer {
	return &ProductTransformer{
		cdn: cdn,
	}
}

func (p *ProductTransformer) Transform(product *domain.Product) *Product {
	var imageURL *string
	if product.ImageFilePath != nil {
		//todo: temporary comment out cdn make
		//nolint:all
		// url := p.cdn.GenerateFullURL(*product.ImageFilePath)
		imageURL = product.ImageFilePath
	}

	var brandTitle *string
	if product.Brand != nil {
		brandTitle = &product.Brand.Title
	}

	categories := make([]Category, 0, len(product.Category))
	for _, category := range product.Category {
		categories = append(categories, Category{
			ID:    hashids.Encode([]int{int(category.ID)}), //nolint:all
			Title: category.Title,
		})
	}

	var discountPercent *uint64
	if product.DiscountPercent != nil {
		// Round to nearest integer as per PM consultation - frontend should not display decimals
		discountPercentUint := uint64(math.Round(*product.DiscountPercent))
		discountPercent = &discountPercentUint
	}

	return &Product{
		//nolint:all
		ID:                     hashids.Encode([]int{int(product.ID)}),
		Title:                  product.Title,
		Description:            product.Description,
		Price:                  product.Price,
		DiscountedPrice:        product.DiscountedPrice,
		VoucherDiscountedPrice: product.VoucherDiscountedPrice,
		VoucherDiscountPercent: product.VoucherDiscountPercent,
		DiscountPercent:        discountPercent,
		ImageURL:               imageURL,
		Stock:                  product.Stock,
		Categories:             categories,
		BrandTitle:             brandTitle,
	}
}

func (p *ProductTransformer) TransformMany(products []*domain.Product) []Product {
	result := make([]Product, 0, len(products))
	for _, list := range products {
		result = append(result, *p.Transform(list))
	}

	return result
}

// ProductWithSelectedResponse represents a product with a selected flag
type ProductWithSelectedResponse struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Price           uint64  `json:"price"`
	DiscountedPrice *uint64 `json:"discounted_price"`
	DiscountPercent *uint64 `json:"discount_percent"`
	ImageFilePath   *string `json:"image_file_path"`
	Selected        bool    `json:"selected"`
	Stock           uint64  `json:"stock"`
}

type ProductWithSelectedTransformer struct {
	cdn *cdnmaker.CDNMaker
}

func NewProductWithSelectedTransformer(cdn *cdnmaker.CDNMaker) *ProductWithSelectedTransformer {
	return &ProductWithSelectedTransformer{
		cdn: cdn,
	}
}

// TransformProductWithSelected transforms a product with selected flag to response
func (p *ProductWithSelectedTransformer) TransformProductWithSelected(product *domain.ProductWithSelected) *ProductWithSelectedResponse {
	if product == nil {
		return &ProductWithSelectedResponse{}
	}

	productTransformer := ProductTransformer{
		cdn: p.cdn,
	}
	productResponse := productTransformer.Transform(&product.Product)

	var discountPercent *uint64
	if product.DiscountPercent != nil {
		// Round to nearest integer as per PM consultation - frontend should not display decimals
		discountPercentUint := uint64(math.Round(*product.DiscountPercent))
		discountPercent = &discountPercentUint
	}

	var stock uint64
	if productResponse.Stock != nil {
		stock = *productResponse.Stock
	}

	return &ProductWithSelectedResponse{
		ID:              productResponse.ID,
		Title:           productResponse.Title,
		Description:     productResponse.Description,
		Price:           productResponse.Price,
		DiscountedPrice: productResponse.DiscountedPrice,
		DiscountPercent: discountPercent,
		ImageFilePath:   productResponse.ImageURL,
		Selected:        product.Selected,
		Stock:           stock,
	}
}

// TransformManyProductsWithSelected transforms multiple products with selected flag to responses
func (p *ProductWithSelectedTransformer) TransformManyProductsWithSelected(products []*domain.ProductWithSelected) []ProductWithSelectedResponse {
	if products == nil {
		return nil
	}

	response := make([]ProductWithSelectedResponse, 0, len(products))
	for _, product := range products {
		response = append(response, *p.TransformProductWithSelected(product))
	}

	return response
}

package catalog

import (
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
)

type catalogResponse struct {
	Products []productResponse `json:"products"`
	Total    int64             `json:"total"`
	Offset   int               `json:"offset"`
	Limit    int               `json:"limit"`
}

type productResponse struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category,omitempty"`
}

type productDetailResponse struct {
	Code     string            `json:"code"`
	Price    float64           `json:"price"`
	Category string            `json:"category,omitempty"`
	Variants []variantResponse `json:"variants"`
}

type variantResponse struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

func toCatalogResponse(products []models.Product, total int64, offset, limit int) catalogResponse {
	responses := make([]productResponse, len(products))
	for i, p := range products {
		responses[i] = toProductResponse(p)
	}
	return catalogResponse{
		Products: responses,
		Total:    total,
		Offset:   offset,
		Limit:    limit,
	}
}

func toProductResponse(p models.Product) productResponse {
	pr := productResponse{
		Code:  p.Code,
		Price: p.Price.InexactFloat64(),
	}
	if p.Category != nil {
		pr.Category = p.Category.Name
	}
	return pr
}

func toProductDetailResponse(p models.Product) productDetailResponse {
	variants := make([]variantResponse, len(p.Variants))
	for i, v := range p.Variants {
		price := v.Price
		if price.Equal(decimal.Zero) {
			price = p.Price
		}
		variants[i] = variantResponse{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: price.InexactFloat64(),
		}
	}

	detail := productDetailResponse{
		Code:     p.Code,
		Price:    p.Price.InexactFloat64(),
		Variants: variants,
	}
	if p.Category != nil {
		detail.Category = p.Category.Name
	}
	return detail
}

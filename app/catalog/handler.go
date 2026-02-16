package catalog

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type ProductRepository interface {
	GetProducts(filter models.ProductFilter) ([]models.Product, int64, error)
	GetProductByCode(code string) (*models.Product, error)
}

type CatalogHandler struct {
	repo ProductRepository
}

func NewCatalogHandler(r ProductRepository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

// HandleGet lists products with optional filtering and pagination.
//
//	@Summary		List products
//	@Description	Returns a paginated list of products, optionally filtered by category and max price.
//	@Tags			catalog
//	@Produce		json
//	@Param			offset			query		int		false	"Pagination offset"		default(0)
//	@Param			limit			query		int		false	"Pagination limit (1-100)"	default(10)
//	@Param			category		query		string	false	"Filter by category code"
//	@Param			price_less_than	query		number	false	"Filter by maximum price"
//	@Success		200				{object}	catalogResponse
//	@Failure		400				{object}	api.errorBody
//	@Failure		500				{object}	api.errorBody
//	@Router			/catalog [get]
func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := parsePaginationParams(r)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	filter := models.ProductFilter{
		Offset: offset,
		Limit:  limit,
	}

	if cat := r.URL.Query().Get("category"); cat != "" {
		filter.CategoryCode = cat
	}

	if priceStr := r.URL.Query().Get("price_less_than"); priceStr != "" {
		price, err := decimal.NewFromString(priceStr)
		if err != nil {
			api.ErrorResponse(w, http.StatusBadRequest, "invalid price_less_than parameter")
			return
		}
		filter.MaxPrice = &price
	}

	products, total, err := h.repo.GetProducts(filter)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, toCatalogResponse(products, total, offset, limit))
}

// HandleGetByCode returns a single product with its variants.
//
//	@Summary		Get product by code
//	@Description	Returns a single product with its variants. Variants with zero price inherit the product price.
//	@Tags			catalog
//	@Produce		json
//	@Param			code	path		string	true	"Product code"
//	@Success		200		{object}	productDetailResponse
//	@Failure		404		{object}	api.errorBody
//	@Failure		500		{object}	api.errorBody
//	@Router			/catalog/{code} [get]
func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	product, err := h.repo.GetProductByCode(code)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			api.ErrorResponse(w, http.StatusNotFound, "product not found")
			return
		}
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, toProductDetailResponse(*product))
}

func parsePaginationParams(r *http.Request) (offset, limit int, err error) {
	offset = 0
	limit = 10

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return 0, 0, errors.New("invalid offset parameter")
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			return 0, 0, errors.New("invalid limit parameter: must be between 1 and 100")
		}
	}

	return offset, limit, nil
}

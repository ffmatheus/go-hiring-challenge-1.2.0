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

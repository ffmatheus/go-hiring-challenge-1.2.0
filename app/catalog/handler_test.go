package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"

	"github.com/mytheresa/go-hiring-challenge/app/catalog/mocks"
	"github.com/mytheresa/go-hiring-challenge/models"
)

func TestHandleGet_Success(t *testing.T) {
	cat := &models.Category{Code: "shoes", Name: "Shoes"}
	products := []models.Product{
		{Code: "PROD-1", Price: decimal.NewFromFloat(99.99), Category: cat},
		{Code: "PROD-2", Price: decimal.NewFromFloat(49.50), Category: cat},
	}

	repo := mocks.NewMockProductRepository(t)
	repo.EXPECT().GetProducts(mock.AnythingOfType("models.ProductFilter")).Return(products, int64(2), nil)

	handler := NewCatalogHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
	w := httptest.NewRecorder()

	handler.HandleGet(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp catalogResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
	if len(resp.Products) != 2 {
		t.Errorf("expected 2 products, got %d", len(resp.Products))
	}
	if resp.Offset != 0 {
		t.Errorf("expected offset 0, got %d", resp.Offset)
	}
	if resp.Limit != 10 {
		t.Errorf("expected default limit 10, got %d", resp.Limit)
	}
}

func TestHandleGet_WithPagination(t *testing.T) {
	repo := mocks.NewMockProductRepository(t)
	repo.EXPECT().GetProducts(models.ProductFilter{Offset: 5, Limit: 20}).Return(nil, int64(0), nil)

	handler := NewCatalogHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/catalog?offset=5&limit=20", nil)
	w := httptest.NewRecorder()

	handler.HandleGet(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestHandleGet_WithCategoryFilter(t *testing.T) {
	expectedFilter := models.ProductFilter{Offset: 0, Limit: 10, CategoryCode: "boots"}
	repo := mocks.NewMockProductRepository(t)
	repo.EXPECT().GetProducts(expectedFilter).Return(nil, int64(0), nil)

	handler := NewCatalogHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/catalog?category=boots", nil)
	w := httptest.NewRecorder()

	handler.HandleGet(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestHandleGet_WithPriceFilter(t *testing.T) {
	price := decimal.NewFromInt(100)
	expectedFilter := models.ProductFilter{Offset: 0, Limit: 10, MaxPrice: &price}
	repo := mocks.NewMockProductRepository(t)
	repo.EXPECT().GetProducts(expectedFilter).Return(nil, int64(0), nil)

	handler := NewCatalogHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/catalog?price_less_than=100", nil)
	w := httptest.NewRecorder()

	handler.HandleGet(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestHandleGet_InvalidPriceFilter(t *testing.T) {
	repo := mocks.NewMockProductRepository(t)

	handler := NewCatalogHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/catalog?price_less_than=abc", nil)
	w := httptest.NewRecorder()

	handler.HandleGet(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandleGet_InvalidOffset(t *testing.T) {
	repo := mocks.NewMockProductRepository(t)

	handler := NewCatalogHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/catalog?offset=-1", nil)
	w := httptest.NewRecorder()

	handler.HandleGet(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandleGet_InvalidLimit(t *testing.T) {
	tests := []struct {
		name  string
		limit string
	}{
		{"zero", "0"},
		{"too large", "101"},
		{"non-numeric", "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockProductRepository(t)

			handler := NewCatalogHandler(repo)
			req := httptest.NewRequest(http.MethodGet, "/catalog?limit="+tt.limit, nil)
			w := httptest.NewRecorder()

			handler.HandleGet(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d", w.Code)
			}
		})
	}
}

func TestHandleGet_RepoError(t *testing.T) {
	repo := mocks.NewMockProductRepository(t)
	repo.EXPECT().GetProducts(mock.AnythingOfType("models.ProductFilter")).Return(nil, int64(0), errors.New("db error"))

	handler := NewCatalogHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
	w := httptest.NewRecorder()

	handler.HandleGet(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestHandleGetByCode_Success(t *testing.T) {
	cat := &models.Category{Code: "shoes", Name: "Shoes"}
	product := &models.Product{
		Code:     "PROD-1",
		Price:    decimal.NewFromFloat(99.99),
		Category: cat,
		Variants: []models.Variant{
			{Name: "Size 42", SKU: "PROD-1-42", Price: decimal.NewFromFloat(109.99)},
			{Name: "Size 43", SKU: "PROD-1-43", Price: decimal.Zero},
		},
	}

	repo := mocks.NewMockProductRepository(t)
	repo.EXPECT().GetProductByCode("PROD-1").Return(product, nil)

	handler := NewCatalogHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog/{code}", handler.HandleGetByCode)

	req := httptest.NewRequest(http.MethodGet, "/catalog/PROD-1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp productDetailResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Code != "PROD-1" {
		t.Errorf("expected code PROD-1, got %s", resp.Code)
	}
	if resp.Category != "Shoes" {
		t.Errorf("expected category Shoes, got %s", resp.Category)
	}
	if len(resp.Variants) != 2 {
		t.Fatalf("expected 2 variants, got %d", len(resp.Variants))
	}
	if resp.Variants[1].Price != 99.99 {
		t.Errorf("expected variant with zero price to inherit product price 99.99, got %f", resp.Variants[1].Price)
	}
	if resp.Variants[0].Price != 109.99 {
		t.Errorf("expected variant price 109.99, got %f", resp.Variants[0].Price)
	}
}

func TestHandleGetByCode_NotFound(t *testing.T) {
	repo := mocks.NewMockProductRepository(t)
	repo.EXPECT().GetProductByCode("UNKNOWN").Return(nil, models.ErrNotFound)

	handler := NewCatalogHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog/{code}", handler.HandleGetByCode)

	req := httptest.NewRequest(http.MethodGet, "/catalog/UNKNOWN", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestHandleGetByCode_RepoError(t *testing.T) {
	repo := mocks.NewMockProductRepository(t)
	repo.EXPECT().GetProductByCode("PROD-1").Return(nil, errors.New("db error"))

	handler := NewCatalogHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog/{code}", handler.HandleGetByCode)

	req := httptest.NewRequest(http.MethodGet, "/catalog/PROD-1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

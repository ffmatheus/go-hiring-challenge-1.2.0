package categories

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/mytheresa/go-hiring-challenge/app/categories/mocks"
	"github.com/mytheresa/go-hiring-challenge/models"
)

func TestHandleGetCategories_Success(t *testing.T) {
	cats := []models.Category{
		{ID: 1, Code: "shoes", Name: "Shoes"},
		{ID: 2, Code: "boots", Name: "Boots"},
	}

	repo := mocks.NewMockCategoryRepository(t)
	repo.EXPECT().GetAllCategories().Return(cats, nil)

	handler := NewCategoriesHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	handler.HandleGetCategories(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp categoriesResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if len(resp.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(resp.Categories))
	}
	if resp.Categories[0].Code != "shoes" {
		t.Errorf("expected code 'shoes', got %q", resp.Categories[0].Code)
	}
	if resp.Categories[1].Name != "Boots" {
		t.Errorf("expected name 'Boots', got %q", resp.Categories[1].Name)
	}
}

func TestHandleGetCategories_Empty(t *testing.T) {
	repo := mocks.NewMockCategoryRepository(t)
	repo.EXPECT().GetAllCategories().Return([]models.Category{}, nil)

	handler := NewCategoriesHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	handler.HandleGetCategories(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp categoriesResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if len(resp.Categories) != 0 {
		t.Errorf("expected 0 categories, got %d", len(resp.Categories))
	}
}

func TestHandleGetCategories_RepoError(t *testing.T) {
	repo := mocks.NewMockCategoryRepository(t)
	repo.EXPECT().GetAllCategories().Return(nil, errors.New("db error"))

	handler := NewCategoriesHandler(repo)
	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	handler.HandleGetCategories(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestHandleCreateCategory_Success(t *testing.T) {
	repo := mocks.NewMockCategoryRepository(t)
	repo.EXPECT().CreateCategory(mock.AnythingOfType("*models.Category")).
		Run(func(cat *models.Category) { cat.ID = 1 }).
		Return(nil)

	handler := NewCategoriesHandler(repo)
	body := `{"code":"sandals","name":"Sandals"}`
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreateCategory(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var resp categoryResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Code != "sandals" {
		t.Errorf("expected code 'sandals', got %q", resp.Code)
	}
	if resp.Name != "Sandals" {
		t.Errorf("expected name 'Sandals', got %q", resp.Name)
	}
}

func TestHandleCreateCategory_InvalidBody(t *testing.T) {
	repo := mocks.NewMockCategoryRepository(t)

	handler := NewCategoriesHandler(repo)
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader("not json"))
	w := httptest.NewRecorder()

	handler.HandleCreateCategory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandleCreateCategory_MissingCode(t *testing.T) {
	repo := mocks.NewMockCategoryRepository(t)

	handler := NewCategoriesHandler(repo)
	body := `{"name":"Sandals"}`
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.HandleCreateCategory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandleCreateCategory_MissingName(t *testing.T) {
	repo := mocks.NewMockCategoryRepository(t)

	handler := NewCategoriesHandler(repo)
	body := `{"code":"sandals"}`
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.HandleCreateCategory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandleCreateCategory_RepoError(t *testing.T) {
	repo := mocks.NewMockCategoryRepository(t)
	repo.EXPECT().CreateCategory(mock.AnythingOfType("*models.Category")).Return(errors.New("duplicate"))

	handler := NewCategoriesHandler(repo)
	body := `{"code":"sandals","name":"Sandals"}`
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.HandleCreateCategory(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

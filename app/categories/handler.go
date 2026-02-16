package categories

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CategoryRepository interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(category *models.Category) error
}

type CategoriesHandler struct {
	repo CategoryRepository
}

func NewCategoriesHandler(r CategoryRepository) *CategoriesHandler {
	return &CategoriesHandler{
		repo: r,
	}
}

func (h *CategoriesHandler) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.repo.GetAllCategories()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, toCategoriesResponse(cats))
}

func (h *CategoriesHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "code is required")
		return
	}

	if req.Name == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "name is required")
		return
	}

	category := &models.Category{
		Code: req.Code,
		Name: req.Name,
	}

	if err := h.repo.CreateCategory(category); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toCategoryResponse(*category))
}

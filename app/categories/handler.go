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

// HandleGetCategories returns all categories.
//
//	@Summary		List categories
//	@Description	Returns all available product categories.
//	@Tags			categories
//	@Produce		json
//	@Success		200	{object}	categoriesResponse
//	@Failure		500	{object}	api.errorBody
//	@Router			/categories [get]
func (h *CategoriesHandler) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.repo.GetAllCategories()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, toCategoriesResponse(cats))
}

// HandleCreateCategory creates a new category.
//
//	@Summary		Create category
//	@Description	Creates a new product category. Both code and name are required.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createCategoryRequest	true	"Category to create"
//	@Success		201		{object}	categoryResponse
//	@Failure		400		{object}	api.errorBody
//	@Failure		500		{object}	api.errorBody
//	@Router			/categories [post]
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

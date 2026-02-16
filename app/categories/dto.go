package categories

import "github.com/mytheresa/go-hiring-challenge/models"

type categoriesResponse struct {
	Categories []categoryResponse `json:"categories"`
}

type categoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type createCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func toCategoriesResponse(cats []models.Category) categoriesResponse {
	responses := make([]categoryResponse, len(cats))
	for i, c := range cats {
		responses[i] = toCategoryResponse(c)
	}
	return categoriesResponse{Categories: responses}
}

func toCategoryResponse(c models.Category) categoryResponse {
	return categoryResponse{
		Code: c.Code,
		Name: c.Name,
	}
}

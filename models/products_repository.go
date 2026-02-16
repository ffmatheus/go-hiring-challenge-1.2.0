package models

import (
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductFilter struct {
	CategoryCode string
	MaxPrice     *decimal.Decimal
	Offset       int
	Limit        int
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetProducts(filter ProductFilter) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := r.db.Model(&Product{})

	if filter.CategoryCode != "" {
		query = query.Where("category_id = (SELECT id FROM categories WHERE code = ?)", filter.CategoryCode)
	}

	if filter.MaxPrice != nil {
		query = query.Where("price < ?", filter.MaxPrice)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Category").Preload("Variants").
		Offset(filter.Offset).Limit(filter.Limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductsRepository) GetProductByCode(code string) (*Product, error) {
	var product Product
	err := r.db.Preload("Category").Preload("Variants").
		Where("code = ?", code).
		First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

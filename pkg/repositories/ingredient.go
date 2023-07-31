package repositories

import (
	"github.com/cvele/recipe/pkg/models"
	"github.com/jinzhu/gorm"
)

type GormIngredientRepository struct {
	db *gorm.DB
}

func NewGormIngredientRepository(db *gorm.DB) *GormIngredientRepository {
	return &GormIngredientRepository{
		db: db,
	}
}

func (r *GormIngredientRepository) FindByID(id uint) (*models.Ingredient, error) {
	var ingredient models.Ingredient
	if err := r.db.First(&ingredient, id).Error; err != nil {
		return nil, err
	}
	return &ingredient, nil
}

package repositories

import (
	"github.com/cvele/recipe/pkg/models"
	"github.com/jinzhu/gorm"
)

type GormRecipeIngredientRepository struct {
	db *gorm.DB
}

func NewGormRecipeIngredientRepository(db *gorm.DB) *GormRecipeIngredientRepository {
	return &GormRecipeIngredientRepository{
		db: db,
	}
}

func (r *GormRecipeIngredientRepository) FindByRecipeID(id uint) ([]*models.RecipeIngredient, error) {
	var recipeIngredients []*models.RecipeIngredient
	if err := r.db.Where("recipe_id = ?", id).Find(&recipeIngredients).Error; err != nil {
		return nil, err
	}
	return recipeIngredients, nil
}

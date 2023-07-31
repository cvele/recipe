package repositories

import "github.com/cvele/recipe/pkg/models"

type IngredientRepository interface {
	FindByID(id uint) (*models.Ingredient, error)
}

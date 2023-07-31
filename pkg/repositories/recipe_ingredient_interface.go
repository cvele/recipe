package repositories

import "github.com/cvele/recipe/pkg/models"

type RecipeIngredientRepository interface {
	FindByRecipeID(id uint) ([]*models.RecipeIngredient, error)
}

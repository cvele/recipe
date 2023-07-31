package repositories

import (
	"github.com/cvele/recipe/pkg/models"
)

type RecipeRepositoryInterface interface {
	GetAllRecipes() ([]models.Recipe, error)
	GetRecipeByID(id uint) (*models.Recipe, error)
	CreateRecipe(recipe *models.Recipe) error
	UpdateRecipe(recipe *models.Recipe) error
	DeleteRecipe(id uint) error
	GetRandomRecipeByType(mealType models.MealType) (*models.Recipe, error)
	GetRecipesByType(mealType models.MealType) ([]models.Recipe, error)
}

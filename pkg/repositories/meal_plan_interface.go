package repositories

import "github.com/cvele/recipe/pkg/models"

type MealPlanRepository interface {
	Create(mealPlan *models.MealPlan) error
	FindByUserID(userID uint) ([]*models.MealPlan, error)
}

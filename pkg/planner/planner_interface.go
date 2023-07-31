package planner

import (
	"time"

	"github.com/cvele/recipe/pkg/models"
)

type MealPlanner interface {
	CreateMealPlans(
		startDate time.Time,
		endDate time.Time,
		mealTime time.Time,
		mealType models.MealType,
	) ([]models.MealPlan, error)
}

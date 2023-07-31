package models

import "time"

type MealType int

const (
	Breakfast MealType = iota
	Lunch
	Dinner
	Snack
)

type MealPlanParams struct {
	StartDate       time.Time
	EndDate         time.Time
	TargetBudget    float64
	MaxBudget       float64
	TargetNutrients NutritionalValues
	MaxNutrients    NutritionalValues
	MinNutrients    NutritionalValues
	Servings        int
	MealType        MealType
}

func (m MealType) String() string {
	switch m {
	case Breakfast:
		return "breakfast"
	case Lunch:
		return "lunch"
	case Dinner:
		return "dinner"
	case Snack:
		return "snack"
	default:
		return "unknown"
	}
}

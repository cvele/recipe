package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type RecipeIngredient struct {
	gorm.Model
	RecipeID     uint       `json:"recipe_id" gorm:"not null"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID"`
	IngredientID uint       `json:"ingredient_id" gorm:"not null"`
	Quantity     float64    `json:"quantity" gorm:"not null"`
	Unit         string     `json:"unit" gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

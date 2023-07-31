package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Recipe struct {
	gorm.Model
	ID                uint                `gorm:"primary_key"`
	Title             string              `json:"title" gorm:"type:varchar(100);not null"`
	Description       string              `json:"description" gorm:"type:text;not null"`
	Servings          int                 `json:"servings" gorm:"not null"`
	PreparationTime   int                 `json:"preparation_time" gorm:"not null"`
	RecipeIngredients *[]RecipeIngredient `json:"ingredients" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsBreakfast       bool                `json:"is_breakfast" gorm:"type:bool"`
	IsLunch           bool                `json:"is_lunch" gorm:"type:bool"`
	IsDinner          bool                `json:"is_dinner" gorm:"type:bool"`
	IsSnack           bool                `json:"is_snack" gorm:"type:bool"`
	Version           int                 `gorm:"not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

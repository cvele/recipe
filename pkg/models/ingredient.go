package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Ingredient struct {
	gorm.Model
	ID           uint    `gorm:"primary_key"`
	Name         string  `json:"name" gorm:"type:varchar(100);not null"`
	PricePerUnit int     `json:"price_per_unit" gorm:"type:int;not null"`        // price per unit in cents
	Unit         float64 `json:"unit" gorm:"type:varchar(32);not null"`          // unit for price per unit for example kg or l
	UnitType     string  `json:"unit_type" gorm:"type:varchar(32);not null"`     // type of the unit for example mass or volume
	Quantity     float64 `json:"quantity" gorm:"type:decimal(10,2);not null"`    // quantity for which the nutrients are given (NutritionalValues)
	QuantityUnit string  `json:"quantity_unit" gorm:"type:varchar(32);not null"` // unit of the quantity for which the nutrients are given
	Nutrients    NutritionalValues
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

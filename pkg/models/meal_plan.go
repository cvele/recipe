package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/cvele/recipe/pkg/units"
	"github.com/jinzhu/gorm"
)

type MealPlan struct {
	gorm.Model
	ID            uint `gorm:"primary_key"`
	UserID        uint
	RecipeID      uint
	Recipe        *Recipe `gorm:"foreignKey:RecipeID"`
	RecipeVersion int     `gorm:"not null"`
	Servings      int     `gorm:"not null"`
	Synced        bool
	MealTime      time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
}

func (mp *MealPlan) AdjustServings(converter units.UnitConverterInterface) error {
	// Defensive check: Ensure that converter is not nil
	if converter == nil {
		return errors.New("nil unit converter")
	}

	// Defensive check: Ensure that mp.Recipe is not nil
	if mp.Recipe == nil {
		return errors.New("nil recipe in meal plan")
	}

	// Defensive check: Ensure that mp.RecipeIngredients slice is not nil
	if mp.Recipe.RecipeIngredients == nil {
		return errors.New("nil recipe ingredients slice in meal plan")
	}

	if mp.Servings <= 0 {
		return errors.New("invalid number of servings in meal plan")
	}

	// Defensive check: Ensure that mp.Recipe.Servings is non-zero
	if mp.Recipe.Servings == 0 {
		return errors.New("invalid recipe servings (zero) in meal plan")
	}

	servingsRatio := float64(mp.Servings) / float64(mp.Recipe.Servings)

	for i, recipeIngredient := range *mp.Recipe.RecipeIngredients {
		ingredient := recipeIngredient.Ingredient
		defaultUnit := ""
		if ingredient.UnitType == "mass" {
			defaultUnit = converter.GetDefaultUnit("mass")
		} else if ingredient.UnitType == "volume" {
			defaultUnit = converter.GetDefaultUnit("volume")
		} else {
			return fmt.Errorf("unsupported unit type")
		}

		// Convert the RecipeIngredient's unit to the default unit
		convertedQuantity, err := converter.ConvertUnits(recipeIngredient.Quantity, recipeIngredient.Unit, defaultUnit, ingredient.UnitType)
		if err != nil {
			return err
		}

		// Adjust the quantity according to the servings ratio
		adjustedQuantity := servingsRatio * convertedQuantity

		(*mp.Recipe.RecipeIngredients)[i].Quantity = adjustedQuantity
		(*mp.Recipe.RecipeIngredients)[i].Unit = defaultUnit
	}

	mp.Recipe.Servings = mp.Servings

	return nil
}

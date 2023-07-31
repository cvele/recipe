package shoppinglist

import (
	"fmt"

	"github.com/cvele/recipe/pkg/models"
	"github.com/cvele/recipe/pkg/units"
)

var _ ShoppingListServiceInterface = &ShoppingListService{}

type ShoppingListService struct {
	UnitConverter units.UnitConverterInterface
	MealPlans     []models.MealPlan
}

func (s *ShoppingListService) GenerateShoppingList() ([]models.ShoppingItem, *int, error) {
	shoppingListMap := make(map[string]models.ShoppingItem)
	totalCost := 0

	for _, mealPlan := range s.MealPlans {
		recipe := mealPlan.Recipe

		if mealPlan.Servings == 0 || recipe.Servings == 0 {
			return nil, nil, fmt.Errorf("servings in meal plan or recipe cannot be 0")
		}

		servingsRatio := float64(mealPlan.Servings) / float64(recipe.Servings)

		for _, recipeIngredient := range *recipe.RecipeIngredients {
			ingredient := recipeIngredient.Ingredient
			defaultUnit := ""

			if ingredient.UnitType == "mass" {
				defaultUnit = s.UnitConverter.GetDefaultUnit("mass")
			} else if ingredient.UnitType == "volume" {
				defaultUnit = s.UnitConverter.GetDefaultUnit("volume")
			} else {
				return nil, nil, fmt.Errorf("unsupported unit type")
			}

			// Convert the RecipeIngredient's unit to the default unit
			convertedQuantity, err := s.UnitConverter.ConvertUnits(recipeIngredient.Quantity, recipeIngredient.Unit, defaultUnit, ingredient.UnitType)
			if err != nil {
				return nil, nil, err
			}

			// Adjust the quantity according to the servings ratio
			adjustedQuantity := servingsRatio * convertedQuantity

			// Calculate the cost
			cost := int(adjustedQuantity * float64(ingredient.PricePerUnit)) // calculate cost in cents
			totalCost += cost

			if item, exists := shoppingListMap[ingredient.Name]; exists {
				item.Quantity += adjustedQuantity
				item.Cost += cost
				shoppingListMap[ingredient.Name] = item
			} else {
				shoppingListMap[ingredient.Name] = models.ShoppingItem{
					Name:     ingredient.Name,
					Quantity: adjustedQuantity,
					Unit:     defaultUnit,
					Cost:     cost,
				}
			}
		}
	}

	shoppingList := make([]models.ShoppingItem, 0, len(shoppingListMap))
	for _, item := range shoppingListMap {
		shoppingList = append(shoppingList, item)
	}

	return shoppingList, &totalCost, nil
}

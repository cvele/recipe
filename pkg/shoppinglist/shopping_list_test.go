package shoppinglist_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/cvele/recipe/pkg/models"
	"github.com/cvele/recipe/pkg/shoppinglist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUnitConverter struct {
	mock.Mock
}

func (m *MockUnitConverter) ConvertUnits(quantity float64, fromUnit string, toUnit string, unitType string) (float64, error) {
	args := m.Called(quantity, fromUnit, toUnit, unitType)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockUnitConverter) GetDefaultUnit(unitType string) string {
	args := m.Called(unitType)
	if args.Get(0) == nil {
		return ""
	}
	return args.Get(0).(string)
}

func (m *MockUnitConverter) GetAvailableUnits(unitType string) []string {
	args := m.Called(unitType)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).([]string)
}

func (m *MockUnitConverter) IsValidUnit(unit string, unitType string) bool {
	args := m.Called(unit, unitType)

	if args.Get(0) == nil {
		return false
	}

	return args.Get(0).(bool)
}

func TestGenerateShoppingList(t *testing.T) {
	mockUnitConverter := new(MockUnitConverter)
	mockUnitConverter.On("GetDefaultUnit", "mass").Return("g")
	mockUnitConverter.On("GetDefaultUnit", "volume").Return("ml")
	mockUnitConverter.On("IsValidUnit", "kg", "mass").Return(true)
	mockUnitConverter.On("ConvertUnits", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(100.0, nil)

	svc := &shoppinglist.ShoppingListService{
		UnitConverter: mockUnitConverter,
		MealPlans: []models.MealPlan{
			{
				Servings: 4,
				Recipe: &models.Recipe{
					Servings: 2,
					RecipeIngredients: &[]models.RecipeIngredient{
						{
							Ingredient: models.Ingredient{
								Name:         "Test Ingredient 1",
								UnitType:     "mass",
								PricePerUnit: 10,
							},
							Quantity: 2,
							Unit:     "kg",
						},
					},
				},
			},
		},
	}

	expectedShoppingItem := models.ShoppingItem{
		Name:     "Test Ingredient 1",
		Quantity: 200.0,
		Unit:     "g",
		Cost:     2000,
	}

	list, totalCost, err := svc.GenerateShoppingList()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 2000, *totalCost)
	assert.Contains(t, list, expectedShoppingItem)

	svc.MealPlans[0].Servings = 0
	list, totalCost, err = svc.GenerateShoppingList()

	assert.NotNil(t, err)
	assert.Nil(t, list)
	assert.Nil(t, totalCost)
	assert.Equal(t, "servings in meal plan or recipe cannot be 0", err.Error())

	// Test scenario where servings in recipe is 0
	svc.MealPlans[0].Servings = 4
	svc.MealPlans[0].Recipe.Servings = 0
	list, totalCost, err = svc.GenerateShoppingList()

	assert.NotNil(t, err)
	assert.Nil(t, list)
	assert.Nil(t, totalCost)
	assert.Equal(t, "servings in meal plan or recipe cannot be 0", err.Error())

	// Test scenario with unsupported unit type
	svc.MealPlans[0].Recipe.Servings = 2
	(*svc.MealPlans[0].Recipe.RecipeIngredients)[0].Ingredient.UnitType = "unsupported"
	list, totalCost, err = svc.GenerateShoppingList()

	assert.NotNil(t, err)
	assert.Nil(t, list)
	assert.Nil(t, totalCost)
	assert.Equal(t, "unsupported unit type", err.Error())
}

func TestGenerateShoppingList_ServingsScaling(t *testing.T) {
	mockUnitConverter := new(MockUnitConverter)
	mockUnitConverter.On("GetDefaultUnit", "mass").Return("g")
	mockUnitConverter.On("ConvertUnits", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(100.0, nil)

	// Test scenario with valid meal plans
	svc := &shoppinglist.ShoppingListService{
		UnitConverter: mockUnitConverter,
		MealPlans: []models.MealPlan{
			{
				Servings: 2,
				Recipe: &models.Recipe{
					Servings: 1,
					RecipeIngredients: &[]models.RecipeIngredient{
						{
							Ingredient: models.Ingredient{
								Name:         "Test Ingredient 1",
								UnitType:     "mass",
								PricePerUnit: 10,
							},
							Quantity: 1,
							Unit:     "kg",
						},
					},
				},
			},
		},
	}

	expectedShoppingItem := models.ShoppingItem{
		Name:     "Test Ingredient 1",
		Quantity: 200.0,
		Unit:     "g",
		Cost:     2000,
	}

	list, totalCost, err := svc.GenerateShoppingList()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 2000, *totalCost)
	assert.Contains(t, list, expectedShoppingItem)

	// Scaling up the servings in the meal plan
	svc.MealPlans[0].Servings = 4

	expectedShoppingItem.Quantity = 400.0
	expectedShoppingItem.Cost = 4000

	list, totalCost, err = svc.GenerateShoppingList()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 4000, *totalCost)
	assert.Contains(t, list, expectedShoppingItem)

	// Scaling down the servings in the meal plan
	svc.MealPlans[0].Servings = 1

	expectedShoppingItem.Quantity = 100.0
	expectedShoppingItem.Cost = 1000

	list, totalCost, err = svc.GenerateShoppingList()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1000, *totalCost)
	assert.Contains(t, list, expectedShoppingItem)
}

func TestGenerateShoppingList_IngredientUnitConversion(t *testing.T) {
	mockUnitConverter := new(MockUnitConverter)
	mockUnitConverter.On("GetDefaultUnit", "mass").Return("g")
	mockUnitConverter.On("GetDefaultUnit", "volume").Return("ml")
	mockUnitConverter.On("ConvertUnits", 1.0, "kg", "g", "mass").Return(1000.0, nil)
	mockUnitConverter.On("ConvertUnits", 1.0, "L", "ml", "volume").Return(1000.0, nil)

	svc := &shoppinglist.ShoppingListService{
		UnitConverter: mockUnitConverter,
		MealPlans: []models.MealPlan{
			{
				Servings: 1,
				Recipe: &models.Recipe{
					Servings: 1,
					RecipeIngredients: &[]models.RecipeIngredient{
						{
							Ingredient: models.Ingredient{
								Name:         "Test Ingredient 1",
								UnitType:     "mass",
								PricePerUnit: 10,
							},
							Quantity: 1,
							Unit:     "kg",
						},
						{
							Ingredient: models.Ingredient{
								Name:         "Test Ingredient 2",
								UnitType:     "volume",
								PricePerUnit: 5,
							},
							Quantity: 1,
							Unit:     "L",
						},
					},
				},
			},
		},
	}

	expectedShoppingItems := []models.ShoppingItem{
		{
			Name:     "Test Ingredient 1",
			Quantity: 1000.0,
			Unit:     "g",
			Cost:     10000,
		},
		{
			Name:     "Test Ingredient 2",
			Quantity: 1000.0,
			Unit:     "ml",
			Cost:     5000,
		},
	}

	list, totalCost, err := svc.GenerateShoppingList()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 15000, *totalCost)
	for _, expectedItem := range expectedShoppingItems {
		assert.Contains(t, list, expectedItem)
	}
}

func TestGenerateShoppingList_CostCalculation(t *testing.T) {
	mockUnitConverter := new(MockUnitConverter)
	mockUnitConverter.On("GetDefaultUnit", "mass").Return("g")
	mockUnitConverter.On("ConvertUnits", 1.0, "kg", "g", "mass").Return(1000.0, nil)

	// Test scenario with valid meal plans
	svc := &shoppinglist.ShoppingListService{
		UnitConverter: mockUnitConverter,
		MealPlans: []models.MealPlan{
			{
				Servings: 1,
				Recipe: &models.Recipe{
					Servings: 1,
					RecipeIngredients: &[]models.RecipeIngredient{
						{
							Ingredient: models.Ingredient{
								Name:         "Test Ingredient 1",
								UnitType:     "mass",
								PricePerUnit: 10,
							},
							Quantity: 1,
							Unit:     "kg",
						},
						{
							Ingredient: models.Ingredient{
								Name:         "Test Ingredient 2",
								UnitType:     "mass",
								PricePerUnit: 20,
							},
							Quantity: 1,
							Unit:     "kg",
						},
					},
				},
			},
		},
	}

	expectedShoppingItems := []models.ShoppingItem{
		{
			Name:     "Test Ingredient 1",
			Quantity: 1000.0,
			Unit:     "g",
			Cost:     10000,
		},
		{
			Name:     "Test Ingredient 2",
			Quantity: 1000.0,
			Unit:     "g",
			Cost:     20000,
		},
	}

	list, totalCost, err := svc.GenerateShoppingList()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 30000, *totalCost)
	for _, expectedItem := range expectedShoppingItems {
		assert.Contains(t, list, expectedItem)
	}
}

func TestGenerateShoppingList_ZeroServings(t *testing.T) {
	mockUnitConverter := new(MockUnitConverter)

	// Test scenario with zero servings in meal plan
	svc := &shoppinglist.ShoppingListService{
		UnitConverter: mockUnitConverter,
		MealPlans: []models.MealPlan{
			{
				Servings: 0,
				Recipe: &models.Recipe{
					Servings: 1,
					RecipeIngredients: &[]models.RecipeIngredient{
						{
							Ingredient: models.Ingredient{
								Name:         "Test Ingredient 1",
								UnitType:     "mass",
								PricePerUnit: 10,
							},
							Quantity: 1,
							Unit:     "kg",
						},
					},
				},
			},
		},
	}

	list, totalCost, err := svc.GenerateShoppingList()

	assert.Nil(t, list)
	assert.Nil(t, totalCost)
	assert.Error(t, err)

	// Test scenario with zero servings in recipe
	svc.MealPlans[0].Servings = 1
	svc.MealPlans[0].Recipe.Servings = 0

	list, totalCost, err = svc.GenerateShoppingList()

	assert.Nil(t, list)
	assert.Nil(t, totalCost)
	assert.Error(t, err)
}

func TestGenerateShoppingList_InvalidUnits(t *testing.T) {
	mockUnitConverter := new(MockUnitConverter)
	mockUnitConverter.On("GetDefaultUnit", "mass").Return("g")

	// Test scenario with unsupported unit type in meal plan
	svc := &shoppinglist.ShoppingListService{
		UnitConverter: mockUnitConverter,
		MealPlans: []models.MealPlan{
			{
				Servings: 1,
				Recipe: &models.Recipe{
					Servings: 1,
					RecipeIngredients: &[]models.RecipeIngredient{
						{
							Ingredient: models.Ingredient{
								Name:         "Test Ingredient 1",
								UnitType:     "unsupported_unit_type",
								PricePerUnit: 10,
							},
							Quantity: 1,
							Unit:     "kg",
						},
					},
				},
			},
		},
	}

	list, totalCost, err := svc.GenerateShoppingList()

	assert.Nil(t, list)
	assert.Nil(t, totalCost)
	assert.Error(t, err)
}

func TestGenerateShoppingList_EmptyMealPlan(t *testing.T) {
	mockUnitConverter := new(MockUnitConverter)

	svc := &shoppinglist.ShoppingListService{
		UnitConverter: mockUnitConverter,
		MealPlans:     []models.MealPlan{},
	}

	list, totalCost, err := svc.GenerateShoppingList()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 0, len(list))
	assert.NotNil(t, totalCost)
	assert.Equal(t, 0, *totalCost)
}

func TestGenerateShoppingList_UnitConversionError(t *testing.T) {
	mockUnitConverter := new(MockUnitConverter)
	mockUnitConverter.On("GetDefaultUnit", "mass").Return("g")
	mockUnitConverter.On("ConvertUnits", 1.0, "kg", "g", "mass").Return(0.0, fmt.Errorf("unit conversion error"))

	svc := &shoppinglist.ShoppingListService{
		UnitConverter: mockUnitConverter,
		MealPlans: []models.MealPlan{
			{
				Servings: 1,
				Recipe: &models.Recipe{
					Servings: 1,
					RecipeIngredients: &[]models.RecipeIngredient{
						{
							Ingredient: models.Ingredient{
								Name:         "Test Ingredient 1",
								UnitType:     "mass",
								PricePerUnit: 10,
							},
							Quantity: 1,
							Unit:     "kg",
						},
					},
				},
			},
		},
	}

	list, totalCost, err := svc.GenerateShoppingList()

	assert.Nil(t, list)
	assert.Nil(t, totalCost)
	assert.Error(t, err)
	assert.Equal(t, "unit conversion error", err.Error())
}

func TestGenerateShoppingList_BoundaryCases(t *testing.T) {
	mockUnitConverter := new(MockUnitConverter)
	mockUnitConverter.On("GetDefaultUnit", "mass").Return("g")
	mockUnitConverter.On("ConvertUnits", math.MaxFloat64, "kg", "g", "mass").Return(math.MaxFloat64, nil)

	svc := &shoppinglist.ShoppingListService{
		UnitConverter: mockUnitConverter,
		MealPlans: []models.MealPlan{
			{
				Servings: 1,
				Recipe: &models.Recipe{
					Servings: 1,
					RecipeIngredients: &[]models.RecipeIngredient{
						{
							Ingredient: models.Ingredient{
								Name:         "Test Ingredient 1",
								UnitType:     "mass",
								PricePerUnit: math.MaxInt,
							},
							Quantity: math.MaxFloat64,
							Unit:     "kg",
						},
					},
				},
			},
		},
	}

	list, totalCost, err := svc.GenerateShoppingList()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.NotNil(t, totalCost)
	assert.Contains(t, list, models.ShoppingItem{
		Name:     "Test Ingredient 1",
		Quantity: math.MaxFloat64,
		Unit:     "g",
		Cost:     math.MaxInt64,
	})
}

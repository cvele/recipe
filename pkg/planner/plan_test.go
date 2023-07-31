package planner_test

import (
	"testing"
	"time"

	"github.com/cvele/recipe/pkg/models"
	"github.com/cvele/recipe/pkg/planner"
	"github.com/cvele/recipe/pkg/units"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RecipeRepositoryMock struct {
	mock.Mock
}

func (m *RecipeRepositoryMock) GetAllRecipes() ([]models.Recipe, error) {
	args := m.Called()
	return args.Get(0).([]models.Recipe), args.Error(1)
}

func (m *RecipeRepositoryMock) GetRecipeByID(id uint) (*models.Recipe, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Recipe), args.Error(1)
}

func (m *RecipeRepositoryMock) CreateRecipe(recipe *models.Recipe) error {
	args := m.Called(recipe)
	return args.Error(0)
}

func (m *RecipeRepositoryMock) UpdateRecipe(recipe *models.Recipe) error {
	args := m.Called(recipe)
	return args.Error(0)
}

func (m *RecipeRepositoryMock) DeleteRecipe(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RecipeRepositoryMock) GetRandomRecipeByType(mealType models.MealType) (*models.Recipe, error) {
	args := m.Called(mealType)
	return args.Get(0).(*models.Recipe), args.Error(1)
}

func (m *RecipeRepositoryMock) GetRecipesByType(mealType models.MealType) ([]models.Recipe, error) {
	args := m.Called(mealType)
	return args.Get(0).([]models.Recipe), args.Error(1)
}

func TestNewGeneticMealPlanner(t *testing.T) {
	mockRepo := new(RecipeRepositoryMock)
	params := models.MealPlanParams{}

	unitConverter := units.NewUnitConverter("kg", "l")
	gmp := planner.NewGeneticMealPlanner(100, 50, 0.7, 0.1, mockRepo, params, unitConverter)

	assert.Equal(t, 100, gmp.PopulationSize())
	assert.Equal(t, 50, gmp.MaxGenerations())
	assert.Equal(t, 0.7, gmp.CrossoverRate())
	assert.Equal(t, 0.1, gmp.MutationRate())
	assert.Equal(t, mockRepo, gmp.RecipeRepo())
	assert.Equal(t, params, gmp.Params())
}

func TestCreateMealPlans(t *testing.T) {
	mockRepo := new(RecipeRepositoryMock)
	params := models.MealPlanParams{}

	mockRepo.On("GetRecipesByType", models.Breakfast).Return([]models.Recipe{
		{
			ID:          1,
			IsBreakfast: true,
			RecipeIngredients: &[]models.RecipeIngredient{
				{
					Ingredient: models.Ingredient{
						ID:           1,
						Name:         "Egg",
						Nutrients:    models.NutritionalValues{},
						PricePerUnit: 1,
					},
					Quantity: 1,
					Unit:     "unit",
				},
			},
		},
	}, nil)

	unitConverter := units.NewUnitConverter("kg", "l")
	gmp := planner.NewGeneticMealPlanner(100, 50, 0.7, 0.1, mockRepo, params, unitConverter)

	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 7) // One week later
	mealTime := time.Now()
	mealType := models.Breakfast

	mealPlans, err := gmp.CreateMealPlans(startDate, endDate, mealTime, mealType)

	assert.NoError(t, err)
	assert.Len(t, mealPlans, 7)

	for _, mealPlan := range mealPlans {
		assert.NotNil(t, mealPlan.Recipe)
		assert.True(t, mealPlan.Recipe.IsBreakfast)
		assert.Equal(t, params.Servings, mealPlan.Servings)
	}
}

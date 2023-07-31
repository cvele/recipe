package repositories

import (
	"fmt"

	"github.com/cvele/recipe/pkg/models"
	"github.com/jinzhu/gorm"
)

type GormRecipeRepository struct {
	db *gorm.DB
}

func NewGormRecipeRepository(db *gorm.DB) *GormRecipeRepository {
	return &GormRecipeRepository{
		db: db,
	}
}

func (r *GormRecipeRepository) GetAllRecipes() ([]models.Recipe, error) {
	var recipes []models.Recipe
	subquery := r.db.Table("recipes").
		Select("id, MAX(version) as max_version").
		Group("id").SubQuery()

	if err := r.db.Table("recipes").
		Joins("JOIN (?) AS latest ON recipes.id = latest.id AND recipes.version = latest.max_version", subquery).
		Find(&recipes).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}

func (r *GormRecipeRepository) CreateRecipe(recipe *models.Recipe) error {
	if err := r.db.Create(recipe).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormRecipeRepository) UpdateRecipe(recipe *models.Recipe) error {
	recipe.Version++
	if err := r.db.Create(recipe).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormRecipeRepository) DeleteRecipe(id uint) error {
	if err := r.db.Where("id = ?", id).Delete(&models.Recipe{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormRecipeRepository) GetRandomRecipeByType(mealType models.MealType) (*models.Recipe, error) {
	var recipe models.Recipe
	if err := r.db.Where("meal_type = ?", mealType).Order("version DESC").Order(gorm.Expr("rand()")).Limit(1).Find(&recipe).Error; err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *GormRecipeRepository) GetRecipesByType(mealType models.MealType) ([]models.Recipe, error) {
	var recipes []models.Recipe
	subquery := r.db.Table("recipes").
		Select("id, MAX(version) as max_version").
		//mealType is safe
		Where(fmt.Sprintf("is_%s = ?", mealType.String()), true).
		Group("id").SubQuery()

	query := r.db.Table("recipes").
		Joins("JOIN (?) AS latest ON recipes.id = latest.id AND recipes.version = latest.max_version", subquery).
		Where(fmt.Sprintf("is_%s = ?", mealType.String()), true).
		Order("version DESC")

	if err := query.Find(&recipes).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}

func (r *GormRecipeRepository) GetRecipeByID(id uint) (*models.Recipe, error) {
	var recipe models.Recipe
	if err := r.db.Where("id = ?", id).Limit(1).Find(&recipe).Error; err != nil {
		return nil, err
	}
	return &recipe, nil
}

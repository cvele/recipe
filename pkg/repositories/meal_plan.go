package repositories

import (
	"github.com/cvele/recipe/pkg/models"
	"github.com/jinzhu/gorm"
)

type GormMealPlanRepository struct {
	db *gorm.DB
}

func NewGormMealPlanRepository(db *gorm.DB) *GormMealPlanRepository {
	return &GormMealPlanRepository{
		db: db,
	}
}

func (r *GormMealPlanRepository) Create(mealPlan *models.MealPlan) error {
	if err := r.db.Create(mealPlan).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormMealPlanRepository) FindByUserID(userID uint) ([]*models.MealPlan, error) {
	var mealPlans []*models.MealPlan
	if err := r.db.Where("user_id = ?", userID).Find(&mealPlans).Error; err != nil {
		return nil, err
	}
	return mealPlans, nil
}

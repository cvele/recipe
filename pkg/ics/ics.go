package ics

import (
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/cvele/recipe/pkg/models"
)

func GenerateMealPlanICal(mealPlans []models.MealPlan) (string, error) {
	cal := ics.NewCalendar()

	// Add an event for each meal in the meal plan
	for _, m := range mealPlans {
		event := cal.AddEvent(m.Recipe.Title)

		event.SetCreatedTime(m.CreatedAt)
		event.SetDtStampTime(m.CreatedAt)
		event.SetModifiedAt(m.UpdatedAt)

		event.SetStartAt(m.MealTime)
		event.SetEndAt(m.MealTime.Add(time.Hour)) // @TODO: for now assuming each meal is an hour long

		event.SetSummary("Meal Plan for Recipe " + fmt.Sprint(m.RecipeID))

		// You can also add a Description with more details about the meal plan
		event.SetDescription(m.Recipe.Description)
	}

	return cal.Serialize(), nil
}

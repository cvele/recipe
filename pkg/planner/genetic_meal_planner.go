package planner

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/cvele/recipe/pkg/models"
	"github.com/cvele/recipe/pkg/repositories"
	"github.com/cvele/recipe/pkg/units"
)

var _ MealPlanner = (*GeneticMealPlanner)(nil)

type GeneticMealPlanner struct {
	populationSize    int
	maxGenerations    int
	crossoverRate     float64
	mutationRate      float64
	recipeRepo        repositories.RecipeRepositoryInterface
	params            models.MealPlanParams
	population        []models.MealPlan
	fitnessValues     []float64
	bestMealPlan      models.MealPlan
	bestFitness       float64
	currentGeneration int
	unitConverter     units.UnitConverterInterface
}

func NewGeneticMealPlanner(
	populationSize int,
	maxGenerations int,
	crossoverRate float64,
	mutationRate float64,
	recipeRepo repositories.RecipeRepositoryInterface,
	params models.MealPlanParams,
	unitConverter units.UnitConverterInterface,
) *GeneticMealPlanner {
	return &GeneticMealPlanner{
		populationSize: populationSize,
		maxGenerations: maxGenerations,
		crossoverRate:  crossoverRate,
		mutationRate:   mutationRate,
		recipeRepo:     recipeRepo,
		params:         params,
		unitConverter:  unitConverter,
	}
}

func (g *GeneticMealPlanner) CreateMealPlans(
	startDate time.Time,
	endDate time.Time,
	mealTime time.Time,
	mealType models.MealType,
) ([]models.MealPlan, error) {
	var mealPlans []models.MealPlan
	duration := endDate.Sub(startDate)
	days := int(duration.Round(time.Hour*24).Hours() / 24)
	for i := 0; i < days; i++ {
		err := g.initializePopulation(mealType)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize population: %v", err)
		}

		for g.currentGeneration = 0; g.currentGeneration < g.maxGenerations; g.currentGeneration++ {
			g.calculateFitness()
			g.updateBest()
			if g.terminate() {
				break
			}

			g.selectNewPopulation()
			g.crossover()
		}

		mealPlans = append(mealPlans, g.bestMealPlan)
	}

	return mealPlans, nil
}

func (g *GeneticMealPlanner) initializePopulation(mealType models.MealType) error {
	recipes, err := g.recipeRepo.GetRecipesByType(mealType)
	if err != nil {
		return err
	}
	if len(recipes) == 0 {
		return errors.New("no recipes available for this meal type")
	}

	populationMap := sync.Map{}

	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.NumCPU()) // Semaphore

	for i := 0; i < g.populationSize; i++ {
		wg.Add(1)
		go func(i int) {
			sem <- struct{}{} // Acquire a token
			defer wg.Done()
			defer func() { <-sem }() // Release the token

			randomIndex := rand.Intn(len(recipes))
			recipe := recipes[randomIndex]

			mealPlan := models.MealPlan{
				Recipe:   &recipe,
				RecipeID: recipe.ID,
				Servings: g.params.Servings,
			}

			// Adjust servings
			mealPlan.AdjustServings(g.unitConverter)

			populationMap.Store(i, mealPlan)
		}(i)
	}
	wg.Wait()

	g.population = make([]models.MealPlan, g.populationSize)
	populationMap.Range(func(key, value interface{}) bool {
		index := key.(int)
		mealPlan := value.(models.MealPlan)
		g.population[index] = mealPlan
		return true
	})

	return nil
}

func (g *GeneticMealPlanner) calculateFitness() {
	fitnessValuesMap := sync.Map{}

	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.NumCPU()) // Semaphore

	for i, recipe := range g.population {
		wg.Add(1)
		go func(i int, recipe models.MealPlan) {
			sem <- struct{}{} // Acquire a token
			defer wg.Done()
			defer func() { <-sem }() // Release the token

			totalNutrition := models.NutritionalValues{}
			totalCost := 0.0
			scalingFactor := float64(g.params.Servings) / float64(recipe.Servings)

			for _, ingredient := range *recipe.Recipe.RecipeIngredients {
				scaledQuantity := ingredient.Quantity * scalingFactor

				totalNutrition.Calories += ingredient.Ingredient.Nutrients.Calories * scaledQuantity
				totalNutrition.Protein += ingredient.Ingredient.Nutrients.Protein * scaledQuantity
				totalNutrition.Fat += ingredient.Ingredient.Nutrients.Fat * scaledQuantity
				totalNutrition.Carbs += ingredient.Ingredient.Nutrients.Carbs * scaledQuantity
				totalNutrition.Fiber += ingredient.Ingredient.Nutrients.Fiber * scaledQuantity
				totalNutrition.Sugar += ingredient.Ingredient.Nutrients.Sugar * scaledQuantity

				totalCost += float64(ingredient.Ingredient.PricePerUnit) * scaledQuantity
			}

			fitness := 0.0
			minTarget := g.params.MinNutrients
			maxTarget := g.params.MaxNutrients
			target := g.params.TargetNutrients

			fitness += calculateNutrientFitness(totalNutrition.Calories, target.Calories, minTarget.Calories, maxTarget.Calories)
			fitness += calculateNutrientFitness(totalNutrition.Protein, target.Protein, minTarget.Protein, maxTarget.Protein)
			fitness += calculateNutrientFitness(totalNutrition.Fat, target.Fat, minTarget.Fat, maxTarget.Fat)
			fitness += calculateNutrientFitness(totalNutrition.Carbs, target.Carbs, minTarget.Carbs, maxTarget.Carbs)
			fitness += calculateNutrientFitness(totalNutrition.Fiber, target.Fiber, minTarget.Fiber, maxTarget.Fiber)
			fitness += calculateNutrientFitness(totalNutrition.Sugar, target.Sugar, minTarget.Sugar, maxTarget.Sugar)

			fitness -= totalCost // Consider reducing cost as improving fitness

			fitnessValuesMap.Store(i, fitness)
		}(i, recipe)
	}
	wg.Wait()

	g.fitnessValues = make([]float64, g.populationSize)
	fitnessValuesMap.Range(func(key, value interface{}) bool {
		index := key.(int)
		fitness := value.(float64)
		g.fitnessValues[index] = fitness
		return true
	})
}

func calculateNutrientFitness(value float64, target float64, minTarget float64, maxTarget float64) float64 {
	if value < minTarget {
		return (minTarget - value) * 2 // Penalize more heavily for falling below minimum
	} else if value > maxTarget {
		return (value - maxTarget) * 2 // Penalize more heavily for exceeding maximum
	} else {
		return math.Abs(target - value) // Encourage matching target
	}
}

func (g *GeneticMealPlanner) updateBest() {
	bestFitnessIndex := 0
	for i, fitness := range g.fitnessValues {
		if fitness < g.fitnessValues[bestFitnessIndex] {
			bestFitnessIndex = i
		}
	}

	if g.bestFitness > g.fitnessValues[bestFitnessIndex] || g.currentGeneration == 0 {
		g.bestFitness = g.fitnessValues[bestFitnessIndex]
		g.bestMealPlan = g.population[bestFitnessIndex]
	}
}

func (g *GeneticMealPlanner) terminate() bool {
	const improvementThreshold = 0.01
	const stagnantGenerations = 10

	if g.currentGeneration >= g.maxGenerations {
		return true
	}

	// If the best fitness is stagnant for a certain number of generations
	if g.currentGeneration > stagnantGenerations {
		improved := false
		for i := 0; i < stagnantGenerations; i++ {
			if math.Abs(g.fitnessValues[g.currentGeneration-i]-g.fitnessValues[g.currentGeneration-i-1]) > improvementThreshold {
				improved = true
				break
			}
		}
		if !improved {
			return true
		}
	}

	return false
}

func (g *GeneticMealPlanner) selectNewPopulation() {
	newPopulation := make([]models.MealPlan, g.populationSize)

	for i := 0; i < g.populationSize; i++ {
		index1, index2 := rand.Intn(g.populationSize), rand.Intn(g.populationSize)
		// Ensure index1 and index2 are different
		for index1 == index2 {
			index2 = rand.Intn(g.populationSize)
		}

		competitor1 := g.population[index1]
		competitor2 := g.population[index2]
		fitness1 := g.fitnessValues[index1]
		fitness2 := g.fitnessValues[index2]

		if fitness1 < fitness2 {
			newPopulation[i] = competitor1
		} else {
			newPopulation[i] = competitor2
		}
	}

	g.population = newPopulation
}

func (g *GeneticMealPlanner) crossover() {
	crossoverLimit := g.populationSize
	if g.populationSize%2 != 0 {
		crossoverLimit = g.populationSize - 1
	}

	populationMap := sync.Map{}

	var wg sync.WaitGroup
	sem := make(chan struct{}, runtime.NumCPU()) // Semaphore

	for i := 0; i < crossoverLimit; i += 2 {
		wg.Add(1)
		go func(i int) {
			sem <- struct{}{} // Acquire a token
			defer wg.Done()
			defer func() { <-sem }() // Release the token

			if rand.Float64() < g.crossoverRate {
				parent1 := g.population[rand.Intn(g.populationSize)]
				parent2 := g.population[rand.Intn(g.populationSize)]

				if rand.Intn(2) == 0 {
					populationMap.Store(i, parent2)
					populationMap.Store(i+1, parent1)
				}
			}
		}(i)
	}
	wg.Wait()

	populationMap.Range(func(key, value interface{}) bool {
		index := key.(int)
		mealPlan := value.(models.MealPlan)
		g.population[index] = mealPlan
		return true
	})
}

func (g *GeneticMealPlanner) PopulationSize() int {
	return g.populationSize
}

func (g *GeneticMealPlanner) MaxGenerations() int {
	return g.maxGenerations
}

func (g *GeneticMealPlanner) CrossoverRate() float64 {
	return g.crossoverRate
}

func (g *GeneticMealPlanner) MutationRate() float64 {
	return g.mutationRate
}

func (g *GeneticMealPlanner) RecipeRepo() repositories.RecipeRepositoryInterface {
	return g.recipeRepo
}

func (g *GeneticMealPlanner) Params() models.MealPlanParams {
	return g.params
}

func (g *GeneticMealPlanner) Population() []models.MealPlan {
	return g.population
}

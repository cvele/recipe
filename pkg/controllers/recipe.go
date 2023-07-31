package controllers

import (
	"net/http"
	"strconv"

	"github.com/cvele/recipe/pkg/models"
	"github.com/cvele/recipe/pkg/repositories"
	"github.com/gin-gonic/gin"
)

type RecipeController struct {
	repo repositories.RecipeRepositoryInterface
}

func NewRecipeController(repo repositories.RecipeRepositoryInterface) *RecipeController {
	return &RecipeController{repo: repo}
}

func (rc *RecipeController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/recipes", rc.getAllRecipes)
	r.POST("/recipes", rc.createRecipe)
	r.GET("/recipes/:id", rc.getRecipeByID)
	r.PUT("/recipes/:id", rc.updateRecipe)
	r.DELETE("/recipes/:id", rc.deleteRecipe)
}

func (rc *RecipeController) getAllRecipes(c *gin.Context) {
	recipes, err := rc.repo.GetAllRecipes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipes)
}

func (rc *RecipeController) getRecipeByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	recipe, err := rc.repo.GetRecipeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

func (rc *RecipeController) createRecipe(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := rc.repo.CreateRecipe(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, recipe)
}

func (rc *RecipeController) updateRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = uint(id)
	err = rc.repo.UpdateRecipe(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

func (rc *RecipeController) deleteRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	err = rc.repo.DeleteRecipe(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

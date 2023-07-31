package shoppinglist

import (
	"github.com/cvele/recipe/pkg/models"
)

type ShoppingListServiceInterface interface {
	GenerateShoppingList() ([]models.ShoppingItem, *int, error)
}

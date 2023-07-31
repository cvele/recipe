package models

type ShoppingItem struct {
	Name     string
	Quantity float64
	Unit     string
	Cost     int `json:"cost"`
}

package models

type NutritionalValues struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
	Fiber    float64 `json:"fiber"`
	Sugar    float64 `json:"sugar"`
}

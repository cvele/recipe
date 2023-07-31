package units

import "errors"

type UnitConverter struct {
	massConversionRates   map[string]map[string]float64
	volumeConversionRates map[string]map[string]float64
	defaultMassUnit       string
	defaultVolumeUnit     string
}

func NewUnitConverter(defaultMassUnit string, defaultVolumeUnit string) *UnitConverter {
	return &UnitConverter{
		defaultMassUnit:   defaultMassUnit,
		defaultVolumeUnit: defaultVolumeUnit,
		massConversionRates: map[string]map[string]float64{
			"kg": {
				"g":  1000,
				"lb": 2.20462,
				"oz": 35.274,
			},
			"g": {
				"kg": 0.001,
				"lb": 0.00220462,
				"oz": 0.035274,
			},
			"lb": {
				"g":  453.592,
				"kg": 0.453592,
				"oz": 16,
			},
			"oz": {
				"g":  28.3495,
				"kg": 0.0283495,
				"lb": 0.0625,
			},
		},
		volumeConversionRates: map[string]map[string]float64{

			"l": {
				"ml":    1000,
				"fl-oz": 33.814,
				"cups":  4.22675,
				"pt":    2.11338,
				"qt":    1.05669,
				"gal":   0.264172,
				"tsp":   202.884,
				"tbsp":  67.628,
			},
			"ml": {
				"l":     0.001,
				"fl-oz": 0.033814,
				"cups":  0.00422675,
				"pt":    0.00211338,
				"qt":    0.00105669,
				"gal":   0.000264172,
				"tsp":   0.202884,
				"tbsp":  0.067628,
			},
			"fl-oz": {
				"l":    0.0295735,
				"ml":   29.5735,
				"cups": 0.125,
				"pt":   0.0625,
				"qt":   0.03125,
				"gal":  0.0078125,
				"tsp":  6,
				"tbsp": 2,
			},
			"cups": {
				"l":     0.236588,
				"ml":    236.588,
				"fl-oz": 8,
				"pt":    0.5,
				"qt":    0.25,
				"gal":   0.0625,
				"tsp":   48,
				"tbsp":  16,
			},
			"pt": {
				"l":     0.473176,
				"ml":    473.176,
				"fl-oz": 16,
				"cups":  2,
				"qt":    0.5,
				"gal":   0.125,
				"tsp":   96,
				"tbsp":  32,
			},
			"qt": {
				"l":     0.946353,
				"ml":    946.353,
				"fl-oz": 32,
				"cups":  4,
				"pt":    2,
				"gal":   0.25,
				"tsp":   192,
				"tbsp":  64,
			},
			"gal": {
				"l":     3.78541,
				"ml":    3785.41,
				"fl-oz": 128,
				"cups":  16,
				"pt":    8,
				"qt":    4,
				"tsp":   768,
				"tbsp":  256,
			},
			"tsp": {
				"l":     0.00492892,
				"ml":    4.92892,
				"fl-oz": 0.166667,
				"cups":  0.0208333,
				"pt":    0.0104167,
				"qt":    0.00520833,
				"gal":   0.00130208,
				"tbsp":  0.333333,
			},
			"tbsp": {
				"l":     0.0147868,
				"ml":    14.7868,
				"fl-oz": 0.5,
				"cups":  0.0625,
				"pt":    0.03125,
				"qt":    0.015625,
				"gal":   0.00390625,
				"tsp":   3,
			},
		},
	}
}

func (u *UnitConverter) ConvertUnits(quantity float64, fromUnit string, toUnit string, unitType string) (float64, error) {
	if fromUnit == toUnit {
		return quantity, nil
	}

	var conversionRates map[string]map[string]float64
	if unitType == "mass" {
		conversionRates = u.massConversionRates
	} else if unitType == "volume" {
		conversionRates = u.volumeConversionRates
	} else {
		return 0, errors.New("unknown unit type")
	}

	conversionRate, ok := conversionRates[fromUnit][toUnit]
	if !ok {
		return 0, errors.New("unsupported unit conversion")
	}

	return quantity * conversionRate, nil
}

func (u *UnitConverter) GetAvailableUnits(unitType string) []string {
	var units []string

	switch unitType {
	case "mass":
		for unit := range u.massConversionRates {
			units = append(units, unit)
		}
	case "volume":
		for unit := range u.volumeConversionRates {
			units = append(units, unit)
		}
	default:
		return nil
	}

	return units
}

func (u *UnitConverter) IsValidUnit(unit string, unitType string) bool {
	switch unitType {
	case "mass":
		_, isValid := u.massConversionRates[unit]
		return isValid
	case "volume":
		_, isValid := u.volumeConversionRates[unit]
		return isValid
	default:
		return false
	}
}

func (u *UnitConverter) GetDefaultUnit(unitType string) string {
	if unitType == "mass" {
		return u.defaultMassUnit
	} else if unitType == "volume" {
		return u.defaultVolumeUnit
	}
	return ""
}

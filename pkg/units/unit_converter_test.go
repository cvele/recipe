package units_test

import (
	"testing"

	"github.com/cvele/recipe/pkg/units"
)

func TestUnitConverter_ConvertUnits(t *testing.T) {
	unitConverter := units.NewUnitConverter("kg", "l")

	tests := []struct {
		name      string
		quantity  float64
		fromUnit  string
		toUnit    string
		unitType  string
		want      float64
		expectErr bool
	}{
		{"Mass: g to kg", 1, "kg", "g", "mass", 1000, false},
		{"Mass: kg to g", 1000, "g", "kg", "mass", 1, false},
		{"Mass: kg to lb", 1, "kg", "lb", "mass", 2.20462, false},
		{"Mass: lb to g", 1, "lb", "g", "mass", 453.592, false},
		{"Volume: l to ml", 1, "l", "ml", "volume", 1000, false},
		{"Volume: fl-oz to cups", 1, "fl-oz", "cups", "volume", 0.125, false},
		{"Invalid unit", 1, "kg", "ml", "mass", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unitConverter.ConvertUnits(tt.quantity, tt.fromUnit, tt.toUnit, tt.unitType)

			if (err != nil) != tt.expectErr {
				t.Errorf("ConvertUnits() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if got != tt.want {
				t.Errorf("ConvertUnits() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnitConverter_GetAvailableUnits(t *testing.T) {
	unitConverter := units.NewUnitConverter("kg", "l")

	massUnits := unitConverter.GetAvailableUnits("mass")
	volumeUnits := unitConverter.GetAvailableUnits("volume")

	if len(massUnits) == 0 {
		t.Errorf("GetAvailableUnits() returned empty for mass units")
	}

	if len(volumeUnits) == 0 {
		t.Errorf("GetAvailableUnits() returned empty for volume units")
	}
}

func TestUnitConverter_IsValidUnit(t *testing.T) {
	unitConverter := units.NewUnitConverter("kg", "l")

	tests := []struct {
		name     string
		unit     string
		unitType string
		want     bool
	}{
		{"Valid mass unit", "kg", "mass", true},
		{"Valid volume unit", "l", "volume", true},
		{"Invalid mass unit", "l", "mass", false},
		{"Invalid volume unit", "kg", "volume", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unitConverter.IsValidUnit(tt.unit, tt.unitType); got != tt.want {
				t.Errorf("IsValidUnit() = %v, want %v", got, tt.want)
			}
		})
	}
}

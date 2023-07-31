package units

type UnitConverterInterface interface {
	ConvertUnits(quantity float64, fromUnit string, toUnit string, unitType string) (float64, error)
	GetAvailableUnits(unitType string) []string
	IsValidUnit(unit string, unitType string) bool
	GetDefaultUnit(unitType string) string
}

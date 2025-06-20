package temperature

type Temperature struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func New(city string, celsius float64) *Temperature {
	return &Temperature{
		City:       city,
		Celsius:    celsius,
		Fahrenheit: ToFahrenheit(celsius),
		Kelvin:     ToKelvin(celsius),
	}
}

func ToFahrenheit(celsius float64) float64 {
	return (celsius * 1.8) + 32
}

func ToKelvin(celsius float64) float64 {
	return celsius + 273
}

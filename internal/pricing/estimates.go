package pricing

type FareEstimate struct {
	VehicleType string `json:"vehicle_type"`
	FareVND     int    `json:"fare_vnd"`
}

var vehicleRates = [...]struct {
	vehicleType  string
	rateVNDPerKm int
}{
	{"motorcycle", 5_000},
	{"car", 12_000},
	{"electric_car", 14_000},
}

func EstimateFares(distanceKm float64) []FareEstimate {
	estimates := make([]FareEstimate, 0, len(vehicleRates))
	for _, rate := range vehicleRates {
		estimates = append(estimates, FareEstimate{
			VehicleType: rate.vehicleType,
			FareVND:     CalculateFare(distanceKm, rate.rateVNDPerKm),
		})
	}
	return estimates
}

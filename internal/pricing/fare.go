package pricing

import "math"

// CalculateFare rounds the fare up to the next 1,000 VND.
// distanceKm and rateVNDPerKm must be non-negative.
func CalculateFare(distanceKm float64, rateVNDPerKm int) int {
	rawFareVND := float64(rateVNDPerKm) * distanceKm
	return int(math.Ceil(rawFareVND/1000) * 1000)
}

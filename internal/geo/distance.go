package geo

import "math"

type Coordinate struct {
	Lat float64
	Lng float64
}

func HaversineKm(a, b Coordinate) float64 {
	const earthRadiusKm = 6371.0

	lat1 := a.Lat * math.Pi / 180
	lat2 := b.Lat * math.Pi / 180

	deltaLat := (b.Lat - a.Lat) * math.Pi / 180
	deltaLng := (b.Lng - a.Lng) * math.Pi / 180

	h := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	h = math.Max(0, math.Min(1, h))

	return 2 * earthRadiusKm * math.Asin(math.Sqrt(h))
}

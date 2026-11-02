package geo

import (
	"math"
	"testing"
)

func TestHaversineKm(t *testing.T) {
	for _, tt := range []struct {
		name string
		a, b Coordinate
		want float64
	}{
		{"same point", Coordinate{10, 20}, Coordinate{10, 20}, 0},
		{"one degree at equator", Coordinate{0, 0}, Coordinate{0, 1}, 111.194927},
		{"antipodes", Coordinate{0, 0}, Coordinate{0, 180}, math.Pi * 6371},
		{"near antipodes", Coordinate{45, 10}, Coordinate{-45, -169.999999}, math.Pi * 6371},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := HaversineKm(tt.a, tt.b)
			if math.IsNaN(got) || math.IsInf(got, 0) || math.Abs(got-tt.want) > 0.001 {
				t.Fatalf("got %v km, want approximately %v km", got, tt.want)
			}
			if reverse := HaversineKm(tt.b, tt.a); math.Abs(got-reverse) > 1e-9 {
				t.Fatalf("distance changed when reversing points: %v vs %v", got, reverse)
			}
		})
	}
}

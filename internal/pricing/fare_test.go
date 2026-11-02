package pricing

import "testing"

func TestCalculateFare(t *testing.T) {
	for _, tt := range []struct {
		name       string
		distanceKm float64
		rate       int
		want       int
	}{
		{"motorcycle fractional", 2.3, 5_000, 12_000},
		{"car fractional", 2.3, 12_000, 28_000},
		{"electric car fractional", 2.3, 14_000, 33_000},
		{"already rounded", 2, 5_000, 10_000},
		{"zero distance", 0, 5_000, 0},
		{"small distance", 0.001, 5_000, 1_000},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateFare(tt.distanceKm, tt.rate); got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

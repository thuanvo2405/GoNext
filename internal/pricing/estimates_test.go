package pricing

import (
	"reflect"
	"testing"
)

func TestEstimateFares(t *testing.T) {
	want := []FareEstimate{
		{VehicleType: "motorcycle", FareVND: 12_000},
		{VehicleType: "car", FareVND: 28_000},
		{VehicleType: "electric_car", FareVND: 33_000},
	}
	if got := EstimateFares(2.3); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

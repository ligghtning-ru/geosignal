package geosignal_test

import (
	"fmt"
	"time"

	"github.com/ligghtning-ru/geosignal"
)

func ExampleResolve_regionCluster() {
	display := geosignal.Resolve([]geosignal.Observation{
		{Source: "provider_a", City: "Lviv", Country: "UA", Lat: 49.8397, Lng: 24.0297, HasCoord: true},
		{Source: "provider_b", City: "Vynnyky", Country: "UA", Lat: 49.8156, Lng: 24.1297, HasCoord: true},
	}, geosignal.Observation{City: "Lviv", Country: "UA"})

	fmt.Println(display.Label)
	// Output: Lviv area, UA
}

func ExampleResolveWithAgreement() {
	obs := []geosignal.Observation{
		{Source: "provider_a", City: "Los Angeles", Country: "US", Count: 2},
		{Source: "provider_b", City: "Brooklyn", Country: "US", Count: 1},
	}

	display, agreement := geosignal.ResolveWithAgreement(obs, geosignal.Observation{Country: "US"})

	fmt.Println(display.Label)
	fmt.Println(agreement.ObservedCities)
	// Output:
	// Los Angeles, US
	// [Brooklyn Los Angeles]
}

func ExampleTimezonesCompatible() {
	ok := geosignal.TimezonesCompatible(
		"Europe/Kiev",
		"Europe/Kyiv",
		"UA",
		time.Date(2026, time.May, 30, 12, 0, 0, 0, time.UTC),
	)

	fmt.Println(ok)
	// Output: true
}

package main

import (
	"fmt"

	"github.com/ligghtning-ru/geosignal"
)

func main() {
	obs := []geosignal.Observation{
		{Source: "provider_a", City: "Los Angeles", Country: "US", Count: 2},
		{Source: "provider_b", City: "Brooklyn", Country: "US", Count: 1},
	}

	display := geosignal.Resolve(obs, geosignal.Observation{Country: "US"})
	agreement := geosignal.AgreementForDisplay(display, obs)

	fmt.Println(display.Label)
	fmt.Println(agreement.ObservedCities)
}

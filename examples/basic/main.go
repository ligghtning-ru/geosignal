package main

import (
	"fmt"

	"github.com/ligghtning-ru/geosignal"
)

func main() {
	display := geosignal.Resolve([]geosignal.Observation{
		{Source: "provider_a", City: "Lviv", Country: "UA", Lat: 49.8397, Lng: 24.0297, HasCoord: true},
		{Source: "provider_b", City: "Vynnyky", Country: "UA", Lat: 49.8156, Lng: 24.1297, HasCoord: true},
	}, geosignal.Observation{City: "Lviv", Country: "UA"})

	fmt.Println(display.Label)
}

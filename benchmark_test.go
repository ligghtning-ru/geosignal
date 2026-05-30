package geosignal

import "testing"

func BenchmarkResolveSmall(b *testing.B) {
	obs := []Observation{
		{Source: "a", City: "Lviv", Country: "UA", Lat: 49.8397, Lng: 24.0297, HasCoord: true},
		{Source: "b", City: "Vynnyky", Country: "UA", Lat: 49.8156, Lng: 24.1297, HasCoord: true},
		{Source: "c", City: "Lviv", Country: "UA", Lat: 49.8397, Lng: 24.0297, HasCoord: true},
	}
	primary := Observation{City: "Lviv", Country: "UA"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Resolve(obs, primary)
	}
}

func BenchmarkResolveDisagreement(b *testing.B) {
	obs := []Observation{
		{Source: "a", City: "Dallas", Country: "US", Lat: 32.7767, Lng: -96.7970, HasCoord: true},
		{Source: "b", City: "San Diego", Country: "US", Lat: 32.7157, Lng: -117.1611, HasCoord: true},
		{Source: "c", City: "New Brunswick", Country: "US", Lat: 40.4862, Lng: -74.4518, HasCoord: true},
		{Source: "d", City: "Dallas", Country: "US", Lat: 32.7767, Lng: -96.7970, HasCoord: true},
	}
	primary := Observation{City: "Dallas", Country: "US"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		display, agreement := ResolveWithAgreement(obs, primary)
		if display.Label == "" || agreement.SourceCount == 0 {
			b.Fatal("empty result")
		}
	}
}

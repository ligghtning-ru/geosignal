package geosignal

import "testing"

func TestResolve(t *testing.T) {
	tests := []struct {
		name    string
		obs     []Observation
		primary Observation
		want    Display
	}{
		{
			name: "same city aliases collapse",
			obs: []Observation{
				{Source: "a", City: "Kyiv", Country: "UA"},
				{Source: "b", City: "Kiev", Country: "UA"},
				{Source: "c", City: "Київ", Country: "UA"},
			},
			primary: Observation{Country: "UA"},
			want:    Display{Kind: KindSingleSource, City: "Kyiv", CityLabel: "Kyiv", Country: "UA", Label: "Kyiv, UA"},
		},
		{
			name: "nearby cities become area",
			obs: []Observation{
				{Source: "a", City: "Lviv", Country: "UA", Lat: 49.8397, Lng: 24.0297, HasCoord: true},
				{Source: "b", City: "Vynnyky", Country: "UA", Lat: 49.8156, Lng: 24.1297, HasCoord: true},
			},
			primary: Observation{City: "Lviv", Country: "UA"},
			want:    Display{Kind: KindRegionCluster, City: "Lviv", CityLabel: "Lviv area", Country: "UA", Label: "Lviv area, UA", Derived: true},
		},
		{
			name: "far same-country cities stay country level",
			obs: []Observation{
				{Source: "a", City: "Dallas", Country: "US", Lat: 32.7767, Lng: -96.7970, HasCoord: true},
				{Source: "b", City: "San Diego", Country: "US", Lat: 32.7157, Lng: -117.1611, HasCoord: true},
				{Source: "c", City: "New Brunswick", Country: "US", Lat: 40.4862, Lng: -74.4518, HasCoord: true},
			},
			primary: Observation{City: "Dallas", Country: "US"},
			want:    Display{Kind: KindAmbiguousCity, CityLabel: "City varies by source", Country: "US", Label: "US", Derived: true},
		},
		{
			name: "cross-country conflict never claims city",
			obs: []Observation{
				{Source: "a", City: "São Paulo", Country: "BR", Lat: -23.5505, Lng: -46.6333, HasCoord: true},
				{Source: "b", City: "Los Angeles", Country: "US", Lat: 34.0522, Lng: -118.2437, HasCoord: true},
			},
			primary: Observation{City: "São Paulo", Country: "BR"},
			want:    Display{Kind: KindAmbiguousCountry, CityLabel: "Location varies by source", Label: "Location varies by source", Derived: true},
		},
		{
			name: "plurality same country wins",
			obs: []Observation{
				{Source: "a", City: "Los Angeles", Country: "US", Count: 2},
				{Source: "b", City: "Brooklyn", Country: "US", Count: 1},
			},
			primary: Observation{City: "Brooklyn", Country: "US"},
			want:    Display{Kind: KindPluralitySameCountry, City: "Los Angeles", CityLabel: "Los Angeles", Country: "US", Label: "Los Angeles, US", Derived: true},
		},
		{
			name:    "country-only primary is still useful",
			primary: Observation{Country: "United States"},
			want:    Display{Kind: KindSingleSource, Country: "US", Label: "US"},
		},
		{
			name: "city suffix and unicode alias normalize",
			obs: []Observation{
				{Source: "a", City: "München", Country: "DE"},
				{Source: "b", City: "Munich, Germany", Country: "DE"},
			},
			primary: Observation{Country: "Germany"},
			want:    Display{Kind: KindSingleSource, City: "Munich", CityLabel: "Munich", Country: "DE", Label: "Munich, DE"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Resolve(tt.obs, tt.primary)
			if got != tt.want {
				t.Fatalf("Resolve() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestResolveWithCustomAlias(t *testing.T) {
	got := Resolve([]Observation{
		{Source: "a", City: "Foo City", Country: "ZZ"},
		{Source: "b", City: "Foo", Country: "ZZ"},
	}, Observation{Country: "ZZ"}, WithCityAlias("Foo City", "Foo"))
	if got.Kind != KindSingleSource || got.Label != "Foo, ZZ" {
		t.Fatalf("custom alias not applied: %+v", got)
	}
}

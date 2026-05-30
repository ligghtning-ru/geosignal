package geosignal

// Kind describes how a display label was resolved from source observations.
type Kind string

const (
	// KindSingleSource means the observations point to one stable location, or
	// only country-level data is available.
	KindSingleSource Kind = "single_source"

	// KindRegionCluster means different city names are close enough to show an
	// area label instead of a precise city.
	KindRegionCluster Kind = "region_cluster"

	// KindPluralitySameCountry means one same-country city has a strong enough
	// plurality to be shown as the display city.
	KindPluralitySameCountry Kind = "plurality_same_country"

	// KindAmbiguousCity means the country is stable, but the city is not.
	KindAmbiguousCity Kind = "ambiguous_city"

	// KindAmbiguousCountry means observations disagree on the country.
	KindAmbiguousCountry Kind = "ambiguous_country"
)

// String returns the wire value of the display kind.
func (k Kind) String() string {
	return string(k)
}

// Observation is one geo signal from a provider, database, API, or local source.
// Country should be ISO-3166 alpha-2 when possible; common country names are
// normalized by the default resolver.
type Observation struct {
	Source   string  `json:"source,omitempty"`
	City     string  `json:"city,omitempty"`
	Country  string  `json:"country,omitempty"`
	Timezone string  `json:"timezone,omitempty"`
	Lat      float64 `json:"lat,omitempty"`
	Lng      float64 `json:"lng,omitempty"`
	HasCoord bool    `json:"has_coord,omitempty"`
	Count    int     `json:"count,omitempty"`
}

// Display is the user-facing location decision.
type Display struct {
	Kind      Kind   `json:"kind"`
	City      string `json:"city,omitempty"`
	CityLabel string `json:"city_label,omitempty"`
	Country   string `json:"country,omitempty"`
	Label     string `json:"label,omitempty"`
	Derived   bool   `json:"derived,omitempty"`
}

// Agreement describes how observations relate to a resolved display.
type Agreement struct {
	SourceCount         int      `json:"source_count"`
	AgreeingSources     []string `json:"agreeing_sources,omitempty"`
	DisagreeingSources  []string `json:"disagreeing_sources,omitempty"`
	ObservedCountries   []string `json:"observed_countries,omitempty"`
	ObservedCities      []string `json:"observed_cities,omitempty"`
	CountryDisagreement bool     `json:"country_disagreement,omitempty"`
	CityDisagreement    bool     `json:"city_disagreement,omitempty"`
}

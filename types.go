package geosignal

type Kind string

const (
	KindSingleSource         Kind = "single_source"
	KindRegionCluster        Kind = "region_cluster"
	KindPluralitySameCountry Kind = "plurality_same_country"
	KindAmbiguousCity        Kind = "ambiguous_city"
	KindAmbiguousCountry     Kind = "ambiguous_country"
)

type Observation struct {
	Source   string
	City     string
	Country  string
	Timezone string
	Lat      float64
	Lng      float64
	HasCoord bool
	Count    int
}

type Display struct {
	Kind      Kind
	City      string
	CityLabel string
	Country   string
	Label     string
	Derived   bool
}

type Agreement struct {
	SourceCount         int
	AgreeingSources     []string
	DisagreeingSources  []string
	ObservedCountries   []string
	ObservedCities      []string
	CountryDisagreement bool
	CityDisagreement    bool
}

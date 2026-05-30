package geosignal

import "strings"

const (
	defaultSameCityRadiusKM = 5.0
	defaultRegionRadiusKM   = 40.0
)

type Option func(*config)

type config struct {
	sameCityRadiusKM float64
	regionRadiusKM   float64
	cityAliases      map[string]string
	countryAliases   map[string]string
	citySuffixes     []string
	timezoneAliases  map[string]string
	countryTimezones map[string]map[string]struct{}
}

// WithSameCityRadiusKM changes the maximum distance for collapsing different
// city names into one precise city.
func WithSameCityRadiusKM(km float64) Option {
	return func(c *config) {
		if km > 0 {
			c.sameCityRadiusKM = km
		}
	}
}

// WithRegionRadiusKM changes the maximum distance for showing an "area" label.
func WithRegionRadiusKM(km float64) Option {
	return func(c *config) {
		if km > 0 {
			c.regionRadiusKM = km
		}
	}
}

// WithCityAlias adds or overrides a city spelling alias.
func WithCityAlias(alias, canonical string) Option {
	return func(c *config) {
		alias = cityKeyRaw(alias)
		canonical = strings.TrimSpace(canonical)
		if alias != "" && canonical != "" {
			c.cityAliases[alias] = canonical
		}
	}
}

// WithCityAliases adds or overrides multiple city spelling aliases.
func WithCityAliases(aliases map[string]string) Option {
	return func(c *config) {
		for alias, canonical := range aliases {
			alias = cityKeyRaw(alias)
			canonical = strings.TrimSpace(canonical)
			if alias != "" && canonical != "" {
				c.cityAliases[alias] = canonical
			}
		}
	}
}

// WithCountryAlias adds or overrides a country name to ISO-2 mapping.
func WithCountryAlias(alias, iso2 string) Option {
	return func(c *config) {
		alias = countryNameKey(alias)
		iso2 = strings.ToUpper(strings.TrimSpace(iso2))
		if alias != "" && len(iso2) == 2 {
			c.countryAliases[alias] = iso2
		}
	}
}

// WithCountryAliases adds or overrides multiple country name mappings.
func WithCountryAliases(aliases map[string]string) Option {
	return func(c *config) {
		for alias, iso2 := range aliases {
			alias = countryNameKey(alias)
			iso2 = strings.ToUpper(strings.TrimSpace(iso2))
			if alias != "" && len(iso2) == 2 {
				c.countryAliases[alias] = iso2
			}
		}
	}
}

// WithCitySuffix adds a suffix stripped during city normalization.
func WithCitySuffix(suffix string) Option {
	return func(c *config) {
		suffix = cityKeyRaw(suffix)
		if suffix != "" {
			c.citySuffixes = append(c.citySuffixes, suffix)
		}
	}
}

func defaultConfig() config {
	return config{
		sameCityRadiusKM: defaultSameCityRadiusKM,
		regionRadiusKM:   defaultRegionRadiusKM,
		cityAliases:      defaultCityAliases,
		countryAliases:   defaultCountryAliases,
		citySuffixes:     defaultCitySuffixes,
		timezoneAliases:  defaultTimezoneAliases,
		countryTimezones: defaultCountryTimezones,
	}
}

func newConfig(opts []Option) config {
	if len(opts) == 0 {
		return defaultConfig()
	}
	c := config{
		sameCityRadiusKM: defaultSameCityRadiusKM,
		regionRadiusKM:   defaultRegionRadiusKM,
		cityAliases:      copyStringMap(defaultCityAliases),
		countryAliases:   copyStringMap(defaultCountryAliases),
		citySuffixes:     append([]string(nil), defaultCitySuffixes...),
		timezoneAliases:  copyStringMap(defaultTimezoneAliases),
		countryTimezones: copyTimezoneSets(defaultCountryTimezones),
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&c)
		}
	}
	return c
}

func copyStringMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func copyTimezoneSets(in map[string]map[string]struct{}) map[string]map[string]struct{} {
	out := make(map[string]map[string]struct{}, len(in))
	for country, set := range in {
		next := make(map[string]struct{}, len(set))
		for key := range set {
			next[key] = struct{}{}
		}
		out[country] = next
	}
	return out
}

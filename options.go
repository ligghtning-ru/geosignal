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

func WithSameCityRadiusKM(km float64) Option {
	return func(c *config) {
		if km > 0 {
			c.sameCityRadiusKM = km
		}
	}
}

func WithRegionRadiusKM(km float64) Option {
	return func(c *config) {
		if km > 0 {
			c.regionRadiusKM = km
		}
	}
}

func WithCityAlias(alias, canonical string) Option {
	return func(c *config) {
		alias = cityKeyRaw(alias)
		canonical = strings.TrimSpace(canonical)
		if alias != "" && canonical != "" {
			c.cityAliases[alias] = canonical
		}
	}
}

func WithCountryAlias(alias, iso2 string) Option {
	return func(c *config) {
		alias = countryNameKey(alias)
		iso2 = strings.ToUpper(strings.TrimSpace(iso2))
		if alias != "" && len(iso2) == 2 {
			c.countryAliases[alias] = iso2
		}
	}
}

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
		cityAliases:      copyStringMap(defaultCityAliases),
		countryAliases:   copyStringMap(defaultCountryAliases),
		citySuffixes:     append([]string(nil), defaultCitySuffixes...),
		timezoneAliases:  copyStringMap(defaultTimezoneAliases),
		countryTimezones: copyTimezoneSets(defaultCountryTimezones),
	}
}

func newConfig(opts []Option) config {
	c := defaultConfig()
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

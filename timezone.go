package geosignal

import (
	"strings"
	"time"
)

var defaultTimezoneAliases = map[string]string{
	"europe/kiev": "Europe/Kyiv",
}

var defaultCountryTimezones = map[string]map[string]struct{}{
	"AM": timezoneSet("Asia/Yerevan"),
	"AZ": timezoneSet("Asia/Baku"),
	"BY": timezoneSet("Europe/Minsk"),
	"GE": timezoneSet("Asia/Tbilisi"),
	"JP": timezoneSet("Asia/Tokyo"),
	"KR": timezoneSet("Asia/Seoul"),
	"TR": timezoneSet("Europe/Istanbul", "Asia/Istanbul"),
	"UA": timezoneSet("Europe/Kyiv", "Europe/Kiev", "Europe/Uzhgorod", "Europe/Zaporozhye", "Europe/Simferopol"),
}

type TimezoneOption func(*timezoneConfig)

type timezoneConfig struct {
	aliases          map[string]string
	countryTimezones map[string]map[string]struct{}
}

func WithTimezoneAlias(alias, canonical string) TimezoneOption {
	return func(c *timezoneConfig) {
		alias = timezoneKey(alias)
		canonical = strings.TrimSpace(canonical)
		if alias != "" && canonical != "" {
			c.aliases[alias] = canonical
		}
	}
}

func WithTimezoneAliases(aliases map[string]string) TimezoneOption {
	return func(c *timezoneConfig) {
		for alias, canonical := range aliases {
			alias = timezoneKey(alias)
			canonical = strings.TrimSpace(canonical)
			if alias != "" && canonical != "" {
				c.aliases[alias] = canonical
			}
		}
	}
}

func CanonicalTimezone(name string, opts ...TimezoneOption) string {
	c := newTimezoneConfig(opts)
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	if canonical := c.aliases[timezoneKey(name)]; canonical != "" {
		return canonical
	}
	return name
}

func TimezonesCompatible(a, b, country string, at time.Time, opts ...TimezoneOption) bool {
	c := newTimezoneConfig(opts)
	a = CanonicalTimezone(a, WithTimezoneAliases(c.aliases))
	b = CanonicalTimezone(b, WithTimezoneAliases(c.aliases))
	if a == "" || b == "" {
		return false
	}
	if strings.EqualFold(a, b) || sameOffsetSignature(a, b, at, c) {
		return true
	}
	return countryTimezoneCompatible(country, a, b, c)
}

func TimezoneOffsetMinutes(name string, at time.Time, opts ...TimezoneOption) (int, bool) {
	name = CanonicalTimezone(name, opts...)
	if name == "" {
		return 0, false
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return 0, false
	}
	_, offsetSeconds := at.In(loc).Zone()
	return offsetSeconds / 60, true
}

func newTimezoneConfig(opts []TimezoneOption) timezoneConfig {
	c := timezoneConfig{
		aliases:          copyStringMap(defaultTimezoneAliases),
		countryTimezones: copyTimezoneSets(defaultCountryTimezones),
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&c)
		}
	}
	return c
}

func sameOffsetSignature(a, b string, at time.Time, c timezoneConfig) bool {
	for _, point := range signaturePoints(at) {
		ao, ok := TimezoneOffsetMinutes(a, point, WithTimezoneAliases(c.aliases))
		if !ok {
			return false
		}
		bo, ok := TimezoneOffsetMinutes(b, point, WithTimezoneAliases(c.aliases))
		if !ok || ao != bo {
			return false
		}
	}
	return true
}

func signaturePoints(at time.Time) []time.Time {
	year := at.UTC().Year()
	return []time.Time{
		at,
		time.Date(year, time.January, 15, 12, 0, 0, 0, time.UTC),
		time.Date(year, time.April, 15, 12, 0, 0, 0, time.UTC),
		time.Date(year, time.July, 15, 12, 0, 0, 0, time.UTC),
		time.Date(year, time.October, 15, 12, 0, 0, 0, time.UTC),
	}
}

func countryTimezoneCompatible(country, a, b string, c timezoneConfig) bool {
	allowed, ok := c.countryTimezones[strings.ToUpper(strings.TrimSpace(country))]
	if !ok {
		return false
	}
	_, aOK := allowed[timezoneKey(a)]
	_, bOK := allowed[timezoneKey(b)]
	return aOK && bOK
}

func timezoneSet(items ...string) map[string]struct{} {
	out := make(map[string]struct{}, len(items))
	for _, item := range items {
		if key := timezoneKey(item); key != "" {
			out[key] = struct{}{}
		}
	}
	return out
}

func timezoneKey(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

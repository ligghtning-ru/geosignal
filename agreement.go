package geosignal

import (
	"sort"
	"strings"
)

func AgreementForDisplay(display Display, observations []Observation, opts ...Option) Agreement {
	c := newConfig(opts)
	rows := normalizeAgreementObservations(observations, c)
	out := Agreement{}
	displayCountry := normalizeCountry(display.Country, c)
	displayCity := cityKey(display.City, c)
	countries := map[string]struct{}{}
	cities := map[string]string{}
	agreeing := map[string]struct{}{}
	disagreeing := map[string]struct{}{}

	for _, item := range rows {
		if item.Country != "" {
			countries[item.Country] = struct{}{}
		}
		if item.City != "" {
			cities[cityKey(item.City, c)] = item.City
		}
		source := strings.TrimSpace(item.Source)
		if source == "" {
			continue
		}
		out.SourceCount++
		if display.Kind == KindAmbiguousCountry {
			continue
		}

		countryMatches := displayCountry == "" || item.Country == "" || item.Country == displayCountry
		cityMatches := displayCity == "" || item.City == "" || cityKey(item.City, c) == displayCity
		switch {
		case displayCity != "" && countryMatches && cityMatches:
			agreeing[source] = struct{}{}
		case displayCity == "" && displayCountry != "" && countryMatches:
			agreeing[source] = struct{}{}
		case !countryMatches || (displayCity != "" && !cityMatches):
			disagreeing[source] = struct{}{}
		}
	}

	out.ObservedCountries = sortedKeys(countries)
	out.ObservedCities = sortedValues(cities)
	out.CountryDisagreement = len(out.ObservedCountries) > 1
	out.CityDisagreement = len(out.ObservedCities) > 1
	out.AgreeingSources = sortedKeys(agreeing)
	out.DisagreeingSources = sortedKeys(disagreeing)
	return out
}

func normalizeAgreementObservations(observations []Observation, c config) []Observation {
	out := make([]Observation, 0, len(observations))
	for _, item := range observations {
		item = normalizeObservation(item, c)
		if item.City == "" && item.Country == "" {
			continue
		}
		out = append(out, item)
	}
	return out
}

func sortedKeys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for key := range m {
		if key != "" {
			out = append(out, key)
		}
	}
	sort.Strings(out)
	return out
}

func sortedValues(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for _, value := range m {
		if value != "" {
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}

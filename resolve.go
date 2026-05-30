package geosignal

import (
	"sort"
	"strings"
)

// Resolve turns source observations into a conservative user-facing display.
// It prefers honest broader labels over precise city labels that are not
// supported by the observations.
func Resolve(observations []Observation, primary Observation, opts ...Option) Display {
	c := newConfig(opts)
	rows := normalizeObservations(observations, c)
	primary = normalizeObservation(primary, c)
	if hasCountryConflict(rows) || primaryCountryConflicts(rows, primary) {
		return Display{
			Kind:      KindAmbiguousCountry,
			CityLabel: "Location varies by source",
			Label:     "Location varies by source",
			Derived:   true,
		}
	}

	country := firstNonEmpty(primary.Country, commonCountry(rows))
	cityRows := rowsWithCity(rows)
	if len(cityRows) == 0 {
		return singleDisplay(primary.City, country)
	}
	if sameCity(cityRows) {
		city := cityRows[0].City
		if primary.City != "" && cityKey(primary.City, c) == cityKey(city, c) {
			city = primary.City
		}
		return singleDisplay(city, country)
	}
	if display, ok := sameCoordinateDisplay(cityRows, primary, country, c); ok {
		return display
	}
	if display, ok := regionClusterDisplay(cityRows, primary, country, c); ok {
		return display
	}
	if display, ok := pluralityDisplay(cityRows, country); ok {
		return display
	}
	return ambiguousCityDisplay(country)
}

func singleDisplay(city, country string) Display {
	city = strings.TrimSpace(city)
	country = strings.TrimSpace(country)
	if city == "" {
		if country == "" {
			return Display{}
		}
		return Display{Kind: KindSingleSource, Country: country, Label: country}
	}
	return Display{
		Kind:      KindSingleSource,
		City:      city,
		CityLabel: city,
		Country:   country,
		Label:     joinLabel(city, country),
	}
}

func ambiguousCityDisplay(country string) Display {
	country = strings.TrimSpace(country)
	label := country
	if label == "" {
		label = "City varies by source"
	}
	return Display{
		Kind:      KindAmbiguousCity,
		CityLabel: "City varies by source",
		Country:   country,
		Label:     label,
		Derived:   true,
	}
}

func hasCountryConflict(rows []Observation) bool {
	country := ""
	for _, item := range rows {
		if item.Country == "" {
			continue
		}
		if country == "" {
			country = item.Country
			continue
		}
		if country != item.Country {
			return true
		}
	}
	return false
}

func primaryCountryConflicts(rows []Observation, primary Observation) bool {
	if primary.Country == "" {
		return false
	}
	observed := commonCountry(rows)
	return observed != "" && observed != primary.Country
}

func commonCountry(rows []Observation) string {
	country := ""
	for _, item := range rows {
		if item.Country == "" {
			continue
		}
		if country == "" {
			country = item.Country
			continue
		}
		if country != item.Country {
			return ""
		}
	}
	return country
}

func rowsWithCity(rows []Observation) []Observation {
	out := make([]Observation, 0, len(rows))
	for _, item := range rows {
		if item.City != "" {
			out = append(out, item)
		}
	}
	return out
}

func sameCity(rows []Observation) bool {
	if len(rows) == 0 {
		return false
	}
	key := cityKeyRaw(rows[0].City)
	if key == "" {
		return false
	}
	for _, item := range rows[1:] {
		if cityKeyRaw(item.City) != key {
			return false
		}
	}
	return true
}

func sameCoordinateDisplay(rows []Observation, primary Observation, country string, c config) (Display, bool) {
	if len(rows) < 2 || !allHaveCoords(rows) || !allPairsWithin(rows, c.sameCityRadiusKM) {
		return Display{}, false
	}
	anchor := clusterAnchor(rows, primary)
	if anchor == "" {
		return Display{}, false
	}
	out := singleDisplay(anchor, country)
	out.Derived = true
	return out, true
}

func regionClusterDisplay(rows []Observation, primary Observation, country string, c config) (Display, bool) {
	withCoords := observationsWithCoords(rows)
	if len(withCoords) < 2 || !allPairsWithin(withCoords, c.regionRadiusKM) {
		return Display{}, false
	}
	anchor := clusterAnchor(rows, primary)
	if anchor == "" {
		return Display{}, false
	}
	cityLabel := anchor + " area"
	return Display{
		Kind:      KindRegionCluster,
		City:      anchor,
		CityLabel: cityLabel,
		Country:   country,
		Label:     joinLabel(cityLabel, country),
		Derived:   true,
	}, true
}

func clusterAnchor(rows []Observation, primary Observation) string {
	if primary.City != "" {
		key := cityKeyRaw(primary.City)
		for _, item := range rows {
			if cityKeyRaw(item.City) == key {
				return item.City
			}
		}
	}
	best := Observation{}
	for _, item := range rows {
		if best.City == "" ||
			item.Count > best.Count ||
			(item.Count == best.Count && strings.ToLower(item.City) < strings.ToLower(best.City)) {
			best = item
		}
	}
	return best.City
}

func pluralityDisplay(rows []Observation, country string) (Display, bool) {
	if len(rows) < 2 {
		return Display{}, false
	}
	rows = sortedRows(rows)
	total := 0
	for _, item := range rows {
		total += item.Count
	}
	hasStrongRatio := rows[0].Count >= 2 && rows[0].Count >= rows[1].Count*2
	hasMajority := rows[0].Count >= 3 && rows[0].Count > total/2
	if !hasStrongRatio && !hasMajority {
		return Display{}, false
	}
	return Display{
		Kind:      KindPluralitySameCountry,
		City:      rows[0].City,
		CityLabel: rows[0].City,
		Country:   country,
		Label:     joinLabel(rows[0].City, country),
		Derived:   true,
	}, true
}

func sortedRows(rows []Observation) []Observation {
	out := append([]Observation(nil), rows...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return strings.ToLower(out[i].City) < strings.ToLower(out[j].City)
	})
	return out
}

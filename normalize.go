package geosignal

import "strings"

var defaultCityAliases = map[string]string{
	"moscow":           "Moscow",
	"moskva":           "Moscow",
	"москва":           "Moscow",
	"kyiv":             "Kyiv",
	"kiev":             "Kyiv",
	"київ":             "Kyiv",
	"киев":             "Kyiv",
	"odesa":            "Odesa",
	"odessa":           "Odesa",
	"одеса":            "Odesa",
	"одесса":           "Odesa",
	"saint petersburg": "Saint Petersburg",
	"sankt petersburg": "Saint Petersburg",
	"st petersburg":    "Saint Petersburg",
	"spb":              "Saint Petersburg",
	"lviv":             "Lviv",
	"lvov":             "Lviv",
	"львів":            "Lviv",
	"львов":            "Lviv",
	"dnipro":           "Dnipro",
	"dnepr":            "Dnipro",
	"dnipropetrovsk":   "Dnipro",
	"new york":         "New York",
	"new york city":    "New York",
	"nyc":              "New York",
	"los angeles":      "Los Angeles",
	"la":               "Los Angeles",
	"sao paulo":        "São Paulo",
	"são paulo":        "São Paulo",
	"munich":           "Munich",
	"munchen":          "Munich",
	"münchen":          "Munich",
	"vienna":           "Vienna",
	"wien":             "Vienna",
	"warsaw":           "Warsaw",
	"warszawa":         "Warsaw",
	"prague":           "Prague",
	"praha":            "Prague",
	"lisbon":           "Lisbon",
	"lisboa":           "Lisbon",
	"milan":            "Milan",
	"milano":           "Milan",
	"rome":             "Rome",
	"roma":             "Rome",
	"athens":           "Athens",
	"athina":           "Athens",
	"tokyo":            "Tokyo",
	"東京":               "Tokyo",
	"seoul":            "Seoul",
	"서울":               "Seoul",
	"beijing":          "Beijing",
	"北京":               "Beijing",
	"shanghai":         "Shanghai",
	"上海":               "Shanghai",
	"hong kong":        "Hong Kong",
	"香港":               "Hong Kong",
}

var defaultCountryAliases = map[string]string{
	"united states":              "US",
	"usa":                        "US",
	"u.s.":                       "US",
	"uk":                         "GB",
	"united kingdom":             "GB",
	"great britain":              "GB",
	"russia":                     "RU",
	"russian federation":         "RU",
	"россия":                     "RU",
	"ukraine":                    "UA",
	"украина":                    "UA",
	"україна":                    "UA",
	"germany":                    "DE",
	"deutschland":                "DE",
	"brazil":                     "BR",
	"brasil":                     "BR",
	"south korea":                "KR",
	"korea, republic of":         "KR",
	"republic of korea":          "KR",
	"china":                      "CN",
	"people's republic of china": "CN",
}

var defaultCitySuffixes = []string{
	" city",
	" town",
	" metropolis",
	", russia",
	", россия",
	", usa",
	", united states",
	", uk",
	", united kingdom",
	", germany",
	", deutschland",
	", ukraine",
	", украина",
	", україна",
	", brazil",
	", brasil",
}

func normalizeObservations(rows []Observation, c config) []Observation {
	merged := map[string]Observation{}
	countryOnly := make([]Observation, 0)
	for _, item := range rows {
		item = normalizeObservation(item, c)
		if item.City == "" && item.Country == "" {
			continue
		}
		if item.City == "" {
			countryOnly = append(countryOnly, item)
			continue
		}
		key := cityKey(item.City, c) + "\x00" + item.Country
		if existing, ok := merged[key]; ok {
			existing.Count += item.Count
			if !existing.HasCoord && item.HasCoord {
				existing.Lat = item.Lat
				existing.Lng = item.Lng
				existing.HasCoord = true
			}
			merged[key] = existing
			continue
		}
		merged[key] = item
	}
	out := make([]Observation, 0, len(merged)+len(countryOnly))
	for _, item := range merged {
		out = append(out, item)
	}
	out = append(out, countryOnly...)
	return sortedRows(out)
}

func normalizeObservation(item Observation, c config) Observation {
	item.Source = strings.TrimSpace(item.Source)
	item.City = normalizeCity(item.City, c)
	item.Country = normalizeCountry(item.Country, c)
	item.Timezone = canonicalTimezone(item.Timezone, c.timezoneAliases)
	if item.Count <= 0 {
		item.Count = 1
	}
	return item
}

func normalizeCity(city string, c config) string {
	city = strings.Join(strings.Fields(strings.TrimSpace(city)), " ")
	if city == "" {
		return ""
	}
	if open := strings.Index(city, "("); open > 0 && strings.HasSuffix(city, ")") {
		city = strings.TrimSpace(city[:open])
	}
	city = stripCitySuffix(city, c.citySuffixes)
	if alias := c.cityAliases[cityKeyRaw(city)]; alias != "" {
		return alias
	}
	return city
}

func stripCitySuffix(city string, suffixes []string) string {
	key := cityKeyRaw(city)
	for _, suffix := range suffixes {
		if strings.HasSuffix(key, suffix) {
			return strings.TrimSpace(city[:len(city)-len(suffix)])
		}
	}
	return city
}

func normalizeCountry(country string, c config) string {
	country = strings.Join(strings.Fields(strings.TrimSpace(country)), " ")
	if country == "" {
		return ""
	}
	if len(country) == 2 {
		return strings.ToUpper(country)
	}
	if alias := c.countryAliases[countryNameKey(country)]; alias != "" {
		return alias
	}
	return strings.ToUpper(country)
}

func cityKey(city string, c config) string {
	return cityKeyRaw(normalizeCity(city, c))
}

func cityKeyRaw(city string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(city)), " "))
}

func countryNameKey(country string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(country)), " "))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func joinLabel(parts ...string) string {
	clean := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			clean = append(clean, part)
		}
	}
	return strings.Join(clean, ", ")
}

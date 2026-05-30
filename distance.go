package geosignal

import "math"

func allHaveCoords(rows []Observation) bool {
	for _, item := range rows {
		if !item.HasCoord {
			return false
		}
	}
	return len(rows) > 0
}

func observationsWithCoords(rows []Observation) []Observation {
	out := make([]Observation, 0, len(rows))
	for _, item := range rows {
		if item.HasCoord {
			out = append(out, item)
		}
	}
	return out
}

func allPairsWithin(rows []Observation, km float64) bool {
	if len(rows) < 2 {
		return false
	}
	for i := 0; i < len(rows); i++ {
		for j := i + 1; j < len(rows); j++ {
			if DistanceKM(rows[i].Lat, rows[i].Lng, rows[j].Lat, rows[j].Lng) > km {
				return false
			}
		}
	}
	return true
}

// DistanceKM returns the great-circle distance between two coordinates.
func DistanceKM(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKM = 6371.0
	toRad := func(v float64) float64 { return v * math.Pi / 180 }
	dLat := toRad(lat2 - lat1)
	dLng := toRad(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*math.Sin(dLng/2)*math.Sin(dLng/2)
	return earthRadiusKM * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

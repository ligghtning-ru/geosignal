package geosignal

import (
	"testing"
	"time"
)

func TestTimezonesCompatible(t *testing.T) {
	at := time.Date(2026, time.May, 30, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		a       string
		b       string
		country string
		want    bool
	}{
		{name: "kyiv legacy alias", a: "Europe/Kiev", b: "Europe/Kyiv", country: "UA", want: true},
		{name: "single-country variants", a: "Europe/Uzhgorod", b: "Europe/Kyiv", country: "UA", want: true},
		{name: "same dst signature", a: "Europe/Vienna", b: "Europe/Berlin", country: "AT", want: true},
		{name: "different continents", a: "Europe/Kyiv", b: "America/Los_Angeles", country: "UA", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TimezonesCompatible(tt.a, tt.b, tt.country, at)
			if got != tt.want {
				t.Fatalf("TimezonesCompatible() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanonicalTimezone(t *testing.T) {
	if got := CanonicalTimezone("Europe/Kiev"); got != "Europe/Kyiv" {
		t.Fatalf("CanonicalTimezone() = %q", got)
	}
}

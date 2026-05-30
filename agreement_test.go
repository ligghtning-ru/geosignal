package geosignal

import "testing"

func TestAgreementForDisplay(t *testing.T) {
	obs := []Observation{
		{Source: "a", City: "Moscow", Country: "RU"},
		{Source: "b", City: "Москва", Country: "RU"},
		{Source: "c", City: "Kazan", Country: "RU"},
	}
	display := Resolve(obs, Observation{City: "Moscow", Country: "RU"})
	got := AgreementForDisplay(display, obs)
	if got.SourceCount != 3 {
		t.Fatalf("SourceCount = %d, want 3", got.SourceCount)
	}
	if !contains(got.AgreeingSources, "a") || !contains(got.AgreeingSources, "b") {
		t.Fatalf("alias agreement not detected: %+v", got)
	}
	if !contains(got.DisagreeingSources, "c") {
		t.Fatalf("disagreement not detected: %+v", got)
	}
	if !got.CityDisagreement || got.CountryDisagreement {
		t.Fatalf("unexpected disagreement flags: %+v", got)
	}
}

func contains(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

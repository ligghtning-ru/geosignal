package geosignal

import (
	"encoding/json"
	"os"
	"testing"
)

type regressionCase struct {
	Name         string        `json:"name"`
	Primary      Observation   `json:"primary"`
	Observations []Observation `json:"observations"`
	Want         struct {
		Kind  Kind   `json:"kind"`
		Label string `json:"label"`
		City  string `json:"city"`
	} `json:"want"`
}

func TestRegressionCases(t *testing.T) {
	raw, err := os.ReadFile("testdata/cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []regressionCase
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("empty regression corpus")
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			got := Resolve(tc.Observations, tc.Primary)
			if got.Kind != tc.Want.Kind || got.Label != tc.Want.Label || got.City != tc.Want.City {
				t.Fatalf("Resolve() = %+v, want kind=%q label=%q city=%q", got, tc.Want.Kind, tc.Want.Label, tc.Want.City)
			}
		})
	}
}

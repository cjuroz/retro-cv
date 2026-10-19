package profile

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestProfileJSONUnmarshal(t *testing.T) {
	input := `{
		"name": "Cristian Uroz",
		"headline": "Go Engineer",
		"summary": "Backend engineering and cloud infrastructure",
		"experience": [
			{
				"company": "Example Corp",
				"role": "Staff Engineer",
				"period": "2024-2026",
				"description": "Built backend services",
				"technologies": ["Go", "Docker", "AWS"]
			}
		]
	}`

	var got Profile
	if err := json.Unmarshal([]byte(input), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if got.Name != "Cristian Uroz" {
		t.Errorf("Name = %q, want %q", got.Name, "Cristian Uroz")
	}

	if got.Headline != "Go Engineer" {
		t.Errorf("Headline = %q, want %q", got.Headline, "Go Engineer")
	}

	if len(got.Experience) != 1 {
		t.Fatalf("len(Experience) = %d, want 1", len(got.Experience))
	}

	exp := got.Experience[0]
	if exp.Company != "Example Corp" {
		t.Errorf("Company = %q, want %q", exp.Company, "Example Corp")
	}

	wantTechnologies := []string{"Go", "Docker", "AWS"}
	if !reflect.DeepEqual(exp.Technologies, wantTechnologies) {
		t.Errorf("Technologies = %v, want %v",
			exp.Technologies, wantTechnologies)
	}
}

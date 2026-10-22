package profile

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestExperienceJSON(t *testing.T) {
	data := `{
		"company": "Trend Micro",
		"summary": "A global cybersecurity company.",
		"website": "https://www.trendmicro.com",
		"location": "Taipei, Taiwan",
		"role": "Lead DevOps & Cloud Solutions Architect",
		"period": "2015 - 2020",
		"tasks": [
			"Led cloud infrastructure engineering.",
			"Standardized CI/CD pipelines."
		],
		"technologies": ["AWS", "Docker", "Kubernetes"]
	}`

	var got Experience
	if err := json.Unmarshal([]byte(data), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if got.Company != "Trend Micro" {
		t.Errorf("Company = %q, want %q", got.Company, "Trend Micro")
	}
	if got.Summary != "A global cybersecurity company." {
		t.Errorf("Summary = %q, want %q", got.Summary, "A global cybersecurity company.")
	}
	if got.Website != "https://www.trendmicro.com" {
		t.Errorf("Website = %q, want %q", got.Website, "https://www.trendmicro.com")
	}
	if got.Location != "Taipei, Taiwan" {
		t.Errorf("Location = %q, want %q", got.Location, "Taipei, Taiwan")
	}
	if got.Role != "Lead DevOps & Cloud Solutions Architect" {
		t.Errorf("Role = %q, want %q", got.Role, "Lead DevOps & Cloud Solutions Architect")
	}
	if got.Period != "2015 - 2020" {
		t.Errorf("Period = %q, want %q", got.Period, "2015 - 2020")
	}

	if len(got.Tasks) != 2 {
		t.Fatalf("len(Description) = %d, want 2", len(got.Tasks))
	}
	for i, description := range got.Tasks {
		if strings.TrimSpace(description) == "" {
			t.Errorf("Description[%d] should not be empty", i)
		}
	}

	if len(got.Technologies) != 3 {
		t.Errorf("len(Technologies) = %d, want 3", len(got.Technologies))
	}
}

func TestExperienceWebsiteOptional(t *testing.T) {
	data := `{
		"company": "Example Company",
		"summary": "An example company.",
		"location": "Taipei, Taiwan",
		"role": "Software Engineer",
		"period": "2020 - 2024",
		"description": ["Built backend services."],
		"technologies": ["Go"]
	}`

	var got Experience
	if err := json.Unmarshal([]byte(data), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if got.Website != "" {
		t.Errorf("Website = %q, want empty string when omitted", got.Website)
	}
}

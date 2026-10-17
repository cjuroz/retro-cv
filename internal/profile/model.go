package profile

type Experience struct {
	Company      string   `json:"company"`
	Role         string   `json:"role"`
	Period       string   `json:"period"`
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
}

type Profile struct {
	Name       string       `json:"name"`
	Headline   string       `json:"headline"`
	Summary    string       `json:"summary"`
	Experience []Experience `json:"experience"`
}

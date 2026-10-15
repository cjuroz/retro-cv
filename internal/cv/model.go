package cv

type Experience struct {
	Company      string
	Role         string
	Period       string
	Description  string
	Technologies []string
}

type Profile struct {
	Name       string
	Headline   string
	Summary    string
	Experience []Experience
}

package profile

type Profile struct {
	Name       string       `json:"name"`
	NameZh     string       `json:"name_zh,omitempty"`
	ResumeURL  string       `json:"resume_url"`
	Address    string       `json:"address"`
	Email      string       `json:"email"`
	WorkStatus string       `json:"work_status"`
	Headline   string       `json:"headline"`
	Summary    string       `json:"summary"`
	Socials    Socials      `json:"socials"`
	Skills     []string     `json:"skills"`
	Expertises []string     `json:"expertises"`
	Education  []Education  `json:"education"`
	Languages  []Language   `json:"languages"`
	Experience []Experience `json:"experience"`
}

type Experience struct {
	Company      string   `json:"company"`
	Summary      string   `json:"summary"`
	Website      string   `json:"website,omitempty"`
	Location     string   `json:"location"`
	Role         string   `json:"role"`
	Period       string   `json:"period"`
	Tasks        []string `json:"tasks"`
	Technologies []string `json:"technologies"`
}

type Language struct {
	Language string `json:"language"`
	Level    string `json:"level"`
}

type Education struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Period      string `json:"period"`
}

type Socials struct {
	LinkedIn string `json:"linkedin"`
	GitHub   string `json:"github"`
}

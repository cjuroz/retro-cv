package profile

type Profile struct {
	Name        string       `json:"name"`
	NameZh      string       `json:"name_zh,omitempty"`
	ResumeURL   string       `json:"resume_url"`
	Address     string       `json:"address"`
	Email       string       `json:"email"`
	WorkStatus  string       `json:"work_status"`
	LinkedIn    string       `json:"linkedin"`
	GitHub      string       `json:"github"`
	Headline    string       `json:"headline"`
	Summary     string       `json:"summary"`
	Skills      []string     `json:"skills"`
	Expertises  []string     `json:"expertises"`
	Educations  []Education  `json:"educations"`
	Languages   []Language   `json:"languages"`
	Experiences []Experience `json:"experiences"`
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

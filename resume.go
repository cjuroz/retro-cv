package main

import (
	"fmt"
	"strings"

	"github.com/cjuroz/retro-cv/internal/profile"
)

func generateResumeMarkdown(p profile.Profile) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# %s", p.Name)
	if p.NameZh != "" {
		fmt.Fprintf(&b, " %s", p.NameZh)
	}
	b.WriteString("\n\n")

	fmt.Fprintf(&b, "**ONLINE RESUME (GO RETRO-WEBSITE) ->** <%s>\n\n", p.ResumeURL)
	fmt.Fprintf(&b, "%s | [%s](mailto:%s)  \n", p.Address, p.Email, p.Email)
	fmt.Fprintf(&b, "Work Status: %s  \n", p.WorkStatus)
	fmt.Fprintf(&b, "[LinkedIn](%s) | [GitHub](%s)\n\n", p.Socials.LinkedIn, p.Socials.GitHub)

	b.WriteString("## PROFESSIONAL SUMMARY\n\n")
	fmt.Fprintf(&b, "### %s\n\n", p.Headline)
	fmt.Fprintf(&b, "%s\n\n", p.Summary)

	b.WriteString("## CORE EXPERTISE\n\n")
	for _, expertise := range p.Expertises {
		fmt.Fprintf(&b, "- %s\n", expertise)
	}
	b.WriteString("\n")

	b.WriteString("## TECHNICAL SKILLS\n\n")
	for _, skill := range p.Skills {
		fmt.Fprintf(&b, "- %s\n", skill)
	}
	b.WriteString("\n")

	b.WriteString("## WORK EXPERIENCE\n\n")
	for _, exp := range p.Experience {
		fmt.Fprintf(&b, "### %s @ %s | %s %s", exp.Role, exp.Company, exp.Location, exp.Period)

		if exp.Website != "" {
			fmt.Fprintf(&b, " | [%s](%s)",
				strings.TrimPrefix(exp.Website, "https://"), exp.Website)
		}

		b.WriteString("\n\n")

		fmt.Fprintf(&b, "%s\n\n", exp.Summary)

		for _, item := range exp.Tasks {
			fmt.Fprintf(&b, "- %s\n", item)
		}
		b.WriteString("\n")

		if len(exp.Technologies) > 0 {
			fmt.Fprintf(&b, "**Technologies:** %s\n\n",
				strings.Join(exp.Technologies, ", "))
		}
	}

	b.WriteString("## EDUCATION & LANGUAGES\n\n")
	for _, edu := range p.Education {
		fmt.Fprintf(&b, "**%s** %s  \n %s\n\n", edu.Degree, edu.Period, edu.Institution)
	}

	languages := make([]string, 0, len(p.Languages))
	for _, lang := range p.Languages {
		languages = append(languages, fmt.Sprintf("**%s:** %s", lang.Language, lang.Level))
	}
	fmt.Fprintf(&b, "%s\n\n", strings.Join(languages, " | "))

	return b.String()
}

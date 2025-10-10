package analyze

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func DetectName(html string) (string, string) {
	html = strings.ReplaceAll(html, "\u00A0", " ")
	re := regexp.MustCompile(`Agenda de l'utilisateur\s+([A-Z'-]+)\s+([A-Za-z'-]+)`)
	match := re.FindStringSubmatch(html)
	if len(match) < 3 {
		return "", ""
	}

	lastName := cases.Title(language.French).String(strings.ToLower(match[1]))
	firstName := cases.Title(language.French).String(match[2])

	return firstName, lastName
}

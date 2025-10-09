package analyze

import (
	"regexp"
	"strings"
)

func DetectName(html string) (string, string) {
	html = strings.ReplaceAll(html, "\u00A0", " ")
	re := regexp.MustCompile(`Agenda de l'utilisateur\s+([A-Z'-]+)\s+([A-Za-z'-]+)`)
	match := re.FindStringSubmatch(html)
	if len(match) < 3 {
		return "", ""
	}

	lastName := strings.Title(match[1])
	firstName := strings.Title(strings.ToLower(match[2]))

	return firstName, lastName
}

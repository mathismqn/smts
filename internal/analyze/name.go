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

	return match[2], match[1]
}

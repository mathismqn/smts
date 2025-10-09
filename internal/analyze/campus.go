package analyze

import (
	"regexp"
	"strings"
)

func DetectCampus(html string) string {
	counts := map[string]int{
		"Brest":  strings.Count(html, "BR-"),
		"Nantes": strings.Count(html, "NA-"),
		"Rennes": strings.Count(html, "RE-"),
	}

	if counts["Brest"] == 0 && counts["Nantes"] == 0 && counts["Rennes"] == 0 {
		reCampus := regexp.MustCompile(`(?i)(Brest|Nantes|Rennes)`)
		if match := reCampus.FindString(html); match != "" {
			return strings.Title(strings.ToLower(match))
		}
		return "unknown"
	}

	maxCampus := "unknown"
	maxCount := 0
	for campus, count := range counts {
		if count > maxCount {
			maxCount = count
			maxCampus = campus
		}
	}

	return maxCampus
}

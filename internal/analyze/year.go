package analyze

import "strings"

func DetectYear(html string) int {
	counts := map[int]int{
		1: strings.Count(html, "A1"),
		2: strings.Count(html, "A2"),
		3: strings.Count(html, "A3"),
	}

	counts[1] += strings.Count(html, "FIP A1") + strings.Count(html, "FIPA1")
	counts[2] += strings.Count(html, "FIP A2") + strings.Count(html, "FIPA2")
	counts[3] += strings.Count(html, "FIP A3") + strings.Count(html, "FIPA3")

	maxYear := 0
	maxCount := 0
	for year, count := range counts {
		if count > maxCount {
			maxCount = count
			maxYear = year
		}
	}

	return maxYear
}

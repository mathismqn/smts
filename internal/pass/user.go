package pass

import (
	"fmt"
	"smts/internal/analyze"
)

type User struct {
	FirstName string
	LastName  string
	Campus    string
	Year      int
}

func (*Client) getUserInfo(html string) (*User, error) {
	campus := analyze.DetectCampus(html)
	if campus == "unknown" {
		return nil, fmt.Errorf("could not detect campus")
	}
	year := analyze.DetectYear(html)
	if year == 0 {
		return nil, fmt.Errorf("could not detect year")
	}
	firstName, lastName := analyze.DetectName(html)
	if firstName == "" || lastName == "" {
		return nil, fmt.Errorf("could not detect name")
	}

	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Campus:    campus,
		Year:      year,
	}, nil
}

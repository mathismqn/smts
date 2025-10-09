package auth

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

func (s *Session) getUserInfo(body string) (*User, error) {
	campus := analyze.DetectCampus(body)
	if campus == "unknown" {
		return nil, fmt.Errorf("could not detect campus")
	}
	year := analyze.DetectYear(body)
	if year == 0 {
		return nil, fmt.Errorf("could not detect year")
	}
	firstName, lastName := analyze.DetectName(body)
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

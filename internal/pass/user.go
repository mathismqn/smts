package pass

import (
	"fmt"
	"smts/internal/analyze"
)

type User struct {
	FirstName string
	LastName  string
	Campus    string
}

func (*Client) getUserInfo(html string) (*User, error) {
	campus := analyze.DetectCampus(html)
	if campus == "unknown" {
		return nil, fmt.Errorf("could not detect campus")
	}
	firstName, lastName := analyze.DetectName(html)
	if firstName == "" || lastName == "" {
		return nil, fmt.Errorf("could not detect name")
	}

	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Campus:    campus,
	}, nil
}

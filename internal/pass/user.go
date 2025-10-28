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

func (c *Client) parseUserInfo(html string) (*User, error) {
	campus := analyze.DetectCampus(html)
	if campus == "unknown" {
		return nil, fmt.Errorf("user information not found")
	}
	firstName, lastName := analyze.DetectName(html)
	if firstName == "" || lastName == "" {
		return nil, fmt.Errorf("user information not found")
	}

	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Campus:    campus,
	}, nil
}

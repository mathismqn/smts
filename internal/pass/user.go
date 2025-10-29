package pass

import (
	"smts/internal/analyze"
)

type User struct {
	FirstName string
	LastName  string
	Campus    string
}

func (c *Client) parseUserInfo(html string) (*User, error) {
	campus := analyze.DetectCampus(html)
	firstName, lastName := analyze.DetectName(html)

	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Campus:    campus,
	}, nil
}

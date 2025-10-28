package creds

import (
	"encoding/json"

	"github.com/zalando/go-keyring"
)

const service = "imta-prod"

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func New(username, password string) *Config {
	return &Config{
		Username: username,
		Password: password,
	}
}

func (c *Config) Save() error {
	confJSON, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err := keyring.Set(service, "creds", string(confJSON)); err != nil {
		return err
	}

	return nil
}

func Load() (*Config, error) {
	confJSON, err := keyring.Get(service, "creds")
	if err != nil {
		return nil, err
	}

	var creds Config
	if err := json.Unmarshal([]byte(confJSON), &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

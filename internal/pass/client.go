package pass

import (
	"net/http"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{httpClient: httpClient}
}

func (c *Client) Authenticate() error {
	samlReq, err := c.getSAMLRequest()
	if err != nil {
		return err
	}

	consentURL, err := samlReq.submit(c.httpClient)
	if err != nil {
		return err
	}

	samlResp, err := c.getSAMLResponse(consentURL)
	if err != nil {
		return err
	}

	return samlResp.submit(c.httpClient)
}

package pass

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"smts/internal/sso"
	"strings"
)

type Client struct {
	httpClient *http.Client
	User       *User
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{httpClient: httpClient}
}

func (c *Client) Login(samlResp *sso.SAMLResponse) error {
	form := url.Values{}
	form.Add("SAMLResponse", samlResp.Value)

	req, err := http.NewRequest("POST", samlResp.URL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	cookies := c.httpClient.Jar.Cookies(req.URL)
	cookieHeader := ""
	for i, c := range cookies {
		if i > 0 {
			cookieHeader += "; "
		}
		cookieHeader += fmt.Sprintf("%s=%s", c.Name, c.Value)
	}
	cookieHeader += `; AuthenticationProvider={"Selected":"SAMLv2ProviderConfiguration"}`
	req.Header.Set("Cookie", cookieHeader)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !strings.Contains(string(body), "Bandeau.aspx") {
		return fmt.Errorf("unexpected error")
	}

	return nil
}

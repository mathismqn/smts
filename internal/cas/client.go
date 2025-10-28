package cas

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{httpClient: httpClient}
}

func (c *Client) Login(username, password string) error {
	resp, err := c.httpClient.Get("https://cas.imt-atlantique.fr/cas/login")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	reqURL := doc.Find("form").AttrOr("action", "")
	lt := doc.Find("input[name='lt']").AttrOr("value", "")
	exec := doc.Find("input[name='execution']").AttrOr("value", "")
	if reqURL == "" || lt == "" || exec == "" {
		return fmt.Errorf("unexpected error")
	}

	reqURL = "https://cas.imt-atlantique.fr" + reqURL
	resp, err = c.httpClient.PostForm(reqURL, map[string][]string{
		"username":  {username},
		"password":  {password},
		"lt":        {lt},
		"execution": {exec},
		"_eventId":  {"submit"},
		"submit":    {"SE CONNECTER"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if strings.Contains(string(body), "The credentials you provided cannot be determined to be authentic") {
		return fmt.Errorf("invalid username or password")
	}

	return nil
}

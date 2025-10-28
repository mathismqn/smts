package pass

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AgendaSession struct {
	Cookies []*http.Cookie
	URL     string
	User    *User
}

func (c *Client) GetAgendaSession() (*AgendaSession, error) {
	grpID, err := c.getGroupID()
	if err != nil {
		return nil, err
	}

	reqURL, html, err := c.getAgendaURL(grpID)
	if err != nil {
		return nil, err
	}

	user, err := c.parseUserInfo(html)
	if err != nil {
		return nil, err
	}

	return &AgendaSession{
		Cookies: c.httpClient.Jar.Cookies(reqURL),
		URL:     reqURL.String(),
		User:    user,
	}, nil
}

func (c *Client) getGroupID() (string, error) {
	resp, err := c.httpClient.Get("https://pass.imt-atlantique.fr/OpDotNet/Noyau/Bandeau.aspx?")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`var\s+IdGroupe\s*=\s*(\d+);`)
	match := re.FindStringSubmatch(string(body))
	if len(match) <= 1 {
		return "", fmt.Errorf("unexpected error")
	}

	return match[1], nil
}

func (c *Client) getAgendaForm(grpID string) (string, error) {
	resp, err := c.httpClient.Get("https://pass.imt-atlantique.fr/OpDotNet/Noyau/Content.aspx?groupe=" + grpID)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	iframeSrc := doc.Find("iframe").AttrOr("src", "")
	if iframeSrc == "" {
		return "", fmt.Errorf("unexpected error")
	}

	resp, err = c.httpClient.Get("https://pass.imt-atlantique.fr" + iframeSrc)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) getAgendaURL(grpID string) (*url.URL, string, error) {
	html, err := c.getAgendaForm(grpID)
	if err != nil {
		return nil, "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, "", err
	}

	reqURL := doc.Find("form").AttrOr("action", "")
	if reqURL == "" {
		return nil, "", fmt.Errorf("unexpected error")
	}
	if strings.HasPrefix(reqURL, "/") {
		reqURL = "https://pass.imt-atlantique.fr" + reqURL
	}

	formData := url.Values{}
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		name, exists := s.Attr("name")
		if !exists || name == "" {
			return
		}

		value := s.AttrOr("value", "")
		formData.Add(name, value)
	})

	resp, err := c.httpClient.PostForm(reqURL, formData)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return resp.Request.URL, string(body), nil
}

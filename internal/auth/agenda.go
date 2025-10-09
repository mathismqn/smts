package auth

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (s *Session) GetAgendaSession() ([]*http.Cookie, string, error) {
	grpID, err := s.getGroupID()
	if err != nil {
		return nil, "", err
	}

	reqURL, html, err := s.getAgendaURL(grpID)
	if err != nil {
		return nil, "", err
	}

	s.User, err = s.getUserInfo(html)
	if err != nil {
		return nil, "", err
	}

	return s.client.Jar.Cookies(reqURL), reqURL.String(), nil
}

func (s *Session) getGroupID() (string, error) {
	resp, err := s.client.Get("https://pass.imt-atlantique.fr/OpDotNet/Noyau/Bandeau.aspx?")
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

func (s *Session) getAgendaURL(grpID string) (*url.URL, string, error) {
	resp, err := s.client.Get("https://pass.imt-atlantique.fr/OpDotNet/Noyau/Content.aspx?groupe=" + grpID)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, "", err
	}

	iframeSrc := doc.Find("iframe").AttrOr("src", "")
	if iframeSrc == "" {
		return nil, "", fmt.Errorf("unexpected error")
	}

	resp, err = s.client.Get("https://pass.imt-atlantique.fr" + iframeSrc)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	doc, err = goquery.NewDocumentFromReader(resp.Body)
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

	resp, err = s.client.PostForm(reqURL, formData)
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

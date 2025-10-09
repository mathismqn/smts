package auth

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Session struct {
	client *http.Client
	User   *User
}

func NewSession() *Session {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	return &Session{client: client}
}

func (s *Session) Login(username, password string) error {
	samlReq, err := s.getSAMLRequest()
	if err != nil {
		return err
	}

	consentURL, err := s.authenticate(samlReq, username, password)
	if err != nil {
		return err
	}

	samlResp, err := s.getSAMLResponse(consentURL)
	if err != nil {
		return err
	}

	return s.finalizeLogin(samlResp)
}

func (s *Session) authenticate(samlReq *samlRequest, username, password string) (string, error) {
	resp, err := s.client.PostForm(samlReq.url, map[string][]string{"SAMLRequest": {samlReq.value}})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	lt := doc.Find("input[name='lt']").AttrOr("value", "")
	exec := doc.Find("input[name='execution']").AttrOr("value", "")
	reqURL := "https://cas.imt-atlantique.fr" + doc.Find("form").AttrOr("action", "")
	if lt == "" || exec == "" || reqURL == "" {
		return "", fmt.Errorf("unexpected error")
	}

	resp, err = s.client.PostForm(reqURL, map[string][]string{
		"username":  {username},
		"password":  {password},
		"lt":        {lt},
		"execution": {exec},
		"_eventId":  {"submit"},
		"submit":    {"SE CONNECTER"},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	consentURL := "https://idp.imt-atlantique.fr" + doc.Find("form").AttrOr("action", "")

	if strings.Contains(doc.Text(), "The credentials you provided cannot be determined to be authentic") {
		return "", fmt.Errorf("invalid username or password")
	}
	if !strings.Contains(doc.Text(), "Select an information release consent duration") || consentURL == "" {
		return "", fmt.Errorf("unexpected error")
	}

	return consentURL, nil
}

func (s *Session) finalizeLogin(samlResp *samlResponse) error {
	form := url.Values{}
	form.Add("SAMLResponse", samlResp.value)

	req, err := http.NewRequest("POST", samlResp.url, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	cookies := s.client.Jar.Cookies(req.URL)
	cookieHeader := ""
	for i, c := range cookies {
		if i > 0 {
			cookieHeader += "; "
		}
		cookieHeader += fmt.Sprintf("%s=%s", c.Name, c.Value)
	}
	cookieHeader += `; AuthenticationProvider={"Selected":"SAMLv2ProviderConfiguration"}`
	req.Header.Set("Cookie", cookieHeader)

	resp, err := s.client.Do(req)
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

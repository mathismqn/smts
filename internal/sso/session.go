package sso

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Session struct {
	Client       *http.Client
	SAMLResponse *SAMLResponse
}

func NewSession() *Session {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	return &Session{Client: client}
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

	return s.getSAMLResponse(consentURL)
}

func (s *Session) authenticate(samlReq *SAMLRequest, username, password string) (string, error) {
	resp, err := s.Client.PostForm(samlReq.URL, map[string][]string{"SAMLRequest": {samlReq.Value}})
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

	resp, err = s.Client.PostForm(reqURL, map[string][]string{
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

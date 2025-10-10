package sso

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type SAMLRequest struct {
	URL   string
	Value string
}

type SAMLResponse struct {
	URL   string
	Value string
}

func (s *Session) getSAMLRequest() (*SAMLRequest, error) {
	resp, err := s.Client.Get("https://pass.imt-atlantique.fr/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resp, err = s.Client.Post("https://pass.imt-atlantique.fr/OpDotNet/Noyau/Login.aspx", "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resp, err = s.Client.Get("https://pass.imt-atlantique.fr/OpDotNet/Noyau/Login.aspx?auth=SAMLv2ProviderConfiguration")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	samlReq := &SAMLRequest{}
	samlReq.Value = doc.Find("input[name='SAMLRequest']").AttrOr("value", "")
	samlReq.URL = doc.Find("form").AttrOr("action", "")
	if samlReq.URL == "" || samlReq.Value == "" {
		return nil, fmt.Errorf("unexpected error")
	}

	return samlReq, nil
}

func (s *Session) getSAMLResponse(reqUrl string) error {
	resp, err := s.Client.PostForm(reqUrl, map[string][]string{
		"_shib_idp_consentIds":     {"eduPersonPrincipalName", "givenName", "mail", "surName"},
		"_shib_idp_consentOptions": {"_shib_idp_globalConsent"},
		"_eventId_proceed":         {"Accepter"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	samlResp := &SAMLResponse{}
	samlResp.URL = doc.Find("form").AttrOr("action", "")
	samlResp.Value = doc.Find("input[name='SAMLResponse']").AttrOr("value", "")
	if samlResp.URL == "" || samlResp.Value == "" {
		return fmt.Errorf("unexpected error")
	}

	s.SAMLResponse = samlResp

	return nil
}

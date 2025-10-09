package auth

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type samlRequest struct {
	url   string
	value string
}

type samlResponse struct {
	url   string
	value string
}

func (s *Session) getSAMLRequest() (*samlRequest, error) {
	resp, err := s.client.Get("https://pass.imt-atlantique.fr/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resp, err = s.client.Post("https://pass.imt-atlantique.fr/OpDotNet/Noyau/Login.aspx", "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resp, err = s.client.Get("https://pass.imt-atlantique.fr/OpDotNet/Noyau/Login.aspx?auth=SAMLv2ProviderConfiguration")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	samlReq := &samlRequest{}
	samlReq.value = doc.Find("input[name='SAMLRequest']").AttrOr("value", "")
	samlReq.url = doc.Find("form").AttrOr("action", "")
	if samlReq.url == "" || samlReq.value == "" {
		return nil, fmt.Errorf("unexpected error")
	}

	return samlReq, nil
}

func (s *Session) getSAMLResponse(reqUrl string) (*samlResponse, error) {
	resp, err := s.client.PostForm(reqUrl, map[string][]string{
		"_shib_idp_consentIds":     {"eduPersonPrincipalName", "givenName", "mail", "surName"},
		"_shib_idp_consentOptions": {"_shib_idp_rememberConsent"},
		"_eventId_proceed":         {"Accepter"},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	samlResp := &samlResponse{}
	samlResp.url = doc.Find("form").AttrOr("action", "")
	samlResp.value = doc.Find("input[name='SAMLResponse']").AttrOr("value", "")
	if samlResp.url == "" || samlResp.value == "" {
		return nil, fmt.Errorf("unexpected error")
	}

	return samlResp, nil
}

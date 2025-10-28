package pass

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SAMLRequest struct {
	TargetURL string
	Token     string
}

type SAMLResponse struct {
	TargetURL string
	Token     string
}

func (c *Client) getSAMLRequest() (*SAMLRequest, error) {
	resp, err := c.httpClient.Get("https://pass.imt-atlantique.fr/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resp, err = c.httpClient.Post("https://pass.imt-atlantique.fr/OpDotNet/Noyau/Login.aspx", "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resp, err = c.httpClient.Get("https://pass.imt-atlantique.fr/OpDotNet/Noyau/Login.aspx?auth=SAMLv2ProviderConfiguration")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	samlReq := &SAMLRequest{}
	samlReq.TargetURL = doc.Find("form").AttrOr("action", "")
	samlReq.Token = doc.Find("input[name='SAMLRequest']").AttrOr("value", "")
	if samlReq.TargetURL == "" || samlReq.Token == "" {
		return nil, fmt.Errorf("unexpected error")
	}

	return samlReq, nil
}

func (r *SAMLRequest) submit(httpClient *http.Client) (string, error) {
	resp, err := httpClient.PostForm(r.TargetURL, map[string][]string{"SAMLRequest": {r.Token}})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	consentURL := doc.Find("form").AttrOr("action", "")
	if consentURL == "" {
		return "", fmt.Errorf("unexpected error")
	}

	return "https://idp.imt-atlantique.fr" + consentURL, nil
}

func (c *Client) getSAMLResponse(reqURL string) (*SAMLResponse, error) {
	resp, err := c.httpClient.PostForm(reqURL, map[string][]string{
		"_shib_idp_consentIds":     {"eduPersonPrincipalName", "givenName", "mail", "surName"},
		"_shib_idp_consentOptions": {"_shib_idp_globalConsent"},
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

	samlResp := &SAMLResponse{}
	samlResp.TargetURL = doc.Find("form").AttrOr("action", "")
	samlResp.Token = doc.Find("input[name='SAMLResponse']").AttrOr("value", "")
	if samlResp.TargetURL == "" || samlResp.Token == "" {
		return nil, fmt.Errorf("unexpected error")
	}

	return samlResp, nil
}

func (r *SAMLResponse) submit(httpClient *http.Client) error {
	form := url.Values{}
	form.Add("SAMLResponse", r.Token)

	req, err := http.NewRequest("POST", r.TargetURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	cookies := httpClient.Jar.Cookies(req.URL)
	cookieHeader := ""
	for i, cookie := range cookies {
		if i > 0 {
			cookieHeader += "; "
		}
		cookieHeader += fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
	}
	cookieHeader += `; AuthenticationProvider={"Selected":"SAMLv2ProviderConfiguration"}`
	req.Header.Set("Cookie", cookieHeader)

	resp, err := httpClient.Do(req)
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

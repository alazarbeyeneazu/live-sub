package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func SubLister(fqdn string) []string {
	url := fmt.Sprintf("https://subdomainfinder.c99.nl/scans/2024-01-15/%s", fqdn)
	response, err := http.Get(url)
	if err != nil {
		return []string{}
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []string{}
	}
	htmlResponse := string(body)

	pattern := `onclick="checkStatus\('([^']+)'`

	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatch(htmlResponse, -1)

	var subdomains []string
	for _, match := range matches {
		if len(match) > 1 {
			subdomain := match[1]
			subdomains = append(subdomains, subdomain)
		}
	}

	return subdomains

}

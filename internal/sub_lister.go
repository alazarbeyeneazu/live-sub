package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

func SubLister(fqdn string, startDate int) []string {
	date := time.Now().Add(time.Hour * -time.Duration(startDate) * 24)
	fDate := date.Format("2006-01-02")

	url := fmt.Sprintf("https://subdomainfinder.c99.nl/scans/%s/%s", fDate, fqdn)
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

	pattern := fmt.Sprintf(`href='//([^']*\b%s\b)'`, fqdn)

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

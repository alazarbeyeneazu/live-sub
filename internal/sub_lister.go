package internal

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func SubLister(fqdn string) []string {
	url := fmt.Sprintf("https://subdomainfinder.c99.nl/scans/2024-01-15/%s", fqdn)
	cmd := exec.Command("curl", url)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing curl command:", err)
		os.Exit(1)
	}

	htmlResponse := string(output)

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

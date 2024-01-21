package internal

import (
	"log"
	"os/exec"
	"regexp"
)

func runAmass(fqdn string) (string, error) {
	cmd := exec.Command("amass", "enum", "-d", fqdn)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func extractSubdomains(data string) []string {
	subdomainRegex := regexp.MustCompile(`(?m)([a-zA-Z0-9.-]+\.[a-zA-Z]+)`)
	matches := subdomainRegex.FindAllString(data, -1)

	subdomainsMap := make(map[string]bool)
	for _, match := range matches {
		subdomainsMap[match] = true
	}

	uniqueSubdomains := make([]string, 0, len(subdomainsMap))
	for subdomain := range subdomainsMap {
		uniqueSubdomains = append(uniqueSubdomains, subdomain)
	}

	return uniqueSubdomains
}

func AmassFindSubDomains(fqdn string) []string {
	amassOutput, err := runAmass(fqdn)
	if err != nil {
		log.Fatal("Unable to run the subdomain finder.\nNote that Amass is a prerequisite.", err)
	}
	return extractSubdomains(amassOutput)
}

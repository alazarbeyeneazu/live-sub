package internal

import (
	"log"
	"os/exec"
)

func runSubFinder(fqdn string) (string, error) {
	cmd := exec.Command("subfinder", "-d", fqdn)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func SubFinderFindSubDomains(fqdn string) []string {
	amassOutput, err := runAmass(fqdn)
	if err != nil {
		log.Fatal("Unable to run the subdomain finder.\nNote that subfinder is a prerequisite.", err)
	}
	return extractSubdomains(amassOutput)
}

package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hacker301et/live-sub/models"
)

func checkWebsite(url string, timeout time.Duration) error {
	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Status: %s", response.Status)
	}

	return nil
}

func CheckSubDomain(subs []string, respChannel chan models.ResponseMsg) {

	protocols := []string{"http://", "https://", "http://www.", "https://www."}
	timeout := time.Second * 5 // Adjust the timeout as needed

	for _, sub := range subs {
		for _, protocol := range protocols {
			url := protocol + sub

			err := checkWebsite(url, timeout)
			if err != nil {
				continue
			} else {
				respChannel <- models.ResponseMsg{ToolName: "sub-lister", FQDN: url}
			}
		}
	}
}

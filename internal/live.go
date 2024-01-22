package internal

import (
	"fmt"
	"net/http"
	"sync"
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

func worker(subChannel <-chan string, respChannel chan<- models.ResponseMsg, wg *sync.WaitGroup, timeout time.Duration, toolName string) {
	defer wg.Done()
	protocols := []string{"http://", "https://", "http://www.", "https://www."}

	for sub := range subChannel {
		for _, protocol := range protocols {
			url := protocol + sub
			err := checkWebsite(url, timeout)
			if err != nil {
				continue
			} else {
				respChannel <- models.ResponseMsg{ToolName: toolName, FQDN: url}
				break
			}
		}
	}
}

func CheckSubDomain(subs []string, respChannel chan models.ResponseMsg, toolName string) {
	numWorkers := 25
	timeout := time.Second * 5
	subChannel := make(chan string, len(subs))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(subChannel, respChannel, &wg, timeout, toolName)
	}
	for _, sub := range subs {
		subChannel <- sub
	}

}

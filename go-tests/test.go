package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
)

const (
	ServerURL             = "https://localhost:8081"
	ApiKey                = "Dzqqm3sncCYry01AuHiK7FDcCc35S4IzoOjgm2v8KyBpNlS52DyhMEXiJev6e8bq"
	DefaultOrgName        = "org"
	DefaultParentTeamName = "na"
	TeamNamePrefix        = "team-"
	DefaultClusterName    = "hubcluster"
)

var HttpClient = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

func main() {
	createTeamsParallel(1)
}

func createTeamsParallel(num int) {
	var wg sync.WaitGroup
	for i := 1; i <= num; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			name := fmt.Sprintf("team-%d", idx)
			if err := createTeam(name); err != nil {
				log.Printf("failed to create a team %s: %s", name, err.Error())
			} else {
				log.Printf("created a team %s", name)
			}
		}(i)
	}
	wg.Wait()
}

func createTeam(name string) error {
	url := fmt.Sprintf("%s/api/team/%s/%s/%s", ServerURL, DefaultOrgName, DefaultParentTeamName, name)
	_, err := performRequest("POST", url, nil, nil)
	return err
}

func performRequest(method, urlStr string, body []byte, headers http.Header) (*http.Response, error) {
	requestURL, _ := url.Parse(urlStr)
	req := &http.Request{
		Method: method,
		URL:    requestURL,
		Body:   io.NopCloser(bytes.NewReader(body)),
		Header: map[string][]string{
			"X-Api-Key": {ApiKey},
		},
	}
	for k, v := range headers {
		req.Header[k] = v
	}
	res, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("status code is %d", res.StatusCode)
	}
	return res, nil
}

package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	PipelineNamePrefix = "pipeline-"
)

var HttpClient = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

func main() {
	n := 20
	if _, err := CreateTeam(); err != nil {
		log.Fatal(err)
	}
	if err := CreatePipelines(n); err != nil {
		log.Fatal(err)
	}
	if err := RunPipelines(n); err != nil {
		log.Fatal(err)
	}
}

func CreateTeam() (*http.Response, error) {
	u := "https://localhost:8081/api/team/org/na/team"
	return PerformRequest("POST", u, nil, nil)
}

func CreatePipelines(n int) error {
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("%s%d", PipelineNamePrefix, i)
		u := fmt.Sprintf("https://localhost:8081/api/pipeline/org/team/%s", name)
		body := `{
    "clusterName": "hubcluster",
    "gitRepoSchema": {
        "gitRepo": "https://github.com/morzhanov/atlas-playground.git",
        "pathToRoot": "delay",
        "gitCred": {}
    },
    "reconciliationSchema": {
        "linkedPipelineName": "",
        "linkedArgoCdApplication": ""
    }
}`
		if _, err := PerformRequest("POST", u, []byte(body), nil); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
		log.Printf("created pipeline %s", name)
	}
	return nil
}

func RunPipelines(n int) error {
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("%s%d", PipelineNamePrefix, i)
		u := fmt.Sprintf("https://localhost:8081/api/sync/org/team/%s/ROOT_COMMIT", name)
		if _, err := PerformRequest("POST", u, nil, nil); err != nil {
			return err
		}
		time.Sleep(time.Second * 2)
		log.Printf("run pipeline %s", name)
	}
	return nil
}

func PerformRequest(method, urlStr string, body []byte, headers http.Header) (*http.Response, error) {
	requestURL, _ := url.Parse(urlStr)
	req := &http.Request{
		Method: method,
		URL:    requestURL,
		Body:   io.NopCloser(bytes.NewReader(body)),
		Header: map[string][]string{
			"X-Api-Key": {"Dzqqm3sncCYry01AuHiK7FDcCc35S4IzoOjgm2v8KyBpNlS52DyhMEXiJev6e8bq"},
		},
	}
	for k, v := range headers {
		req.Header[k] = v
	}
	res, err := HttpClient.Do(req)
	if err != nil {
		errRes := ReadErrorResponse(res)
		if errRes != "" {
			log.Printf("error response: %s", errRes)
		}
		return nil, err
	}
	if res.StatusCode >= 300 {
		errRes := ReadErrorResponse(res)
		if errRes != "" {
			log.Printf("error response: %s", errRes)
		}
		return nil, fmt.Errorf("status code is %d", res.StatusCode)
	}
	return res, nil
}

func ReadErrorResponse(res *http.Response) string {
	if res == nil {
		return ""
	}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(res.Body)
	if err != nil {
		return ""
	}
	return string(buf.Bytes())
}

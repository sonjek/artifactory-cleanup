package rest

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

func DoRSTQuery(method string, url string, body string, auth string) ([]gjson.Result, error) {
	client := &http.Client{}

	// Prepare query
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "text/plain")

	if body != "" {
		req.Body = io.NopCloser(strings.NewReader(body))
	}

	// Do query and handle error
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read and process the response body
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	responseString := string(response)
	errors := gjson.Get(responseString, "errors").Array()
	if len(errors) > 0 {
		return nil, fmt.Errorf("error in response: %v", errors)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error in response: %v", responseString)
	}

	return gjson.Get(responseString, "results").Array(), nil
}

func DoAQLQuery(endpoint string, aqlQuery string, auth string) ([]gjson.Result, error) {
	return DoRSTQuery("POST", endpoint, aqlQuery, auth)
}

package client

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

// RestClient is the client making REST calls.
type RestClient struct {
	Ctx                   context.Context
	BaseURL               string
	Verb                  string
	Body                  *bytes.Buffer
	AdditionalQueryParams []KeyValue
	AdditionalHeaders     []KeyValue
}

// NewRestClient is a constructor for creating a new client
func NewRestClient() *RestClient {
	return &RestClient{}
}

// MakeRestCall creates a new http client and makes arequest based on the parameters
func (r *RestClient) MakeRestCall() ([]byte, http.Header, error) {
	restURL, err := url.Parse(r.BaseURL)

	log.Println(restURL)

	if err != nil {
		log.WithFields(log.Fields{"baseURL": r.BaseURL, "verb": r.Verb}).Errorf("Failed to Parse URL with Error : %v", err)
		return nil, nil, err
	}
	queryParameters := restURL.Query()
	for _, kv := range r.AdditionalQueryParams {
		queryParameters.Add(kv.Key, kv.Value)
	}

	restURL.RawQuery = queryParameters.Encode()
	url := restURL.String()

	// create the client
	client := &http.Client{}

	// create a new request and check if there are errors
	request, err := http.NewRequestWithContext(r.Ctx, r.Verb, url, r.Body)
	if err != nil {
		log.WithFields(log.Fields{"url": url}).Errorf("Could not build request. : %v", err)
		return make([]byte, 0), nil, err
	}

	for _, kv := range r.AdditionalHeaders {
		request.Header.Add(kv.Key, kv.Value)
	}

	response, err := client.Do(request)
	if err != nil {
		log.WithFields(log.Fields{"url": url}).Errorf("Could not make request. : %v", err)
		return make([]byte, 0), nil, err
	}

	// save the response header
	h := response.Header

	// read the response body and check for errors
	responseBytes, err := ioutil.ReadAll(response.Body)
	if response.StatusCode <= 200 && response.StatusCode > 300 {
		log.WithFields(log.Fields{"url": url, "status_code": response.StatusCode}).Errorf("Received bad status code. : %v", err)
		return make([]byte, 0), nil, errors.New(response.Status)
	}
	if err != nil {
		log.WithFields(log.Fields{"url": url}).Errorf("Could not read response for body. : %v", err)
		return make([]byte, 0), nil, err
	}
	// return the response body, headers and no error
	return responseBytes, h, nil
}

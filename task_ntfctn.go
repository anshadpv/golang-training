package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func cnsPost(payload Payload, reqId string, cnsCampaignId string, cnsBulkDomain string, cnsBulkApiEndpoint string, httpMethod string) error {
	//Adding reqid and campaignid to payload
	payload.Email.ReqID = reqId
	payload.Email.CampaignID = cnsCampaignId

	// Converting payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload to JSON: %v", err)
	}

	// Printing payload for debugging
	fmt.Println("payload from campaign:", string(payloadBytes))

	// Creating HTTP client
	client := &http.Client{}

	// Creating request
	req, err := http.NewRequest(httpMethod, cnsBulkDomain+cnsBulkApiEndpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Sending request
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer res.Body.Close()

	// Reading response body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// Checking response status code
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Cns Error: %s", string(resBody))
	}

	return nil
}

// func main() {

// 	go func() {
// 		http.HandleFunc("/api/endpoint", func(w http.ResponseWriter, r *http.Request) {
// 		})
// 		http.ListenAndServe(":8880", nil)
// 	}()

// 	payload := map[string]interface{}{
// 		"data_placeholder": "some data",
// 		"campaign_name":    "example campaign",
// 	}
// 	reqId := "123"
// 	cnsCampaignId := "456"
// 	cnsBulkDomain := "http://localhost:8880/"
// 	cnsBulkApiEndpoint := "/api/endpoint"
// 	httpMethod := "POST"

// 	err := CnsPost(payload, reqId, cnsCampaignId, cnsBulkDomain, cnsBulkApiEndpoint, httpMethod)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	}
// }

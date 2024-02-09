package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

type RowElement struct {
	ClientID             string
	DestinationEmail     string
	DestinationMobNumber string
	DailyFrequency       int
}

// payload structure for sending notificatins
type Payload struct {
	Email    EmailPayload    `json:"email,omitempty"`
	Push     PushPayload     `json:"push,omitempty"`
	SMS      SMSPayload      `json:"sms,omitempty"`
	WhatsApp WhatsAppPayload `json:"whatsapp,omitempty"`
}

type DestinationDetails struct {
	ClientID             string                 `json:"client_id,omitempty"`
	DestinationEmail     string                 `json:"destination_email,omitempty"`
	DestinationMobNo     string                 `json:"destination_mobile_no,omitempty"`
	MaxNotificationCount int                    `json:"max_notification_count,omitempty"`
	DataPlaceholder      map[string]interface{} `json:"data_placeholder,omitempty"`
}

// payload structure specifically for email notifications
type EmailPayload struct {
	DestinationDetails []DestinationDetails `json:"destination_details,omitempty"`
	ReqID              string               `json:"req_id,omitempty"`
	CampaignID         string               `json:"campaign_id,omitempty"`
}

// payload structure specifically for push notifications
type PushPayload struct {
	DestinationDetails []DestinationDetails `json:"destination_details,omitempty"`
}

// payload structure specifically for sms notifications
type SMSPayload struct {
	DestinationDetails []DestinationDetails `json:"destination_details,omitempty"`
}

// payload structure specifically for whatsapp notifications
type WhatsAppPayload struct {
	DestinationDetails []DestinationDetails `json:"destination_details,omitempty"`
}

func handleBatch(row []RowElement, payload *Payload, columns []string, requestId, cnsCampaignID string) ([]string, int) {
	var failedRequests []string          //for storing failed request
	partitionID := generatePartitionID() //creating a random partition id
	fmt.Println("partitionID", partitionID)
	fmt.Println("number of columns", len(columns))
	isClientID := contains(columns, "client_id")
	isDestinationEmail := contains(columns, "destination_email")
	isDestinationMobNumber := contains(columns, "destination_mob_number")
	isDailyFrequency := contains(columns, "daily_frequency")

	fmt.Println(isClientID)
	fmt.Println(isDestinationEmail)
	fmt.Println(isDestinationMobNumber)
	fmt.Println(isDailyFrequency)

	if !(isClientID || isDestinationEmail || isDestinationMobNumber) {
		err := errors.New("client id, destination email, and destination mobile number are not present")
		panic(err)
	}

	count := 0
	reqCount := 0

	fmt.Println("Number of", len(row), "clients records", partitionID, len(row))

	rowSubList := makeSubLists(row, 3) // assuming cnsBulkAPIBatchSize is 3
	if rowSubList == nil || len(rowSubList) == 0 {
		fmt.Println("row sub list is empty")
		return failedRequests, count
	}

	for _, rowElement := range rowSubList {
		var clientDetails, destinationEmailDetails, destinationMobDetails []DestinationDetails
		reqCount++
		reqID := fmt.Sprintf("%s_%s_%d", requestId, partitionID, reqCount)
		for _, rowe := range rowElement {
			rowDailyFrequency := 0
			if isDailyFrequency {
				rowDailyFrequency = rowe.DailyFrequency
			}
			placeholder := make(map[string]interface{})
			if len(rowe.ClientID) > 0 && isClientID {
				clientData := DestinationDetails{
					ClientID:             rowe.ClientID,
					MaxNotificationCount: rowDailyFrequency,
					DataPlaceholder:      placeholder,
				}
				clientDetails = append(clientDetails, clientData)
			}
			if len(rowe.DestinationEmail) > 0 && isDestinationEmail {
				destinationEmail := DestinationDetails{
					DestinationEmail:     rowe.DestinationEmail,
					MaxNotificationCount: rowDailyFrequency,
					DataPlaceholder:      placeholder,
				}
				destinationEmailDetails = append(destinationEmailDetails, destinationEmail)
			}
			if len(rowe.DestinationMobNumber) > 0 && isDestinationMobNumber {
				destinationMobNumber := DestinationDetails{
					DestinationMobNo:     rowe.DestinationMobNumber,
					MaxNotificationCount: rowDailyFrequency,
					DataPlaceholder:      placeholder,
				}
				destinationMobDetails = append(destinationMobDetails, destinationMobNumber)
			}
		}

		// Constructing Payload
		if len(payload.Email.DestinationDetails) > 0 {
			// Payload with email destination details
			if len(clientDetails) != 0 {
				payload.Email.DestinationDetails = clientDetails
			} else if len(destinationEmailDetails) != 0 {
				payload.Email.DestinationDetails = destinationEmailDetails
			} else {
				panic(errors.New("clientid and email both empty for Email"))
			}
		}

		if len(payload.Push.DestinationDetails) > 0 {
			// Payload with push destination details
			if len(clientDetails) != 0 {
				payload.Push.DestinationDetails = clientDetails
			} else {
				panic(errors.New("clientid empty for Push"))
			}
		}

		if len(payload.SMS.DestinationDetails) > 0 {
			//payload with SMS destination details
			if len(clientDetails) != 0 {
				payload.SMS.DestinationDetails = clientDetails
			} else if len(destinationMobDetails) != 0 {
				payload.SMS.DestinationDetails = destinationMobDetails
			} else {
				panic(errors.New("clientid and mobile empty for SMS"))
			}
		}
		if len(payload.WhatsApp.DestinationDetails) > 0 {
			//payload with Whatsapp destination details
			if len(clientDetails) != 0 {
				payload.WhatsApp.DestinationDetails = clientDetails
			} else if len(destinationMobDetails) != 0 {
				payload.WhatsApp.DestinationDetails = destinationMobDetails
			} else {
				panic(errors.New("clientid and mobile empty for WhatsApp"))
			}
		}

		fmt.Println("Request Body", payload)

		cnsBulkDomain := "http://localhost:8880/"
		cnsBulkApiEndpoint := "/api/endpoint"
		httpMethod := "POST"

		// Calling cnsPost
		cnsPost(*payload, reqID, cnsCampaignID, cnsBulkDomain, cnsBulkApiEndpoint, httpMethod)
		count += len(rowElement)
	}

	return failedRequests, count
}

func generatePartitionID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var partitionID strings.Builder
	for i := 0; i < 5; i++ {
		partitionID.WriteByte(charset[rand.Intn(len(charset))])
	}
	return partitionID.String()
}

func contains(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

func makeSubLists(row []RowElement, batchSize int) [][]RowElement {
	if len(row) == 0 {
		return nil
	}
	var subLists [][]RowElement
	for i := 0; i < len(row); i += batchSize {
		end := i + batchSize
		if end > len(row) {
			end = len(row)
		}
		subLists = append(subLists, row[i:end])
	}
	return subLists
}

func main() {
	// Example usage
	row := []RowElement{
		{ClientID: "client1", DestinationEmail: "email1@example.com", DestinationMobNumber: "1234567890", DailyFrequency: 2},
		{ClientID: "client2", DestinationEmail: "email2@example.com", DestinationMobNumber: "0987654321", DailyFrequency: 3},
		// Add more rows if necessary
	}
	payload := &Payload{
		Email:    EmailPayload{},
		Push:     PushPayload{},
		SMS:      SMSPayload{},
		WhatsApp: WhatsAppPayload{},
	}
	columns := []string{"client_id", "destination_email", "destination_mob_number", "daily_frequency"}
	requestID := "request123"
	cnsCampaignID := "campaign123"

	failedRequests, count := handleBatch(row, payload, columns, requestID, cnsCampaignID)
	fmt.Println("Failed Requests:", failedRequests)
	fmt.Println("Processed Count:", count)
}

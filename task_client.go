// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"strings"
// 	"sync"
// )

// // type Message struct {
// // 	Text string `json:"text"`
// // }

// func main() {
// 	var wg sync.WaitGroup

// 	for i := 0; i < 10; i++ {
// 		wg.Add(1)
// 		go sendMsg(&wg)
// 	}

// 	wg.Wait()
// }

// func sendMsg(wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	apiURL := "http://localhost:8009/send-message"

// 	message := "Hello API, This is a test message."

// 	requestBody := fmt.Sprintf(`{"text": "%s"}`, message)

// 	resp, err := http.Post(apiURL, "application/json", strings.NewReader(requestBody))
// 	if err != nil {
// 		fmt.Println("Error sending request:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	fmt.Printf("Response Status for '%s': %s\n", message, resp.Status)
// }

package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Message struct {
	Text string
}

func main() {
	messageChannel := make(chan Message, 5)
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go sendMsg(ctx, i, messageChannel, &wg)
	}

	close(messageChannel)

	// Wait for all goroutines to finish
	wg.Wait()
}

func sendMsg(ctx context.Context, messageNumber int, messageChannel chan Message, wg *sync.WaitGroup) {

	defer wg.Done()

	apiURL := "http://localhost:8009/send-message"
	message := Message{Text: fmt.Sprintf("Test message %d", messageNumber)}
	requestBody := fmt.Sprintf(`{"text": "%s"}`, message.Text)

	if requestBody == `{"text": "Test message 4"}` {
		time.Sleep(5 * time.Second)
	}

	select {
	case <-ctx.Done():
		fmt.Printf("Goroutine for message %d canceled\n", messageNumber)
		return
	default:
	}

	resp, err := http.Post(apiURL, "application/json", strings.NewReader(requestBody))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response Status for '%s': %s\n", message.Text, resp.Status)
}

// ctxtimeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
// defer cancel()

// select {
// case <-ctxtimeout.Done():
// 	fmt.Println("Taking more than 3 seconds : ", ctxtimeout.Err())
// }

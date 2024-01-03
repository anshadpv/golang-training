package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/send-message", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Inside the API.")
	})

	// Start the server on port 8009
	fmt.Println("Server listening on :8009")

	// ListenAndServe should be called outside of the loop
	http.ListenAndServe(":8009", nil)
}

// Simple mock server for testing the Marketplace Order API
// This server ONLY accepts array format [...] to match the real API
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/orderextws/public/Orders/createOrder", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request to %s", r.Method, r.URL.Path)
		log.Printf("Headers:")
		for k, v := range r.Header {
			log.Printf("  %s: %v", k, v)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("Body: %s", string(body))

		// STRICT: Only accept direct array [...] format (like the real API)
		var orders []map[string]interface{}
		if err := json.Unmarshal(body, &orders); err != nil {
			log.Printf("ERROR: Body is not a JSON array! Got: %s", string(body[:min(100, len(body))]))
			http.Error(w, "Invalid JSON body - must be an array", http.StatusBadRequest)
			return
		}

		log.Printf("SUCCESS: Parsed as direct array with %d orders", len(orders))

		// Return a mock success response
		response := map[string]interface{}{
			"success":    true,
			"message":    "Order created successfully",
			"orderId":    "ORD-2024-001234",
			"orders":     len(orders),
			"bodyFormat": "array",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		log.Printf("Responded with success")
	})

	log.Println("Mock Marketplace API running on http://localhost:9999 (STRICT array-only mode)")
	log.Fatal(http.ListenAndServe(":9999", nil))
}

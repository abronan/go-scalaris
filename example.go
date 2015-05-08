package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/abronan/go-scalaris/client"
)

var (
	url = "localhost:8000"
)

// Example usage
// Work in progress
func main() {
	log.Printf("Scalaris endpoint: ", url)

	client := &client.Client{Url: url}

	res, err := client.TxWrite("hello", "woot")
	if err != nil {
		log.Fatalf("Err: ", err)
	}
	log.Info(res)

	res, err = client.Read("hello")
	if err != nil {
		log.Fatalf("Err: ", err)
	}
	log.Info(res)

	oldValue := map[string]interface{}{
		"type":  "as_is",
		"value": "woot",
	}
	newValue := map[string]interface{}{
		"type":  "as_is",
		"value": "wootwoot",
	}
	res, err = client.TxTestAndSet("hello", oldValue, newValue)
	if err != nil {
		log.Fatalf("Err: ", err)
	}
	log.Info(res)

	// Test CAS
	res, err = client.TxWrite("genId", 1)
	if err != nil {
		log.Fatalf("Err: ", err)
	}
	log.Info(res)

	// Read back the value
	res, err = client.Read("genId")
	if err != nil {
		log.Fatalf("Err: ", err)
	}
	log.Info(res)

	// Get the counter, needs proper deserialization method to get back the value
	oldCounter := res["result"].(map[string]interface{})["value"].(map[string]interface{})["value"].(float64)
	newCounter := oldCounter + 1
	log.Info("value: ", newCounter)

	// Use Scalaris Master (post REV-7316) to make this working
	res, err = client.TxTestAndSet("genId", oldCounter, newCounter)
	if err != nil {
		log.Fatalf("Err: ", err)
	}
	log.Info(res)
}

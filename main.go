package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type ReadResp struct {
	Status string `json:"status"`
	Value  string `json:"value"`
	Reason string `json:"reason"`
}

type WriteResp struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type Client struct {
	url string
}

var (
	url     = "localhost:8000"
	api     = "/api"
	tx      = "/tx.yaws"
	rdht    = "/rdht.yaws"
	dht_raw = "/dht_raw.yaws"
	pubsub  = "/pubsub.yaws"
	monitor = "/monitor.yaws"
)

func init() {
	env := os.Getenv("SCALARIS_API_ENDPOINT")
	if env != "" {
		url = env
	}
}

// Example program
func main() {
	fmt.Printf("Scalaris endpoint: %s", url)

	client := &Client{url: "192.168.42.10:8000"}

	res, err := client.TxWrite("hello", "woot")
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	log.Println(res)

	res, err = client.Read("hello")
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	log.Println(res)

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
		log.Fatalf("Err: %v", err)
	}
	log.Println(res)

	// Test Id and increment
	res, err = client.TxWrite("genId", 1)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	log.Println(res)

	res, err = client.Read("genId")
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	log.Println(res)
	log.Println("\n")

	// FIXME simplify - Get Value and Increment
	oldCounter := res["result"].(map[string]interface{})["value"].(map[string]interface{})["value"].(float64)
	newCounter := oldCounter + 1
	fmt.Printf("value: %d", newCounter)
	fmt.Printf("\n")

	res, err = client.TestAndSet("genId", oldCounter, newCounter)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	log.Println(res)
	fmt.Printf("\n")
}

// TODO improve with full-fledged config, etc.
func NewClient(url string) (*Client, error) {
	return &Client{url: url}, nil
}

func Call(address string, method string, id interface{}, params []interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(map[string]interface{}{
		"method": method,
		"id":     id,
		"params": params,
	})
	fmt.Printf(string(data))

	if err != nil {
		log.Fatalf("Marshal: %v", err)
		return nil, err
	}

	resp, err := http.Post(address, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Post: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
		return nil, err
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil, err
	}

	return result, nil
}

// TODO allow to retrieve value as <json_value> for test and set operations
func (c *Client) Read(key string) (map[string]interface{}, error) {
	url := fmt.Sprint("http://", c.url, api, tx)
	res, err := Call(url, "read", 1, []interface{}{key})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	return res, nil
}

func write() {

}

func (c *Client) TestAndSet(key string, oldValue interface{}, newValue interface{}) (map[string]interface{}, error) {
	// {"write": {"key": <key>, "old": <json_value>, "new": <json_value>} }
	data := map[string]interface{}{
		"key": key,
		"old": encode_value(oldValue),
		"new": encode_value(newValue),
	}

	url := fmt.Sprint("http://", c.url, api, tx)
	res, err := Call(url, "test_and_set", 1, []interface{}{data})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	return res, nil
}

func (c *Client) Delete(key string) error {
	// {"write": {"key": <key>, "old": <json_value>, "new": <json_value>} }
	data := map[string]interface{}{
		"key": key,
	}

	url := fmt.Sprint("http://", c.url, api, rdht)
	_, err := Call(url, "delete", 1, []interface{}{data})
	if err != nil {
		log.Fatalf("Err: %v", err)
		return err
	}
	return nil
}

//func tx_suite(ops ...Operation) {
// TODO
//}

func (c *Client) TxRead(key string) (map[string]interface{}, error) {
	data := []interface{}{
		// {"read": <key>}
		map[string]interface{}{
			"read": key,
		},
		// {"commit": ""}
		commit(),
	}

	url := fmt.Sprint("http://", c.url, api, tx)
	res, err := Call(url, "read", 1, []interface{}{data})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	return res, nil
}

func (c *Client) TxWrite(key string, value interface{}) (map[string]interface{}, error) {
	data := []interface{}{
		// {"write": {<key>: {"type": "as_is" or "as_bin", "value": <value>} } }
		map[string]interface{}{
			"write": map[string]interface{}{
				key: encode_value(value),
			},
		},
		// {"commit": ""}
		commit(),
	}

	url := fmt.Sprint("http://", c.url, api, tx)
	res, err := Call(url, "req_list", 1, []interface{}{data})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	return res, nil
}

func (c *Client) TxTestAndSet(key string, oldValue interface{}, newValue interface{}) (map[string]interface{}, error) {
	data := []interface{}{
		// {"test_and_set": {"key": <key>, "old": <oldValue>, new: <newValue>} }
		map[string]interface{}{
			"test_and_set": map[string]interface{}{
				"key": key,
				"old": encode_value(oldValue),
				"new": encode_value(newValue),
			},
		},
		// {"commit": ""}
		commit(),
	}
	log.Println(data)

	url := fmt.Sprint("http://", c.url, api, tx)
	res, err := Call(url, "req_list", 1, []interface{}{data})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	return res, nil
}

func encode_value(value interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	switch value := value.(type) {
	default:
		res = map[string]interface{}{
			"type":  "as_is",
			"value": value,
		}
	case []byte:
		res = map[string]interface{}{
			"type":  "as_is",
			"value": value,
		}
	}
	return res
}

func commit() map[string]interface{} {
	res := map[string]interface{}{
		"commit": "",
	}
	return res
}

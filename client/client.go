package client

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
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

// TODO more config parameters + TLS
type Client struct {
	Url string
}

const (
	api     = "/api"
	tx      = "/tx.yaws"
	rdht    = "/rdht.yaws"
	dht_raw = "/dht_raw.yaws"
	pubsub  = "/pubsub.yaws"
	monitor = "/monitor.yaws"
)

var (
	url = "localhost:8000"
)

func init() {
	env := os.Getenv("SCALARIS_API_ENDPOINT")
	if env != "" {
		url = env
	}
}

func NewClient(url string) (*Client, error) {
	return &Client{Url: url}, nil
}

// TODO allow to retrieve value as <json_value> for test and set operations
func (c *Client) Read(key string) (map[string]interface{}, error) {
	url := fmt.Sprint("http://", c.Url, api, tx)
	res, err := Call(url, "read", 1, []interface{}{key})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	return res, nil
}

func write() {

}

// FIXME not working
func (c *Client) TestAndSet(key string, oldValue interface{}, newValue interface{}) (map[string]interface{}, error) {
	// {"write": {"key": <key>, "old": <json_value>, "new": <json_value>} }
	data := map[string]interface{}{
		"key": key,
		"old": encode_value(oldValue),
		"new": encode_value(newValue),
	}

	url := fmt.Sprint("http://", c.Url, api, tx)
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

	url := fmt.Sprint("http://", c.Url, api, rdht)
	_, err := Call(url, "delete", 1, []interface{}{data})
	if err != nil {
		log.Fatalf("Err: %v", err)
		return err
	}
	return nil
}

// TODO Add transaction suites
// func (c *Client) TxSuite(ops ...Operation) {
//
// }

func (c *Client) TxRead(key string) (map[string]interface{}, error) {
	data := []interface{}{
		// {"read": <key>}
		map[string]interface{}{
			"read": key,
		},
		// {"commit": ""}
		commit(),
	}

	url := fmt.Sprint("http://", c.Url, api, tx)
	res, err := Call(url, "read", 1, []interface{}{data})
	if err != nil {
		log.Fatalf("Err: ", err)
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

	url := fmt.Sprint("http://", c.Url, api, tx)
	res, err := Call(url, "req_list", 1, []interface{}{data})
	if err != nil {
		log.Fatalf("Err: ", err)
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

	url := fmt.Sprint("http://", c.Url, api, tx)
	res, err := Call(url, "req_list", 1, []interface{}{data})
	if err != nil {
		log.Fatalf("Err: ", err)
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

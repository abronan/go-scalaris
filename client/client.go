package client

import (
	"fmt"
	"os"
	"time"

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
	defaultURL     = "localhost:8000"
	defaultTimeout = 5 * time.Second
)

func init() {
	env := os.Getenv("SCALARIS_JSON_URL")
	if env != "" {
		defaultURL = env
	}
}

func NewClient(url string) (*Client, error) {
	return &Client{Url: url}, nil
}

// Operation represents a single operation
// to be included in a request list, this
// can be done in the case of a transaction
type Operation map[string]interface{}

// TxSuite is used to atomically execute
// a suite of Operations inside the store
// It is used as such:
//
//    resp, err := scalaris.TxSuite {
//        scalaris.TestAndSet(key, newValue, oldValue),
//        scalaris.Write(key, value),
//        scalaris.Delete(key),
//    }
//
// The resp contains the set of responses
// from all the operations inside the transaction
//
// If the operation fails, TxSuite throws the
//
//
func (c *Client) TxSuite(ops ...Operation) (resp interface{}, err error) {
	var tx []interface{}
	for _, op := range ops {
		tx = append(tx, op)
	}
	tx = append(tx, commit())

	url := fmt.Sprint("http://", c.Url, api, tx)
	res, err := Call(url, "req_list", []interface{}{tx})
	if err != nil {
		log.Fatalf("Err: ", err)
		return nil, ErrAbort
	}
	return res, nil
}

// TODO allow to retrieve value as <json_value> for test and set operations
func (c *Client) Read(key string) (map[string]interface{}, error) {
	url := fmt.Sprint("http://", c.Url, api, tx)
	res, err := Call(url, "read", []interface{}{key})
	if err != nil {
		// TODO Parse error, throw ErrNotFound
		log.Fatalf("Err: %v", err)
		return nil, err
	}
	return res, nil
}

func write() {

}

func (c *Client) Delete(key string) (map[string]interface{}, error) {
	// {"write": {"key": <key>, "old": <json_value>, "new": <json_value>} }
	data := map[string]interface{}{
		"key": key,
	}

	url := fmt.Sprint("http://", c.Url, api, rdht)
	res, err := Call(url, "delete", []interface{}{data})
	if err != nil {
		log.Fatalf("Err: %v", err)
		return res, err
	}
	return res, nil
}

func (c *Client) TxRead(key string) (map[string]interface{}, error) {
	data := []interface{}{
		// {"read": <key>}
		map[string]interface{}{
			"read": key,
		},
		// {"commit": ""}
		commit(),
	}
	log.Debug("Built request: ", data)

	url := fmt.Sprint("http://", c.Url, api, tx)
	res, err := Call(url, "read", []interface{}{data})
	if err != nil {
		// TODO Parse error throw ErrNotFound
		log.Fatalf("Err: ", err)
		return nil, err
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
	log.Debug("Built request: ", data)

	url := fmt.Sprint("http://", c.Url, api, tx)
	res, err := Call(url, "req_list", []interface{}{data})
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
	log.Debug("Built request: ", data)

	url := fmt.Sprint("http://", c.Url, api, tx)
	res, err := Call(url, "req_list", []interface{}{data})
	if err != nil {
		log.Fatalf("Err: ", err)
	}
	return res, nil
}

// TestAndSet is used as an Operation inside
// a TxSuite, it atomically puts a value inside
// the store
func (c *Client) TestAndSet(key string, oldValue interface{}, newValue interface{}) map[string]interface{} {
	return map[string]interface{}{
		"test_and_set": map[string]interface{}{
			"key": key,
			"old": encode_value(oldValue),
			"new": encode_value(newValue),
		},
	}
}

// Write is used as an Operation inside
// a TxSuite, it atomically puts a value inside
// the store
func (c *Client) Write(key string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"write": map[string]interface{}{
			key: encode_value(value),
		},
	}
}

// Commit marks the end of a Transaction Suite
func commit() map[string]interface{} {
	return map[string]interface{}{
		"commit": "",
	}
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

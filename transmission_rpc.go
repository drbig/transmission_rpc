package transmission_rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type Client struct {
	Client   *http.Client // HTTP client
	Address  string       // Transmission RPC server URL
	Endpoint string       // Transmission RPC server endpoint path
	login    string
	password string
	auth     bool
	token    string
}

type request struct {
	Arguments interface{} `json:"arguments"`
	Method    string      `json:"method"`
	Tag       int         `json:"tag"`
}

type Response struct {
	Arguments map[string]interface{} `json:"arguments"` // map of response arguments
	Result    string                 `json:"result"`
	Tag       int                    `json:"tag"`
}

const (
	defaultEndpoint = `/transmission/rpc`
	defaultTries    = 3
)

var (
	tag    = 0
	tagMtx sync.Mutex
)

func getTag() int {
	tagMtx.Lock()
	defer tagMtx.Unlock()
	tag++
	if tag < 0 {
		tag = 0
	}
	return tag
}

func NewClient(address string) *Client {
	return &Client{
		Client:   &http.Client{},
		Address:  address,
		Endpoint: defaultEndpoint,
	}
}

func (c *Client) SetAuth(login, password string) {
	c.login = login
	c.password = password
	c.auth = true
}

func (c *Client) RemAuth() {
	c.login = ""
	c.password = ""
	c.auth = false
}

func (c *Client) RequestRaw(request []byte) ([]byte, error) {
	for t := 0; t < defaultTries; t++ {
		req, err := http.NewRequest("POST", c.Address+c.Endpoint, bytes.NewBuffer(request))
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-Transmission-Session-Id", c.token)
		if c.auth {
			req.SetBasicAuth(c.login, c.password)
		}
		res, err := c.Client.Do(req)
		if err != nil {
			return nil, err
		}
		switch res.StatusCode {
		case 409:
			c.token = res.Header.Get("X-Transmission-Session-Id")
		case 200:
			data, err := ioutil.ReadAll(res.Body)
			if err == nil {
				res.Body.Close()
				return data, nil
			}
		}
		res.Body.Close()
	}
	return nil, fmt.Errorf("Gave up after %d tries", defaultTries)
}

func (c *Client) Request(method string, args interface{}) (*Response, error) {
	req := request{
		Method:    method,
		Arguments: args,
		Tag:       getTag(),
	}
	js, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.RequestRaw(js)
	if err != nil {
		return nil, err
	}
	var res Response
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}
	if res.Tag != req.Tag {
		return &res, fmt.Errorf("Tag mismatch (%d != %d)", res.Tag, req.Tag)
	}
	if res.Result != "success" {
		return &res, fmt.Errorf("Unsuccessful response")
	}
	return &res, nil
}

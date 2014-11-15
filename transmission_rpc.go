package transmission_rpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
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

const (
	defaultEndpoint = `/transmission/rpc`
	defaultTries    = 3
)

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

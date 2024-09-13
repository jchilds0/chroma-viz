package hub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Client struct {
	Address string
	Port    int
	Client  http.Client
}

func (c Client) URL() string {
	return "http://" + c.Address + ":" + strconv.Itoa(c.Port)
}

func (c Client) GetJSON(path string, v any) (err error) {
	res, err := c.Client.Get(c.URL() + path)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("Response: %s", res.Status)
		return
	}

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, v)
	return
}

func (c Client) PutJSON(dir string, v any) (err error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(jsonData)

	res, err := c.Client.Post(c.URL()+dir, "", buf)
	if err != nil {
		return
	}
	defer res.Body.Close()
	return
}

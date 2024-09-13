package hub

import (
	"bytes"
	"chroma-viz/library/templates"
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

func (c *Client) URL() string {
	return "http://" + c.Address + ":" + strconv.Itoa(c.Port)
}

func GetTemplates(c Client) (temps []*templates.Template, err error) {
	res, err := c.Client.Get(c.URL() + "/templates")
	if err != nil {
		return
	}

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, &temps)
	if err != nil {
		return
	}

	for _, t := range temps {
		if t == nil {
			continue
		}

		err = t.Init()
		if err != nil {
			return
		}
	}

	return
}

func GetTemplate(c Client, tempID int) (temp templates.Template, err error) {
	url := fmt.Sprintf("%s/template/%d", c.URL(), tempID)

	res, err := c.Client.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("Server returned: %s", res.Status)
		return
	}

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, &temp)
	if err != nil {
		return
	}

	err = temp.Init()
	return
}

func PutTemplate(c Client, temp *templates.Template) (err error) {
	jsonData, err := json.Marshal(temp)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(jsonData)

	res, err := c.Client.Post(c.URL()+"/template", "", buf)
	if err != nil {
		return
	}
	defer res.Body.Close()

	return
}

func GetTemplateIDs(c Client) (tempIDs map[int]string, err error) {
	res, err := c.Client.Get(c.URL() + "/tempIDs")
	if err != nil {
		return
	}
	defer res.Body.Close()

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, &tempIDs)
	return
}

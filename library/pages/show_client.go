package pages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ShowClient struct {
	Client   http.Client
	ShowAddr string
	ShowPort int
}

func (c *ShowClient) url() string {
	return fmt.Sprintf("http://%s:%d", c.ShowAddr, c.ShowPort)
}

func (c *ShowClient) AddPage(page *Page) (err error) {
	jsonData, err := json.Marshal(page)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(jsonData)
	res, err := c.Client.Post(c.url()+"/page", "", buf)
	if err != nil {
		return
	}
	defer res.Body.Close()

	jsonData, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonData, &page.PageNum)
	return err
}

func (c *ShowClient) GetPage(pageNum int) (*Page, bool) {
	path := fmt.Sprintf("%s/page/%d", c.url(), pageNum)

	res, err := c.Client.Get(path)
	if err != nil {
		log.Println("Error getting page:", err)
		return nil, false
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, false
	}

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error getting page:", err)
		return nil, false
	}

	var page Page
	err = json.Unmarshal(jsonData, &page)
	if err != nil {
		log.Println("Error getting page:", err)
		return nil, false
	}

	return &page, true
}

func (c *ShowClient) GetPages() (map[int]PageData, error) {
	res, err := c.Client.Get(c.url() + "/page")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("No pages: %s", res.Status)
	}

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data map[int]PageData
	err = json.Unmarshal(jsonData, &data)
	return data, err
}

func (c *ShowClient) DeletePage(pageNum int) {
	log.Println("Delete page not implemented")
}

func (c *ShowClient) Clear() {
	log.Println("Clear not implemented")
}

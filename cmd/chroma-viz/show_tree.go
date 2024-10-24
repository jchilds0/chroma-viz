package main

import (
	"bufio"
	"chroma-viz/library/pages"
	"encoding/json"
	"net"

	"github.com/gotk3/gotk3/gtk"
)

const (
	PAGENUM = iota
	TITLE
	TEMPLATE_ID
	TEMPLATE_NAME
	NUM_COL
)

var KEYTITLE = map[int]string{
	PAGENUM:       "Page Num",
	TITLE:         "Title",
	TEMPLATE_ID:   "Template ID",
	TEMPLATE_NAME: "Template Name",
}

type ShowTree interface {
	TreeView() *gtk.TreeView
	SelectedPage() (int, error)
	AddPage(page pages.Page) (err error)
	GetPage(pageNum int) (*pages.Page, bool)
	GetPages() map[int]PageData
	DeletePage(pageNum int)
	Clear()
}

type PageData struct {
	PageNum int
	Title   string
	TempID  int
	Layer   int
}

const (
	NO_MESSAGE = iota
	CREATE_PAGE
	READ_PAGE
	UPDATE_PAGE
	DELETE_PAGE
	GET_PAGES
	RECIEVE_UPDATES
)

type Message struct {
	Type     int
	PageInfo PageData
	PageData map[int]PageData
	Page     pages.Page
}

func sendMessage(conn net.Conn, m Message) (err error) {
	buf, err := json.Marshal(m)
	if err != nil {
		return
	}

	_, err = conn.Write(buf)
	if err != nil {
		return
	}

	_, err = conn.Write([]byte{0})
	return
}

func recvMessage(conn net.Conn) (m Message, err error) {
	buf, err := bufio.NewReader(conn).ReadBytes(0)
	if err != nil {
		return
	}

	err = json.Unmarshal(buf[:len(buf)-1], &m)
	return
}

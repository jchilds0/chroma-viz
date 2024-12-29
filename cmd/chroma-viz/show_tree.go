package main

import (
	"bufio"
	"chroma-viz/library/pages"
	"chroma-viz/library/util"
	"encoding/json"
	"log"
	"net"

	"github.com/gotk3/gotk3/gtk"
)

const (
	PAGENUM = iota
	TITLE
	TEMPLATE_ID
	TEMPLATE_NAME
	TEMPLATE_LAYER
	NUM_COL
)

var KEYTITLE = map[int]string{
	PAGENUM:        "Page Num",
	TITLE:          "Title",
	TEMPLATE_ID:    "Template ID",
	TEMPLATE_NAME:  "Template Name",
	TEMPLATE_LAYER: "Template Layer",
}

func createShowTreeModel(updateTitle func(title string, pageNum int)) (*gtk.TreeView, error) {
	treeView, err := gtk.TreeViewNew()
	if err != nil {
		return nil, err
	}

	treeView.SetReorderable(true)

	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		return nil, err
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(KEYTITLE[PAGENUM], cell, "text", PAGENUM)
	if err != nil {
		return nil, err
	}

	column.SetResizable(true)
	treeView.AppendColumn(column)

	title, err := gtk.CellRendererTextNew()
	if err != nil {
		return nil, err
	}

	title.SetProperty("editable", true)
	title.Connect("edited",
		func(cell *gtk.CellRendererText, path string, text string) {
			list, err := treeView.GetModel()
			if err != nil {
				log.Println("Error editing page:", err)
				return
			}

			model := list.ToTreeModel()
			iter, err := model.GetIterFromString(path)
			if err != nil {
				log.Println("Error editing page:", err)
				return
			}

			pageNum, err := util.ModelGetValue[int](model, iter, PAGENUM)
			if err != nil {
				log.Println("Error editing page:", err)
				return
			}

			updateTitle(text, pageNum)
		})

	column, err = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[TITLE], title, "text", TITLE)
	if err != nil {
		return nil, err
	}

	column.SetExpand(true)
	column.SetResizable(true)
	treeView.AppendColumn(column)

	column, err = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[TEMPLATE_ID], cell, "text", TEMPLATE_ID)
	if err != nil {
		return nil, err
	}

	column.SetResizable(true)
	treeView.AppendColumn(column)

	column, err = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[TEMPLATE_NAME], cell, "text", TEMPLATE_NAME)
	if err != nil {
		return nil, err
	}

	column.SetResizable(true)
	treeView.AppendColumn(column)

	column, err = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[TEMPLATE_LAYER], cell, "text", TEMPLATE_LAYER)
	if err != nil {
		return nil, err
	}

	column.SetResizable(true)
	treeView.AppendColumn(column)
	return treeView, nil
}

type ShowTree interface {
	TreeView() *gtk.TreeView
	SelectedPage() (int, error)
	WritePage(page pages.Page) (err error)
	ReadPage(pageNum int) (*pages.Page, bool)
	GetPages() map[int]PageData
	DeletePage(pageNum int)
	Clear()
}

func NextPageNum(showTree ShowTree) (pageNum int) {
	pageNum = 1

	for _, p := range showTree.GetPages() {
		pageNum = max(p.PageNum+1, pageNum)
	}

	return
}

type PageData struct {
	PageNum  int
	Title    string
	TempID   int
	TempName string
	Layer    int
}

const (
	NO_MESSAGE = iota
	WRITE_PAGE
	READ_PAGE
	UPDATE_PAGE_INFO
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

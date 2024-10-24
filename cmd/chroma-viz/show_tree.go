package main

import (
	"chroma-viz/library/pages"

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
	GetPages() map[int]pages.Page
	DeletePage(pageNum int)
	Clear()
}

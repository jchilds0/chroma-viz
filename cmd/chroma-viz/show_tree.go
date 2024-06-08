package main

import (
	"chroma-viz/library"
	"chroma-viz/library/attribute"
	"chroma-viz/library/pages"
	"chroma-viz/library/props"
	"chroma-viz/library/util"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
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

type ShowTree struct {
	treeView *gtk.TreeView
	treeList *gtk.ListStore
	show     *pages.Show
	columns  map[int]bool
}

func NewShowTree(pageToEditor func(*pages.Page)) *ShowTree {
	var err error
	showTree := &ShowTree{}

	showTree.treeView, err = gtk.TreeViewNew()
	if err != nil {
		log.Fatalf("Error creating show (%s)", err)
	}

	showTree.show = pages.NewShow()
	showTree.treeView.SetReorderable(true)
	showTree.columns = make(map[int]bool, NUM_COL)

	showTree.columns[PAGENUM] = true
	showTree.columns[TITLE] = true
	showTree.columns[TEMPLATE_ID] = true
	showTree.columns[TEMPLATE_NAME] = true

	// create tree columns
	var column *gtk.TreeViewColumn
	for key := 0; key < NUM_COL; key++ {
		if !showTree.columns[key] {
			continue
		}

		switch key {
		case TITLE:
			title, err := gtk.CellRendererTextNew()
			if err != nil {
				log.Fatalf("Error creating show (%s)", err)
			}

			title.SetProperty("editable", true)
			title.Connect("edited",
				func(cell *gtk.CellRendererText, path string, text string) {
					iter, err := showTree.treeList.GetIterFromString(path)
					if err != nil {
						log.Printf("Error editing page (%s)", err)
						return
					}

					model := &showTree.treeList.TreeModel
					pageNum, err := util.ModelGetValue[int](model, iter, PAGENUM)
					if err != nil {
						log.Printf("Error editing page (%s)", err)
						return
					}

					if _, ok := showTree.show.Pages[pageNum]; !ok {
						log.Print("Error getting page")
						return
					}

					showTree.show.Pages[pageNum].Title = text
					showTree.treeList.SetValue(iter, TITLE, text)
				})

			column, err = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[key], title, "text", 1)
			if err != nil {
				log.Fatalf("Error creating show (%s)", err)
			}

			column.SetExpand(true)

		default:
			cell, err := gtk.CellRendererTextNew()
			if err != nil {
				log.Fatalf("Error creating show (%s)", err)
			}

			column, err = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[key], cell, "text", key)
			if err != nil {
				log.Fatalf("Error creating show (%s)", err)
			}

		}

		column.SetResizable(true)
		showTree.treeView.AppendColumn(column)
	}

	showTree.treeList, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_INT, glib.TYPE_STRING)
	if err != nil {
		log.Fatalf("Error creating show (%s)", err)
	}

	showTree.treeView.SetModel(showTree.treeList)

	// send page to editor on double click
	showTree.treeView.Connect("row-activated",
		func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
			iter, err := showTree.treeList.GetIter(path)
			if err != nil {
				log.Printf("Error sending page to editor (%s)", err)
				return
			}

			model := &showTree.treeList.TreeModel
			pageNum, err := util.ModelGetValue[int](model, iter, PAGENUM)
			if err != nil {
				log.Printf("Error editing page (%s)", err)
				return
			}

			pageToEditor(showTree.show.Pages[pageNum])
			SendPreview(showTree.show.Pages[pageNum], library.ANIMATE_ON)
		})

	return showTree
}

func (showTree *ShowTree) ImportPage(page *pages.Page) (err error) {
	if page == nil {
		err = fmt.Errorf("Attempted to import nil page")
		return
	}

	showTree.show.AddPage(page)

	for _, prop := range page.PropMap {
		if prop.PropType != props.CLOCK_PROP {
			continue
		}

		/*
		   Clock requires a way to send updates to viz
		   to animate the clock. We manually add this
		   after parsing the page.
		*/
		attr, ok := prop.Attr["string"]
		if !ok {
			continue
		}

		clockAttr, ok := attr.(*attribute.ClockAttribute)
		if !ok {
			continue
		}

		clockAttr.SetClock(func() { SendEngine(page, library.CONTINUE) })
	}

	iter := showTree.treeList.Append()
	err = showTree.treeList.SetValue(iter, PAGENUM, page.PageNum)
	if err != nil {
		return
	}

	err = showTree.treeList.SetValue(iter, TITLE, page.Title)
	if err != nil {
		return
	}

	err = showTree.treeList.SetValue(iter, TEMPLATE_ID, page.TemplateID)
	if err != nil {
		return
	}

	err = showTree.treeList.SetValue(iter, TEMPLATE_NAME, page.Title)
	return
}

func (showTree *ShowTree) ImportShow(filename string) {
	var show pages.Show
	err := show.ImportShow(filename)
	if err != nil {
		log.Print(err)
	}

	for _, page := range show.Pages {
		err = showTree.ImportPage(page)
		if err != nil {
			log.Print(err)
		}
	}
}

func (ShowTree *ShowTree) Clean() {
	ShowTree.treeList.Clear()
	ShowTree.show.Pages = make(map[int]*pages.Page, 10)
	ShowTree.show.NumPages = 1
}

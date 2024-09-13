package main

import (
	"chroma-viz/library/hub"
	"chroma-viz/library/util"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	PREVIEW = iota
	TOP_LEFT
	LOWER_FRAME
	TICKER
)

type TempTree struct {
	treeView     *gtk.TreeView
	treeList     *gtk.ListStore
	sendTemplate func(int)
}

func NewTempTree(templateToShow func(int)) *TempTree {
	var err error
	temp := &TempTree{sendTemplate: templateToShow}

	temp.treeView, err = gtk.TreeViewNew()
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	// create tree columns
	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute("Name", cell, "text", 0)
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	temp.treeView.AppendColumn(column)
	column, err = gtk.TreeViewColumnNewWithAttribute("Template ID", cell, "text", 1)
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	temp.treeView.AppendColumn(column)

	temp.treeList, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_INT)
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	temp.treeView.SetModel(temp.treeList)

	// send template to show on double click
	temp.treeView.Connect("row-activated",
		func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
			iter, err := temp.treeList.GetIter(path)
			if err != nil {
				log.Fatalf("Error sending template to show (%s)", err)
			}

			model := &temp.treeList.TreeModel
			tempID, err := util.ModelGetValue[int](model, iter, 1)
			if err != nil {
				log.Fatalf("Error sending template to show (%s)", err)
			}

			temp.sendTemplate(tempID)
		})

	return temp
}

func (temp *TempTree) ImportTemplates(c hub.Client) {
	tempids, err := hub.GetTemplateIDs(c)
	if err != nil {
		log.Printf("Error importing templates: %s", err)
		return
	}

	for id, title := range tempids {
		iter := temp.treeList.Append()
		err := temp.treeList.SetValue(iter, 0, title)
		if err != nil {
			log.Printf("Error importing templates: %s", err)
			return
		}

		err = temp.treeList.SetValue(iter, 1, id)
		if err != nil {
			log.Printf("Error importing templates: %s", err)
			return
		}
	}
}

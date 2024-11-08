package main

import (
	"chroma-viz/library/hub"
	"chroma-viz/library/util"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	TEMP_NAME = iota
	TEMP_ID
	TEMP_LAYER
)

type TempTree struct {
	treeView     *gtk.TreeView
	treeList     *gtk.ListStore
	sendTemplate func(int)
}

func NewTempTree(templateToShow func(int)) (*TempTree, error) {
	var err error
	temp := &TempTree{sendTemplate: templateToShow}

	temp.treeView, err = gtk.TreeViewNew()
	if err != nil {
		return nil, err
	}

	titles := []string{"Name", "Template ID", "Layer"}
	for i, name := range titles {
		cell, err := gtk.CellRendererTextNew()
		if err != nil {
			return nil, err
		}

		column, err := gtk.TreeViewColumnNewWithAttribute(name, cell, "text", i)
		if err != nil {
			return nil, err
		}

		temp.treeView.AppendColumn(column)

		if i == TEMP_ID {
			column.SetSortIndicator(true)
			column.SetSortColumnID(1)
		}
	}

	temp.treeList, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_INT, glib.TYPE_INT)
	if err != nil {
		return nil, err
	}

	temp.treeView.SetModel(temp.treeList)

	// send template to show on double click
	temp.treeView.Connect("row-activated",
		func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
			iter, err := temp.treeList.GetIter(path)
			if err != nil {
				log.Println("Error sending template to show:", err)
			}

			model := &temp.treeList.TreeModel
			tempID, err := util.ModelGetValue[int](model, iter, 1)
			if err != nil {
				log.Println("Error sending template to show:", err)
			}

			temp.sendTemplate(tempID)
		})

	return temp, nil
}

func (temp *TempTree) ImportTemplates(c hub.Client) {
	var temps []hub.TemplateHeader

	err := c.GetJSON("/template/list", &temps)
	if err != nil {
		log.Printf("Error importing templates: %s", err)
		return
	}

	for _, header := range temps {
		iter := temp.treeList.Append()
		err := temp.treeList.SetValue(iter, TEMP_NAME, header.Title)
		if err != nil {
			log.Printf("Error importing templates: %s", err)
			continue
		}

		err = temp.treeList.SetValue(iter, TEMP_ID, header.TemplateID)
		if err != nil {
			log.Printf("Error importing templates: %s", err)
			continue
		}

		err = temp.treeList.SetValue(iter, TEMP_LAYER, header.Layer)
		if err != nil {
			log.Printf("Error importing templates: %s", err)
			continue
		}
	}
}

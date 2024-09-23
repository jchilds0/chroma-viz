package main

import (
	"chroma-viz/library/hub"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	PREVIEW = iota
	TOP_LEFT
	LOWER_FRAME
	TICKER
)

type TemplateChooserDialog struct {
	*gtk.Dialog
	treeView *gtk.TreeView
	treeList *gtk.ListStore
}

func NewTemplateChooserDialog(win *gtk.Window) (dialog *TemplateChooserDialog, err error) {
	dialog = &TemplateChooserDialog{}
	dialog.Dialog, err = gtk.DialogNewWithButtons(
		"Import Template", win, gtk.DIALOG_MODAL,
		[]interface{}{"_Close", gtk.RESPONSE_REJECT},
		[]interface{}{"_Open", gtk.RESPONSE_ACCEPT},
	)

	if err != nil {
		return
	}

	dialog.SetResizable(true)
	dialog.SetDecorated(true)
	dialog.SetDefaultSize(600, 300)

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	dialogContent, err := dialog.GetContentArea()
	if err != nil {
		return
	}

	dialogContent.PackStart(box, true, true, 10)

	scroll, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return
	}

	box.PackStart(scroll, true, true, 0)

	dialog.treeView, err = gtk.TreeViewNew()
	if err != nil {
		return
	}

	scroll.Add(dialog.treeView)

	// create tree columns
	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	column, err := gtk.TreeViewColumnNewWithAttribute("Name", cell, "text", 0)
	if err != nil {
		return
	}

	dialog.treeView.AppendColumn(column)
	column, err = gtk.TreeViewColumnNewWithAttribute("Template ID", cell, "text", 1)
	if err != nil {
		return
	}

	dialog.treeView.AppendColumn(column)

	dialog.treeList, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_INT)
	if err != nil {
		return
	}

	dialog.treeView.SetModel(dialog.treeList)
	box.ShowAll()

	return
}

func (dialog *TemplateChooserDialog) ImportTemplates(c hub.Client) (err error) {
	dialog.treeList.Clear()

	var tempids map[int]string

	err = c.GetJSON("/template/list", &tempids)
	if err != nil {
		return
	}

	for id, title := range tempids {
		iter := dialog.treeList.Append()
		err = dialog.treeList.SetValue(iter, 0, title)
		if err != nil {
			return
		}

		err = dialog.treeList.SetValue(iter, 1, id)
		if err != nil {
			return
		}
	}

	return
}

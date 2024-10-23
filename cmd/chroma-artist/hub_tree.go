package main

import (
	"chroma-viz/library/hub"
	"chroma-viz/library/util"
	"fmt"

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
	dialog.treeView.Connect("row-activated", func() {
		dialog.Dialog.Emit("response", glib.TYPE_INT, gtk.RESPONSE_ACCEPT)
	})

	// create tree columns
	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	column, err := gtk.TreeViewColumnNewWithAttribute("Name", cell, "text", 0)
	if err != nil {
		return
	}

	column.SetExpand(true)
	dialog.treeView.AppendColumn(column)

	column, err = gtk.TreeViewColumnNewWithAttribute("Template ID", cell, "text", 1)
	if err != nil {
		return
	}

	column.SetSortIndicator(true)
	column.SetSortColumnID(1)
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

func (dialog *TemplateChooserDialog) SelectedTemplateID() (id int, err error) {
	selection, err := dialog.treeView.GetSelection()
	if err != nil {
		return
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		return 0, fmt.Errorf("No template selected")
	}

	id, err = util.ModelGetValue[int](dialog.treeList.ToTreeModel(), iter, 1)
	return
}

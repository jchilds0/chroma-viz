package main

import (
	"chroma-viz/library/hub"
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

type TemplateChooserDialog struct {
	*gtk.Dialog
	treeView *gtk.TreeView
	treeList *gtk.ListStore
}

func NewTemplateChooserDialog(win *gtk.Window) *TemplateChooserDialog {
	var err error
	dialog := &TemplateChooserDialog{}
	dialog.Dialog, err = gtk.DialogNewWithButtons(
		"Import Template", win, gtk.DIALOG_MODAL,
		[]interface{}{"_Close", gtk.RESPONSE_REJECT},
		[]interface{}{"_Open", gtk.RESPONSE_ACCEPT},
	)

	if err != nil {
		log.Fatal(err)
	}

	dialog.SetResizable(true)
	dialog.SetDecorated(true)
	dialog.SetDefaultSize(600, 300)

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal(err)
	}

	dialogContent, err := dialog.GetContentArea()
	if err != nil {
		log.Fatal(err)
	}

	dialogContent.PackStart(box, true, true, 10)

	scroll, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	box.PackStart(scroll, true, true, 0)

	dialog.treeView, err = gtk.TreeViewNew()
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	scroll.Add(dialog.treeView)

	// create tree columns
	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute("Name", cell, "text", 0)
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	dialog.treeView.AppendColumn(column)
	column, err = gtk.TreeViewColumnNewWithAttribute("Template ID", cell, "text", 1)
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	dialog.treeView.AppendColumn(column)

	dialog.treeList, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_INT)
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	dialog.treeView.SetModel(dialog.treeList)
	box.ShowAll()

	return dialog
}

func (dialog *TemplateChooserDialog) ImportTemplates(c hub.Client) {
	dialog.treeList.Clear()

	tempids, err := hub.GetTemplateIDs(c)
	if err != nil {
		log.Printf("Error importing templates: %s", err)
		return
	}

	for id, title := range tempids {
		iter := dialog.treeList.Append()
		err := dialog.treeList.SetValue(iter, 0, title)
		if err != nil {
			log.Printf("Error importing templates: %s", err)
			return
		}

		err = dialog.treeList.SetValue(iter, 1, id)
		if err != nil {
			log.Printf("Error importing templates: %s", err)
			return
		}
	}

}

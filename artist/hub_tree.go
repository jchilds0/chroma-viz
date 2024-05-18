package artist

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

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

	dialogContent.PackStart(box, true, true, 0)

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

func (dialog *TemplateChooserDialog) ImportTemplates(hub net.Conn) {
	if hub == nil {
		log.Print("Chroma Hub is disconnected")
		return
	}

	dialog.treeList.Clear()

	s := fmt.Sprintf("ver 0 1 tempids;")
	hub.Write([]byte(s))

	buf := bufio.NewReader(hub)
	for {
		s, err := buf.ReadString(';')
		if err != nil {
			return
		}

		s = strings.TrimSuffix(s, ";")
		if s == "EOF" {
			return
		}

		data := strings.Split(s, " ")

		tempID, err := strconv.Atoi(data[0])
		if err != nil {
			log.Printf("Error reading template id (%s)", err)
			continue
		}
		title := strings.Join(data[1:], " ")

		err = dialog.treeList.Set(
			dialog.treeList.Append(),
			[]int{0, 1},
			[]interface{}{title, tempID},
		)
		if err != nil {
			log.Printf("Error adding template to gtk treestore (%s)", err)
			return
		}
	}
}

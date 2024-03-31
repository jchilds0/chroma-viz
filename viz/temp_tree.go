package viz

import (
	"chroma-viz/library/gtk_utils"
	"chroma-viz/library/templates"
	"log"
	"net"

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
	Temps        *templates.Temps
	sendTemplate func(*templates.Template)
}

func NewTempTree(templateToShow func(*templates.Template)) *TempTree {
	var err error
	temp := &TempTree{sendTemplate: templateToShow}

	temp.treeView, err = gtk.TreeViewNew()
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	temp.Temps = templates.NewTemps()

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
			tempID, err := gtk_utils.ModelGetValue[int](model, iter, 1)
			if err != nil {
				log.Fatalf("Error sending template to show (%s)", err)
			}

			temp.sendTemplate(temp.Temps.Temps[tempID])
		})

	return temp
}

func (temp *TempTree) AddTemplate(template *templates.Template) error {
	err := temp.treeList.Set(
		temp.treeList.Append(),
		[]int{0, 1},
		[]interface{}{template.Title, template.TempID},
	)

	return err
}

func (temp *TempTree) ImportTemplates(hub net.Conn) {
	err := temp.Temps.ImportTemplates(hub)
	if err != nil {
		log.Printf("Error importing hub (%s)", err)
	}

	for _, template := range temp.Temps.Temps {
		//log.Printf("Imported Template %d (%s)", template.TempID, template.Title)
		temp.AddTemplate(template)
	}
}

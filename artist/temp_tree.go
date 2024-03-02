package artist

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
    ICON = iota
    NAME
    VISIBLE
    PROP_NUM 
    NUM_COLS
)

type TempTree struct {
    model *gtk.TreeStore
    view *gtk.TreeView
}

func NewTempTree(propToEditor func(propID int)) (*TempTree, error) {
    var err error
    temp := &TempTree{}

    temp.view, err = gtk.TreeViewNew()
    if err != nil {
        return nil, err
    }

    temp.view.Set("reorderable", true)

    cell, err := gtk.CellRendererTextNew()
    if err != nil {
        return nil, err
    }

    column, err := gtk.TreeViewColumnNewWithAttribute("Name", cell, "text", NAME)
    if err != nil {
        return nil, err
    }
    temp.view.AppendColumn(column)

    temp.model, err = gtk.TreeStoreNew(glib.TYPE_OBJECT, glib.TYPE_STRING, glib.TYPE_OBJECT, glib.TYPE_INT)
    if err != nil {
        return nil, err
    }

    temp.view.SetModel(temp.model)

    temp.view.Connect("row-activated", 
        func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
            iter, err := temp.model.GetIter(path)
            if err != nil {
                log.Fatalf("Error sending page to editor (%s)", err)
            }

            id, err := temp.model.GetValue(iter, PROP_NUM)
            if err != nil {
                log.Fatalf("Error sending prop to editor (%s)", err)
            }

            val, err := id.GoValue()
            if err != nil {
                log.Fatalf("Error sending prop to editor (%s)", err)
            }

            propID, ok := val.(int)
            if !ok {
                log.Fatalf("Error sending prop to editor (value not int)")
            }

            propToEditor(propID)
    })

    return temp, nil
}

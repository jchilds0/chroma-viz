package artist

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
    ICON = iota
    NAME
    VISIBLE
    NUM_COLS
)

type TempTree struct {
    model *gtk.TreeStore
    view *gtk.TreeView
}

func NewTempTree() (*TempTree, error) {
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

    temp.model, err = gtk.TreeStoreNew(glib.TYPE_OBJECT, glib.TYPE_STRING, glib.TYPE_OBJECT)
    if err != nil {
        return nil, err
    }

    temp.view.SetModel(temp.model)

    return temp, nil
}

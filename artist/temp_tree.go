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

func NewTempTree(propToEditor func(propID int), updateParent func(propId, parentID int)) (*TempTree, error) {
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

    temp.view.Connect("cursor-changed", 
        func(tree *gtk.TreeView) { 
            selection, err := tree.GetSelection()
            if err != nil {
                log.Fatalf("Error sending prop to editor (%s)", err)
            }

            _, iter, ok := selection.GetSelected()
            if !ok {
                log.Fatalf("Error sending prop to editor (%s)", err)
            }

            propID := getIntFromModel(temp.model, iter, PROP_NUM)

            var parentID int
            var parent gtk.TreeIter
            ok = temp.model.IterParent(&parent, iter)
            if ok {
                parentID = getIntFromModel(temp.model, &parent, PROP_NUM)
            }

            updateParent(propID, parentID)
    })

    temp.view.Connect("row-activated", 
        func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
            iter, err := temp.model.GetIter(path)
            if err != nil {
                log.Fatalf("Error sending prop to editor (%s)", err)
            }

            propID := getIntFromModel(temp.model, iter, PROP_NUM)
            propToEditor(propID)
    })

    return temp, nil
}

func getIntFromModel(model *gtk.TreeStore, iter *gtk.TreeIter, col int) int {
    id, err := model.GetValue(iter, col)
    if err != nil {
        log.Fatalf("Error sending prop to editor (%s)", err)
    }

    val, err := id.GoValue()
    if err != nil {
        log.Fatalf("Error sending prop to editor (%s)", err)
    }

    integer, ok := val.(int)
    if !ok {
        log.Fatalf("Error sending prop to editor (value not int)")
    }

    return integer
}

package artist

import (
	"chroma-viz/library/gtk_utils"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
    PROP_TYPE = iota
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

    typeCell, err := gtk.CellRendererTextNew()
    if err != nil {
        return nil, err
    }

    column, err := gtk.TreeViewColumnNewWithAttribute("Geometry", typeCell, "text", PROP_TYPE)
    if err != nil {
        return nil, err
    }
    temp.view.AppendColumn(column)

    nameCell, err := gtk.CellRendererTextNew()
    if err != nil {
        return nil, err
    }

    nameCell.SetProperty("editable", true)
    nameCell.Connect("edited", func(cell *gtk.CellRendererText, path, text string) {
        iter, err := temp.model.GetIterFromString(path)
        if err != nil {
            log.Printf("Error editing geometry (%s)", err)
            return
        }

        model := &temp.model.TreeModel
        geoID, err := gtk_utils.ModelGetValue[int](model, iter, PROP_NUM)
        if err != nil {
            log.Printf("Error editing geometry (%s)", err)
            return
        }

        geo := template.Geometry[geoID]
        if geo == nil { 
            log.Printf("Error getting geometry %d", geoID)
            return
        }

        geo.Name = text
        temp.model.SetValue(iter, NAME, text)
    })

    column, err = gtk.TreeViewColumnNewWithAttribute("Name", nameCell, "text", NAME)
    if err != nil {
        return nil, err
    }
    temp.view.AppendColumn(column)

    temp.model, err = gtk.TreeStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_OBJECT, glib.TYPE_INT)
    if err != nil {
        return nil, err
    }

    temp.view.SetModel(temp.model)

    temp.view.Connect("row-activated", 
        func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
            iter, err := temp.model.GetIter(path)
            if err != nil {
                log.Printf("Error sending prop to editor (%s)", err)
                return
            }

            model := &temp.model.TreeModel
            propID, err := gtk_utils.ModelGetValue[int](model, iter, PROP_NUM)
            if err != nil {
                log.Printf("Error sending prop to editor (%s)", err)
                return
            }

            propToEditor(propID)
    })

    return temp, nil
}

func (tempView *TempTree) AddRow(iter *gtk.TreeIter, name, propName string, propNum int) {
    tempView.model.SetValue(iter, PROP_TYPE, propName)
    tempView.model.SetValue(iter, NAME, name)
    tempView.model.SetValue(iter, PROP_NUM, propNum)
}

func (tempView *TempTree) Clean() {
    var err error
    tempView.model, err = gtk.TreeStoreNew(
        glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_OBJECT, glib.TYPE_INT)
    if err != nil {
        log.Print(err)
        return
    }

    tempView.view.SetModel(tempView.model)
}

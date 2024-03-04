package attribute

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type graphCell struct {
    *gtk.CellRendererText
    columnNum int
}

func NewGraphCell(i int) *graphCell {
    cell, err := gtk.CellRendererTextNew()
    if err != nil {
        log.Fatalf("Error creating graph cell (%s)", err)
    }

    gCell := &graphCell{CellRendererText: cell, columnNum: i}
    return gCell
}

type ListAttribute struct {
    name          string
    listStore     *gtk.ListStore
}

func NewListAttribute(name string) *ListAttribute {
    var err error
    list := &ListAttribute{name: name}
    list.listStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_INT)
    if err != nil {
        log.Fatalf("Error creating list store")
    }

    return list
}

func (listAttr *ListAttribute) String() (s string) {
    // currently chroma_viz allocates 100 nodes for each graph statically
    s = "num_node=0#" 

    iter, ok := listAttr.listStore.GetIterFirst()
    i := 0 
    for ok {
        // TODO: generic columns
        xVal := getIntFromIter(listAttr.listStore, iter, 0)
        yVal := getIntFromIter(listAttr.listStore, iter, 1)
        
        s = s + fmt.Sprintf("graph_node=%d %d %d#", i, xVal, yVal)
        ok = listAttr.listStore.IterNext(iter)
        i++
    }

    return
}

func getIntFromIter(model *gtk.ListStore, iter *gtk.TreeIter, col int) int {
    row, err := model.GetValue(iter, col)
    if err != nil {
        log.Printf("Error getting graph row (%s)", err)
        return 0
    }

    rowVal, err := row.GoValue()
    if err != nil {
        log.Printf("Error converting row to go val (%s)", err)
        return 0
    }
        
    return rowVal.(int)
}

func (listAttr *ListAttribute) Encode() (s string) {
    iter, ok := listAttr.listStore.GetIterFirst()
    i := 0 
    for ok {
        // TODO: generic columns
        xVal := getIntFromIter(listAttr.listStore, iter, 0)
        yVal := getIntFromIter(listAttr.listStore, iter, 1)
        
        s = s + fmt.Sprintf("node %d %d;", xVal, yVal)
        ok = listAttr.listStore.IterNext(iter)
        i++
    }

    return
}

func (listAttr *ListAttribute) Decode(s string) error {
    line := strings.Split(s, " ")
    if len(line) != 3 {
        return fmt.Errorf("Incorrect list attr string (%s)", line)
    }

    x, err := strconv.Atoi(line[1])
    if err != nil { 
        return err
    }

    y, err := strconv.Atoi(line[2])
    if err != nil { 
        return err
    }

    listAttr.listStore.Set(listAttr.listStore.Append(), []int{0, 1}, []interface{}{x, y})
    return nil
}

func (listAttr *ListAttribute) Update(edit Editor) error {
    _, ok := edit.(*ListEditor) 
    if !ok {
        return fmt.Errorf("ListAttribute.Update requires a ListEditor")
    }

    /*
        No changes required. ListEdit has a pointer to the 
        list store in ListAttribute, which is updated
    */
    return nil
}

type ListEditor struct {
    name        string
    box         *gtk.Box
    treeView    *gtk.TreeView
    listStore   *gtk.ListStore
}

func NewListEditor(name string, columns []string, animate func()) (listEdit *ListEditor, err error) {
    listEdit = &ListEditor{name: name}
    listEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        return
    }

    listEdit.treeView, err = gtk.TreeViewNew()
    if err != nil {
        return 
    }

    listEdit.treeView.SetVisible(true)

    for i, name := range columns {
        gCell := NewGraphCell(i)
        gCell.SetProperty("editable", true)
        gCell.Connect("edited", 
            func(cell *gtk.CellRendererText, path string, text string) {
                if listEdit.listStore == nil {
                    log.Printf("Error editing list attribute (listStore missing)")
                    return
                }

                iter, err := listEdit.listStore.ToTreeModel().GetIterFromString(path)
                if err != nil {
                    log.Printf("Error editing list attribute (%s)", err)
                    return
                }

                id_val, err := strconv.Atoi(text)
                if err != nil {
                    log.Printf("Error editing list attribute (%s)", err)
                    return
                }

                listEdit.listStore.SetValue(iter, gCell.columnNum, id_val)
                animate()
        })
        column, err := gtk.TreeViewColumnNewWithAttribute(name, gCell, "text", i)
        if err != nil {
            log.Printf("Error creating list column (%s)", err)
        }

        listEdit.treeView.AppendColumn(column)
    }

    frame, err := gtk.FrameNew(name)
    if err != nil {
        return
    }

    frame.Set("border-width", 2 * padding)
    frame.Add(listEdit.treeView)
    listEdit.box.PackStart(frame, true, true, 0)
    
    label, err := gtk.LabelNew("Data Rows")
    if err != nil {
        return 
    }

    label.SetVisible(true)
    listEdit.box.PackStart(label, false, false, padding)

    // add rows
    button, err := gtk.ButtonNewWithLabel("+")
    if err != nil {
        return
    }

    button.Connect("clicked", func() { 
        if listEdit.listStore == nil {
            log.Printf("Graph prop editor does not have a list store")
            return
        }

        listEdit.listStore.Append()
    })

    button.SetVisible(true)
    listEdit.box.PackStart(button, false, false, padding)

    // remove rows
    button, err = gtk.ButtonNewWithLabel("-")
    if err != nil {
        log.Printf("Error creating graph table (%s)", err)
    }

    button.Connect("clicked", func() {
        if listEdit.listStore == nil {
            log.Printf("Graph prop editor does not have a list store")
            return
        }

        selection, err := listEdit.treeView.GetSelection()
        if err != nil {
            log.Printf("Error getting current row (%s)", err)
            return
        }

        _, iter, ok := selection.GetSelected()
        if !ok {
            log.Printf("Error getting selected")
            return
        }

        listEdit.listStore.Remove(iter)
    })

    button.SetVisible(true)
    listEdit.box.PackStart(button, false, false, padding)

    return
}

func (listEdit *ListEditor) Update(attr Attribute) error {
    listAttr, ok := attr.(*ListAttribute)
    if !ok {
        return fmt.Errorf("ListEditor.Update requires a ListAttribute")
    }

    listEdit.listStore = listAttr.listStore
    listEdit.treeView.SetModel(listEdit.listStore)
    return nil
}

func (listEdit *ListEditor) Box() *gtk.Box {
    return listEdit.box
}

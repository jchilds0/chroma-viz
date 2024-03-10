package attribute

import (
	"encoding/json"
	"fmt"
	"log"

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
    FileName      string
    ChromaName    string
    Type          int
    NumCols       int
    Selected      bool
    selectedIter  *gtk.TreeIter
    ListStore     *gtk.ListStore
}

func NewListAttribute(file, chroma string, numCols int, selected bool) *ListAttribute {
    var err error
    list := &ListAttribute{
        FileName: file, 
        ChromaName: chroma,
        Type: LIST,
        NumCols: numCols, 
        Selected: selected,
    }

    cols := make([]glib.Type, list.NumCols)
    for i := range cols {
        cols[i] = glib.TYPE_STRING
    }

    list.ListStore, err = gtk.ListStoreNew(cols...)
    if err != nil {
        log.Fatalf("Error creating list store")
    }

    return list
}

func (listAttr *ListAttribute) String() (s string) {
    // currently chroma_engine allocates 100 nodes for each list statically
    if listAttr.Selected {
        // send only the currently selected item from the list
        if listAttr.selectedIter == nil {
            return 
        }

        return listAttr.stringRow(listAttr.selectedIter)
    }

    iter, ok := listAttr.ListStore.GetIterFirst()
    i := 0 
    for ok {
        s = s + listAttr.stringRow(iter)
        ok = listAttr.ListStore.IterNext(iter)
        i++
    }

    s = fmt.Sprintf("num_node=%d#", i) + s
    return
}

func (listAttr *ListAttribute) stringRow(iter *gtk.TreeIter) (s string) {
    s = listAttr.ChromaName + "="
    for j := 0; j < listAttr.NumCols - 1; j++ {
        item := getStringFromIter(listAttr.ListStore, iter, j)
        s = s + item + " "
    }

    item := getStringFromIter(listAttr.ListStore, iter, listAttr.NumCols - 1)
    s = s + item + "#"
    return
}

func getStringFromIter(model *gtk.ListStore, iter *gtk.TreeIter, col int) string {
    row, err := model.GetValue(iter, col)
    if err != nil {
        log.Printf("Error getting graph row (%s)", err)
        return "" 
    }

    rowVal, err := row.GoValue()
    if err != nil {
        log.Printf("Error converting row to go val (%s)", err)
        return ""
    }
        
    return rowVal.(string)
}

type ListAttributeJSON struct {
    ListAttribute
    ListStore       [][]string
    MarshalJSON     struct {}
    UnmarshalJSON   struct {}
}

func (listAttr *ListAttribute) MarshalJSON() (b []byte, err error) {
    listAttrJSON := &ListAttributeJSON{
        ListAttribute: *listAttr,
        ListStore: make([][]string, 0, 10),
    }

    iter, ok := listAttr.ListStore.GetIterFirst()
    for ok {
        row := listAttr.encodeRow(iter)
        ok = listAttr.ListStore.IterNext(iter)
        listAttrJSON.ListStore = append(listAttrJSON.ListStore, row)
    }

    return json.Marshal(listAttrJSON)
}

func (listAttr *ListAttribute) UnmarshalJSON(b []byte) error {
    var listAttrJSON ListAttributeJSON

    err := json.Unmarshal(b, &listAttrJSON)
    if err != nil {
        return err
    }

    listAttr.FileName = listAttrJSON.FileName
    listAttr.ChromaName = listAttrJSON.ChromaName
    listAttr.Type = listAttrJSON.Type
    listAttr.NumCols = listAttrJSON.NumCols
    listAttr.Selected = listAttrJSON.Selected

    cols := make([]glib.Type, listAttr.NumCols)
    for i := range cols {
        cols[i] = glib.TYPE_STRING
    }

    listAttr.ListStore, err = gtk.ListStoreNew(cols...)
    if err != nil {
        log.Fatalf("Error creating list store")
    }

    colIdx := make([]int, listAttr.NumCols)
    for i := range colIdx {
        colIdx[i] = i
    }

    rowInterface := make([]interface{}, listAttr.NumCols)
    for _, row := range listAttrJSON.ListStore {
        for i, col := range row {
            rowInterface[i] = interface{}(col)
        }

        listAttr.ListStore.Set(listAttr.ListStore.Append(), colIdx, rowInterface)
    }

    return nil
}

func (listAttr *ListAttribute) encodeRow(iter *gtk.TreeIter) []string {
    row := make([]string, listAttr.NumCols)
    for j := 0; j < listAttr.NumCols; j++ {
        row[j] = getStringFromIter(listAttr.ListStore, iter, j)
    }

    return row
}

func (listAttr *ListAttribute) Update(edit Editor) error {
    listEdit, ok := edit.(*ListEditor) 
    if !ok {
        return fmt.Errorf("ListAttribute.Update requires a ListEditor")
    }

    if !listAttr.Selected {
        return nil
    }

    selected, err := listEdit.treeView.GetSelection()
    if err != nil {
        return err
    }
    _, listAttr.selectedIter, ok = selected.GetSelected()
    if !ok {
        return fmt.Errorf("Error getting selected iter from tree view selection")
    }
    // Increment selection
    ok = listAttr.ListStore.IterNext(listAttr.selectedIter)
    if !ok {
        // last item in the list
        listAttr.selectedIter, ok = listAttr.ListStore.GetIterFirst()
    }

    if ok {
        selected.SelectIter(listAttr.selectedIter)
    }

    return nil
}

type ListEditor struct {
    name        string
    box         *gtk.Box
    treeView    *gtk.TreeView
    listStore   *gtk.ListStore
}

func NewListEditor(name string, columns []string) *ListEditor {
    var err error
    listEdit := &ListEditor{name: name}
    listEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Print(err)
        return nil
    }

    listEdit.treeView, err = gtk.TreeViewNew()
    if err != nil {
        log.Print(err)
        return nil
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

                listEdit.listStore.SetValue(iter, gCell.columnNum, text)
        })
        column, err := gtk.TreeViewColumnNewWithAttribute(name, gCell, "text", i)
        if err != nil {
            log.Printf("Error creating list column (%s)", err)
        }

        listEdit.treeView.AppendColumn(column)
    }

    frame, err := gtk.FrameNew(name)
    if err != nil {
        log.Print(err)
        return nil
    }

    frame.Set("border-width", 2 * padding)
    frame.Add(listEdit.treeView)
    frame.SetVisible(true)
    
    actionBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil {
        log.Print(err)
        return nil
    }

    actionBox.SetVisible(true)

    label, err := gtk.LabelNew("Data Rows")
    if err != nil {
        log.Print(err)
        return nil 
    }

    label.SetVisible(true)
    label.SetWidthChars(12)
    actionBox.PackStart(label, false, false, padding)

    // add rows
    button, err := gtk.ButtonNewWithLabel("+")
    if err != nil {
        log.Print(err)
        return nil
    }

    button.Connect("clicked", func() { 
        if listEdit.listStore == nil {
            log.Printf("Graph prop editor does not have a list store")
            return
        }

        listEdit.listStore.Append()
    })

    button.SetVisible(true)
    actionBox.PackStart(button, false, false, padding)

    // remove rows
    button, err = gtk.ButtonNewWithLabel("-")
    if err != nil {
        log.Printf("Error creating graph table (%s)", err)
        return nil
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
    actionBox.PackStart(button, false, false, padding)

    listEdit.box.PackStart(actionBox, true, true, 0)
    listEdit.box.PackStart(frame, true, true, 0)
    return listEdit
}

func (listEdit *ListEditor) Update(attr Attribute) error {
    listAttr, ok := attr.(*ListAttribute)
    if !ok {
        return fmt.Errorf("ListEditor.Update requires a ListAttribute")
    }

    listEdit.listStore = listAttr.ListStore
    listEdit.treeView.SetModel(listEdit.listStore)
    return nil
}

func (listEdit *ListEditor) Box() *gtk.Box {
    return listEdit.box
}

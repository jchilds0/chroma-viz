package props

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type GraphCell struct {
    *gtk.CellRendererText
    columnNum int
}

func NewGraphCell(i int) *GraphCell {
    cell, err := gtk.CellRendererTextNew()
    if err != nil {
        log.Fatalf("Error creating graph cell (%s)", err)
    }

    gCell := &GraphCell{CellRendererText: cell, columnNum: i}
    return gCell
}

type GraphEditor struct {
    box *gtk.Box
    treeView *gtk.TreeView
    listStore *gtk.ListStore
    value [2]*gtk.SpinButton
}

func NewGraphEditor(width, height int, animate func()) PropertyEditor {
    var err error
    g := &GraphEditor{}

    g.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Printf("Error creating graph box (%s)", err)
    }

    g.value[0], err = gtk.SpinButtonNewWithRange(-float64(width), float64(width), 1)
    if err != nil { 
        log.Printf("Error creating graph spin button (%s)", err) 
    }

    g.value[1], err = gtk.SpinButtonNewWithRange(-float64(height), float64(height), 1)
    if err != nil { 
        log.Printf("Error creating graph spin button (%s)", err) 
    }

    posBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil {
        log.Printf("Error creating graph box (%s)", err)
    }

    input := IntEditor("x Pos", g.value[0], animate)
    posBox.PackStart(input, false, false, 0)

    input = IntEditor("y Pos", g.value[1], animate)
    posBox.PackStart(input, false, false, 0)

    posBox.SetVisible(true)
    g.box.PackStart(posBox, false, false, padding)

    columns := []string{"x Pos", "y Pos"}

    g.treeView, err = gtk.TreeViewNew()
    if err != nil {
        log.Printf("Error creating list box (%s)", err)
    }

    for i, name := range columns {
        // creating a graph cell just to reference the column num inside the edit callback
        // i dont like this  

        gCell := NewGraphCell(i)
        gCell.SetProperty("editable", true)
        gCell.Connect("edited", 
            func(cell *gtk.CellRendererText, path string, text string) {
                if g.listStore == nil {
                    log.Printf("Error editing graph prop")
                    return
                }

                iter, err := g.listStore.ToTreeModel().GetIterFromString(path)
                if err != nil {
                    log.Printf("Error editing graph prop (%s)", err)
                    return
                }

                id_val, err := strconv.Atoi(text)
                if err != nil {
                    log.Printf("Error editing graph prop (%s)", err)
                    return
                }

                g.listStore.SetValue(iter, gCell.columnNum, id_val)
                animate()
        })
        column, err := gtk.TreeViewColumnNewWithAttribute(name, gCell, "text", i)
        if err != nil {
            log.Printf("Error creating graph list (%s)", err)
        }

        g.treeView.AppendColumn(column)
    }

    frame, err := gtk.FrameNew("Graph Data")
    if err != nil {
        log.Printf("Error creating graph table (%s)", err)
    }

    frame.Set("border-width", 2 * padding)
    frame.Add(g.treeView)
    g.treeView.SetVisible(true)

    g.box.PackStart(frame, true, true, 0)

    label, err := gtk.LabelNew("Data Rows")
    if err != nil {
        log.Printf("Error creating graph table (%s)", err)
    }

    label.SetVisible(true)
    posBox.PackStart(label, false, false, padding)

    button, err := gtk.ButtonNewWithLabel("+")
    if err != nil {
        log.Printf("Error creating graph table (%s)", err)
    }

    button.Connect("clicked", func() { 
        if g.listStore == nil {
            log.Printf("Graph prop editor does not have a list store")
            return
        }

        g.listStore.Append()
    })

    button.SetVisible(true)
    posBox.PackStart(button, false, false, padding)

    button, err = gtk.ButtonNewWithLabel("-")
    if err != nil {
        log.Printf("Error creating graph table (%s)", err)
    }

    button.Connect("clicked", func() {
        if g.listStore == nil {
            log.Printf("Graph prop editor does not have a list store")
            return
        }

        selection, err := g.treeView.GetSelection()
        if err != nil {
            log.Printf("Error getting current row (%s)", err)
            return
        }

        _, iter, ok := selection.GetSelected()
        if !ok {
            log.Printf("Error getting selected")
            return
        }

        g.listStore.Remove(iter)
    })

    button.SetVisible(true)
    posBox.PackStart(button, false, false, padding)

    frame.SetVisible(true)
    g.box.SetVisible(true)
 
    return g
}

func (g *GraphEditor) Box() *gtk.Box {
    return g.box
}

func (gEdit *GraphEditor) Update(g Property) {
    gProp, ok := g.(*GraphProp)
    if !ok {
        log.Printf("GraphEditor.Update requires a GraphProp")
        return
    }

    gEdit.value[0].SetValue(float64(gProp.Value[0]))
    gEdit.value[1].SetValue(float64(gProp.Value[1]))

    gEdit.listStore = gProp.listStore
    gEdit.treeView.SetModel(gEdit.listStore)
}

type GraphProp struct {
    name string
    listStore *gtk.ListStore
    Value [2]int
}

func NewGraphProp(name string) Property {
    var err error
    g := &GraphProp{name: name}

    g.listStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_INT)
    if err != nil {
        log.Printf("Error creating graph list (%s)", err)
    }

    return g
}

func (g *GraphProp) Type() int {
    return GRAPH_PROP 
}

func (g *GraphProp) Name() string {
    return g.name
}

func (g *GraphProp) String() string {
    str := fmt.Sprintf("rel_x=%d#rel_y=%d#num_node=0#", 
        g.Value[0], g.Value[1])

    iter, ok := g.listStore.GetIterFirst()
    i := 0
    for ok {
        xVal := getIntFromIter(g.listStore, iter, 0)
        yVal := getIntFromIter(g.listStore, iter, 1)
        
        str = fmt.Sprintf("%sgraph_node=%d %d %d#", str, i, xVal, yVal)
        ok = g.listStore.IterNext(iter)
        i++
    }

    return str
}

func (g *GraphProp) Encode() string {
    str := fmt.Sprintf("x %d;y %d;",
        g.Value[0], g.Value[1])

    iter, ok := g.listStore.GetIterFirst()
    i := 0
    for ok {
        xVal := getIntFromIter(g.listStore, iter, 0)
        yVal := getIntFromIter(g.listStore, iter, 1)
        
        str = fmt.Sprintf("%snode %d %d;", str, xVal, yVal)
        ok = g.listStore.IterNext(iter)
        i++
    }

    return str
}

func (g *GraphProp) Decode(input string) {
    attrs := strings.Split(input, ";")

    for _, attr := range attrs[1:] {
        line := strings.Split(attr, " ")

        if len(line) < 2 {
            continue
        }

        name := line[0]

        switch (name) {
        case "x":
            value, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Printf("Error decoding graph (%s)", err) 
            }

            g.Value[0] = value
        case "y":
            value, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Printf("Error decoding graph (%s)", err) 
            }

            g.Value[1] = value
        case "node":
            x, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Printf("Error decoding graph (%s)", err) 
            }

            y, err := strconv.Atoi(line[2])
            if err != nil { 
                log.Printf("Error decoding graph (%s)", err) 
            }

            g.listStore.Set(g.listStore.Append(), []int{0, 1}, []interface{}{x, y})
        default:
            log.Printf("Unknown GraphProp attr name (%s)\n", name)
        }
    }
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
func (gProp *GraphProp) Update(g PropertyEditor, action int) {
    gEdit, ok := g.(*GraphEditor)
    if !ok {
        log.Printf("GraphProp.Update requires GraphEditor")
        return
    }

    gProp.Value[0] = gEdit.value[0].GetValueAsInt()
    gProp.Value[1] = gEdit.value[1].GetValueAsInt()
}


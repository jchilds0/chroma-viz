package props

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
    LINE_NUM = iota
    TEXT
    COUNT
)

type TickerEditor struct {
    box         *gtk.Box
    text        *gtk.Entry
    treeView    *gtk.TreeView
    listStore   *gtk.ListStore
    value       [2]*gtk.SpinButton
}

func NewTickerEditor(width, height int, animate func()) PropertyEditor {
    var err error 
    t := &TickerEditor{}

    t.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Printf("Error creating ticker (%s)", err)
    }

    pos, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil {
        log.Printf("Error creating ticker (%s)", err)
    }

    pos.SetVisible(true)
    t.box.PackStart(pos, false, false, padding)

    upper := []int{width, height}
    for i := range t.value {
        t.value[i], err = gtk.SpinButtonNewWithRange(-float64(upper[i]), float64(upper[i]), 1)
        if err != nil { 
            log.Fatalf("Error creating text prop (%s)", err) 
        }

        spin := IntEditor("x Pos", t.value[i], animate)
        pos.PackStart(spin, false, false, 0)
    }

    t.treeView, err = gtk.TreeViewNew()
    if err != nil {
        log.Printf("Error creating ticker (%s)", err)
    }

    cell, err := gtk.CellRendererTextNew()
    if err != nil {
        log.Printf("Error creating ticker (%s)", err)
    }

    column, err := gtk.TreeViewColumnNewWithAttribute("Num", cell, "text", LINE_NUM)
    if err != nil {
        log.Printf("Error creating ticker prop (%s)", err)
    }

    t.treeView.AppendColumn(column)

    cell, err = gtk.CellRendererTextNew()
    if err != nil {
        log.Printf("Error creating ticker (%s)", err)
    }
    
    cell.SetProperty("editable", true)
    cell.Connect("edited", 
        func(cell *gtk.CellRendererText, path string, text string) {
            if t.listStore == nil {
                log.Printf("TickerEditor does not have a list store")
                return
            }

            iter, err := t.listStore.GetIterFromString(path)
            if err != nil {
                log.Printf("Error editing ticker list (%s)", err)
                return
            }

            t.listStore.SetValue(iter, TEXT, text) 
    })

    column, err = gtk.TreeViewColumnNewWithAttribute("Text", cell, "text", TEXT)
    if err != nil {
        log.Printf("Error creating ticker prop (%s)", err)
    }

    t.treeView.AppendColumn(column)
    t.treeView.SetVisible(true)

    label, err := gtk.LabelNew("Ticker Rows")
    if err != nil {
        log.Printf("Error creating graph table (%s)", err)
    }

    label.SetVisible(true)
    pos.PackStart(label, false, false, padding)

    button, err := gtk.ButtonNewWithLabel("+")
    if err != nil {
        log.Printf("Error creating graph table (%s)", err)
    }

    button.Connect("clicked", func() { 
        t.listStore.Append()
    })
    button.SetVisible(true)
    pos.PackStart(button, false, false, padding)

    button, err = gtk.ButtonNewWithLabel("-")
    if err != nil {
        log.Printf("Error creating graph table (%s)", err)
    }

    button.Connect("clicked", func() {
        selection, err := t.treeView.GetSelection()
        if err != nil {
            log.Printf("Error getting current row (%s)", err)
            return
        }

        _, iter, ok := selection.GetSelected()
        if !ok {
            log.Printf("Error getting selected")
            return
        }

        t.listStore.Remove(iter)
    })

    button.SetVisible(true)
    pos.PackStart(button, false, false, padding)

    frame, err := gtk.FrameNew("Ticker Text")
    if err != nil {
        log.Printf("Error creating graph table (%s)", err)
    }

    frame.Set("border-width", 2 * padding)
    frame.Add(t.treeView)
    frame.SetVisible(true)

    t.box.PackStart(frame, true, true, 0)
    t.box.SetVisible(true)

    return t
}

func (tickEdit *TickerEditor) Box() *gtk.Box {
    return tickEdit.box
}

func (tickEdit *TickerEditor) Update(tick Property) {
    tickProp, ok := tick.(*TickerProp)
    if !ok {
        log.Printf("TickerEditor.Update requires a tickerProp property")
        return
    }

    tickEdit.listStore = tickProp.listStore
    tickEdit.treeView.SetModel(tickEdit.listStore)

    tickEdit.value[0].SetValue(float64(tickProp.value[0]))
    tickEdit.value[1].SetValue(float64(tickProp.value[1]))
}

type TickerProp struct {
    name string
    text string 
    value [2]int
    listStore *gtk.ListStore
}

func NewTickerProp(name string) Property {
    var err error

    t := &TickerProp{
        name: name, 
    }

    t.listStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING)
    if err != nil {
        log.Printf("Error creating ticker (%s)", err)
    }

    t.listStore.Connect("row-inserted", 
        func(model *gtk.ListStore, path *gtk.TreePath, iter *gtk.TreeIter) {
            iter, ok := model.GetIterFirst()
            i := 1

            for ok {
                model.SetValue(iter, LINE_NUM, i)
                ok = model.IterNext(iter)
                i++
            }
    })

    t.listStore.Connect("row-deleted", 
        func(model *gtk.ListStore, path *gtk.TreePath) {
            iter, ok := model.GetIterFirst()
            i := 1

            for ok {
                model.SetValue(iter, LINE_NUM, i)
                ok = model.IterNext(iter)
                i++
            }
    })

    return t
}

func (t *TickerProp) Type() int {
    return TICKER_PROP
}

func (t *TickerProp) Name() string {
    return t.name
}

func(t *TickerProp) String() string {
    return fmt.Sprintf("string=%s#rel_x=%d#rel_y=%d#", 
        t.text, t.value[0], t.value[1])
}

func (t *TickerProp) Encode() string {
    str := fmt.Sprintf("x %d;y %d;", t.value[0], t.value[1])
    iter, ok := t.listStore.GetIterFirst()

    for ok {
        g_text, err := t.listStore.GetValue(iter, TEXT)
        if err != nil {
            log.Printf("Error exporting ticker prop (%s)", err)
            return ""
        }

        text, err := g_text.GoValue()
        if err != nil {
            log.Printf("Error exporting ticker prop (%s)", err)
            return ""
        }

        str = fmt.Sprintf("%stext %s;", str, text)
        ok = t.listStore.IterNext(iter)
    }

    return str
}

func (t *TickerProp) Decode(input string) {
    attrs := strings.Split(input, ";")

    for _, attr := range attrs[1:] {
        line := strings.Split(attr, " ")
        name := line[0]

        switch (name) {
        case "x":
            value, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Fatalf("Error decoding ticker prop (%s)", err) 
            }

            t.value[0] = value
        case "y":
            value, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Printf("Error decoding ticker prop (%s)", err) 
            }

            t.value[1] = value
        case "text":
            t.listStore.Set(
                t.listStore.Append(), 
                []int{LINE_NUM, TEXT}, 
                []interface{}{0, strings.TrimPrefix(attr, "text ")})
        case "":
        default:
            log.Printf("Unknown TickerProp attr name (%s)\n", name)
        }
    }

    t.listStore.Append()
}

func (tickProp *TickerProp) Update(t PropertyEditor, action int) {
    tickEdit, ok := t.(*TickerEditor)
    if !ok {
        log.Printf("TickerProp.Update requires a TickerEditor")
        return
    }

    switch action {
    case ANIMATE_ON, CONTINUE:
        tickProp.value[0] = tickEdit.value[0].GetValueAsInt()
        tickProp.value[1] = tickEdit.value[1].GetValueAsInt()

        // Get text from selection
        selection, err := tickEdit.treeView.GetSelection()
        if err != nil {
            log.Printf("Error getting ticker selection (%s)", err)
            return 
        }

        _, iter, ok := selection.GetSelected()
        if !ok {
            log.Printf("Error getting selected")
            return
        }

        new_text, err := tickProp.listStore.GetValue(iter, TEXT)
        if err != nil {
            log.Printf("Error getting selection value (%s)", err)
            return
        }

        text, err := new_text.GetString()
        if err != nil {
            text = ""
        }

        tickProp.text = text

        // Increment selection
        ok = tickProp.listStore.IterNext(iter)
        if !ok {
            // last item in the list
            iter, ok = tickProp.listStore.GetIterFirst()
        }

        if ok {
            selection.SelectIter(iter)
        }
    case ANIMATE_OFF:
    default:
        log.Printf("Unknown action")
    }
}

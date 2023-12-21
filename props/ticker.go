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

type TickerProp struct {
    box *gtk.Box
    listStore *gtk.ListStore
    treeView *gtk.TreeView
    name string
    text *TextProp
}

func NewTickerProp(width, height int, animate func(), name string) Property {
    var err error

    t := &TickerProp{
        name: name, 
        text: NewTextProp(width, height, animate, name),
    }

    t.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Printf("Error creating ticker (%s)", err)
    }

    t.text.pos.Unparent()

    t.box.PackStart(t.text.pos, false, false, padding)

    t.treeView, err = gtk.TreeViewNew()
    if err != nil {
        log.Printf("Error creating ticker (%s)", err)
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
    t.treeView.SetModel(t.listStore)
    t.treeView.SetVisible(true)

    label, err := gtk.LabelNew("Ticker Rows")
    if err != nil {
        log.Printf("Error creating graph table (%s)", err)
    }

    label.SetVisible(true)
    t.text.pos.PackStart(label, false, false, padding)

    button, err := gtk.ButtonNewWithLabel("+")
    if err != nil {
        log.Printf("Error creating graph table (%s)", err)
    }

    button.Connect("clicked", func() { 
        t.listStore.Append()
    })
    button.SetVisible(true)
    t.text.pos.PackStart(button, false, false, padding)

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
    t.text.pos.PackStart(button, false, false, padding)

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

func (t *TickerProp) Tab() *gtk.Box {
    return t.box
}

func (t *TickerProp) Name() string {
    return t.name
}

func(t *TickerProp) String() string {
        return t.text.String()
}

func (t *TickerProp) Encode() string {
    str := t.text.Encode()
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
                log.Fatalf("Error decoding text prop (%s)", err) 
            }

            t.text.value[0].SetValue(float64(value))
        case "y":
            value, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Printf("Error decoding text prop (%s)", err) 
            }

            t.text.value[1].SetValue(float64(value))
        case "string":
            t.text.entry.SetText(strings.TrimPrefix(attr, "string "))
        case "text":
            t.listStore.Set(
                t.listStore.Append(), 
                []int{LINE_NUM, TEXT}, 
                []interface{}{0, strings.TrimPrefix(attr, "text ")})
        case "":
        default:
            log.Printf("Unknown TextProp attr name (%s)\n", name)
        }
    }

    t.listStore.Append()
}

func (t *TickerProp) Update(action int) {
    switch action {
    case ANIMATE_ON, CONTINUE:
        // Get text from selection
        selection, err := t.treeView.GetSelection()
        if err != nil {
            log.Printf("Error getting ticker selection (%s)", err)
            return 
        }

        _, iter, ok := selection.GetSelected()
        if !ok {
            log.Printf("Error getting selected")
            return
        }

        new_text, err := t.listStore.GetValue(iter, TEXT)
        if err != nil {
            log.Printf("Error getting selection value (%s)", err)
            return
        }

        text, err := new_text.GetString()
        if err != nil {
            text = ""
        }

        // Update TextProp text
        t.text.entry.SetText(text)

        // Increment selection
        ok = t.listStore.IterNext(iter)
        if !ok {
            // last item in the list
            iter, ok = t.listStore.GetIterFirst()
        }

        if ok {
            selection.SelectIter(iter)
        }
    case ANIMATE_OFF:
    default:
        log.Printf("Unknown action")
    }
}

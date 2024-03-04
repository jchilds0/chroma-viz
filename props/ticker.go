package props

import (
	"chroma-viz/attribute"

	"github.com/gotk3/gotk3/gtk"
)

const (
    LINE_NUM = iota
    TEXT
    COUNT
)

type TickerEditor struct {
    box         *gtk.Box
    edit        map[string]attribute.Editor
}

func NewTickerEditor(width, height int, animate func()) (t *TickerEditor, err error) {
    t = &TickerEditor{}
    t.edit = make(map[string]attribute.Editor, 5)

    t.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        return
    }

    columns := []string{"Num", "Text"}
    t.edit["text"], err = attribute.NewListEditor("Ticker", columns, animate)
    if err != nil {
        return 
    }

    t.edit["x"], err = attribute.NewIntEditor("x", 0, float64(width), animate)
    if err != nil {
        return
    }

    t.edit["y"], err = attribute.NewIntEditor("y", 0, float64(height), animate)
    if err != nil {
        return
    }

    t.box.SetVisible(true)
    for _, attr := range t.edit {
        t.box.PackStart(attr.Box(), false, false, padding)
    }

    return
}

func (tickEdit *TickerEditor) Box() *gtk.Box {
    return tickEdit.box
}

func (tickEdit *TickerEditor) Editors() map[string]attribute.Editor {
    return tickEdit.edit
}

type TickerProp struct {
    name      string
    attrs     map[string]attribute.Attribute
    visible   map[string]bool
}

func NewTickerProp(name string, visible map[string]bool) *TickerProp {
    t := &TickerProp{name: name, visible: visible}

    t.attrs = make(map[string]attribute.Attribute, 5)
    t.attrs["x"] = attribute.NewIntAttribute("rel_x")
    t.attrs["y"] = attribute.NewIntAttribute("rel_y")
    t.attrs["text"] = attribute.NewListAttribute("ticker")

    return t
}

func (t *TickerProp) Type() int {
    return TICKER_PROP
}

func (t *TickerProp) Name() string {
    return t.name
}

func (t *TickerProp) Visible() map[string]bool {
    return t.visible
}

func (t *TickerProp) Attributes() map[string]attribute.Attribute {
    return t.attrs
}

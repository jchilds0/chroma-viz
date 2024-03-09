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

    columns := []string{"Text"}
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
    t.box.PackStart(t.edit["x"].Box(), false, false, padding)
    t.box.PackStart(t.edit["y"].Box(), false, false, padding)
    t.box.PackStart(t.edit["text"].Box(), false, false, padding)

    return
}

func (tickEdit *TickerEditor) Box() *gtk.Box {
    return tickEdit.box
}

func (tickEdit *TickerEditor) Editors() map[string]attribute.Editor {
    return tickEdit.edit
}


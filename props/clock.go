package props

import (
	"chroma-viz/attribute"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type ClockEditor struct {
    box             *gtk.Box
    edit            map[string]attribute.Editor
}

func NewClockEditor(width, height int, animate, cont func()) (clockEdit *ClockEditor, err error) {
    clockEdit = &ClockEditor{}

    clockEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    clockEdit.box.SetVisible(true)

    clockEdit.edit = make(map[string]attribute.Editor, 5)
    clockEdit.edit["clock"], err = attribute.NewClockEditor("Time", animate, cont)
    if err != nil {
        return 
    }

    clockEdit.edit["x"], err = attribute.NewIntEditor("x", 0, float64(width), animate)
    if err != nil {
        return
    }

    clockEdit.edit["y"], err = attribute.NewIntEditor("y", 0, float64(height), animate)
    if err != nil {
        return
    }

    clockEdit.edit["color"], err = attribute.NewColorEditor("Color", animate)
    if err != nil {
        return
    }

    clockEdit.box.PackStart(clockEdit.edit["x"].Box(), false, false, padding)
    clockEdit.box.PackStart(clockEdit.edit["y"].Box(), false, false, padding)
    clockEdit.box.PackStart(clockEdit.edit["clock"].Box(), false, false, padding)
    clockEdit.box.PackStart(clockEdit.edit["color"].Box(), false, false, padding)

    return 
}

func (clock *ClockEditor) Box() *gtk.Box {
    return clock.box
}

func (clock *ClockEditor) Editors() map[string]attribute.Editor {
    return clock.edit
}


package props

import (
	"chroma-viz/attribute"

	"github.com/gotk3/gotk3/gtk"
)


type RectEditor struct {
    box     *gtk.Box
    edit    map[string]attribute.Editor
}

func NewRectEditor(width, height int, animate func()) (rectEdit *RectEditor, err error) {
    rectEdit = &RectEditor{}
    rectEdit.edit = make(map[string]attribute.Editor, 5) 

    rectEdit.edit["x"], err = attribute.NewIntEditor("x", 0, float64(width), animate)
    if err != nil {
        return
    }

    rectEdit.edit["y"], err = attribute.NewIntEditor("y", 0, float64(height), animate)
    if err != nil {
        return
    }

    rectEdit.edit["width"], err = attribute.NewIntEditor("Width", 0, float64(width), animate)
    if err != nil {
        return
    }

    rectEdit.edit["height"], err = attribute.NewIntEditor("Height", 0, float64(height), animate)
    if err != nil {
        return
    }

    rectEdit.edit["color"], err = attribute.NewColorEditor("Color", animate)
    if err != nil {
        return
    }

    rectEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        return 
    }

    rectEdit.box.SetVisible(true)
    rectEdit.box.PackStart(rectEdit.edit["x"].Box(), false, false, padding)
    rectEdit.box.PackStart(rectEdit.edit["y"].Box(), false, false, padding)
    rectEdit.box.PackStart(rectEdit.edit["width"].Box(), false, false, padding)
    rectEdit.box.PackStart(rectEdit.edit["height"].Box(), false, false, padding)
    rectEdit.box.PackStart(rectEdit.edit["color"].Box(), false, false, padding)

    return 
}

func (rectEdit *RectEditor) Box() *gtk.Box {
    return rectEdit.box
}

func (rectEdit *RectEditor) Editors() map[string]attribute.Editor {
    return rectEdit.edit
}


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

    labels := []string{"x", "y", "Width", "Height"}
    upper := []int{width, height, width, height}

    for i := range labels {
        rectEdit.edit[labels[i]], err = attribute.NewIntEditor(labels[i], 0, float64(upper[i]), animate)

        if err != nil {
            return
        }
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
    for _, attr := range rectEdit.edit {
        rectEdit.box.PackStart(attr.Box(), false, false, padding)
    }

    return 
}

func (rectEdit *RectEditor) Box() *gtk.Box {
    return rectEdit.box
}

func (rectEdit *RectEditor) Editors() map[string]attribute.Editor {
    return rectEdit.edit
}

type RectProp struct {
    name        string 
    attrs       map[string]attribute.Attribute
    visible     map[string]bool
}

func NewRectProp(name string) *RectProp {
    rect := &RectProp{name: name}
    rect.attrs = make(map[string]attribute.Attribute, 5)
    rect.visible = make(map[string]bool, 5)

    rect.attrs["x"] = attribute.NewIntAttribute("x")
    rect.attrs["y"] = attribute.NewIntAttribute("y")
    rect.attrs["width"] = attribute.NewIntAttribute("width")
    rect.attrs["height"] = attribute.NewIntAttribute("height")
    rect.attrs["color"] = attribute.NewColorAttribute()

    return rect
}

func (rect *RectProp) Name() string {
    return rect.name
}

func (rect *RectProp) Type() int {
    return RECT_PROP
}

func (rect *RectProp) Visible() map[string]bool {
    return rect.visible
}

func (rect *RectProp) Attributes() map[string]attribute.Attribute {
    return rect.attrs
}

package props

import (
	"chroma-viz/attribute"

	"github.com/gotk3/gotk3/gtk"
)

type TextEditor struct {
    box *gtk.Box
    edit map[string]attribute.Editor
}

func NewTextEditor(width, height int, animate func()) (textEdit *TextEditor, err error) {
    textEdit = &TextEditor{}

    textEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        return
    }

    textEdit.box.SetVisible(true)
    textEdit.edit = make(map[string]attribute.Editor, 5)

    textEdit.edit["string"], err = attribute.NewStringEditor("Text", animate)
    if err != nil {
        return
    }

    textEdit.edit["x"], err = attribute.NewIntEditor("x", 0, float64(width), animate)
    if err != nil {
        return
    }

    textEdit.edit["y"], err = attribute.NewIntEditor("y", 0, float64(height), animate)
    if err != nil {
        return
    }

    textEdit.edit["color"], err = attribute.NewColorEditor("Color", animate)
    if err != nil {
        return
    }

    textEdit.box.PackStart(textEdit.edit["x"].Box(), false, false, padding)
    textEdit.box.PackStart(textEdit.edit["y"].Box(), false, false, padding)
    textEdit.box.PackStart(textEdit.edit["string"].Box(), false, false, padding)
    textEdit.box.PackStart(textEdit.edit["color"].Box(), false, false, padding)

    return 
}

func (text *TextEditor) Box() *gtk.Box {
    return text.box
}

func (text *TextEditor) Editors() map[string]attribute.Editor {
    return text.edit
}

type TextProp struct {
    name        string
    attrs       map[string]attribute.Attribute
    visible     map[string]bool
}

func NewTextProp(name string, visible map[string]bool) *TextProp {
    text := &TextProp{name: name}
    text.attrs = make(map[string]attribute.Attribute, 5)
    text.visible = visible

    text.attrs["x"] = attribute.NewIntAttribute("rel_x")
    text.attrs["y"] = attribute.NewIntAttribute("rel_y")
    text.attrs["string"] = attribute.NewStringAttribute("string")
    text.attrs["color"] = attribute.NewColorAttribute()

    return text
}

func (text *TextProp) Type() int {
    return TEXT_PROP
}

func (text *TextProp) Name() string {
    return text.name
}

func (text *TextProp) Visible() map[string]bool {
    return text.visible
}

func (text *TextProp) Attributes() map[string]attribute.Attribute {
    return text.attrs
}

package props

import (
	"chroma-viz/attribute"

	"github.com/gotk3/gotk3/gtk"
)

type CircleEditor struct {
    box *gtk.Box
    edit map[string]attribute.Editor
}

func NewCircleEditor(width, height int, animate func()) (circleEdit *CircleEditor, err error) {
    circleEdit = &CircleEditor{}
    circleEdit.edit = make(map[string]attribute.Editor, 10)

    circleEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        return
    }

    name := []string{
        "x", 
        "y", 
        "inner_radius", 
        "outer_radius", 
        "start_angle", 
        "end_angle",
    }

    labels := []string{
        "Center x", 
        "Center y", 
        "Inner Radius", 
        "Outer Radius", 
        "Start Angle", 
        "End Angle",
    }
    upper := []int{
        width, 
        height, 
        width, 
        width, 
        360, 
        360,
    }   

    for i := range labels {
        circleEdit.edit[name[i]], err = attribute.NewIntEditor(labels[i], 0, float64(upper[i]), animate)

        if err != nil { 
            return
        }
    }

    circleEdit.edit["color"], err = attribute.NewColorEditor("Color", animate)
    if err != nil {
        return
    }

    circleEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        return 
    }

    circleEdit.box.SetVisible(true)
    for _, attr := range circleEdit.edit {
        circleEdit.box.PackStart(attr.Box(), false, false, padding)
    }

    return 
}

func (circleEdit *CircleEditor) Box() *gtk.Box {
    return circleEdit.box
}

func (circleEdit *CircleEditor) Editors() map[string]attribute.Editor {
    return circleEdit.edit
}

type CircleProp struct {
    name      string 
    attrs     map[string]attribute.Attribute
    visible   map[string]bool
}

func NewCircleProp(name string, visible map[string]bool) *CircleProp {
    circle := &CircleProp{name: name, visible: visible}
    circle.attrs = make(map[string]attribute.Attribute, 10)

    circle.attrs["x"] = attribute.NewIntAttribute("rel_x")
    circle.attrs["y"] = attribute.NewIntAttribute("rel_y")
    circle.attrs["inner_radius"] = attribute.NewIntAttribute("inner_radius")
    circle.attrs["outer_radius"] = attribute.NewIntAttribute("outer_radius")
    circle.attrs["start_angle"] = attribute.NewIntAttribute("start_angle")
    circle.attrs["end_angle"] = attribute.NewIntAttribute("end_angle")

    return circle
}

func (circle *CircleProp) Type() int {
    return CIRCLE_PROP
}

func (circle *CircleProp) Name() string {
    return circle.name
}

func (circle *CircleProp) Visible() map[string]bool {
    return circle.visible
}

func (circle *CircleProp) Attributes() map[string]attribute.Attribute {
    return circle.attrs
}

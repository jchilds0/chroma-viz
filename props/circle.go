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
    circleEdit.box.PackStart(circleEdit.edit["x"].Box(), false, false, padding)
    circleEdit.box.PackStart(circleEdit.edit["y"].Box(), false, false, padding)
    circleEdit.box.PackStart(circleEdit.edit["inner_radius"].Box(), false, false, padding)
    circleEdit.box.PackStart(circleEdit.edit["outer_radius"].Box(), false, false, padding)
    circleEdit.box.PackStart(circleEdit.edit["start_angle"].Box(), false, false, padding)
    circleEdit.box.PackStart(circleEdit.edit["end_angle"].Box(), false, false, padding)
    circleEdit.box.PackStart(circleEdit.edit["color"].Box(), false, false, padding)

    return 
}

func (circleEdit *CircleEditor) Box() *gtk.Box {
    return circleEdit.box
}

func (circleEdit *CircleEditor) Editors() map[string]attribute.Editor {
    return circleEdit.edit
}


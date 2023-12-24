package props

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type CircleEditor struct {
    box *gtk.Box
    value [6]*gtk.SpinButton
}

func NewCircleEditor(width, height int, animate func()) PropertyEditor {
    var err error
    circle := &CircleEditor{}

    circle.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Printf("Error creating circle box (%s)", err) 
    }

    labels := []string{"Center x", "Center y", "Inner Radius", 
        "Outer Radius", "Start Angle", "End Angle"}
    upper := []int{width, height, width, width, 360, 360}   

    for i := range circle.value {
        circle.value[i], err = gtk.SpinButtonNewWithRange(-float64(upper[i]), float64(upper[i]), 1)
        if err != nil { 
            log.Printf("Error creating circle spin button (%s)", err) 
        }

        input := IntEditor(labels[i], circle.value[i], animate)
        circle.box.PackStart(input, false, false, padding)
    }

    circle.value[2].SetValue(1.0)
    circle.value[5].SetValue(360.0)
    circle.box.SetVisible(true)

    return circle
}

func (circle *CircleEditor) Box() *gtk.Box {
    return circle.box
}

func (circleEdit *CircleEditor) Update(circle Property) {
    circleProp, ok := circle.(*CircleProp)
    if !ok {
        log.Printf("CircleEditor.Update requires CircleProp")
        return
    }

    for i := range circleEdit.value {
        circleEdit.value[i].SetValue(float64(circleProp.value[i]))
    }
}

type CircleProp struct {
    name string 
    value [6]int 
}

func NewCircleProp(name string) Property {
    circle := &CircleProp{name: name}
    return circle
}

func (circle *CircleProp) Type() int {
    return CIRCLE_PROP
}

func (circle *CircleProp) Name() string {
    return circle.name
}

func (circle *CircleProp) String() string {
    return fmt.Sprintf("rel_x=%d#rel_y=%d#inner_radius=%d#" + 
        "outer_radius=%d#start_angle=%d#end_angle=%d#", 
        circle.value[0], circle.value[1], circle.value[2],
        circle.value[3], circle.value[4], circle.value[5])
}

func (circle *CircleProp) Encode() string {
    return fmt.Sprintf("x %d;y %d;inner_radius %d;outer_radius %d;" +
        "start_angle %d;end_angle %d;", 
        circle.value[0], circle.value[1], circle.value[2],
        circle.value[3], circle.value[4], circle.value[5])
}

func (circle *CircleProp) Decode(input string) {
    attrs := strings.Split(input, ";")

    for _, attr := range attrs[1:] {
        line := strings.Split(attr, " ")

        if len(line) != 2 {
            continue
        }

        name := line[0]
        value, err := strconv.Atoi(line[1])
        if err != nil { 
            log.Printf("Error circle decode (%s)", err) 
        }

        switch (name) {
        case "x":
            circle.value[0] = value
        case "y":
            circle.value[1] = value
        case "inner_radius":
            circle.value[2] = value
        case "outer_radius":
            circle.value[3] = value
        case "start_angle":
            circle.value[4] = value
        case "end_angle":
            circle.value[5] = value
        default:
            log.Printf("Unknown CircleProp attr name (%s)\n", name)
        }
    }
}

func (circleProp *CircleProp) Update(circle PropertyEditor, action int) {
    circleEdit, ok := circle.(*CircleEditor)
    if !ok {
        log.Printf("CircleProp.Update requires CircleEditor")
        return
    }

    for i := range circleEdit.value {
        circleProp.value[i] = circleEdit.value[i].GetValueAsInt()
    }
}

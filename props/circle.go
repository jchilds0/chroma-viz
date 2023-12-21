package props

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type CircleProp struct {
    name string 
    box *gtk.Box;
    value [6]*gtk.SpinButton
}

func NewCircleProp(width, height int, animate func(), name string) Property {
    var err error
    circle := &CircleProp{name: name}

    circle.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Printf("Error creating circle box (%s)", err) 
    }

    labels := []string{"Center x", "Center y", "Inner Radius", 
        "Outer Radius", "Start Angle", "End Angle"}
    upper := []int{width, height, width, width, 360, 360}   

    for i := range circle.value {
        circle.value[i], err = gtk.SpinButtonNewWithRange(float64(0), float64(upper[i]), 1)
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

func (circle *CircleProp) Tab() *gtk.Box {
    return circle.box
}

func (circle *CircleProp) Name() string {
    return circle.name
}

func (circle *CircleProp) String() string {
    return fmt.Sprintf("center_x=%d#center_y=%d#inner_radius=%d#" + 
        "outer_radius=%d#start_angle=%d#end_angle=%d#", 
        circle.value[0].GetValueAsInt(),
        circle.value[1].GetValueAsInt(),
        circle.value[2].GetValueAsInt(),
        circle.value[3].GetValueAsInt(),
        circle.value[4].GetValueAsInt(),
        circle.value[5].GetValueAsInt())
}

func (circle *CircleProp) Encode() string {
    return fmt.Sprintf("x %d;y %d;inner_radius %d;outer_radius %d;" +
        "start_angle %d;end_angle %d;", 
        circle.value[0].GetValueAsInt(),
        circle.value[1].GetValueAsInt(),
        circle.value[2].GetValueAsInt(),
        circle.value[3].GetValueAsInt(),
        circle.value[4].GetValueAsInt(),
        circle.value[5].GetValueAsInt())
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
            circle.value[0].SetValue(float64(value))
        case "y":
            circle.value[1].SetValue(float64(value))
        case "inner_radius":
            circle.value[2].SetValue(float64(value))
        case "outer_radius":
            circle.value[3].SetValue(float64(value))
        case "start_angle":
            circle.value[4].SetValue(float64(value))
        case "end_angle":
            circle.value[5].SetValue(float64(value))
        default:
            log.Printf("Unknown CircleProp attr name (%s)\n", name)
        }
    }
}

func (circle *CircleProp) Update(action int) {
    switch action {
    case ANIMATE_ON:
    case CONTINUE:
    case ANIMATE_OFF:
    default:
        log.Printf("Unknown action")
    }
}

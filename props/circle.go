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
    value [3]*gtk.SpinButton
}

func NewCircleProp(width, height int, animate func(), name string) Property {
    var err error
    circle := &CircleProp{name: name}

    circle.value[0], err = gtk.SpinButtonNewWithRange(float64(0), float64(width), 1)
    if err != nil { 
        log.Printf("Error creating circle spin button (%s)", err) 
    }

    circle.value[1], err = gtk.SpinButtonNewWithRange(float64(0), float64(height), 1)
    if err != nil { 
        log.Printf("Error creating circle spin button (%s)", err) 
    }

    circle.value[2], err = gtk.SpinButtonNewWithRange(float64(0), float64(width), 1)
    if err != nil { 
        log.Printf("Error creating circle spin button (%s)", err) 
    }

    circle.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Printf("Error creating circle box (%s)", err) 
    }

    labels := [3]string{"Center x", "Center y", "Radius"}

    circle.value[2].SetValue(1.0)
    for i := range circle.value {
        input := IntEditor(labels[i], circle.value[i], animate)
        circle.box.PackStart(input, false, false, padding)
    }

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
    return fmt.Sprintf("center_x=%d#center_y=%d#radius=%d#", 
        circle.value[0].GetValueAsInt(),
        circle.value[1].GetValueAsInt(),
        circle.value[2].GetValueAsInt())
}

func (circle *CircleProp) Encode() string {
    return fmt.Sprintf("x %d;y %d;radius %d;", 
        circle.value[0].GetValueAsInt(),
        circle.value[1].GetValueAsInt(),
        circle.value[2].GetValueAsInt())
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
        case "radius":
            circle.value[2].SetValue(float64(value))
        default:
            log.Printf("Unknown CircleProp attr name (%s)\n", name)
        }
    }
}


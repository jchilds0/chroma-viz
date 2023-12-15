package props

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type RectProp struct {
    box *gtk.Box
    value [4]*gtk.SpinButton
    num int
}

func NewRectProp(num, width, height int, animate func()) Property {
    var err error
    rect := &RectProp{num: num}

    rect.value[0], err = gtk.SpinButtonNewWithRange(float64(0), float64(width), 1)
    if err != nil { 
        log.Printf("Error creating rect spin button (%s)", err) 
    }

    rect.value[1], err = gtk.SpinButtonNewWithRange(float64(0), float64(height), 1)
    if err != nil { 
        log.Printf("Error creating rect spin button (%s)", err) 
    }

    rect.value[2], err = gtk.SpinButtonNewWithRange(float64(0), float64(width), 1)
    if err != nil { 
        log.Printf("Error creating rect spin button (%s)", err) 
    }

    rect.value[3], err = gtk.SpinButtonNewWithRange(float64(0), float64(height), 1)
    if err != nil { 
        log.Printf("Error creating rect spin button (%s)", err) 
    }

    rect.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Printf("Error creating rect box (%s)", err) 
    }

    labels := [4]string{"x Pos", "y Pos", "Width", "Height"}

    for i := range rect.value {
        input := IntEditor(labels[i], rect.value[i], animate)
        rect.box.PackStart(input, false, false, padding)
    }

    rect.box.SetVisible(true)
    return rect
}

func (rect *RectProp) Tab() *gtk.Box {
    return rect.box
}

func (rect *RectProp) String() string {
    return fmt.Sprintf("rect=%d#pos_x=%d#pos_y=%d#width=%d#height=%d#", 
        rect.num,
        rect.value[0].GetValueAsInt(),
        rect.value[1].GetValueAsInt(),
        rect.value[2].GetValueAsInt(),
        rect.value[3].GetValueAsInt())
}

func (rect *RectProp) Encode() string {
    return fmt.Sprintf("x %d;y %d;width %d;height %d;", 
        rect.value[0].GetValueAsInt(),
        rect.value[1].GetValueAsInt(),
        rect.value[2].GetValueAsInt(),
        rect.value[3].GetValueAsInt())
}

func (rect *RectProp) Decode(input string) {
    attrs := strings.Split(input, ";")

    for _, attr := range attrs[1:] {
        line := strings.Split(attr, " ")

        if len(line) != 2 {
            continue
        }

        name := line[0]
        value, err := strconv.Atoi(line[1])
        if err != nil { 
            log.Printf("Error decoding rect (%s)", err) 
        }


        switch (name) {
        case "x":
            rect.value[0].SetValue(float64(value))
        case "y":
            rect.value[1].SetValue(float64(value))
        case "width":
            rect.value[2].SetValue(float64(value))
        case "height":
            rect.value[3].SetValue(float64(value))
        default:
            log.Printf("Unknown RectProp attr name (%s)\n", name)
        }
    }
}


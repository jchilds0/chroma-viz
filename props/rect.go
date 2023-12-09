package props

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type RectProp struct {
    value [4]*gtk.SpinButton
    input [4]*gtk.Box
    box *gtk.Box
}

func NewRectProp(width, height int, animate func()) Property {
    rect := &RectProp{}

    rect.value[0], _ = gtk.SpinButtonNewWithRange(float64(0), float64(width), 1)
    rect.value[1], _ = gtk.SpinButtonNewWithRange(float64(0), float64(height), 1)
    rect.value[2], _ = gtk.SpinButtonNewWithRange(float64(0), float64(width), 1)
    rect.value[3], _ = gtk.SpinButtonNewWithRange(float64(0), float64(height), 1)

    rect.input[0] = IntEditor("x Pos", 0, width, rect.value[0], animate)
    rect.input[1] = IntEditor("y Pos", 0, height, rect.value[1], animate)
    rect.input[2] = IntEditor("width", 0, width, rect.value[2], animate)
    rect.input[3] = IntEditor("height", 0, height, rect.value[3], animate)

    rect.box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

    for _, in := range rect.input {
        rect.box.PackStart(in, false, false, padding)
    }

    rect.box.SetVisible(true)

    return rect
}

func (rect *RectProp) Tab() *gtk.Box {
    return rect.box
}

func (rect *RectProp) String() string {
    return fmt.Sprintf("pos_x#%d#pos_y#%d#width#%d#height#%d#", 
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
        value, _ := strconv.Atoi(line[1])

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


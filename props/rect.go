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
    input [4]*gtk.Box
    num int
}

func NewRectProp(num, width, height int, animate func()) Property {
    rect := &RectProp{num: num}

    rect.value[0], _ = gtk.SpinButtonNewWithRange(float64(0), float64(width), 1)
    rect.value[1], _ = gtk.SpinButtonNewWithRange(float64(0), float64(height), 1)
    rect.value[2], _ = gtk.SpinButtonNewWithRange(float64(0), float64(width), 1)
    rect.value[3], _ = gtk.SpinButtonNewWithRange(float64(0), float64(height), 1)

    rect.box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    labels := [4]string{"x Pos", "y Pos", "Width", "Height"}

    for i := range rect.input {
        rect.input[i] = IntEditor(labels[i], rect.value[i], animate)
        rect.box.PackStart(rect.input[i], false, false, padding)
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


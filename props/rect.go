package props

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type RectEditor struct {
    box *gtk.Box
    value [4]*gtk.SpinButton
}

func NewRectEditor(width, height int, animate func()) PropertyEditor {
    var err error
    rect := &RectEditor{}

    rect.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Printf("Error creating rect box (%s)", err) 
    }

    rect.box.SetVisible(true)

    upper := []int{width, height, width, height}
    labels := [4]string{"x Pos", "y Pos", "Width", "Height"}
    for i := range rect.value {
        rect.value[i], err = gtk.SpinButtonNewWithRange(-float64(upper[i]), float64(upper[i]), 1)
        if err != nil { 
            log.Printf("Error creating rect spin button (%s)", err) 
        }

        input := IntEditor(labels[i], rect.value[i], animate)
        rect.box.PackStart(input, false, false, padding)
    }

    return rect
}

func (rectEdit *RectEditor) Box() *gtk.Box {
    return rectEdit.box
}

func (rectEdit *RectEditor) Update(rect Property) {
    rectProp, ok := rect.(*RectProp)
    if !ok {
        log.Printf("RectEditor.Update requires RectProp")
        return 
    }

    rectEdit.value[0].SetValue(float64(rectProp.value[0]))
    rectEdit.value[1].SetValue(float64(rectProp.value[1]))
    rectEdit.value[2].SetValue(float64(rectProp.value[2]))
    rectEdit.value[3].SetValue(float64(rectProp.value[3]))
}

type RectProp struct {
    name string
    value [4]int
}

func NewRectProp(name string) Property {
    rect := &RectProp{name: name}
    return rect
}

func (rect *RectProp) Type() int {
    return RECT_PROP
}

func (rect *RectProp) Name() string {
    return rect.name
}

func (rect *RectProp) String() string {
    return fmt.Sprintf("rel_x=%d#rel_y=%d#width=%d#height=%d#", 
        rect.value[0], rect.value[1], rect.value[2], rect.value[3])
}

func (rect *RectProp) Encode() string {
    return fmt.Sprintf("x %d;y %d;width %d;height %d;", 
        rect.value[0], rect.value[1], rect.value[2], rect.value[3])
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
            rect.value[0] = value
        case "y":
            rect.value[1] = value
        case "width":
            rect.value[2] = value
        case "height":
            rect.value[3] = value
        default:
            log.Printf("Unknown RectProp attr name (%s)\n", name)
        }
    }
}

func (rectProp *RectProp) Update(rect PropertyEditor, action int) {
    rectEdit, ok := rect.(*RectEditor)
    if !ok {
        log.Printf("RectProp.Update requires RectEditor")
        return
    }

    rectProp.value[0] = rectEdit.value[0].GetValueAsInt()
    rectProp.value[1] = rectEdit.value[1].GetValueAsInt()
    rectProp.value[2] = rectEdit.value[2].GetValueAsInt()
    rectProp.value[3] = rectEdit.value[3].GetValueAsInt()
}

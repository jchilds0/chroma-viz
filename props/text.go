package props

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type TextProp struct {
    entry *gtk.Entry
    input *gtk.Box
    box *gtk.Box
    num int
    x_spin *gtk.SpinButton
    y_spin *gtk.SpinButton
}

func NewTextProp(num, width, height int, animate func()) *TextProp {
    text := &TextProp{num: num}
    text.input, text.entry = TextEditor("Text: ", animate)

    text.box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    text.box.SetVisible(true)
    text.box.PackStart(text.input, false, false, padding)

    text.x_spin, _ = gtk.SpinButtonNewWithRange(float64(0), float64(width), 1)
    text.y_spin, _ = gtk.SpinButtonNewWithRange(float64(0), float64(height), 1)

    text.box.PackStart(IntEditor("x Pos", text.x_spin, animate), false, false, padding)
    text.box.PackStart(IntEditor("y Pos", text.y_spin, animate), false, false, padding)

    return text
}

func (text *TextProp) Tab() *gtk.Box {
    return text.box 
}

func (text *TextProp) String() string {
    entryText, _ := text.entry.GetText()

    return fmt.Sprintf("text=%d#string=%s#pos_x=%d#pos_y=%d#", 
        text.num, entryText, text.x_spin.GetValueAsInt(), text.y_spin.GetValueAsInt())
}
 
func (text *TextProp) Encode() string {
    entryText, _ := text.entry.GetText()
    return fmt.Sprintf("string %s;x %d;y %d;", 
        entryText, text.x_spin.GetValueAsInt(), text.y_spin.GetValueAsInt())
}

func (text *TextProp) Decode(input string) {
    attrs := strings.Split(input, ";")

    for _, attr := range attrs[1:] {
        line := strings.Split(attr, " ")
        name := line[0]

        switch (name) {
        case "x":
            value, _ := strconv.Atoi(line[1])
            text.x_spin.SetValue(float64(value))
        case "y":
            value, _ := strconv.Atoi(line[1])
            text.y_spin.SetValue(float64(value))
        case "string":
            text.entry.SetText(strings.TrimPrefix(attr, "string "))
        case "":
        default:
            log.Printf("Unknown TextProp attr name (%s)\n", name)
        }
    }
}


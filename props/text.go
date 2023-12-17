package props

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type TextProp struct {
    name string
    entry *gtk.Entry
    input *gtk.Box
    box *gtk.Box
    x_spin *gtk.SpinButton
    y_spin *gtk.SpinButton
}

func NewTextProp(width, height int, animate func(), name string) *TextProp {
    var err error
    text := &TextProp{name: name}
    text.input, text.entry = TextEditor("Text: ", animate)

    text.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Fatalf("Error creating text prop (%s)", err) 
    }

    text.box.SetVisible(true)
    text.box.PackStart(text.input, false, false, padding)

    text.x_spin, err = gtk.SpinButtonNewWithRange(float64(0), float64(width), 1)
    if err != nil { 
        log.Fatalf("Error creating text prop (%s)", err) 
    }

    text.y_spin, err = gtk.SpinButtonNewWithRange(float64(0), float64(height), 1)
    if err != nil { 
        log.Fatalf("Error creating text prop (%s)", err) 
    }

    text.box.PackStart(IntEditor("x Pos", text.x_spin, animate), false, false, padding)
    text.box.PackStart(IntEditor("y Pos", text.y_spin, animate), false, false, padding)

    return text
}

func (text *TextProp) Tab() *gtk.Box {
    return text.box 
}

func (text *TextProp) Name() string {
    return text.name
}

func (text *TextProp) String() string {
    entryText, err := text.entry.GetText()
    if err != nil { 
        log.Fatalf("Error creating text string (%s)", err) 
    }


    return fmt.Sprintf("string=%s#pos_x=%d#pos_y=%d#", 
        entryText, text.x_spin.GetValueAsInt(), text.y_spin.GetValueAsInt())
}
 
func (text *TextProp) Encode() string {
    entryText, err := text.entry.GetText()
    if err != nil { 
        log.Fatalf("Error encoding text prop (%s)", err) 
    }

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
            value, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Fatalf("Error decoding text prop (%s)", err) 
            }

            text.x_spin.SetValue(float64(value))
        case "y":
            value, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Printf("Error decoding text prop (%s)", err) 
            }

            text.y_spin.SetValue(float64(value))
        case "string":
            text.entry.SetText(strings.TrimPrefix(attr, "string "))
        case "":
        default:
            log.Printf("Unknown TextProp attr name (%s)\n", name)
        }
    }
}


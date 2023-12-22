package props

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type TextEditor struct {
    box *gtk.Box
    entry *gtk.Entry
    input *gtk.Box
    value [2]*gtk.SpinButton
}

func NewTextEditor(width, height int, animate func()) PropertyEditor {
    var err error
    text := &TextEditor{}

    text.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Fatalf("Error creating text prop (%s)", err) 
    }

    text.box.SetVisible(true)

    text.input, text.entry = StringEditor("Text: ", animate)
    text.box.PackStart(text.input, false, false, padding)

    text.value[0], err = gtk.SpinButtonNewWithRange(float64(0), float64(width), 1)
    if err != nil { 
        log.Fatalf("Error creating text prop (%s)", err) 
    }

    text.value[1], err = gtk.SpinButtonNewWithRange(float64(0), float64(height), 1)
    if err != nil { 
        log.Fatalf("Error creating text prop (%s)", err) 
    }

    text.box.PackStart(IntEditor("x Pos", text.value[0], animate), false, false, padding)
    text.box.PackStart(IntEditor("y Pos", text.value[1], animate), false, false, padding)

    return text
}

func (text *TextEditor) Box() *gtk.Box {
    return text.box
}

func (textEdit *TextEditor) Update(text Property) {
    textProp, ok := text.(*TextProp)
    if !ok {
        log.Printf("TextEditor.Update requires TextProp")
        return
    }

    textEdit.value[0].SetValue(float64(textProp.value[0]))
    textEdit.value[1].SetValue(float64(textProp.value[1]))
    textEdit.entry.SetText(textProp.str)
}

type TextProp struct {
    name string
    str string
    value [2]int
}

func NewTextProp(name string) *TextProp {
    text := &TextProp{name: name}
    return text
}

func (text *TextProp) Type() int {
    return TEXT_PROP
}

func (text *TextProp) Name() string {
    return text.name
}

func (text *TextProp) String() string {
    return fmt.Sprintf("string=%s#pos_x=%d#pos_y=%d#", 
        text.str, text.value[0], text.value[1])
}
 
func (text *TextProp) Encode() string {
    return fmt.Sprintf("string %s;x %d;y %d;", 
        text.str, text.value[0], text.value[1])
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

            text.value[0] = value
        case "y":
            value, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Printf("Error decoding text prop (%s)", err) 
            }

            text.value[1] = value
        case "string":
            text.str = strings.TrimPrefix(attr, "string ")
        case "":
        default:
            log.Printf("Unknown TextProp attr name (%s)\n", name)
        }
    }
}

func (textProp *TextProp) Update(text PropertyEditor, action int) {
    var err error
    textEdit, ok := text.(*TextEditor)
    if !ok {
        log.Printf("TextProp.Update requires TextEditor")
        return
    }

    textProp.value[0] = textEdit.value[0].GetValueAsInt()
    textProp.value[1] = textEdit.value[1].GetValueAsInt()
    textProp.str, err = textEdit.entry.GetText()
    if err != nil {
        log.Printf("Error getting text from editor entry (%s)", err)
        textProp.str = ""
    }
}

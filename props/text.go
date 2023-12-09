package props

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type TextProp struct {
    text []*gtk.Entry
    input []*gtk.Box
    numLines int
    box *gtk.Box
}

func NewTextProp(numLines int, animate func()) *TextProp {
    text := &TextProp{numLines: numLines}
    text.text = make([]*gtk.Entry, numLines)
    text.input = make([]*gtk.Box, numLines)

    for i := range text.text {
        text.input[i], text.text[i] = TextEditor("Line " + strconv.Itoa(i + 1), animate)
    }

    text.box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    text.box.SetVisible(true)

    for _, in := range text.input {
        text.box.PackStart(in, false, false, padding)
    }

    return text
}

func (text *TextProp) Tab() *gtk.Box {
    return text.box 
}

func (text *TextProp) String() string {
    str := ""
    for i, entry := range text.text {
        entryText, _ := entry.GetText()
        str = str + fmt.Sprintf("text%d#%s#", i, entryText)
    }

    return str
}
 
func (text *TextProp) Encode() string {
    str := ""
    for _, entry := range text.text {
        entryText, _ := entry.GetText()
        str = str + fmt.Sprintf("%s;", entryText)
    }

    return str
}

func (text *TextProp) Decode(input string) {
    strings := strings.Split(input, ";")

    for i := range text.text {
        text.text[i].SetText(strings[i + 1])
    }
}


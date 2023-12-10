package props 

import (
	"github.com/gotk3/gotk3/gtk"
)

const padding = 10

const (
    START = iota
    PAUSE
    STOP
)

func IntEditor(name string, spin *gtk.SpinButton, animate func()) *gtk.Box {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    box.SetVisible(true)

    label, _ := gtk.LabelNew(name)
    label.SetVisible(true)
    label.SetWidthChars(7)
    box.PackStart(label, false, false, uint(padding))

    spin.SetVisible(true)
    spin.SetValue(0)
    box.PackStart(spin, false, false, 0)

    spin.Connect("value-changed", animate)
    return box
}

func TextEditor(name string, animate func()) (*gtk.Box, *gtk.Entry) {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    box.SetVisible(true)

    label, _ := gtk.LabelNew(name)
    label.SetVisible(true)
    label.SetWidthChars(7)
    box.PackStart(label, false, false, uint(padding))

    buf, _ := gtk.EntryBufferNew("", 0)
    text, _ := gtk.EntryNewWithBuffer(buf)
    text.SetVisible(true)
    box.PackStart(text, false, false, 0)

    text.Connect("changed", animate)

    return box, text
}

type Property interface {
    Tab() *gtk.Box
    String() string
    Encode() string
    Decode(string)
}


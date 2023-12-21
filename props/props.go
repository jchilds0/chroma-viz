package props

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

const padding = 10

const (
    START = iota
    PAUSE
    STOP
)

const (
    END_OF_CONN = iota + 1
    END_OF_MESSAGE
    ANIMATE_ON
    CONTINUE
    ANIMATE_OFF
)

type Property interface {
    Tab() *gtk.Box
    Name() string
    String() string
    Encode() string
    Decode(string)
    Update(int)
}


func IntEditor(name string, spin *gtk.SpinButton, animate func()) *gtk.Box {
    box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil { 
        log.Printf("Error creating int editor (%s)", err) 
    }

    box.SetVisible(true)

    label, err := gtk.LabelNew(name)
    if err != nil { 
        log.Printf("Error creating int editor (%s)", err) 
    }

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
    box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil { 
        log.Printf("Error creating text editor (%s)", err) 
    }

    box.SetVisible(true)

    label, err := gtk.LabelNew(name)
    if err != nil { 
        log.Printf("Error creating text editor (%s)", err) 
    }

    label.SetVisible(true)
    label.SetWidthChars(7)
    box.PackStart(label, false, false, uint(padding))

    buf, err := gtk.EntryBufferNew("", 0)
    if err != nil { 
        log.Printf("Error creating text editor (%s)", err) 
    }

    text, _ := gtk.EntryNewWithBuffer(buf)
    if err != nil { 
        log.Printf("Error creating text editor (%s)", err) 
    }

    text.SetVisible(true)
    box.PackStart(text, false, false, 0)

    text.Connect("changed", animate)

    return box, text
}


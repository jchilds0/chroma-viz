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

const (
    RECT_PROP = iota
    TEXT_PROP
    CIRCLE_PROP 
    GRAPH_PROP
    TICKER_PROP
    CLOCK_PROP
    NUM_PROPS
)

type PropertyEditor interface {
    Update(Property)
    Box() *gtk.Box
}

type Property interface {
    Name() string
    Type() int
    String() string
    Encode() string
    Decode(string)
    Update(PropertyEditor, int)
}

func NewPropertyEditor(typed int, animate, cont func()) PropertyEditor {
    switch (typed) {
    case RECT_PROP:
        return NewRectEditor(1920, 1080, animate)
    case TEXT_PROP:
        return NewTextEditor(1920, 1080, animate)
    case CIRCLE_PROP:
        return NewCircleEditor(1920, 1080, animate)
    case GRAPH_PROP:
        return NewGraphEditor(1920, 1080, animate)
    case TICKER_PROP:
        return NewTickerEditor(1920, 1080, animate)
    case CLOCK_PROP:
        return NewClockEditor(1920, 1080, animate, cont)
    default:
        log.Printf("Unknown Prop %d", typed)
        return nil
    }
}

func NewProperty(typed int, name string) Property {
    switch (typed) {
    case RECT_PROP:
        return NewRectProp(name)
    case TEXT_PROP:
        return NewTextProp(name)
    case CIRCLE_PROP:
        return NewCircleProp(name)
    case GRAPH_PROP:
        return NewGraphProp(name)
    case TICKER_PROP:
        return NewTickerProp(name)
    case CLOCK_PROP:
        return NewClockProp(name)
    default:
        log.Printf("Unknown Prop %d", typed)
        return nil
    }
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

func StringEditor(name string, animate func()) (*gtk.Box, *gtk.Entry) {
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


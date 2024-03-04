package props

import (
	"chroma-viz/attribute"
	"fmt"
	"log"
	"strings"

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
    IMAGE_PROP
    NUM_PROPS
)

var StringToProp map[string]int = map[string]int{
    "rect": RECT_PROP,
    "text": TEXT_PROP,
    "circle": CIRCLE_PROP,
    "graph": GRAPH_PROP,
    "ticker": TICKER_PROP,
    "clock": CLOCK_PROP,
    "image": IMAGE_PROP,
}

type PropertyEditor interface {
    Box() *gtk.Box
    Editors() map[string]attribute.Editor
}

type Property interface {
    Name() string
    Type() int
    Visible() map[string]bool
    Attributes() map[string]attribute.Attribute
}

func NewPropertyEditor(typed int, animate, cont func()) (PropertyEditor, error) {
    width := 1920
    height := 1080 

    switch (typed) {
    case RECT_PROP:
        return NewRectEditor(width, height, animate)
    case TEXT_PROP:
        return NewTextEditor(width, height, animate)
    case CIRCLE_PROP:
        return NewCircleEditor(width, height, animate)
    case GRAPH_PROP:
        return NewGraphEditor(width, height, animate)
    case TICKER_PROP:
        return NewTickerEditor(width, height, animate)
    case CLOCK_PROP:
        return NewClockEditor(width, height, animate, cont)
    case IMAGE_PROP:
        return NewImageEditor(width, height, animate)
    default:
        return nil, fmt.Errorf("Unknown Prop %d", typed)
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
    case IMAGE_PROP:
        return NewImageProp(name)
    default:
        log.Printf("Unknown Prop %d", typed)
        return nil
    }
}

func PropToString(prop Property) (s string) {
    for _, attr := range prop.Attributes() {
        s = s + attr.String()
    }
    return
}

func EncodeProp(prop Property) (s string) {
    for _, attr := range prop.Attributes() {
        s = s + attr.Encode()
    }
    return
}

func DecodeProp(prop Property, s string) {
    attrs := strings.Split(s, ";")
    props := prop.Attributes()

    for _, attr := range attrs[1:] {
        name := strings.Split(attr, " ")[0]
        if props[name] == nil {
            log.Printf("Error prop %s missing prop attr %s", prop.Name(), name)
            continue
        }

        err := props[name].Decode(attr)
        if err != nil {
            log.Printf("Error decoding prop %s in %s", name, prop.Name())
        }
    }
}

func UpdateEditor(propEdit PropertyEditor, prop Property) {
    attrs := prop.Attributes()
    visible := prop.Visible()

    for name, edit := range propEdit.Editors() {
        edit.Box().SetVisible(visible[name])
        edit.Update(attrs[name])
    }
}

func UpdateProp(prop Property, propEdit PropertyEditor) {
    editors := propEdit.Editors()

    for name, attr := range prop.Attributes() {
        attr.Update(editors[name])
    }
}

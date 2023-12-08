package gui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

const padding = 10

type Property interface {
    Editor(string) *gtk.Box
    String() string
    Copy() Property
    Encode(string) []byte
}

type IntProp struct {
    name       string
    value      int
    lowerBound int
    upperBound int
}

func NewIntProp(name string, lowerBound, upperBound int) *IntProp {
    return &IntProp{name: name, lowerBound: lowerBound, upperBound: upperBound}
}

func (i *IntProp) Editor(name string) *gtk.Box {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    box.SetVisible(true)

    label, _ := gtk.LabelNew(name)
    label.SetVisible(true)
    label.SetWidthChars(7)
    box.PackStart(label, false, false, uint(padding))

    spin, _ := gtk.SpinButtonNewWithRange(float64(i.lowerBound), float64(i.upperBound), 1)
    spin.SetVisible(true)
    spin.SetValue(float64(i.value))
    spin.Connect("value-changed", func(spin *gtk.SpinButton) { i.value = spin.GetValueAsInt() } )
    box.PackStart(spin, false, false, 0)

    return box
}

func (i *IntProp) String() string {
    return fmt.Sprintf("%s#%d#", i.name, i.value)
}

func (i *IntProp) Copy() Property {
    newProp := NewIntProp(i.name, i.lowerBound, i.upperBound)
    newProp.value = i.value

    return newProp
}

func (i *IntProp) Encode(name string) []byte {
    return []byte(fmt.Sprintf("type int;name %s;value %d;\n", name, i.value))
}

type StrProp struct {
    name    string
    value   string
}

func NewStrProp(name string) *StrProp {
    return &StrProp{name: name}
}

func (str *StrProp) Editor(name string) *gtk.Box {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    box.SetVisible(true)

    label, _ := gtk.LabelNew(name)
    label.SetVisible(true)
    label.SetWidthChars(7)
    box.PackStart(label, false, false, uint(padding))

    buf, _ := gtk.EntryBufferNew(str.value, len(str.value))
    text, _ := gtk.EntryNewWithBuffer(buf)
    text.SetVisible(true)
    box.PackStart(text, false, false, 0)

    text.Connect("activate", 
        func(e *gtk.Entry) { 
            str.value, _ = e.GetText() 
        })
    return box
}

func (str *StrProp) String() string {
    return fmt.Sprintf("%s#%s#", str.name, str.value)
}

func (str *StrProp) Copy() Property {
    newProp := NewStrProp(str.name)
    newProp.value = str.value
    return newProp
}

func (str *StrProp) Encode(name string) []byte {
    return []byte(fmt.Sprintf("type string;name %s;value %s;\n", name, str.value))
}

func Decode(page *Page, input string) {
    props := strings.Split(input, ";")
    //log.Printf("%v\n", props)

    typed := strings.TrimPrefix(props[0], "type ")
    name := strings.TrimPrefix(props[1], "name ")
    prop, ok := page.props[name]
    if ok == false {
        log.Printf("property %s does not exist\n", name)
        return
    }

    switch (typed) {
    case "int":
        value := parse_int_value(props[2], "value")
        prop.(*IntProp).value = value
    case "string":
        value := strings.TrimPrefix(props[2], "value ")
        prop.(*StrProp).value = value
    }
}

func parse_int_value(input string, name string) int {
    value, err := strconv.Atoi(strings.TrimLeft(input, name + " "))
    if err != nil {
        log.Println(err)
    }

    return value
}

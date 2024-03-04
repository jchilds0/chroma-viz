package attribute

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type StringAttribute struct {
    name    string
    value   string
}

func NewStringAttribute(name string) *StringAttribute {
    stringAttr := &StringAttribute{name: name}
    return stringAttr
}

func (stringAttr *StringAttribute) String() string {
    return fmt.Sprintf("%s=%s#", stringAttr.name, stringAttr.value)
}

func (stringAttr *StringAttribute) Encode() string {
    return fmt.Sprintf("%s %s;", stringAttr.name, stringAttr.value)
}

func (stringAttr *StringAttribute) Decode(s string) error {
    line := strings.Split(s, " ")
    stringAttr.value = line[1]

    return nil
}

func (stringAttr *StringAttribute) Update(edit Editor) error {
    var err error
    stringEdit, ok := edit.(*StringEditor)
    if !ok {
        return fmt.Errorf("StringAttribute.Update requires StringEditor") 
    }

    stringAttr.value, err = stringEdit.Entry.GetText()
    return err
}

type StringEditor struct {
    box       *gtk.Box
    Entry     *gtk.Entry
}

func NewStringEditor(name string, animate func()) (stringEdit *StringEditor, err error) {
    stringEdit = &StringEditor{}

    stringEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil {
        return 
    }

    stringEdit.box.SetVisible(true)
    label, err := gtk.LabelNew(name)
    if err != nil { 
        return
    }

    label.SetVisible(true)
    label.SetWidthChars(12)
    stringEdit.box.PackStart(label, false, false, uint(padding))

    buf, err := gtk.EntryBufferNew("", 0)
    if err != nil { 
        return
    }

    stringEdit.Entry, err = gtk.EntryNewWithBuffer(buf)
    if err != nil { 
        return
    }

    stringEdit.Entry.SetVisible(true)
    stringEdit.box.PackStart(stringEdit.Entry, false, false, 0)
    stringEdit.Entry.Connect("changed", animate)
    return
}

func (stringEdit *StringEditor) Update(attr Attribute) error {
    stringAttr, ok := attr.(*StringAttribute)
    if !ok {
        return fmt.Errorf("StringEditor.Update requires StringAttribute") 
    }

    stringEdit.Entry.SetText(stringAttr.value)
    return nil
}

func (stringEdit *StringEditor) Box() *gtk.Box {
    return stringEdit.box
}

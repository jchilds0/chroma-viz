package attribute

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type StringAttribute struct {
    Name    string
    Type    int
    Value   string
}

func NewStringAttribute(name string) *StringAttribute {
    stringAttr := &StringAttribute{
        Name: name,
        Type: STRING,
    }

    return stringAttr
}

func (stringAttr *StringAttribute) String() string {
    return fmt.Sprintf("%s=%s#", stringAttr.Name, stringAttr.Value)
}

func (stringAttr *StringAttribute) Encode() string {
    return fmt.Sprintf("{'name': '%s', 'value': '%s'}", 
        stringAttr.Name, stringAttr.Value)
}

func (stringAttr *StringAttribute) Decode(value string) {
    stringAttr.Value = value
}

func (stringAttr *StringAttribute) Copy(attr Attribute) {
    stringAttrCopy, ok := attr.(*StringAttribute)
    if !ok {
        log.Print("Attribute not StringAttribute")
        return
    }

    stringAttr.Value = stringAttrCopy.Value
}

func (stringAttr *StringAttribute) Update(edit Editor) error {
    var err error
    stringEdit, ok := edit.(*StringEditor)
    if !ok {
        return fmt.Errorf("StringAttribute.Update requires StringEditor") 
    }

    stringAttr.Value, err = stringEdit.Entry.GetText()
    return err
}

type StringEditor struct {
    box       *gtk.Box
    Entry     *gtk.Entry
    name      string
}

func NewStringEditor(name string) *StringEditor {
    var err error
    stringEdit := &StringEditor{name: name}

    stringEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil {
        log.Print(err)
        return nil
    }

    stringEdit.box.SetVisible(true)
    label, err := gtk.LabelNew(name)
    if err != nil { 
        log.Print(err)
        return nil
    }

    label.SetVisible(true)
    label.SetWidthChars(12)
    stringEdit.box.PackStart(label, false, false, uint(padding))

    buf, err := gtk.EntryBufferNew("", 0)
    if err != nil { 
        log.Print(err)
        return nil
    }

    stringEdit.Entry, err = gtk.EntryNewWithBuffer(buf)
    if err != nil { 
        log.Print(err)
        return nil
    }

    stringEdit.Entry.SetVisible(true)
    stringEdit.box.PackStart(stringEdit.Entry, false, false, 0)

    return stringEdit
}

func (stringEdit *StringEditor) Name() string {
    return stringEdit.name
}

func (stringEdit *StringEditor) Update(attr Attribute) error {
    stringAttr, ok := attr.(*StringAttribute)
    if !ok {
        return fmt.Errorf("StringEditor.Update requires StringAttribute") 
    }

    stringEdit.Entry.SetText(stringAttr.Value)
    return nil
}

func (stringEdit *StringEditor) Box() *gtk.Box {
    return stringEdit.box
}

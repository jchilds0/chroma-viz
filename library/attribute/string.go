package attribute

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type StringAttribute struct {
	Name  string
	Value string
}

func (stringAttr *StringAttribute) Encode() string {
	return fmt.Sprintf("%s=%s#", stringAttr.Name, stringAttr.Value)
}

func (stringAttr *StringAttribute) UpdateAttribute(stringEdit StringEditor) (err error) {
	stringAttr.Value, err = stringEdit.Entry.GetText()
	return
}

type StringEditor struct {
	Box   *gtk.Box
	Entry *gtk.Entry
	Name  string
}

func NewStringEditor(name string) (stringEdit *StringEditor, err error) {
	stringEdit = &StringEditor{Name: name}

	stringEdit.Box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	stringEdit.Box.SetVisible(true)
	label, err := gtk.LabelNew(name)
	if err != nil {
		return
	}

	label.SetVisible(true)
	label.SetWidthChars(12)
	stringEdit.Box.PackStart(label, false, false, uint(padding))

	buf, err := gtk.EntryBufferNew("", 0)
	if err != nil {
		return
	}

	stringEdit.Entry, err = gtk.EntryNewWithBuffer(buf)
	if err != nil {
		return
	}

	stringEdit.Entry.SetVisible(true)
	stringEdit.Box.PackStart(stringEdit.Entry, false, false, 0)

	return
}

func (stringEdit *StringEditor) UpdateEditor(stringAttr StringAttribute) error {
	stringEdit.Entry.SetText(stringAttr.Value)
	return nil
}

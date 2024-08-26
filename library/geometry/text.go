package geometry

import (
	"chroma-viz/library/attribute"
	"strings"
)

type Text struct {
	Geometry

	String attribute.StringAttribute
	Color  attribute.ColorAttribute
	Scale  attribute.FloatAttribute
}

func (t *Text) UpdateGeometry(tEdit *TextEditor) (err error) {
	return
}

func (t *Text) EncodeEngine(b strings.Builder) {

}

type TextEditor struct {
	GeometryEditor

	String attribute.StringEditor
	Color  attribute.ColorEditor
	Scale  attribute.FloatEditor
}

func NewTextEditor() *TextEditor {
	return nil
}

func (tEdit *TextEditor) UpdateEditor(t *Text) (err error) {
	return
}

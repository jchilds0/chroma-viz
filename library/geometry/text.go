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

func NewText(geo Geometry) *Text {
	text := &Text{
		Geometry: geo,
	}

	text.String.Name = "string"
	text.Color.Name = "color"
	text.Scale.Name = "scale"
	return text
}

func (t *Text) UpdateGeometry(tEdit *TextEditor) (err error) {
	return
}

func (t *Text) EncodeEngine(b strings.Builder) {

}

func (t *Text) EncodeJSON(b strings.Builder) {

}

type TextEditor struct {
	GeometryEditor

	String attribute.StringEditor
	Color  attribute.ColorEditor
	Scale  attribute.FloatEditor
}

func NewTextEditor() (*TextEditor, error) {
	return nil, nil
}

func (tEdit *TextEditor) UpdateEditor(t *Text) (err error) {
	return
}

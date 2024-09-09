package geometry

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/util"
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

	text.String.Name = ATTR_STRING
	text.Color.Name = ATTR_COLOR
	text.Scale.Name = ATTR_SCALE
	text.Color.Alpha = 1.0

	return text
}

func (t *Text) UpdateGeometry(tEdit *TextEditor) (err error) {
	err = t.Geometry.UpdateGeometry(&tEdit.GeometryEditor)
	if err != nil {
		return
	}

	err = t.String.UpdateAttribute(tEdit.String)
	if err != nil {
		return
	}

	err = t.Color.UpdateAttribute(tEdit.Color)
	if err != nil {
		return
	}

	err = t.Scale.UpdateAttribute(tEdit.Scale)
	return
}

func (t *Text) Encode(b *strings.Builder) {
	t.Geometry.Encode(b)

	util.EngineAddKeyValue(b, t.String.Name, t.String.Value)
	util.EngineAddKeyValue(b, t.Scale.Name, t.Scale.Value)
	util.EngineAddKeyValue(b, t.Color.Name, t.Color.ToString())
}

type TextEditor struct {
	GeometryEditor

	String *attribute.StringEditor
	Color  *attribute.ColorEditor
	Scale  *attribute.FloatEditor
}

func NewTextEditor() (tEdit *TextEditor, err error) {
	geo, err := NewGeometryEditor()
	if err != nil {
		return
	}

	tEdit = &TextEditor{
		GeometryEditor: *geo,
	}

	tEdit.String, err = attribute.NewStringEditor(Attrs[ATTR_STRING])
	if err != nil {
		return
	}

	tEdit.Color, err = attribute.NewColorEditor(Attrs[ATTR_COLOR])
	if err != nil {
		return
	}

	tEdit.Scale, err = attribute.NewFloatEditor(Attrs[ATTR_SCALE], 0.0, 1.0, 0.01)
	if err != nil {
		return
	}

	geo.ScrollBox.PackStart(tEdit.String.Box, false, false, padding)
	geo.ScrollBox.PackStart(tEdit.Scale.Box, false, false, padding)
	geo.ScrollBox.PackStart(tEdit.Color.Box, false, false, padding)

	return
}

func (tEdit *TextEditor) UpdateEditor(t *Text) (err error) {
	err = tEdit.GeometryEditor.UpdateEditor(&t.Geometry)
	if err != nil {
		return
	}

	err = tEdit.String.UpdateEditor(&t.String)
	if err != nil {
		return
	}

	err = tEdit.Scale.UpdateEditor(&t.Scale)
	if err != nil {
		return
	}

	err = tEdit.Color.UpdateEditor(&t.Color)
	return
}

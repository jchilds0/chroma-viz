package geometry

import (
	"chroma-viz/library/attribute"
	"strings"
)

type Ticker struct {
	Geometry

	String attribute.ListAttribute
	Scale  attribute.FloatAttribute
	Color  attribute.ColorAttribute
}

func NewTicker(geo Geometry) *Ticker {
	ticker := &Ticker{
		Geometry: geo,
	}

	ticker.String.Name = ATTR_STRING
	ticker.Scale.Name = ATTR_SCALE
	ticker.Color.Name = ATTR_COLOR
	return ticker
}

func (t *Ticker) UpdateGeometry(tEdit *TickerEditor) (err error) {
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

func (t *Ticker) EncodeEngine(b strings.Builder) {

}

type TickerEditor struct {
	GeometryEditor

	String *attribute.ListEditor
	Scale  *attribute.FloatEditor
	Color  *attribute.ColorEditor
}

func NewTickerEditor() (tEdit *TickerEditor, err error) {
	geo, err := NewGeometryEditor()
	if err != nil {
		return
	}

	tEdit = &TickerEditor{
		GeometryEditor: *geo,
	}

	tEdit.String, err = attribute.NewListEditor("List", 1)
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

func (tEdit *TickerEditor) UpdateEditor(t *Ticker) (err error) {
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

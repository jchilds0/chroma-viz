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

func (t *Ticker) UpdateGeometry(tEdit *TickerEditor) (err error) {
	return
}

func (t *Ticker) EncodeEngine(b strings.Builder) {

}

func (t *Ticker) EncodeJSON(b strings.Builder) {

}

type TickerEditor struct {
	GeometryEditor

	String attribute.ListEditor
	Scale  attribute.FloatEditor
	Color  attribute.ColorEditor
}

func NewTickerEditor() (*TickerEditor, error) {
	return nil, nil
}

func (tEdit *TickerEditor) UpdateEditor(t *Ticker) (err error) {
	return
}

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

func (t *Ticker) UpdateGeometry(tEdit *TextEditor) (err error) {
	return
}

func (t *Ticker) EncodeEngine(b strings.Builder) {

}

type TickerEditor struct {
	Geometry

	String attribute.ListEditor
	Scale  attribute.FloatEditor
	Color  attribute.ColorEditor
}

func NewTickerEditor() *TickerEditor {
	return nil
}

func (tEdit *TickerEditor) UpdateEditor(t *Ticker) (err error) {
	return
}

package geometry

import (
	"chroma-viz/library/attribute"
	"strings"
)

type Clock struct {
	Geometry

	String attribute.ClockAttribute
	Color  attribute.ColorAttribute
	Scale  attribute.FloatAttribute
}

func (c *Clock) UpdateGeometry(cEdit *ClockEditor) (err error) {
	return
}

func (c *Clock) EncodeEngine(b strings.Builder) {

}

func (c *Clock) EncodeJSON(b strings.Builder) {

}

type ClockEditor struct {
	GeometryEditor

	String attribute.ClockEditor
	Color  attribute.ColorEditor
	Scale  attribute.FloatEditor
}

func NewClockEditor() (*ClockEditor, error) {
	return nil, nil
}

func (cEdit *ClockEditor) UpdateEditor(c *Clock) (err error) {
	return nil
}

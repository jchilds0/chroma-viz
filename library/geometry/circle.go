package geometry

import (
	"chroma-viz/library/attribute"
	"strings"
)

type Circle struct {
	Geometry

	InnerRadius attribute.IntAttribute
	OuterRadius attribute.IntAttribute
	StartAngle  attribute.IntAttribute
	EndAngle    attribute.IntAttribute
	Color       attribute.ColorAttribute
}

func (c *Circle) UpdateGeometry(cEdit *CircleEditor) (err error) {
	return
}

func (c *Circle) EncodeEngine(b strings.Builder) {
}

type CircleEditor struct {
	GeometryEditor

	InnerRadius attribute.IntEditor
	OuterRadius attribute.IntEditor
	StartAngle  attribute.IntEditor
	EndAngle    attribute.IntEditor
	Color       attribute.ColorEditor
}

func NewCircleEditor() *CircleEditor {
	return nil
}

func (cEdit *CircleEditor) UpdateEditor(c *Circle) (err error) {
	return
}

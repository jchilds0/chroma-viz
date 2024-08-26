package geometry

import (
	"chroma-viz/library/attribute"
	"strings"
)

type Rectangle struct {
	Geometry

	Width    attribute.IntAttribute
	Height   attribute.IntAttribute
	Rounding attribute.IntAttribute
	Color    attribute.ColorAttribute
}

func (r *Rectangle) UpdateGeometry(rEdit *RectangleEditor) (err error) {
	return
}

func (r *Rectangle) EncodeEngine(b strings.Builder) {

}

type RectangleEditor struct {
	GeometryEditor

	Width    attribute.IntEditor
	Height   attribute.IntEditor
	Rounding attribute.IntEditor
	Color    attribute.ColorAttribute
}

func NewRectangleEditor() *RectangleEditor {
	return nil
}

func (rEdit *RectangleEditor) UpdateEditor(r *Rectangle) (err error) {
	return
}

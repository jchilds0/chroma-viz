package geometry

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/util"
	"strings"
)

type Rectangle struct {
	Geometry

	Width    attribute.IntAttribute
	Height   attribute.IntAttribute
	Rounding attribute.IntAttribute
	Color    attribute.ColorAttribute
}

func NewRectangle(geo Geometry) *Rectangle {
	rect := &Rectangle{
		Geometry: geo,
	}

	rect.Width.Name = ATTR_WIDTH
	rect.Height.Name = ATTR_HEIGHT
	rect.Rounding.Name = ATTR_ROUND
	rect.Color.Name = ATTR_COLOR
	rect.Color.Alpha = 1.0

	return rect
}

func (r *Rectangle) UpdateGeometry(rEdit *RectangleEditor) (err error) {
	err = r.Geometry.UpdateGeometry(&rEdit.GeometryEditor)
	if err != nil {
		return
	}

	err = r.Width.UpdateAttribute(rEdit.Width)
	if err != nil {
		return
	}

	err = r.Height.UpdateAttribute(rEdit.Height)
	if err != nil {
		return
	}

	err = r.Rounding.UpdateAttribute(rEdit.Rounding)
	if err != nil {
		return
	}

	err = r.Color.UpdateAttribute(rEdit.Color)
	return
}

func (r *Rectangle) Encode(b *strings.Builder) {
	r.Geometry.Encode(b)

	util.EngineAddKeyValue(b, r.Width.Name, r.Width.Value)
	util.EngineAddKeyValue(b, r.Height.Name, r.Height.Value)
	util.EngineAddKeyValue(b, r.Rounding.Name, r.Rounding.Value)
	util.EngineAddKeyValue(b, r.Color.Name, r.Color.ToString())
}

type RectangleEditor struct {
	GeometryEditor

	Width    *attribute.IntEditor
	Height   *attribute.IntEditor
	Rounding *attribute.IntEditor
	Color    *attribute.ColorEditor
}

func NewRectangleEditor() (rectEdit *RectangleEditor, err error) {
	geoEdit, err := NewGeometryEditor()
	if err != nil {
		return
	}

	rectEdit = &RectangleEditor{
		GeometryEditor: *geoEdit,
	}

	rectEdit.Width, err = attribute.NewIntEditor(Attrs[ATTR_WIDTH])
	if err != nil {
		return
	}

	rectEdit.Height, err = attribute.NewIntEditor(Attrs[ATTR_HEIGHT])
	if err != nil {
		return
	}

	rectEdit.Rounding, err = attribute.NewIntEditor(Attrs[ATTR_ROUND])
	if err != nil {
		return
	}

	rectEdit.Color, err = attribute.NewColorEditor(Attrs[ATTR_COLOR])
	if err != nil {
		return
	}

	geoEdit.ScrollBox.PackStart(rectEdit.Width.Box, false, false, padding)
	geoEdit.ScrollBox.PackStart(rectEdit.Height.Box, false, false, padding)
	geoEdit.ScrollBox.PackStart(rectEdit.Rounding.Box, false, false, padding)
	geoEdit.ScrollBox.PackStart(rectEdit.Color.Box, false, false, padding)

	return
}

func (rEdit *RectangleEditor) UpdateEditor(r *Rectangle) (err error) {
	err = rEdit.GeometryEditor.UpdateEditor(&r.Geometry)
	if err != nil {
		return
	}

	err = rEdit.Width.UpdateEditor(&r.Width)
	if err != nil {
		return
	}

	err = rEdit.Height.UpdateEditor(&r.Height)
	if err != nil {
		return
	}

	err = rEdit.Rounding.UpdateEditor(&r.Rounding)
	if err != nil {
		return
	}

	err = rEdit.Color.UpdateEditor(&r.Color)
	return

}

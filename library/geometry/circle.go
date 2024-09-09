package geometry

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/util"
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

func NewCircle(geo Geometry) *Circle {
	circle := &Circle{
		Geometry: geo,
	}

	circle.InnerRadius.Name = ATTR_INNER_RADIUS
	circle.OuterRadius.Name = ATTR_OUTER_RADIUS
	circle.StartAngle.Name = ATTR_START_ANGLE
	circle.EndAngle.Name = ATTR_END_ANGLE
	circle.Color.Name = ATTR_COLOR
	circle.Color.Alpha = 1.0

	return circle
}

func (c *Circle) UpdateGeometry(cEdit *CircleEditor) (err error) {
	err = c.Geometry.UpdateGeometry(&cEdit.GeometryEditor)
	if err != nil {
		return
	}

	err = c.InnerRadius.UpdateAttribute(cEdit.InnerRadius)
	if err != nil {
		return
	}

	err = c.OuterRadius.UpdateAttribute(cEdit.OuterRadius)
	if err != nil {
		return
	}

	err = c.StartAngle.UpdateAttribute(cEdit.StartAngle)
	if err != nil {
		return
	}

	err = c.EndAngle.UpdateAttribute(cEdit.EndAngle)
	if err != nil {
		return
	}

	err = c.Color.UpdateAttribute(cEdit.Color)
	return
}

func (c *Circle) Encode(b *strings.Builder) {
	c.Geometry.Encode(b)

	util.EngineAddKeyValue(b, c.InnerRadius.Name, c.InnerRadius.Value)
	util.EngineAddKeyValue(b, c.OuterRadius.Name, c.OuterRadius.Value)
	util.EngineAddKeyValue(b, c.StartAngle.Name, c.StartAngle.Value)
	util.EngineAddKeyValue(b, c.EndAngle.Name, c.EndAngle.Value)
	util.EngineAddKeyValue(b, c.Color.Name, c.Color.ToString())
}

type CircleEditor struct {
	GeometryEditor

	InnerRadius *attribute.IntEditor
	OuterRadius *attribute.IntEditor
	StartAngle  *attribute.IntEditor
	EndAngle    *attribute.IntEditor
	Color       *attribute.ColorEditor
}

func NewCircleEditor() (cEdit *CircleEditor, err error) {
	geoEdit, err := NewGeometryEditor()
	if err != nil {
		return
	}

	cEdit = &CircleEditor{
		GeometryEditor: *geoEdit,
	}

	cEdit.InnerRadius, err = attribute.NewIntEditor(Attrs[ATTR_INNER_RADIUS])
	if err != nil {
		return
	}

	cEdit.OuterRadius, err = attribute.NewIntEditor(Attrs[ATTR_OUTER_RADIUS])
	if err != nil {
		return
	}

	cEdit.StartAngle, err = attribute.NewIntEditor(Attrs[ATTR_START_ANGLE])
	if err != nil {
		return
	}

	cEdit.EndAngle, err = attribute.NewIntEditor(Attrs[ATTR_END_ANGLE])
	if err != nil {
		return
	}

	cEdit.Color, err = attribute.NewColorEditor(Attrs[ATTR_COLOR])
	if err != nil {
		return
	}

	geoEdit.ScrollBox.PackStart(cEdit.InnerRadius.Box, false, false, padding)
	geoEdit.ScrollBox.PackStart(cEdit.OuterRadius.Box, false, false, padding)
	geoEdit.ScrollBox.PackStart(cEdit.StartAngle.Box, false, false, padding)
	geoEdit.ScrollBox.PackStart(cEdit.EndAngle.Box, false, false, padding)
	geoEdit.ScrollBox.PackStart(cEdit.Color.Box, false, false, padding)

	return
}

func (cEdit *CircleEditor) UpdateEditor(c *Circle) (err error) {
	err = cEdit.GeometryEditor.UpdateEditor(&c.Geometry)
	if err != nil {
		return
	}

	err = cEdit.InnerRadius.UpdateEditor(&c.InnerRadius)
	if err != nil {
		return
	}

	err = cEdit.OuterRadius.UpdateEditor(&c.OuterRadius)
	if err != nil {
		return
	}

	err = cEdit.OuterRadius.UpdateEditor(&c.OuterRadius)
	if err != nil {
		return
	}

	err = cEdit.StartAngle.UpdateEditor(&c.StartAngle)
	if err != nil {
		return
	}

	err = cEdit.EndAngle.UpdateEditor(&c.EndAngle)
	if err != nil {
		return
	}

	err = cEdit.Color.UpdateEditor(&c.Color)
	if err != nil {
		return
	}

	return
}

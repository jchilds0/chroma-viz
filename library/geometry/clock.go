package geometry

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/util"
	"strings"
)

type Clock struct {
	Geometry

	Clock attribute.ClockAttribute
	Color attribute.ColorAttribute
	Scale attribute.FloatAttribute
}

func NewClock(geo Geometry) (c *Clock) {
	c = &Clock{
		Geometry: geo,
	}

	c.Clock.Name = ATTR_STRING
	c.Color.Name = ATTR_COLOR
	c.Scale.Name = ATTR_SCALE
	c.Color.Alpha = 1.0

	return c
}

func (c *Clock) UpdateGeometry(cEdit *ClockEditor) (err error) {
	err = c.Geometry.UpdateGeometry(&cEdit.GeometryEditor)
	if err != nil {
		return
	}

	err = c.Clock.UpdateAttribute(cEdit.Clock)
	if err != nil {
		return
	}

	err = c.Color.UpdateAttribute(cEdit.Color)
	if err != nil {
		return
	}

	err = c.Scale.UpdateAttribute(cEdit.Scale)
	return
}

func (c *Clock) Encode(b *strings.Builder) {
	c.Geometry.Encode(b)

	util.EngineAddKeyValue(b, c.Clock.Name, c.Clock.CurrentTime)
	util.EngineAddKeyValue(b, c.Color.Name, c.Color.ToString())
	util.EngineAddKeyValue(b, c.Scale.Name, c.Scale.Value)
}

type ClockEditor struct {
	GeometryEditor

	Clock *attribute.ClockEditor
	Color *attribute.ColorEditor
	Scale *attribute.FloatEditor
}

func NewClockEditor() (cEdit *ClockEditor, err error) {
	geoEdit, err := NewGeometryEditor()
	if err != nil {
		return
	}

	cEdit = &ClockEditor{
		GeometryEditor: *geoEdit,
	}

	cEdit.Clock, err = attribute.NewClockEditor("Time")
	if err != nil {
		return
	}

	cEdit.Color, err = attribute.NewColorEditor(Attrs[ATTR_COLOR])
	if err != nil {
		return
	}

	cEdit.Scale, err = attribute.NewFloatEditor(Attrs[ATTR_SCALE], 0.0, 1.0, 0.01)
	if err != nil {
		return
	}

	geoEdit.ScrollBox.PackStart(cEdit.Clock.Box, false, false, padding)
	geoEdit.ScrollBox.PackStart(cEdit.Color.Box, false, false, padding)
	geoEdit.ScrollBox.PackStart(cEdit.Scale.Box, false, false, padding)

	return
}

func (cEdit *ClockEditor) UpdateEditor(c *Clock) (err error) {
	err = cEdit.GeometryEditor.UpdateEditor(&c.Geometry)
	if err != nil {
		return
	}

	err = cEdit.Clock.UpdateEditor(&c.Clock)
	if err != nil {
		return
	}

	err = cEdit.Color.UpdateEditor(&c.Color)
	if err != nil {
		return
	}

	err = cEdit.Scale.UpdateEditor(&c.Scale)
	return
}

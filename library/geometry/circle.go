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

func NewCircle(geo Geometry) Circle {
	circle := Circle{
		Geometry: geo,
	}

	circle.InnerRadius.Name = "inner_radius"
	circle.OuterRadius.Name = "outer_radius"
	circle.StartAngle.Name = "start_angle"
	circle.EndAngle.Name = "end_angle"
	circle.Color.Name = "color"
	return circle
}

func (c *Circle) UpdateGeometry(cEdit *CircleEditor) (err error) {
	err = c.Geometry.UpdateGeometry(&cEdit.GeometryEditor)
	if err != nil {
		return
	}

	err = c.InnerRadius.UpdateAttribute(&cEdit.InnerRadius)
	if err != nil {
		return
	}

	err = c.OuterRadius.UpdateAttribute(&cEdit.OuterRadius)
	if err != nil {
		return
	}

	err = c.StartAngle.UpdateAttribute(&cEdit.StartAngle)
	if err != nil {
		return
	}

	err = c.EndAngle.UpdateAttribute(&cEdit.EndAngle)
	if err != nil {
		return
	}

	err = c.Color.UpdateAttribute(&cEdit.Color)
	return
}

func (c *Circle) EncodeEngine(b strings.Builder) {

}

func (c *Circle) EncodeJSON(b strings.Builder) {

}

type CircleEditor struct {
	GeometryEditor

	InnerRadius attribute.IntEditor
	OuterRadius attribute.IntEditor
	StartAngle  attribute.IntEditor
	EndAngle    attribute.IntEditor
	Color       attribute.ColorEditor
}

func NewCircleEditor() (*CircleEditor, error) {
	return nil, nil
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

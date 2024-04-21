package templates

import (
	"chroma-viz/library/attribute"
)

type Geometry struct {
	GeoID   int64
	Name    string
	GeoType int
	RelX    int
	RelY    int
	Color   [4]byte
	Parent  int
}

type GeometryEncoder interface {
	Encode() map[string]attribute.Attribute
}

func NewGeometry(name string, geoType, rel_x, rel_y int, r, g, b, a byte, parent int) *Geometry {
	geo := &Geometry{
		Name:    name,
		GeoType: geoType,
		RelX:    rel_x,
		RelY:    rel_y,
		Color:   [4]byte{r, g, b, a},
		Parent:  parent,
	}

	return geo
}

func (geo *Geometry) Encode() map[string]attribute.Attribute {
	p := make(map[string]attribute.Attribute, 10)

	x := attribute.NewIntAttribute("rel_x")
	x.Value = geo.RelX
	p[x.Name] = x

	y := attribute.NewIntAttribute("rel_y")
	y.Value = geo.RelY
	p[y.Name] = y

	color := attribute.NewColorAttribute("color")
	color.Red = float64(geo.Color[0]) / 255
	color.Green = float64(geo.Color[1]) / 255
	color.Blue = float64(geo.Color[2]) / 255
	color.Alpha = float64(geo.Color[3]) / 255
	p[color.Name] = color

	return p
}

type Rectangle struct {
	Geometry
	Width    int
	Height   int
	Rounding int
}

func NewRectangle(geo Geometry, width, height, rounding int) *Rectangle {
	rect := &Rectangle{
		Geometry: geo,
		Width:    width,
		Height:   height,
		Rounding: rounding,
	}

	return rect
}

func (rect *Rectangle) Encode() map[string]attribute.Attribute {
	p := rect.Geometry.Encode()

	width := attribute.NewIntAttribute("width")
	width.Value = rect.Width
	p[width.Name] = width

	height := attribute.NewIntAttribute("height")
	height.Value = rect.Height
	p[height.Name] = height

	rounding := attribute.NewIntAttribute("rounding")
	rounding.Value = rect.Rounding
	p[rounding.Name] = rounding

	return p
}

type Text struct {
	Geometry
	Text string
}

func NewText(geo Geometry, text string) *Text {
	t := &Text{
		Geometry: geo,
		Text:     text,
	}

	return t
}

func (text *Text) Encode() map[string]attribute.Attribute {
	p := text.Geometry.Encode()

	t := attribute.NewStringAttribute("string")
	t.Value = text.Text
	p["string"] = t

	return p
}

type Circle struct {
	Geometry
	InnerRadius int
	OuterRadius int
	StartAngle  int
	EndAngle    int
}

func NewCircle(geo Geometry, innerRadius, outerRadius, startAngle, endAngle int) *Circle {
	c := &Circle{
		Geometry:    geo,
		InnerRadius: innerRadius,
		OuterRadius: outerRadius,
		StartAngle:  startAngle,
		EndAngle:    endAngle,
	}

	return c
}

func (circle *Circle) Encode() map[string]attribute.Attribute {
	p := circle.Geometry.Encode()

	innerRadius := attribute.NewIntAttribute("inner_radius")
	innerRadius.Value = circle.InnerRadius
	p[innerRadius.Name] = innerRadius

	outerRadius := attribute.NewIntAttribute("outer_radius")
	outerRadius.Value = circle.OuterRadius
	p[outerRadius.Name] = outerRadius

	startAngle := attribute.NewIntAttribute("start_angle")
	startAngle.Value = circle.StartAngle
	p[startAngle.Name] = startAngle

	endAngle := attribute.NewIntAttribute("end_angle")
	endAngle.Value = circle.EndAngle
	p[endAngle.Name] = endAngle

	return p
}

package templates

import (
	"strconv"
	"strings"
)

const (
	GEO_RECT = iota
	GEO_CIRCLE
	GEO_TEXT
)

var GeoName = map[int]string{
	GEO_RECT:   "rectangle",
	GEO_CIRCLE: "circle",
	GEO_TEXT:   "text",
}

type Geometry struct {
	GeoNum   int
	Name     string
	GeoType  int
	PropType int
	RelX     int
	RelY     int
	Parent   int
}

func NewGeometry(geoID int, name string, propType, geoType, rel_x, rel_y int, parent int) *Geometry {
	geo := &Geometry{
		GeoNum:   geoID,
		Name:     name,
		GeoType:  geoType,
		PropType: propType,
		RelX:     rel_x,
		RelY:     rel_y,
		Parent:   parent,
	}

	return geo
}

func (geo *Geometry) Geom() *Geometry {
	return geo
}

func (geo *Geometry) Attributes() map[string]string {
	p := make(map[string]string, 10)

	p["rel_x"] = strconv.Itoa(geo.RelX)
	p["rel_y"] = strconv.Itoa(geo.RelY)
	p["parent"] = strconv.Itoa(geo.Parent)

	return p
}

func EncodeGeometry(geo Geometry, attr map[string]string) string {
	var b strings.Builder

	b.WriteString("{")
	b.WriteString("'id': ")
	b.WriteString(strconv.Itoa(int(geo.GeoNum)))
	b.WriteString(", ")

	b.WriteString("'name': '")
	b.WriteString(geo.Name)
	b.WriteString("', ")

	b.WriteString("'prop_type': ")
	b.WriteString(strconv.Itoa(geo.PropType))
	b.WriteString(", ")

	b.WriteString("'geo_type': ")
	b.WriteString(strconv.Itoa(geo.GeoType))
	b.WriteString(", ")

	// TODO: Visible

	first := true
	b.WriteString("'attr': [")
	for name, value := range attr {
		if !first {
			b.WriteString(",")
		}

		first = false
		b.WriteString("{'name': '")
		b.WriteString(name)
		b.WriteString("', 'value': '")
		b.WriteString(value)
		b.WriteString("'}")
	}

	b.WriteString("]}")

	return b.String()
}

type Rectangle struct {
	Geometry
	Width    int
	Height   int
	Rounding int
	Color    string
}

func NewRectangle(geo Geometry, width, height, rounding int, color string) *Rectangle {
	rect := &Rectangle{
		Geometry: geo,
		Width:    width,
		Height:   height,
		Rounding: rounding,
		Color:    color,
	}

	return rect
}

func (rect *Rectangle) Attributes() map[string]string {
	p := rect.Geometry.Attributes()

	p["width"] = strconv.Itoa(rect.Width)
	p["height"] = strconv.Itoa(rect.Height)
	p["rounding"] = strconv.Itoa(rect.Rounding)
	p["color"] = rect.Color

	return p
}

type Text struct {
	Geometry
	Text  string
	Color string
}

func NewText(geo Geometry, text, color string) *Text {
	t := &Text{
		Geometry: geo,
		Text:     text,
		Color:    color,
	}

	return t
}

func (text *Text) Attributes() map[string]string {
	p := text.Geometry.Attributes()

	p["string"] = text.Text
	p["color"] = text.Color

	return p
}

type Circle struct {
	Geometry
	InnerRadius int
	OuterRadius int
	StartAngle  int
	EndAngle    int
	Color       string
}

func NewCircle(geo Geometry, innerRadius, outerRadius, startAngle, endAngle int, color string) *Circle {
	c := &Circle{
		Geometry:    geo,
		InnerRadius: innerRadius,
		OuterRadius: outerRadius,
		StartAngle:  startAngle,
		EndAngle:    endAngle,
		Color:       color,
	}

	return c
}

func (circle *Circle) Attributes() map[string]string {
	p := circle.Geometry.Attributes()

	p["inner_radius"] = strconv.Itoa(circle.InnerRadius)
	p["outer_radius"] = strconv.Itoa(circle.OuterRadius)
	p["start_angle"] = strconv.Itoa(circle.StartAngle)
	p["end_angle"] = strconv.Itoa(circle.EndAngle)
	p["color"] = circle.Color

	return p
}

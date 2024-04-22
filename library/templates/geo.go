package templates

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	GEO_RECT = iota
	GEO_CIRCLE
	GEO_TEXT
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

type IGeometry interface {
	Geom() Geometry
	Attributes() map[string]string
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

func (geo Geometry) Geom() Geometry {
	return geo
}

func (geo *Geometry) Attributes() map[string]string {
	p := make(map[string]string, 10)

	p["rel_x"] = strconv.Itoa(geo.RelX)
	p["rel_y"] = strconv.Itoa(geo.RelY)
	p["parent"] = strconv.Itoa(geo.Parent)

	p["color"] = fmt.Sprintf("%f %f %f %f", geo.Color[0]/255,
		geo.Color[1]/255, geo.Color[2]/255, geo.Color[3]/255)

	return p
}

func EncodeGeometry(geo IGeometry) string {
	var b strings.Builder

	geom := geo.Geom()

	b.WriteString("{")
	b.WriteString("'id': ")
	b.WriteString(strconv.Itoa(int(geom.GeoID)))
	b.WriteString(", ")

	b.WriteString("'name': '")
	b.WriteString(geom.Name)
	b.WriteString("', ")

	b.WriteString("'prop_type': '")
	b.WriteString("', ")

	b.WriteString("'geo_type': '")
	b.WriteString("', ")

	// TODO: Visible

	first := true
	b.WriteString("'attr': [")
	for name, value := range geo.Attributes() {
		if !first {
			b.WriteString(",")
		}

		first = false
		b.WriteString("{'")
		b.WriteString(name)
		b.WriteString("': '")
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

func (rect *Rectangle) Attributes() map[string]string {
	p := rect.Geometry.Attributes()

	p["width"] = strconv.Itoa(rect.Width)
	p["height"] = strconv.Itoa(rect.Height)
	p["rounding"] = strconv.Itoa(rect.Rounding)

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

func (text *Text) Attributes() map[string]string {
	p := text.Geometry.Attributes()

	p["string"] = text.Text

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

func (circle *Circle) Attributes() map[string]string {
	p := circle.Geometry.Attributes()

	p["inner_radius"] = strconv.Itoa(circle.InnerRadius)
	p["outer_radius"] = strconv.Itoa(circle.OuterRadius)
	p["start_angle"] = strconv.Itoa(circle.StartAngle)
	p["end_angle"] = strconv.Itoa(circle.EndAngle)

	return p
}

package templates

import (
	"fmt"
	"strconv"
	"strings"
)

type AttrEntry struct {
	Name  string
	Value string
	Index string
}

const (
	GEO_RECT   = "rect"
	GEO_CIRCLE = "circle"
	GEO_TEXT   = "text"
	GEO_IMAGE  = "image"
	GEO_POLY   = "polygon"
)

type Geometry struct {
	GeoNum   int
	Name     string
	GeoType  string
	PropType string
	RelX     int
	RelY     int
	Parent   int
	Mask     int
}

func (geo *Geometry) Geom() *Geometry {
	return geo
}

func (geo *Geometry) Attributes() []AttrEntry {
	p := make([]AttrEntry, 10)

	p = append(p, AttrEntry{Name: "rel_x", Value: strconv.Itoa(geo.RelX)})
	p = append(p, AttrEntry{Name: "rel_y", Value: strconv.Itoa(geo.RelY)})
	p = append(p, AttrEntry{Name: "parent", Value: strconv.Itoa(geo.Parent)})
	p = append(p, AttrEntry{Name: "mask", Value: strconv.Itoa(geo.Mask)})

	return p
}

func EncodeGeometry(geo Geometry, attr []AttrEntry) string {
	var b strings.Builder

	b.WriteString("{")
	b.WriteString("'id': ")
	b.WriteString(strconv.Itoa(int(geo.GeoNum)))
	b.WriteString(", ")

	b.WriteString("'name': '")
	b.WriteString(geo.Name)
	b.WriteString("', ")

	b.WriteString("'prop_type': '")
	b.WriteString(geo.PropType)
	b.WriteString("', ")

	b.WriteString("'geo_type': '")
	b.WriteString(geo.GeoType)
	b.WriteString("', ")

	// TODO: Visible

	first := true
	b.WriteString("'attr': [")
	for _, entry := range attr {
		if !first {
			b.WriteString(",")
		}

		first = false
		b.WriteString("{'name': '")
		b.WriteString(entry.Name)

		if entry.Index != "" {
			b.WriteString(", 'index': ")
			b.WriteString(entry.Index)
		}

		b.WriteString("', 'value': '")
		b.WriteString(entry.Value)

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

func (rect *Rectangle) Attributes() []AttrEntry {
	p := rect.Geometry.Attributes()

	p = append(p, AttrEntry{Name: "width", Value: strconv.Itoa(rect.Width)})
	p = append(p, AttrEntry{Name: "height", Value: strconv.Itoa(rect.Height)})
	p = append(p, AttrEntry{Name: "rounding", Value: strconv.Itoa(rect.Rounding)})
	p = append(p, AttrEntry{Name: "color", Value: rect.Color})

	return p
}

type Text struct {
	Geometry
	Text  string
	Scale float64
	Color string
}

func NewText(geo Geometry, text, color string, scale float64) *Text {
	t := &Text{
		Geometry: geo,
		Text:     text,
		Scale:    scale,
		Color:    color,
	}

	return t
}

func (text *Text) Attributes() []AttrEntry {
	p := text.Geometry.Attributes()

	p = append(p, AttrEntry{Name: "string", Value: text.Text})
	p = append(p, AttrEntry{Name: "scale", Value: strconv.FormatFloat(text.Scale, 'f', 10, 64)})
	p = append(p, AttrEntry{Name: "color", Value: text.Color})

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

func (circle *Circle) Attributes() []AttrEntry {
	p := circle.Geometry.Attributes()

	p = append(p, AttrEntry{Name: "inner_radius", Value: strconv.Itoa(circle.InnerRadius)})
	p = append(p, AttrEntry{Name: "outer_radius", Value: strconv.Itoa(circle.OuterRadius)})
	p = append(p, AttrEntry{Name: "start_angle", Value: strconv.Itoa(circle.StartAngle)})
	p = append(p, AttrEntry{Name: "end_angle", Value: strconv.Itoa(circle.EndAngle)})
	p = append(p, AttrEntry{Name: "color", Value: circle.Color})

	return p
}

type Asset struct {
	Geometry
	Dir   string
	Name  string
	ID    int
	Scale float64
}

func NewAsset(geo Geometry, name, dir string, id int, scale float64) *Asset {
	a := &Asset{
		Geometry: geo,
		Dir:      dir,
		Name:     name,
		ID:       id,
		Scale:    scale,
	}

	return a
}

func (a *Asset) Attributes() []AttrEntry {
	p := a.Geometry.Attributes()

	p = append(p, AttrEntry{Name: "image_id", Value: strconv.Itoa(a.ID)})
	p = append(p, AttrEntry{Name: "scale", Value: strconv.FormatFloat(a.Scale, 'f', 10, 64)})

	return p
}

type Point struct {
	PointIndex int
	PosX       int
	PosY       int
}

type Polygon struct {
	Geometry
	Points map[int]Point
}

func NewPolygon(geo Geometry, numPoints int) *Polygon {
	p := &Polygon{
		Geometry: geo,
		Points:   make(map[int]Point, numPoints),
	}

	return p
}

func (p *Polygon) Attributes() []AttrEntry {
	attr := p.Attributes()

	for _, point := range p.Points {
		value := fmt.Sprintf("%d %d", point.PosX, point.PosY)
		attr = append(attr, AttrEntry{
			Name:  "point",
			Index: strconv.Itoa(point.PointIndex),
			Value: value,
		})
	}

	return attr
}

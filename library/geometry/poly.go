package geometry

import (
	"chroma-viz/library/attribute"
	"strings"
)

type Polygon struct {
	Geometry

	Polygon attribute.PolygonAttribute
	Color   attribute.ColorAttribute
}

func NewPolygon(geo Geometry) *Polygon {
	poly := &Polygon{
		Geometry: geo,
	}

	poly.Polygon.PosX = make(map[int]int, 128)
	poly.Polygon.PosY = make(map[int]int, 128)

	return poly
}

func (p *Polygon) UpdateGeometry(pEdit *PolygonEditor) (err error) {
	return
}

func (p *Polygon) Encode(b *strings.Builder) {
	p.Geometry.Encode(b)

	p.Polygon.Encode(b)
}

type PolygonEditor struct {
	GeometryEditor

	Poly attribute.PolygonAttribute
}

func NewPolygonEditor() (*PolygonEditor, error) {
	return nil, nil
}

func (pEdit *PolygonEditor) UpdateEditor(p *Polygon) (err error) {
	return
}

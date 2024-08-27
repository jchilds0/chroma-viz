package geometry

import (
	"chroma-viz/library/attribute"
	"strings"
)

type Polygon struct {
	Geometry
	Polygon attribute.PolygonAttribute
}

func NewPolygon(geo Geometry, numPoints int) Polygon {
	poly := Polygon{
		Geometry: geo,
	}

	poly.Polygon.PosX = make([]int, 0, numPoints)
	poly.Polygon.PosY = make([]int, 0, numPoints)
	return poly
}

func (p *Polygon) UpdateGeometry(pEdit *PolygonEditor) (err error) {
	return
}

func (p *Polygon) EncodeEngine(b strings.Builder) {

}

func (p *Polygon) EncodeJSON(b strings.Builder) {

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

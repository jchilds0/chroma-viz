package geometry

import (
	"chroma-viz/library/attribute"
	"strings"
)

type Polygon struct {
	Geometry

	Poly attribute.PolyAttribute
}

func (p *Polygon) UpdateGeometry(pEdit *PolygonEditor) (err error) {
	return
}

func (p *Polygon) EncodeEngine(b strings.Builder) {

}

type PolygonEditor struct {
	GeometryEditor

	Poly attribute.PolyAttribute
}

func NewPolygonEditor() *PolygonEditor {
	return nil
}

func (pEdit *PolygonEditor) UpdateEditor(p *Polygon) (err error) {
	return
}

package geometry

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/util"
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

	poly.Polygon.NumPoints = ATTR_NUM_POINTS
	poly.Polygon.Points = ATTR_POINT
	poly.Color.Name = ATTR_COLOR
	poly.Polygon.PosX = make(map[int]int, 128)
	poly.Polygon.PosY = make(map[int]int, 128)

	return poly
}

func (p *Polygon) UpdateGeometry(pEdit *PolygonEditor) (err error) {
	err = p.Geometry.UpdateGeometry(&pEdit.GeometryEditor)
	if err != nil {
		return
	}

	err = p.Color.UpdateAttribute(pEdit.Color)
	if err != nil {
		return
	}

	err = p.Polygon.UpdateAttribute(pEdit.Poly)
	return
}

func (p *Polygon) Encode(b *strings.Builder) {
	p.Geometry.Encode(b)

	util.EngineAddKeyValue(b, ATTR_COLOR_R, p.Color.Red)
	util.EngineAddKeyValue(b, ATTR_COLOR_G, p.Color.Green)
	util.EngineAddKeyValue(b, ATTR_COLOR_B, p.Color.Blue)
	util.EngineAddKeyValue(b, ATTR_COLOR_A, p.Color.Alpha)
	p.Polygon.Encode(b)
}

func (p *Polygon) GetGeometry() *Geometry {
	return &p.Geometry
}

type PolygonEditor struct {
	GeometryEditor

	Poly  *attribute.PolygonEditor
	Color *attribute.ColorEditor
}

func NewPolygonEditor() (pEdit *PolygonEditor, err error) {
	geoEdit, err := NewGeometryEditor()
	if err != nil {
		return
	}

	pEdit = &PolygonEditor{
		GeometryEditor: *geoEdit,
	}

	pEdit.Poly, err = attribute.NewPolygonEditor("Polygon")
	if err != nil {
		return
	}

	pEdit.Color, err = attribute.NewColorEditor(Attrs[ATTR_COLOR])
	if err != nil {
		return
	}

	pEdit.ScrollBox.PackStart(pEdit.Color.Box, false, false, padding)
	pEdit.ScrollBox.PackStart(pEdit.Poly.Box, true, true, padding)
	return
}

func (pEdit *PolygonEditor) UpdateEditor(p *Polygon) (err error) {
	err = pEdit.GeometryEditor.UpdateEditor(&p.Geometry)
	if err != nil {
		return
	}

	err = pEdit.Color.UpdateEditor(&p.Color)
	if err != nil {
		return
	}

	err = pEdit.Poly.UpdateEditor(&p.Polygon)
	return
}

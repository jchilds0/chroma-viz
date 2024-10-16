package geometry

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/util"
	"strings"
)

type List struct {
	Geometry

	String attribute.ListAttribute
	Scale  attribute.FloatAttribute
	Color  attribute.ColorAttribute
}

func NewList(geo Geometry) *List {
	ticker := &List{
		Geometry: geo,
	}

	ticker.String.Name = ATTR_STRING
	ticker.Scale.Name = ATTR_SCALE
	ticker.Color.Name = ATTR_COLOR
	ticker.Color.Alpha = 1.0

	return ticker
}

func (l *List) AddRow(index int, row attribute.ListRow) {
	if l.String.Rows == nil {
		l.String.Rows = make(map[int]attribute.ListRow, 128)
	}

	l.String.Rows[index] = row
}

func (l *List) UpdateGeometry(lEdit *ListEditor) (err error) {
	err = l.Geometry.UpdateGeometry(&lEdit.GeometryEditor)
	if err != nil {
		return
	}

	err = l.String.UpdateAttribute(lEdit.String)
	if err != nil {
		return
	}

	err = l.Color.UpdateAttribute(lEdit.Color)
	if err != nil {
		return
	}

	err = l.Scale.UpdateAttribute(lEdit.Scale)
	return
}

func (l *List) Encode(b *strings.Builder) {
	l.Geometry.Encode(b)

	l.String.Encode(b)
	util.EngineAddKeyValue(b, l.Scale.Name, l.Scale.Value)
	util.EngineAddKeyValue(b, ATTR_COLOR_R, l.Color.Red)
	util.EngineAddKeyValue(b, ATTR_COLOR_G, l.Color.Green)
	util.EngineAddKeyValue(b, ATTR_COLOR_B, l.Color.Blue)
	util.EngineAddKeyValue(b, ATTR_COLOR_A, l.Color.Alpha)
}

type ListEditor struct {
	GeometryEditor

	String *attribute.ListEditor
	Scale  *attribute.FloatEditor
	Color  *attribute.ColorEditor
}

func NewListEditor() (lEdit *ListEditor, err error) {
	geo, err := NewGeometryEditor()
	if err != nil {
		return
	}

	lEdit = &ListEditor{
		GeometryEditor: *geo,
	}

	lEdit.String, err = attribute.NewListEditor("List", 1)
	if err != nil {
		return
	}

	lEdit.Color, err = attribute.NewColorEditor(Attrs[ATTR_COLOR])
	if err != nil {
		return
	}

	lEdit.Scale, err = attribute.NewFloatEditor(Attrs[ATTR_SCALE], 0.0, 1.0, 0.01)
	if err != nil {
		return
	}

	lEdit.ScrollBox.PackStart(lEdit.Scale.Box, false, false, padding)
	lEdit.ScrollBox.PackStart(lEdit.Color.Box, false, false, padding)
	lEdit.ScrollBox.PackStart(lEdit.String.Box, true, true, padding)

	return
}

func (tEdit *ListEditor) UpdateEditor(t *List) (err error) {
	err = tEdit.GeometryEditor.UpdateEditor(&t.Geometry)
	if err != nil {
		return
	}

	err = tEdit.String.UpdateEditor(&t.String)
	if err != nil {
		return
	}

	err = tEdit.Scale.UpdateEditor(&t.Scale)
	if err != nil {
		return
	}

	err = tEdit.Color.UpdateEditor(&t.Color)
	return
}

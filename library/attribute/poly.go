package attribute

import (
	"chroma-viz/library/parser"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type PolygonAttribute struct {
	Name string
	PosX []int
	PosY []int
}

func (polyAttr *PolygonAttribute) Encode(b *strings.Builder) {
	parser.EngineAddKeyValue(b, "num_points", len(polyAttr.PosX))
}

func (polyAttr *PolygonAttribute) UpdateAttribute(polyEditor *PolygonEditor) error {
	return nil
}

type PolygonEditor struct {
	Name   string
	Box    *gtk.Box
	points *gtk.ListStore
}

func NewPolygonEditor(name string) (polyEdit *PolygonEditor, err error) {
	return
}

func (polyEdit *PolygonEditor) UpdateEditor(polyAttr *PolygonAttribute) error {
	return nil
}

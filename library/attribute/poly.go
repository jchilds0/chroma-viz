package attribute

import (
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type PolygonAttribute struct {
	Name string
	PosX []int
	PosY []int
}

func (polyAttr *PolygonAttribute) Encode(b strings.Builder) {
	b.WriteString("num_points")
	b.WriteRune('=')
	b.WriteString(strconv.Itoa(len(polyAttr.PosX)))
	b.WriteRune('#')
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

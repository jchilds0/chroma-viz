package attribute

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type PolygonAttribute struct {
	Name string
	PosX []int
	PosY []int
}

func (polyAttr *PolygonAttribute) Encode() string {
	return fmt.Sprintf("num_points=%d", len(polyAttr.PosX))
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

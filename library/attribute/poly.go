package attribute

import (
	"chroma-viz/library/parser"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type PolygonAttribute struct {
	Name string
	PosX map[int]int
	PosY map[int]int
}

func (polyAttr *PolygonAttribute) AddPoint(index, posX, posY int) {
	if polyAttr.PosX == nil {
		polyAttr.PosX = make(map[int]int, 128)
	}

	if polyAttr.PosY == nil {
		polyAttr.PosY = make(map[int]int, 128)
	}

	polyAttr.PosX[index] = posX
	polyAttr.PosY[index] = posY
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
	polyEdit = &PolygonEditor{
		Name: name,
	}

	polyEdit.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	polyEdit.Box.SetVisible(true)
	polyEdit.points, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		return
	}

	treeView, err := gtk.TreeViewNew()
	if err != nil {
		return
	}

	treeView.SetModel(polyEdit.points)
	cell, _ := NewListCell(0)
	cell.editableCell(polyEdit.points)
	column, _ := gtk.TreeViewColumnNewWithAttribute("Rows", cell, "text", 0)
	treeView.AppendColumn(column)

	return
}

func (polyEdit *PolygonEditor) UpdateEditor(polyAttr *PolygonAttribute) error {
	return nil
}

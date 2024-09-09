package attribute

import (
	"chroma-viz/library/util"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type PolygonAttribute struct {
	NumPoints string
	Points    string
	PosX      map[int]int
	PosY      map[int]int
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
	util.EngineAddKeyValue(b, polyAttr.NumPoints, len(polyAttr.PosX))

	for i := range polyAttr.PosX {
		point := fmt.Sprintf("%d %d %d", i-1, polyAttr.PosX[i], polyAttr.PosY[i])
		util.EngineAddKeyValue(b, polyAttr.Points, point)
	}
}

func (polyAttr *PolygonAttribute) UpdateAttribute(polyEditor *PolygonEditor) (err error) {
	iter, ok := polyEditor.points.GetIterFirst()
	model := polyEditor.points.ToTreeModel()
	polyAttr.PosX = make(map[int]int, polyEditor.numPoints)
	polyAttr.PosY = make(map[int]int, polyEditor.numPoints)

	var i, posX, posY int
	var s string
	for ok {
		i, err = util.ModelGetValue[int](model, iter, 0)
		if err != nil {
			return
		}

		s, err = util.ModelGetValue[string](model, iter, 1)
		if err != nil {
			return
		}

		posX, err = strconv.Atoi(s)
		if err != nil {
			return
		}

		polyAttr.PosX[i] = posX

		s, err = util.ModelGetValue[string](model, iter, 2)
		if err != nil {
			return
		}

		posY, err = strconv.Atoi(s)
		if err != nil {
			return
		}

		polyAttr.PosY[i] = posY

		ok = polyEditor.points.IterNext(iter)
	}

	return
}

type PolygonEditor struct {
	Name      string
	Box       *gtk.Box
	points    *gtk.ListStore
	numPoints int
}

func NewPolygonEditor(name string) (polyEdit *PolygonEditor, err error) {
	polyEdit = &PolygonEditor{
		Name:      name,
		numPoints: 0,
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
	treeView.SetVisible(true)

	{
		// Columns
		cell, _ := gtk.CellRendererTextNew()
		column, _ := gtk.TreeViewColumnNewWithAttribute("Num", cell, "text", 0)
		treeView.AppendColumn(column)

		cell1, _ := NewListCell(1)
		cell1.editableCell(polyEdit.points)
		column, _ = gtk.TreeViewColumnNewWithAttribute("Pos X", cell1, "text", 1)
		treeView.AppendColumn(column)

		cell2, _ := NewListCell(2)
		cell2.editableCell(polyEdit.points)
		column, _ = gtk.TreeViewColumnNewWithAttribute("Pos Y", cell2, "text", 2)
		treeView.AppendColumn(column)
	}

	frame, err := gtk.FrameNew(name)
	if err != nil {
		return
	}

	frame.Set("border-width", 2*padding)
	frame.Add(treeView)
	frame.SetVisible(true)

	actionBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	actionBox.SetVisible(true)

	label, err := gtk.LabelNew("Data Rows")
	if err != nil {
		return
	}

	label.SetVisible(true)
	label.SetWidthChars(12)
	actionBox.PackStart(label, false, false, padding)

	// add rows
	button, err := gtk.ButtonNewWithLabel("+")
	if err != nil {
		return
	}

	button.Connect("clicked", func() {
		iter := polyEdit.points.Append()
		polyEdit.points.SetValue(iter, 0, polyEdit.numPoints+1)
		polyEdit.points.SetValue(iter, 1, 0)
		polyEdit.points.SetValue(iter, 2, 0)

		polyEdit.numPoints++
	})

	button.SetVisible(true)
	actionBox.PackStart(button, false, false, padding)

	// remove rows
	button, err = gtk.ButtonNewWithLabel("-")
	if err != nil {
		return
	}

	button.Connect("clicked", func() {
		selection, err := treeView.GetSelection()
		if err != nil {
			log.Printf("Error getting current row (%s)", err)
			return
		}

		_, iter, ok := selection.GetSelected()
		if !ok {
			log.Printf("Error getting selected")
			return
		}

		polyEdit.points.Remove(iter)
	})

	button.SetVisible(true)
	actionBox.PackStart(button, false, false, padding)

	polyEdit.Box.PackStart(actionBox, false, false, 0)
	polyEdit.Box.PackStart(frame, true, true, 0)

	return
}

func (polyEdit *PolygonEditor) UpdateEditor(polyAttr *PolygonAttribute) error {
	polyEdit.points.Clear()
	polyEdit.numPoints = len(polyAttr.PosX)

	for i := range polyAttr.NumPoints {
		if _, ok := polyAttr.PosX[i]; !ok {
			continue
		}

		iter := polyEdit.points.Append()
		polyEdit.points.SetValue(iter, 0, i)
		polyEdit.points.SetValue(iter, 1, polyAttr.PosX[i])
		polyEdit.points.SetValue(iter, 2, polyAttr.PosY[i])
	}

	return nil
}

package attribute

import (
	"chroma-viz/library/util"
	"fmt"
	"log"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type ListRow []string

type ListAttribute struct {
	Name        string
	NumColumns  int
	Selected    bool
	SelectedRow int
	Rows        []ListRow
}

func NewListAttribute(name string, numCols int, selected bool) (list *ListAttribute, err error) {
	list = &ListAttribute{
		Name:       name,
		NumColumns: numCols,
		Selected:   selected,
	}

	cols := make([]glib.Type, list.NumColumns)
	for i := range cols {
		cols[i] = glib.TYPE_STRING
	}

	list.Rows = make([]ListRow, 0, 10)
	return
}

func (listAttr *ListAttribute) Encode(b *strings.Builder) {
	for i := range listAttr.Rows {
		listAttr.stringRow(b, i)
	}
}

func (listAttr *ListAttribute) stringRow(b *strings.Builder, rowIndex int) {
	b.WriteString(listAttr.Name)
	b.WriteRune('=')
	defer b.WriteRune('#')

	row := listAttr.Rows[rowIndex]

	if len(row) == 0 {
		return
	}

	b.WriteString(row[0])

	if len(row) == 1 {
		return
	}

	for _, elem := range row[1:len(row)] {
		b.WriteRune(' ')
		b.WriteString(elem)
	}

	return
}

func (listAttr *ListAttribute) UpdateAttribute(listEdit *ListEditor) (err error) {
	if listAttr.NumColumns != listEdit.NumColumns {
		return fmt.Errorf(
			"Incorrect number of columns in list editor: "+
				"List Attribute %d, List Editor %d",
			listAttr.NumColumns, listEdit.NumColumns)
	}

	listAttr.Rows = listAttr.Rows[:0]
	model := listEdit.listStore.ToTreeModel()

	iter, ok := listEdit.listStore.GetIterFirst()
	for ok {
		row := make(ListRow, listAttr.NumColumns)
		listAttr.Rows = append(listAttr.Rows, row)

		for i := range listEdit.NumColumns {
			row[i], err = util.ModelGetValue[string](model, iter, i)
			if err != nil {
				return
			}
		}
	}

	return
}

type graphCell struct {
	*gtk.CellRendererText
	columnNum int
}

func NewGraphCell(i int) (gCell *graphCell, err error) {
	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		err = fmt.Errorf("Error creating graph cell (%s)", err)
	}

	gCell = &graphCell{CellRendererText: cell, columnNum: i}
	return
}

func (gCell *graphCell) editableCell(list *gtk.ListStore) {
	gCell.SetProperty("editable", true)

	gCell.Connect("edited", func(cell *gtk.CellRendererText, path string, text string) {
		iter, err := list.ToTreeModel().GetIterFromString(path)
		if err != nil {
			log.Printf("Error editing list attribute (%s)", err)
			return
		}

		list.SetValue(iter, gCell.columnNum, text)
	})
}

type ListEditor struct {
	Name       string
	NumColumns int
	Box        *gtk.Box
	treeView   *gtk.TreeView
	listStore  *gtk.ListStore
}

func NewListEditor(name string, numCols int) (listEdit *ListEditor, err error) {
	listEdit = &ListEditor{
		Name:       name,
		NumColumns: numCols,
	}

	listEdit.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	listEdit.treeView, err = gtk.TreeViewNew()
	if err != nil {
		return
	}

	colTypes := make([]glib.Type, numCols)
	for i := range colTypes {
		colTypes[i] = glib.TYPE_STRING
	}

	listEdit.listStore, err = gtk.ListStoreNew(colTypes...)
	if err != nil {
		return
	}

	listEdit.treeView.SetVisible(true)
	listEdit.treeView.SetVExpand(true)

	var gCell *graphCell
	var column *gtk.TreeViewColumn
	for i := range numCols {
		gCell, err = NewGraphCell(i)
		if err != nil {
			return
		}

		gCell.editableCell(listEdit.listStore)
		column, err = gtk.TreeViewColumnNewWithAttribute("", gCell, "text", i)
		if err != nil {
			return
		}

		listEdit.treeView.AppendColumn(column)
	}

	frame, err := gtk.FrameNew(name)
	if err != nil {
		return
	}

	frame.Set("border-width", 2*padding)
	frame.Add(listEdit.treeView)
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
		listEdit.listStore.Append()
	})

	button.SetVisible(true)
	actionBox.PackStart(button, false, false, padding)

	// remove rows
	button, err = gtk.ButtonNewWithLabel("-")
	if err != nil {
		return
	}

	button.Connect("clicked", func() {
		selection, err := listEdit.treeView.GetSelection()
		if err != nil {
			log.Printf("Error getting current row (%s)", err)
			return
		}

		_, iter, ok := selection.GetSelected()
		if !ok {
			log.Printf("Error getting selected")
			return
		}

		listEdit.listStore.Remove(iter)
	})

	button.SetVisible(true)
	actionBox.PackStart(button, false, false, padding)

	listEdit.Box.PackStart(actionBox, false, false, 0)
	listEdit.Box.PackStart(frame, true, true, 0)
	return
}

func (listEdit *ListEditor) UpdateEditor(listAttr *ListAttribute) (err error) {
	if listEdit.NumColumns != listAttr.NumColumns {
		return fmt.Errorf(
			"Incorrect number of columns in list attribute: "+
				"List Editor %d, List Attribute %d",
			listEdit.NumColumns, listAttr.NumColumns)
	}

	listEdit.listStore.Clear()
	if listAttr.Rows == nil {
		return
	}

	for _, row := range listAttr.Rows {
		if row == nil {
			continue
		}

		iter := listEdit.listStore.Append()

		for i, elem := range row {
			listEdit.listStore.SetValue(iter, i, elem)
		}
	}

	return
}

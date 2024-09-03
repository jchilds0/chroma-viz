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

func NewListRow(s string) ListRow {
	return strings.Split(s, ",")
}

func (row *ListRow) ToString() string {
	if row == nil {
		return ""
	}

	return strings.Join(*row, ",")
}

type ListAttribute struct {
	Name        string
	Selected    bool
	SelectedRow int
	Header      ListRow
	Rows        map[int]ListRow
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
	model := listEdit.listStore.ToTreeModel()
	listAttr.Rows = make(map[int]ListRow, 128)

	iter, ok := listEdit.listStore.GetIterFirst()
	for i := 1; ok; i++ {
		row := make(ListRow, listEdit.NumColumns)
		listAttr.Rows[i] = row

		for j := range listEdit.NumColumns {
			row[j], err = util.ModelGetValue[string](model, iter, j)
			if err != nil {
				return
			}
		}

		ok = listEdit.listStore.IterNext(iter)
	}

	return
}

type listCell struct {
	*gtk.CellRendererText
	columnNum int
}

func NewListCell(i int) (cell *listCell, err error) {
	c, err := gtk.CellRendererTextNew()
	if err != nil {
		err = fmt.Errorf("Error creating graph cell (%s)", err)
	}

	cell = &listCell{CellRendererText: c, columnNum: i}
	return
}

func (cell *listCell) editableCell(list *gtk.ListStore) {
	cell.SetProperty("editable", true)

	cell.Connect("edited", func(c *gtk.CellRendererText, path string, text string) {
		iter, err := list.ToTreeModel().GetIterFromString(path)
		if err != nil {
			log.Printf("Error editing list attribute (%s)", err)
			return
		}

		list.SetValue(iter, cell.columnNum, text)
	})
}

type ListEditor struct {
	Name       string
	NumColumns int
	Box        *gtk.Box
	treeView   *gtk.TreeView
	listStore  *gtk.ListStore
}

func NewListEditor(name string, numColumns int) (listEdit *ListEditor, err error) {
	listEdit = &ListEditor{
		Name:       name,
		NumColumns: numColumns,
	}

	listEdit.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	listEdit.Box.SetVisible(true)

	listEdit.treeView, err = gtk.TreeViewNew()
	if err != nil {
		return
	}

	listEdit.treeView.SetVisible(true)
	listEdit.treeView.SetVExpand(true)

	colTypes := make([]glib.Type, numColumns)
	for i := range numColumns {
		colTypes[i] = glib.TYPE_STRING
	}

	listEdit.listStore, err = gtk.ListStoreNew(colTypes...)
	if err != nil {
		return
	}

	listEdit.treeView.SetModel(listEdit.listStore)

	var cell *listCell
	var column *gtk.TreeViewColumn
	for i := range numColumns {
		cell, err = NewListCell(i)
		if err != nil {
			return
		}

		cell.editableCell(listEdit.listStore)
		column, err = gtk.TreeViewColumnNewWithAttribute("Rows", cell, "text", i)
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

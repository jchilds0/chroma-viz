package attribute

import (
	"fmt"
	"log"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

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

type ListRow []string

type ListAttribute struct {
	Name        string
	NumCols     int
	Selected    bool
	SelectedRow int
	Rows        []ListRow
}

func NewListAttribute(name string, numCols int, selected bool) (list *ListAttribute, err error) {
	list = &ListAttribute{
		Name:     name,
		NumCols:  numCols,
		Selected: selected,
	}

	cols := make([]glib.Type, list.NumCols)
	for i := range cols {
		cols[i] = glib.TYPE_STRING
	}

	list.Rows = make([]ListRow, 0, 10)
	return
}

func (listAttr *ListAttribute) EncodeEngine() (s string) {
	return
}

func (listAttr *ListAttribute) stringRow(row int) (s string, err error) {
	s = listAttr.Name + "=" + strings.Join(listAttr.Rows[row], " ") + "#"
	return
}

func (listAttr *ListAttribute) EncodeJSON() (s string) {
	// currently chroma_engine allocates 100 nodes for each list statically
	if listAttr.Selected {
		// send only the currently selected item from the list
		row, _ := listAttr.encodeRow(listAttr.SelectedRow)
		return strings.Join(row, " ")
	}

	s = fmt.Sprintf("{'name': 'num_node', 'value': '%d'}", listAttr.NumCols-1)
	for _, row := range listAttr.Rows {
		s += fmt.Sprintf(",{'name': '%s', 'value': '%s'}",
			listAttr.Name, strings.Join(row, " "))
	}

	return
}

func (listAttr *ListAttribute) encodeRow(rowIndex int) (row []string, err error) {
	return
}

func (listAttr *ListAttribute) UpdateAttribute(listEdit *ListEditor) (err error) {
	return
}

type ListEditor struct {
	Name      string
	Box       *gtk.Box
	treeView  *gtk.TreeView
	listStore *gtk.ListStore
}

func NewListEditor(name string, columns []string) (listEdit *ListEditor, err error) {
	listEdit = &ListEditor{Name: name}
	listEdit.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	listEdit.treeView, err = gtk.TreeViewNew()
	if err != nil {
		return
	}

	listEdit.treeView.SetVisible(true)
	listEdit.treeView.SetVExpand(true)

	var gCell *graphCell
	var column *gtk.TreeViewColumn
	for i, name := range columns {
		gCell, err = NewGraphCell(i)
		if err != nil {
			return
		}

		gCell.SetProperty("editable", true)
		gCell.Connect("edited",
			func(cell *gtk.CellRendererText, path string, text string) {
				if listEdit.listStore == nil {
					log.Printf("Error editing list attribute (listStore missing)")
					return
				}

				iter, err := listEdit.listStore.ToTreeModel().GetIterFromString(path)
				if err != nil {
					log.Printf("Error editing list attribute (%s)", err)
					return
				}

				listEdit.listStore.SetValue(iter, gCell.columnNum, text)
			})
		column, err = gtk.TreeViewColumnNewWithAttribute(name, gCell, "text", i)
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
		if listEdit.listStore == nil {
			log.Printf("Graph prop editor does not have a list store")
			return
		}

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
		if listEdit.listStore == nil {
			log.Printf("Graph prop editor does not have a list store")
			return
		}

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
	return nil
}

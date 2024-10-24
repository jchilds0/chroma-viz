package main

import (
	"chroma-viz/library"
	"chroma-viz/library/pages"
	"chroma-viz/library/util"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type ExternalShow struct {
	treeView *gtk.TreeView
	treeList *gtk.ListStore

	addr  string
	port  int
	pages map[int]pages.Page
}

func NewExternalShow(addr string, port int, pageToEditor func(*pages.Page) error) *ExternalShow {
	var err error
	show := &ExternalShow{}

	show.treeView, err = gtk.TreeViewNew()
	if err != nil {
		log.Fatalln("Error creating show:", err)
	}

	show.treeView.SetReorderable(true)

	show.treeList, err = gtk.ListStoreNew(
		glib.TYPE_INT,
		glib.TYPE_STRING,
		glib.TYPE_INT,
		glib.TYPE_STRING,
	)
	if err != nil {
		log.Fatalln("Error creating show:", err)
	}

	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatalln("Error creating show:", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(KEYTITLE[PAGENUM], cell, "text", PAGENUM)
	if err != nil {
		log.Fatalln("Error creating show:", err)
	}

	column.SetResizable(true)
	show.treeView.AppendColumn(column)

	title, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatalln("Error creating show:", err)
	}

	title.SetProperty("editable", true)
	title.Connect("edited",
		func(cell *gtk.CellRendererText, path string, text string) {
			iter, err := show.treeList.GetIterFromString(path)
			if err != nil {
				log.Println("Error editing page:", err)
				return
			}

			model := show.treeList.ToTreeModel()
			pageNum, err := util.ModelGetValue[int](model, iter, PAGENUM)
			if err != nil {
				log.Println("Error editing page:", err)
				return
			}

			show.UpdatePageTitle(pageNum, text)
			show.treeList.SetValue(iter, TITLE, text)
		})

	column, err = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[TITLE], title, "text", TITLE)
	if err != nil {
		log.Fatalln("Error creating show:", err)
	}

	column.SetExpand(true)
	column.SetResizable(true)
	show.treeView.AppendColumn(column)

	column, err = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[TEMPLATE_ID], cell, "text", TEMPLATE_ID)
	if err != nil {
		log.Fatalln("Error creating show:", err)
	}

	column.SetResizable(true)
	show.treeView.AppendColumn(column)

	column, err = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[TEMPLATE_NAME], cell, "text", TEMPLATE_NAME)
	if err != nil {
		log.Fatalln("Error creating show:", err)
	}

	column.SetResizable(true)
	show.treeView.AppendColumn(column)

	// send page to editor on double click
	show.treeView.Connect("row-activated",
		func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
			iter, err := show.treeList.GetIter(path)
			if err != nil {
				log.Println("Error sending page to editor:", err)
				return
			}

			model := show.treeList.ToTreeModel()
			pageNum, err := util.ModelGetValue[int](model, iter, PAGENUM)
			if err != nil {
				log.Println("Error sending page to editor:", err)
				return
			}

			page, ok := show.GetPage(pageNum)
			if !ok {
				log.Println("Missing page", pageNum)
				return
			}

			pageToEditor(page)
			SendPreview(page, library.ANIMATE_ON)
		})

	return show
}

func (show *ExternalShow) UpdatePageTitle(pageNum int, title string) {

}

func (show *ExternalShow) TreeView() *gtk.TreeView {

}

func (show *ExternalShow) SelectedPage() (int, error) {

}

func (show *ExternalShow) AddPage(page pages.Page) (err error) {

}

func (show *ExternalShow) GetPage(pageNum int) (*pages.Page, bool) {

}

func (show *ExternalShow) GetPages() map[int]pages.Page {
	return show.pages
}

func (show *ExternalShow) DeletePage(pageNum int) {

}

func (show *ExternalShow) Clear() {

}

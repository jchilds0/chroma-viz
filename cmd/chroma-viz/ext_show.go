package main

import (
	"chroma-viz/library"
	"chroma-viz/library/pages"
	"chroma-viz/library/util"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type ExternalShow struct {
	treeView *gtk.TreeView
	treeList *gtk.ListStore

	server net.Conn // listen for page updates
	conn   net.Conn // send and recieve pages
	pages  map[int]pages.Page
}

func NewExternalShow(addr string, port int, pageToEditor func(*pages.Page) error) *ExternalShow {
	var err error
	show := &ExternalShow{}

	show.server, err = net.Dial("tcp", addr+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}

	m := Message{
		Type: RECIEVE_UPDATES,
	}

	err = sendMessage(show.server, m)
	if err != nil {
		log.Fatal(err)
	}

	show.conn, err = net.Dial("tcp", addr+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}

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

	show.treeView.SetModel(show.treeList)

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

	pages := show.GetPages()
	for _, page := range pages {
		show.addRow(page)
	}

	return show
}

func (show *ExternalShow) addRow(page PageData) {
	iter := show.treeList.Append()
	err := show.treeList.SetValue(iter, PAGENUM, page.PageNum)
	if err != nil {
		return
	}

	err = show.treeList.SetValue(iter, TITLE, page.Title)
	if err != nil {
		return
	}

	err = show.treeList.SetValue(iter, TEMPLATE_ID, page.TempID)
	if err != nil {
		return
	}

	err = show.treeList.SetValue(iter, TEMPLATE_NAME, page.Title)
	return
}

func (show *ExternalShow) UpdatePageTitle(pageNum int, title string) {
	m := Message{
		Type: UPDATE_PAGE,
		PageInfo: PageData{
			PageNum: pageNum,
			Title:   title,
		},
	}

	err := sendMessage(show.conn, m)
	if err != nil {
		log.Println("Update Page Title", err)
	}
}

func (show *ExternalShow) TreeView() *gtk.TreeView {
	return show.treeView
}

func (show *ExternalShow) SelectedPage() (pageNum int, err error) {
	selection, err := show.treeView.GetSelection()
	if err != nil {
		return
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		err = fmt.Errorf("Error getting selection iter")
		return
	}

	model := show.treeList.ToTreeModel()
	pageNum, err = util.ModelGetValue[int](model, iter, PAGENUM)
	return
}

func (show *ExternalShow) AddPage(page pages.Page) (err error) {
	req := Message{
		Type: CREATE_PAGE,
		Page: page,
	}

	err = sendMessage(show.conn, req)
	return
}

func (show *ExternalShow) GetPage(pageNum int) (*pages.Page, bool) {
	req := Message{
		Type: READ_PAGE,
		PageInfo: PageData{
			PageNum: pageNum,
		},
	}

	err := sendMessage(show.conn, req)
	if err != nil {
		log.Println("Error getting page", pageNum, err)
		return nil, false
	}

	res, err := recvMessage(show.conn)
	if err != nil {
		log.Println("Error getting page", pageNum, err)
		return nil, false
	}

	return &res.Page, true
}

func (show *ExternalShow) GetPages() (pages map[int]PageData) {
	pages = make(map[int]PageData)

	req := Message{
		Type: GET_PAGES,
	}

	err := sendMessage(show.conn, req)
	if err != nil {
		log.Println("Error getting pages:", err)
		return
	}

	res, err := recvMessage(show.conn)
	if err != nil {
		log.Println("Error getting pages:", err)
		return
	}

	return res.PageData
}

func (show *ExternalShow) DeletePage(pageNum int) {
	m := Message{
		Type: DELETE_PAGE,
		PageInfo: PageData{
			PageNum: pageNum,
		},
	}

	err := sendMessage(show.conn, m)
	if err != nil {
		log.Printf("Error deleting page %d: %s", pageNum, err)
		return
	}
}

func (show *ExternalShow) Clear() {
}

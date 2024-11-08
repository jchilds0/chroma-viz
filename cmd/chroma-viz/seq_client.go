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

type SequencerClient struct {
	rows     map[int]*gtk.TreeIter
	treeView *gtk.TreeView
	treeList *gtk.ListStore

	server net.Conn // listen for page updates
	conn   net.Conn // send and recieve pages
}

func NewSequencerClient(addr string, port int, pageToEditor func(*pages.Page) error) (*SequencerClient, error) {
	var err error
	show := &SequencerClient{
		rows: make(map[int]*gtk.TreeIter),
	}

	show.server, err = net.Dial("tcp", addr+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	m := Message{
		Type: RECIEVE_UPDATES,
	}

	err = sendMessage(show.server, m)
	if err != nil {
		return nil, err
	}

	show.conn, err = net.Dial("tcp", addr+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	show.treeList, err = gtk.ListStoreNew(
		glib.TYPE_INT,
		glib.TYPE_STRING,
		glib.TYPE_INT,
		glib.TYPE_STRING,
		glib.TYPE_STRING,
	)
	if err != nil {
		return nil, err
	}

	show.treeView, err = createShowTreeModel(func(text string, pageNum int) {
		row, ok := show.rows[pageNum]
		if !ok {
			log.Printf("Missing iter for page %d", pageNum)
			return
		}

		show.UpdatePageTitle(pageNum, text)
		show.treeList.SetValue(row, TITLE, text)
	})
	if err != nil {
		return nil, err
	}

	show.treeView.SetModel(show.treeList)

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

			page, ok := show.ReadPage(pageNum)
			if !ok {
				log.Println("Missing page", pageNum)
				return
			}

			SendPreview(page, library.ANIMATE_ON)
			pageToEditor(page)
		})

	pages := show.GetPages()
	for _, page := range pages {
		show.addRow(page)
	}

	go show.pageUpdates()

	return show, nil
}

func (show *SequencerClient) pageUpdates() {
	for {
		m, err := recvMessage(show.server)
		if err != nil {
			log.Println("Error receiving page update", err)
			break
		}

		show.addRow(m.PageInfo)
	}
}

func (show *SequencerClient) addRow(page PageData) {
	if _, ok := show.rows[page.PageNum]; !ok {
		show.rows[page.PageNum] = show.treeList.Append()
	}

	iter := show.rows[page.PageNum]
	err := show.treeList.SetValue(iter, PAGENUM, page.PageNum)
	if err != nil {
		return
	}

	err = show.treeList.SetValue(iter, TITLE, page.Title)
	if err != nil {
		return
	}

	if page.TempID != 0 {
		err = show.treeList.SetValue(iter, TEMPLATE_ID, page.TempID)
		if err != nil {
			return
		}
	}

	if page.TempName != "" {
		err = show.treeList.SetValue(iter, TEMPLATE_NAME, "")
		if err != nil {
			return
		}
	}

	return
}

func (show *SequencerClient) UpdatePageTitle(pageNum int, title string) {
	m := Message{
		Type: UPDATE_PAGE_INFO,
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

func (show *SequencerClient) TreeView() *gtk.TreeView {
	return show.treeView
}

func (show *SequencerClient) SelectedPage() (pageNum int, err error) {
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

func (show *SequencerClient) WritePage(page *pages.Page) (err error) {
	req := Message{
		Type: WRITE_PAGE,
		Page: page,
	}

	err = sendMessage(show.conn, req)
	return
}

func (show *SequencerClient) ReadPage(pageNum int) (*pages.Page, bool) {
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

	return res.Page, true
}

func (show *SequencerClient) GetPages() (pages map[int]PageData) {
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

func (show *SequencerClient) DeletePage(pageNum int) {
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

func (show *SequencerClient) Clear() {
}

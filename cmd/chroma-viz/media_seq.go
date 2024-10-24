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

type MediaSequencer struct {
	treeView *gtk.TreeView
	treeList *gtk.ListStore

	numPages int
	conn     net.Listener
	pages    map[int]pages.Page
	clients  map[string]net.Conn
}

func NewMediaSequencer(port int, pageToEditor func(*pages.Page) error) *MediaSequencer {
	var err error
	show := &MediaSequencer{
		pages:   make(map[int]pages.Page, 1024),
		clients: make(map[string]net.Conn, 64),
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
				log.Println("Error editing page:", err)
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

	show.conn, err = net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}

	go show.listen()

	return show
}

func (show *MediaSequencer) listen() {
	log.Println("Media Sequencer listening on:", show.conn.Addr())

	for {
		client, err := show.conn.Accept()
		if err != nil {
			log.Print(err)
		}

		log.Println("Connected to client", client.RemoteAddr())
		go show.handleConn(client)
	}

}

func (show *MediaSequencer) handleConn(client net.Conn) {
	defer client.Close()

	for {
		req, err := recvMessage(client)
		if err != nil {
			log.Println(client.RemoteAddr(), err)
			break
		}

		if req.Type == NO_MESSAGE {
			continue
		}

		log.Println(client.RemoteAddr(), "Message:", req)

		switch req.Type {
		case CREATE_PAGE:
		case READ_PAGE:
			res := Message{
				Type: READ_PAGE,
			}

			page, ok := show.GetPage(req.PageInfo.PageNum)
			if ok {
				res.Page = *page

			}

			err = sendMessage(client, res)
			if err != nil {
				log.Println(err)
				continue
			}

		case UPDATE_PAGE:
		case DELETE_PAGE:
		case GET_PAGES:
			res := Message{
				Type:     GET_PAGES,
				PageData: show.GetPages(),
			}

			err = sendMessage(client, res)
			if err != nil {
				log.Println(err)
				continue
			}

		case RECIEVE_UPDATES:
			show.clients[client.RemoteAddr().String()] = client
		}
	}

	log.Print("Closing connection to", client.RemoteAddr().String())
}

func (show *MediaSequencer) TreeView() *gtk.TreeView {
	return show.treeView
}

func (show *MediaSequencer) SelectedPage() (pageNum int, err error) {
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

func (show *MediaSequencer) AddPage(page pages.Page) (err error) {
	show.pages[page.PageNum] = page

	for _, geo := range page.Clock {
		if geo == nil {
			continue
		}

		/*
		   Clock requires a way to send updates to viz
		   to animate the clock. We manually add this
		   after parsing the page.
		*/
		geo.Clock.SetClock(func() { SendEngine(&page, library.CONTINUE) })
	}

	iter := show.treeList.Append()
	err = show.treeList.SetValue(iter, PAGENUM, page.PageNum)
	if err != nil {
		return
	}

	err = show.treeList.SetValue(iter, TITLE, page.Title)
	if err != nil {
		return
	}

	err = show.treeList.SetValue(iter, TEMPLATE_ID, page.TemplateID)
	if err != nil {
		return
	}

	err = show.treeList.SetValue(iter, TEMPLATE_NAME, page.Title)
	return
}

func (show *MediaSequencer) GetPages() map[int]PageData {
	pageData := make(map[int]PageData, len(show.pages))

	for _, page := range show.pages {
		pageData[page.PageNum] = PageData{
			PageNum: page.PageNum,
			TempID:  page.TemplateID,
			Title:   page.Title,
			Layer:   page.Layer,
		}
	}

	return pageData
}

func (show *MediaSequencer) GetPage(pageNum int) (*pages.Page, bool) {
	page, ok := show.pages[pageNum]
	return &page, ok
}

func (show *MediaSequencer) UpdatePageTitle(pageNum int, title string) {
	page, ok := show.pages[pageNum]
	if !ok {
		return
	}

	page.Title = title
	show.pages[pageNum] = page
}

func (show *MediaSequencer) DeletePage(pageNum int) {
	log.Fatal("Not implemented")
}

func (show *MediaSequencer) Clear() {
	show.treeList.Clear()
	clear(show.pages)
}

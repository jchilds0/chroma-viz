package viz

import (
	"chroma-viz/library/editor"
	"chroma-viz/library/gtk_utils"
	"chroma-viz/library/shows"
	"chroma-viz/library/tcp"
	"chroma-viz/library/templates"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var conn *GuiConn = NewGuiConn()

type GuiConn struct {
	hub  *tcp.Connection
	eng  []*tcp.Connection
	prev []*tcp.Connection
}

func NewGuiConn() *GuiConn {
	gui := &GuiConn{}
	gui.eng = make([]*tcp.Connection, 0, 10)
	gui.prev = make([]*tcp.Connection, 0, 10)

	return gui
}

func SendPreview(page tcp.Animator, action int) {
	if page == nil {
		log.Println("SendPreview recieved nil page")
		return
	}

	for _, c := range conn.prev {
		if c == nil {
			continue
		}

		c.SetPage <- page
		c.SetAction <- action
	}
}

func SendEngine(page tcp.Animator, action int) {
	if page == nil {
		log.Println("SendEngine recieved nil page")
		return
	}

	for _, c := range conn.eng {
		if c == nil {
			continue
		}

		c.SetPage <- page
		c.SetAction <- action
	}
}

/* 
    A hook which is run after the viz TempTree and 
    ShowTree are initialised. This allows a test to 
    to call the import methods of these structs
*/
var importHook = func(hub net.Conn, temp *TempTree, show *ShowTree) {}

func VizGui(app *gtk.Application) {
	win, err := gtk.ApplicationWindowNew(app)
	if err != nil {
		log.Fatal(err)
	}

	win.SetDefaultSize(800, 600)
	win.SetTitle("Chroma Viz")

	preview := setupPreviewWindow(conf.PreviewDirectory, conf.PreviewName)

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal(err)
	}

	win.Add(box)

    edit := editor.NewEditor(SendEngine, SendPreview)
    showTree := NewShowTree(func(page *shows.Page) { edit.SetPage(page) })
    tempTree := NewTempTree(conn.hub.Conn, 
        func(temp *templates.Template) {
            page := showTree.show.AddPage(temp.Title, temp)
            showTree.ImportPage(page)
        },
    )

	edit.AddAction("Take On", true, func() { SendEngine(edit.Page, tcp.ANIMATE_ON) })
	edit.AddAction("Continue", true, func() { SendEngine(edit.Page, tcp.CONTINUE) })
	edit.AddAction("Take Off", true, func() { SendEngine(edit.Page, tcp.ANIMATE_OFF) })
	edit.AddAction("Save", false, func() {
		edit.UpdateProps()
		SendPreview(edit.Page, tcp.ANIMATE_ON)
	})
	edit.PageEditor()

    start := time.Now()
	tempTree.ImportTemplates(conn.hub.Conn)
    end := time.Now()
    elapsed := end.Sub(start)

    log.Printf("Imported Graphics Hub in %s", elapsed)
	importHook(conn.hub.Conn, tempTree, showTree)

	/* Menu layout */
	builder, err := gtk.BuilderNew()
    if err := builder.AddFromFile(conf.InstallDirectory + "gtk/viz-menu.ui"); err != nil {
		log.Fatal(err)
	}

	menu, err := builder.GetObject("menubar")
	if err != nil {
		log.Fatal(err)
	}

	app.SetMenubar(menu.(*glib.MenuModel))

	importShow := glib.SimpleActionNew("import_show", nil)
	importShow.Connect("activate", func() {
		err := guiImportShow(win, showTree)
		if err != nil {
			log.Printf("Error importing show (%s)", err)
		}
	})
	app.AddAction(importShow)

	exportShow := glib.SimpleActionNew("export_show", nil)
	exportShow.Connect("activate", func() {
		err := guiExportShow(win, showTree)
		if err != nil {
			log.Printf("Error exporting show (%s)", err)
		}
	})
	app.AddAction(exportShow)

	importPage := glib.SimpleActionNew("import_page", nil)
	importPage.Connect("activate", func() {
		err := guiImportPage(win, showTree)
		if err != nil {
			log.Printf("Error importing page (%s)", err)
		}
	})
	app.AddAction(importPage)

	exportPage := glib.SimpleActionNew("export_page", nil)
	exportPage.Connect("activate", func() {
		err := guiExportPage(win, showTree)
		if err != nil {
			log.Printf("Error exporting page (%s)", err)
		}
	})
	app.AddAction(exportPage)

	deletePage := glib.SimpleActionNew("delete_page", nil)
	deletePage.Connect("activate", func() {
		err := guiDeletePage(showTree)
		if err != nil {
			log.Printf("Error deleting page (%s)", err)
		}
	})
	app.AddAction(deletePage)

	/* Body layout */
	builder, err = gtk.BuilderNew()
	if err := builder.AddFromFile(conf.InstallDirectory + "gtk/viz-gui.ui"); err != nil {
		log.Fatal(err)
	}

	body, err := gtk_utils.BuilderGetObject[*gtk.Paned](builder, "body")
	if err != nil {
		log.Fatal(err)
	}

	box.PackStart(body, true, true, 0)

	tempScroll, err := gtk_utils.BuilderGetObject[*gtk.ScrolledWindow](builder, "templates-win")
	if err != nil {
		log.Fatal(err)
	}

	tempScroll.Add(tempTree.treeView)

	showScroll, err := gtk_utils.BuilderGetObject[*gtk.ScrolledWindow](builder, "show-win")
	if err != nil {
		log.Fatal(err)
	}

	showScroll.Add(showTree.treeView)

	editBox, err := gtk_utils.BuilderGetObject[*gtk.Box](builder, "edit")
	if err != nil {
		log.Fatal(err)
	}

	editBox.PackStart(edit.Box, true, true, 0)

	prevBox, err := gtk_utils.BuilderGetObject[*gtk.Box](builder, "preview")
	if err != nil {
		log.Fatal(err)
	}

	prevBox.PackStart(preview, true, true, 0)

	/* Lower Bar layout */
	lowerBox, err := gtk.ActionBarNew()
	if err != nil {
		log.Fatal(err)
	}

	box.PackEnd(lowerBox, false, false, 0)

	button, err := gtk.ButtonNew()
	if err != nil {
		log.Fatal(err)
	}

	lowerBox.PackEnd(button)
	button.SetLabel("Exit")
	button.Connect("clicked", func() { gtk.MainQuit() })

	for _, c := range conn.eng {
		if c == nil {
			continue
		}

		eng := NewEngineWidget(c)
		lowerBox.PackStart(eng.button)
	}

	for _, c := range conn.prev {
		if c == nil {
			continue
		}

		eng := NewEngineWidget(c)
		lowerBox.PackStart(eng.button)
	}

	win.ShowAll()

}

func guiImportShow(win *gtk.ApplicationWindow, show *ShowTree) error {
	dialog, err := gtk.FileChooserDialogNewWith2Buttons(
		"Import Show", win, gtk.FILE_CHOOSER_ACTION_OPEN,
		"_Cancel", gtk.RESPONSE_CANCEL, "_Open", gtk.RESPONSE_ACCEPT)
	if err != nil {
		return err
	}
	defer dialog.Destroy()

	res := dialog.Run()
	if res == gtk.RESPONSE_ACCEPT {
		filename := dialog.GetFilename()
		show.ImportShow(filename)
	}

	return nil
}

func guiExportShow(win *gtk.ApplicationWindow, showTree *ShowTree) error {
	dialog, err := gtk.FileChooserDialogNewWith2Buttons(
		"Save Show", win, gtk.FILE_CHOOSER_ACTION_SAVE,
		"_Cancel", gtk.RESPONSE_CANCEL, "_Save", gtk.RESPONSE_ACCEPT)
	if err != nil {
		return err
	}
	defer dialog.Destroy()

	dialog.SetCurrentName(".show")
	res := dialog.Run()
	if res == gtk.RESPONSE_ACCEPT {
		filename := dialog.GetFilename()
		showTree.show.ExportShow(filename)
	}

	return nil
}

func guiImportPage(win *gtk.ApplicationWindow, showTree *ShowTree) error {
	dialog, err := gtk.FileChooserDialogNewWith2Buttons(
		"Import Page", win, gtk.FILE_CHOOSER_ACTION_OPEN,
		"_Cancel", gtk.RESPONSE_CANCEL, "_Open", gtk.RESPONSE_ACCEPT)
	if err != nil {
		return err
	}
	defer dialog.Destroy()

	res := dialog.Run()
	if res == gtk.RESPONSE_ACCEPT {
		filename := dialog.GetFilename()

		page := &shows.Page{}
		err := page.ImportPage(filename)
		if err != nil {
			return err
		}

		showTree.ImportPage(page)
	}

	return nil
}

func guiExportPage(win *gtk.ApplicationWindow, showTree *ShowTree) error {
	selection, err := showTree.treeView.GetSelection()
	if err != nil {
		return err
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		return fmt.Errorf("Error getting selected iter")
	}

	model := &showTree.treeList.TreeModel
	title, err := gtk_utils.ModelGetValue[string](model, iter, TITLE)
	if err != nil {
		return err
	}

	pageNum, err := gtk_utils.ModelGetValue[int](model, iter, PAGENUM)
	if err != nil {
		return err
	}

	dialog, err := gtk.FileChooserDialogNewWith2Buttons(
		"Save Page", win, gtk.FILE_CHOOSER_ACTION_SAVE,
		"_Cancel", gtk.RESPONSE_CANCEL, "_Save", gtk.RESPONSE_ACCEPT)
	if err != nil {
		return err
	}
	defer dialog.Destroy()

	dialog.SetCurrentName(title + ".json")
	res := dialog.Run()
	if res == gtk.RESPONSE_ACCEPT {
		filename := dialog.GetFilename()

		page := showTree.show.Pages[pageNum]
		if page == nil {
			return fmt.Errorf("Page %d does not exist", pageNum)
		}

		err := shows.ExportPage(page, filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func guiDeletePage(show *ShowTree) error {
	selection, err := show.treeView.GetSelection()
	if err != nil {
		return err
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		return fmt.Errorf("Error getting selection iter")
	}

	model := &show.treeList.TreeModel
	pageNum, err := gtk_utils.ModelGetValue[int](model, iter, PAGENUM)
	if err != nil {
		return err
	}

	show.treeList.Remove(iter)
	show.show.Pages[pageNum] = nil
	return nil
}


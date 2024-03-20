package viz

import (
	"chroma-viz/library/editor"
	"chroma-viz/library/shows"
	"chroma-viz/library/tcp"
	"chroma-viz/library/templates"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var conn *GuiConn = NewGuiConn()

type GuiConn struct {
    hub   *tcp.Connection
    eng   []*tcp.Connection
    prev  []*tcp.Connection
}

func NewGuiConn() *GuiConn {
    gui := &GuiConn{}
    gui.eng = make([]*tcp.Connection, 0, 10)
    gui.prev = make([]*tcp.Connection, 0, 10)

    return gui
}

func AddGraphicsHub(addr string, port int) {
    conn.hub = tcp.NewConnection("Hub", addr, port)
    conn.hub.Connect()
}

func AddConnection(name, conn_type, ip string, port int) error {
    if conn_type == "engine" {
        conn.eng = append(conn.eng, tcp.NewConnection(name, ip, port))
        return nil
    } else if conn_type == "preview" {
        conn.prev = append(conn.prev, tcp.NewConnection(name, ip, port))
        return nil
    }

    return fmt.Errorf("Unknown connection type %s", conn_type)
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

func CloseConn() {
    for _, c := range conn.eng {
        if c == nil {
            continue
        }

        if c.IsConnected() {
            c.CloseConn()
            log.Printf("Closed %s\n", c.Name)
        }
    }

    for _, c := range conn.prev {
        if c == nil {
            continue
        }

        if c.IsConnected() {
            c.CloseConn()
            log.Printf("Closed %s\n", c.Name)
        }
    }
}

func VizGui(app *gtk.Application) {
    win, err := gtk.ApplicationWindowNew(app)
    if err != nil {
        log.Fatal(err)
    }

    win.SetDefaultSize(800, 600)
    win.SetTitle("Chroma Viz")

    edit := editor.NewEditor(SendEngine, SendPreview)
    edit.AddAction("Take On", true, func() { SendEngine(edit.Page, tcp.ANIMATE_ON) })
    edit.AddAction("Continue", true, func() { SendEngine(edit.Page, tcp.CONTINUE) })
    edit.AddAction("Take Off", true, func() { SendEngine(edit.Page, tcp.ANIMATE_OFF) })
    edit.AddAction("Save", false, func() { 
        edit.UpdateProps()
        SendPreview(edit.Page, tcp.ANIMATE_ON) 
    })
    edit.PageEditor()

    cont := func(page *shows.Page) { SendEngine(page, tcp.CONTINUE) }

    showTree := NewShowTree(func(page *shows.Page) { edit.SetPage(page) })
    tempTree := NewTempTree(func(temp *templates.Template) { 
        page := showTree.show.AddPage(temp.Title, temp, cont)
        showTree.ImportPage(page) 
    })

    tempTree.ImportTemplates(conn.hub.Conn)

    preview := setup_preview_window()

    //testGui(tempView, showView)

    box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Fatal(err)
    }

    win.Add(box)

    /* Menu layout */
    builder, err := gtk.BuilderNew()
    if err := builder.AddFromFile("./gtk/viz-menu.ui"); err != nil {
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
    if err := builder.AddFromFile("./gtk/viz-gui.ui"); err != nil {
        log.Fatal(err)
    }

    body, err := builder.GetObject("body")
    if err != nil {
        log.Fatal(err)
    }
    
    bodyBox, ok := body.(*gtk.Paned)
    if !ok {
        log.Fatal("viz-gui.ui object 'body' is not a gtk.Paned")
    }

    box.PackStart(bodyBox, true, true, 0)

    tempScroll, err := gtkGetObject[*gtk.ScrolledWindow](builder, "templates-win")
    if err != nil {
        log.Fatal(err)
    }

    tempScroll.Add(tempTree.treeView)

    showScroll, err := gtkGetObject[*gtk.ScrolledWindow](builder, "show-win")
    if err != nil {
        log.Fatal(err)
    }

    showScroll.Add(showTree.treeView)

    editBox, err := gtkGetObject[*gtk.Box](builder, "edit")
    if err != nil {
        log.Fatal(err)
    }

    editBox.PackStart(edit.Box, true, true, 0)

    prevBox, err := gtkGetObject[*gtk.Box](builder, "preview")
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
    cont := func(page *shows.Page) { SendEngine(page, tcp.CONTINUE) }

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
        show.ImportShow(filename, cont)
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

        showTree.show.NumPages++
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

    id, err := showTree.treeList.GetValue(iter, TITLE)
    if err != nil { 
        return err 
    }

    val, err := id.GoValue()
    if err != nil { 
        return err 
    }

    title := val.(string)

    id, err = showTree.treeList.GetValue(iter, PAGENUM)
    if err != nil { 
        return err 
    }

    val, err = id.GoValue()
    if err != nil { 
        return err 
    }

    pageNum, err := strconv.Atoi(val.(string))
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

    show.treeList.Remove(iter)
    id, err := show.treeList.GetValue(iter, PAGENUM)
    if err != nil { 
        return err 
    }

    val, err := id.GoValue()
    if err != nil { 
        return err 
    }

    pageNum, err := strconv.Atoi(val.(string))
    if err != nil { 
        return err 
    }

    show.show.Pages[pageNum] = nil
    return nil
}

func testGui(tempTree *TempTree, showTree *ShowTree) {
    cont := func(page *shows.Page) { SendEngine(page, tcp.CONTINUE) }
    num_temps := 10000
    num_props := 100
    num_pages := 1000

    log.Printf("Testing with %d Templates, %d Properties, %d Pages\n", num_temps, num_props, num_pages)

    start := time.Now()
    for i := 1; i < num_temps; i++ {
        //page, _ := tempTree.AddTemplate("Template", i, LOWER_FRAME, num_props)

        // for j := 0; j < num_props; j++ {
        //     page.AddProp("Background", j, props.RECT_PROP)
        // }
    }

    t := time.Now()
    elapsed := t.Sub(start)
    log.Printf("Built Templates in %s\n", elapsed)

    start = time.Now()
    for i := 0; i < num_pages; i++ {
        index := rand.Int() % (num_temps - 1) + 1
        temp := tempTree.Temps.Temps[index]
        page := showTree.show.AddPage(temp.Title, temp, cont)
        showTree.ImportPage(page)
    }

    t = time.Now()
    elapsed = t.Sub(start)
    log.Printf("Built Show in %s\n", elapsed)
}

func gtkGetObject[T any](builder *gtk.Builder, name string) (obj T, err error) {
    gtkObject, err := builder.GetObject(name)
    if err != nil {
        return 
    }

    goObj, ok := gtkObject.(T)
    if !ok {
        err = fmt.Errorf("viz-gui.ui object '%s' is type %v", name, reflect.TypeOf(goObj))
        return 
    }

    return goObj, nil
}

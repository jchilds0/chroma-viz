package viz

import (
	"chroma-viz/editor"
	"chroma-viz/shows"
	"chroma-viz/tcp"
	"chroma-viz/templates"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
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

func SendPreview(page *shows.Page, action int) {
    for _, c := range conn.prev {
        if c == nil {
            continue
        }

        c.SetPage <- page
        c.SetAction <- action 
    }
}

func SendEngine(page *shows.Page, action int) {
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
    show := NewShowTree(func(page *shows.Page) { edit.SetPage(page) })
    temp := NewTempTree(func(temp *templates.Template) { show.ImportPage(temp.Title, temp) })

    temp.ImportTemplates(conn.hub.Conn)

    //testGui(tempView, showView)

    box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Fatal(err)
    }

    win.Add(box)

    /* Menu layout */
    builder, err := gtk.BuilderNew()
    if err := builder.AddFromFile("./gtk/menus.ui"); err != nil {
        log.Fatal(err)
    }

    menu, err := builder.GetObject("menubar")
    if err != nil {
        log.Fatal(err)
    }

    app.SetMenubar(menu.(*glib.MenuModel))

    importShow := glib.SimpleActionNew("import_show", nil)
    importShow.Connect("activate", func() { 
        err := guiImportShow(win, show, temp) 
        if err != nil {
            log.Printf("Error importing show (%s)", err)
        }
    })
    app.AddAction(importShow)

    exportShow := glib.SimpleActionNew("export_show", nil)
    exportShow.Connect("activate", func() { 
        err := guiExportShow(win, show) 
        if err != nil {
            log.Printf("Error exporting show (%s)", err)
        }
    })
    app.AddAction(exportShow)

    importPage := glib.SimpleActionNew("import_page", nil)
    importPage.Connect("activate", func() { 
        err := guiImportPage(win, temp, show) 
        if err != nil {
            log.Printf("Error importing page (%s)", err)
        }
    })
    app.AddAction(importPage)

    exportPage := glib.SimpleActionNew("export_page", nil)
    exportPage.Connect("activate", func() { 
        err := guiExportPage(win, show) 
        if err != nil {
            log.Printf("Error exporting page (%s)", err)
        }
    })
    app.AddAction(exportPage)

    deletePage := glib.SimpleActionNew("delete_page", nil)
    deletePage.Connect("activate", func() { 
        err := guiDeletePage(show) 
        if err != nil {
            log.Printf("Error deleting page (%s)", err)
        }
    })
    app.AddAction(deletePage)

    /* Body layout */
    bodyBox, err := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
    if err != nil {
        log.Fatal(err)
    }

    box.PackStart(bodyBox, true, true, 0)

    leftBox, err := gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
    if err != nil {
        log.Fatal(err)
    }

    rightBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Fatal(err)
    }

    bodyBox.Pack1(leftBox, true, true)
    bodyBox.Pack2(rightBox, true, true)

    /* left */
    leftBox.SetHExpand(true)

    /* templates */
    templates, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Print(err)
    }

    leftBox.Pack1(templates, true, true)

    header1, err := gtk.HeaderBarNew()
    if err != nil {
        log.Fatal(err)
    }

    header1.SetTitle("Templates")
    templates.PackStart(header1, false, false, 0)

    scroll1, err := gtk.ScrolledWindowNew(nil, nil)
    if err != nil {
        log.Fatal(err)
    }

    templates.PackStart(scroll1, true, true, 0)
    scroll1.Add(temp.treeView)

    /* show */
    shows, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Fatal(err)
    }

    leftBox.Pack2(shows, true, true)

    header2, err := gtk.HeaderBarNew()
    if err != nil {
        log.Fatal(err)
    }

    header2.SetTitle("Show")
    shows.PackStart(header2, false, false, 0)
    scroll2, err := gtk.ScrolledWindowNew(nil, nil)
    if err != nil {
        log.Fatal(err)
    }

    shows.PackStart(scroll2, true, true, 0)
    scroll2.Add(show.treeView)

    /* right */
    rightBox.PackStart(edit.Box(), true, true, 0)

    preview := setup_preview_window()
    rightBox.PackEnd(preview, false, false, 0)

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

func guiImportShow(win *gtk.ApplicationWindow, show *ShowTree, temp *TempTree) error {
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
        show.ImportShow(temp, filename)
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

/*
    Simple struct for encoding/decoding json format of a page
    TODO: change the property format to facilitate direct encoding
*/
type GuiPage struct {
    Title     string
    TempID    int
}

func guiImportPage(win *gtk.ApplicationWindow, tempTree *TempTree, showTree *ShowTree) error {
    dialog, err := gtk.FileChooserDialogNewWith2Buttons(
        "Import Page", win, gtk.FILE_CHOOSER_ACTION_OPEN, 
        "_Cancel", gtk.RESPONSE_CANCEL, "_Open", gtk.RESPONSE_ACCEPT)
    if err != nil {
        return err
    }
    defer dialog.Destroy()

    res := dialog.Run()
    if res == gtk.RESPONSE_ACCEPT {
        page := &GuiPage{}
        filename := dialog.GetFilename()

        buf, err := os.ReadFile(filename)
        if err != nil {
            return err
        }

        err = json.Unmarshal(buf, page)
        if err != nil {
            return err 
        }

        err = showTree.ImportPage(page.Title, tempTree.Temps.Temps[page.TempID])
        if err != nil {
            return err 
        }
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
        file, err := os.Create(filename)
        if err != nil {
            return err
        }
        defer file.Close()

        page := showTree.show.Pages[pageNum]
        if page == nil {
            return fmt.Errorf("Page %d does not exist", pageNum)
        }

        guiPage := &GuiPage{
            Title: title, 
            TempID: page.TemplateID,
        }

        buf, err := json.Marshal(guiPage)
        if err != nil {
            return err
        }

        _, err = file.Write(buf)
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
    return nil
}

func testGui(tempTree *TempTree, showTree *ShowTree) {
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
        showTree.ImportPage(temp.Title, temp)
    }

    t = time.Now()
    elapsed = t.Sub(start)
    log.Printf("Built Show in %s\n", elapsed)
}

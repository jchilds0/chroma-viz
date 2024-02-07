package viz

import (
	"chroma-viz/props"
	"chroma-viz/shows"
	"chroma-viz/tcp"
	"chroma-viz/templates"
	"fmt"
	"log"
	"math/rand"
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

    editView := NewEditor()
    showView := NewShow(func(page *shows.Page) { editView.SetPage(page) })
    tempView := NewTempTree(func(temp *templates.Template) {showView.NewShowPage(temp)})

    err = ImportTemplates(conn.hub.Conn, tempView)
    if err != nil {
        log.Printf("Error importing hub (%s)", err)
    } else {
        log.Println("Graphics hub imported")
    }

    //showView.ImportShow(tempView, "/home/josh/Documents/projects/chroma-viz/shows/simple.show")
    //testGui(tempView, showView)

    box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Fatal(err)
    }

    win.Add(box)

    /* Menu layout */
    builder, err := gtk.BuilderNew()
    if err := builder.AddFromFile("/home/josh/Documents/projects/chroma-viz/gtk/menus.ui"); err != nil {
        log.Fatal(err)
    }

    menu, err := builder.GetObject("menubar")
    if err != nil {
        log.Fatal(err)
    }

    app.SetMenubar(menu.(*glib.MenuModel))

    importAction := glib.SimpleActionNew("import_show", nil)
    importAction.Connect("activate", func() { guiImportShow(win, showView, tempView) })
    app.AddAction(importAction)

    exportAction := glib.SimpleActionNew("export_show", nil)
    exportAction.Connect("activate", func() { guiExportShow(win, showView) })
    app.AddAction(exportAction)

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
    scroll1.Add(tempView)

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
    scroll2.Add(showView)

    /* right */
    rightBox.PackStart(editView.Box(), true, true, 0)

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

func guiImportShow(win *gtk.ApplicationWindow, show *ShowTree, temp *TempTree) {
    dialog, err := gtk.FileChooserDialogNewWith2Buttons(
        "Import Show", win, gtk.FILE_CHOOSER_ACTION_OPEN, 
        "_Cancel", gtk.RESPONSE_CANCEL, "_Open", gtk.RESPONSE_ACCEPT)

    if err != nil {
        log.Print(err)
    }

    res := dialog.Run()

    if res == gtk.RESPONSE_ACCEPT {
        filename := dialog.GetFilename()
        err := show.ImportShow(temp, filename)

        if err != nil {
            log.Printf("Error importing show (%s)", err)
        }
    }
    dialog.Destroy()
}

func guiExportShow(win *gtk.ApplicationWindow, show *ShowTree) {
    dialog, err := gtk.FileChooserDialogNewWith2Buttons(
        "Save Show", win, gtk.FILE_CHOOSER_ACTION_SAVE, 
        "_Cancel", gtk.RESPONSE_CANCEL, "_Save", gtk.RESPONSE_ACCEPT)

    if err != nil {
        log.Print(err)
    }

    dialog.SetCurrentName(".show")
    res := dialog.Run()
    if res == gtk.RESPONSE_ACCEPT {
        filename := dialog.GetFilename()
        show.ExportShow(filename)
    }

    dialog.Destroy()
}

func testGui(temp *TempTree, show *ShowTree) {
    num_temps := 10000
    num_props := 100
    num_pages := 1000

    log.Printf("Testing with %d Templates, %d Properties, %d Pages\n", num_temps, num_props, num_pages)

    start := time.Now()
    for i := 1; i < num_temps; i++ {
        page, _ := temp.AddTemplate("Template", i, LOWER_FRAME, num_props)

        for j := 0; j < num_props; j++ {
            page.AddProp("Background", j, props.RECT_PROP)
        }
    }

    t := time.Now()
    elapsed := t.Sub(start)
    log.Printf("Built Templates in %s\n", elapsed)

    start = time.Now()
    for i := 0; i < num_pages; i++ {
        index := rand.Int() % (num_temps - 1) + 1
        show.NewShowPage(temp.Temps[index])
    }

    t = time.Now()
    elapsed = t.Sub(start)
    log.Printf("Built Show in %s\n", elapsed)
}

package artist

import (
	"chroma-viz/tcp"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var conn map[string]*tcp.Connection

func InitConnections(){
    conn = make(map[string]*tcp.Connection)
}

func AddConnection(name string, ip string, port int) {
    conn[name] = tcp.NewConnection(name, ip, port)
}

func CloseConn() {
    for name, c := range conn {
        if c.IsConnected() {
            c.CloseConn()
            log.Printf("Closed %s\n", name)
        }
    }
}

func ArtistGui(app *gtk.Application) {
    win, err := gtk.ApplicationWindowNew(app)
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    win.SetDefaultSize(800, 600)
    win.SetTitle("Chroma Artist")

    temp, err := NewTempTree()
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    win.Add(box)

    /* Menu layout */
    builder, err := gtk.BuilderNew()
    if err := builder.AddFromFile("/home/josh/Documents/projects/chroma-viz/gtk/menus.ui"); err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    menu, err := builder.GetObject("menubar")
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    app.SetMenubar(menu.(*glib.MenuModel))

    /* Body layout */
    bodyBox, err := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    box.PackStart(bodyBox, true, true, 0)

    leftBox, err := gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    rightBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    bodyBox.Pack1(leftBox, true, true)
    bodyBox.Pack2(rightBox, true, true)

    /* left */
    leftBox.SetHExpand(true)

    /* template */
    templates, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    leftBox.Pack1(templates, true, true)

    header1, err := gtk.HeaderBarNew()
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    header1.SetTitle("Template")
    templates.PackStart(header1, false, false, 0)

    tempActions, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    templates.PackStart(tempActions, false, false, 10)

    geoBox, err := gtk.ComboBoxTextNew()
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    tempActions.PackStart(geoBox, false, false, 10)
    geoBox.AppendText("Rectangle")
    geoBox.AppendText("Circle")
    geoBox.AppendText("Text")

    button1, err := gtk.ButtonNewWithLabel("Add Geometry")
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    tempActions.PackStart(button1, false, false, 10)

    button1.Connect("clicked", func() {
        name := geoBox.GetActiveText()
        if name == "" {
            log.Print("No geometry selected")
            return
        }

        temp.model.SetValue(temp.model.Append(nil), NAME, name)
    })


    button2, err := gtk.ButtonNewWithLabel("Remove Geometry")

    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    tempActions.PackStart(button2, false, false, 10)
    button2.Connect("clicked", func() {
        selection, err := temp.view.GetSelection()
        if err != nil {
            log.Printf("Error getting selected")
            return
        }

        _, iter, ok := selection.GetSelected()
        if !ok {
            log.Printf("No geometry selected")
            return
        }
        temp.model.Remove(iter)
    })

    scroll1, err := gtk.ScrolledWindowNew(nil, nil)
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    templates.PackStart(scroll1, true, true, 0)
    scroll1.Add(temp.view)

    /* right */
    //rightBox.PackStart(editView.Box(), true, true, 0)

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

    for name, render := range conn {
        eng := NewEngineWidget(name, render)
        lowerBox.PackStart(eng.button)
    }

    win.ShowAll()
}

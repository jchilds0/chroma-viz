package artist

import (
	"chroma-viz/hub"
	"chroma-viz/library/attribute"
	"chroma-viz/library/editor"
	"chroma-viz/library/gtk_utils"
	"chroma-viz/library/props"
	"chroma-viz/library/tcp"
	"chroma-viz/library/templates"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var conn map[string]*tcp.Connection

func InitConnections() {
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

func SendPreview(page tcp.Animator, action int) {
    if page == nil {
        log.Println("SendPreview recieved nil page")
        return
    }

    for _, c := range conn{
        if c == nil {
            continue
        }

        c.SetPage <- page
        c.SetAction <- action 
    }
}

var template = ArtistPage()
var geo_count []int
var geoms map[int]*geom

func ArtistGui(app *gtk.Application) {
    win, err := gtk.ApplicationWindowNew(app)
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    win.SetDefaultSize(800, 600)
    win.SetTitle("Chroma Artist")

    port := 9100
    geo := []string{"rect", "text", "circle", "graph", "image", "ticker", "clock"}
    geo_count = []int{10, 10, 10, 10, 10, 10, 10}
    geoms = make(map[int]*geom, len(geo))

    index := 1
    for i, name := range geo {
        geoms[props.StringToProp[name]] = newGeom(index, geo_count[i])
        index += geo_count[i]
    }

    hub.GenerateTemplateHub(geo, geo_count, "artist/artist.json")
    hub.ImportArchive("artist/artist.json")
    go hub.StartHub(port, -1)

    editView := editor.NewEditor(func(page tcp.Animator, action int) {}, SendPreview)
    editView.AddAction("Save", true, func() { 
        editView.UpdateProps()
        SendPreview(editView.Page, tcp.ANIMATE_ON) 
    })
    editView.PropertyEditor()
    editView.Page = template

    preview := setup_preview_window(port)

    tempView, err := NewTempTree(
        func(propID int) { 
            prop := template.Geometry[propID]
            editView.SetProperty(prop) 
        },
    )

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
    if err := builder.AddFromFile("./gtk/artist-menu.ui"); err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    menu, err := builder.GetObject("menubar")
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    app.SetMenubar(menu.(*glib.MenuModel))

    importPage := glib.SimpleActionNew("import_page", nil)
    importPage.Connect("activate", func() { 
        err := guiImportPage(win, tempView)
        if err != nil {
            log.Print(err)
        }
    })
    app.AddAction(importPage)

    exportPage := glib.SimpleActionNew("export_page", nil)
    exportPage.Connect("activate", func() { 
        err := guiExportPage(win, tempView)
        if err != nil {
            log.Print(err)
        }
    })
    app.AddAction(exportPage)


    /* Body layout */
    builder, err = gtk.BuilderNew()
    if err := builder.AddFromFile("./gtk/artist-gui.ui"); err != nil {
        log.Fatal(err)
    }

    body, err := gtk_utils.BuilderGetObject[*gtk.Paned](builder, "body")
    if err != nil {
        log.Fatal(err)
    }

    box.PackStart(body, true, true, 0)

    title, err := gtk_utils.BuilderGetObject[*gtk.Entry](builder, "title")
    if err != nil {
        log.Fatal(err)
    }

    title.Connect("changed", func(entry *gtk.Entry) {
        text, err := entry.GetText()
        if err != nil {
            log.Print(err)
            return
        }

        template.Title = text
    })

    tempid, err := gtk_utils.BuilderGetObject[*gtk.Entry](builder, "tempid")
    if err != nil {
        log.Fatal(err)
    }

    tempid.Connect("changed", func(entry *gtk.Entry) {
        text, err := entry.GetText()
        if err != nil {
            log.Print(err)
            return
        }

        id, err := strconv.Atoi(text)
        if err != nil {
            log.Print(err)
            return
        }

        template.TempID = id
    })

    geoSelector, err := gtk_utils.BuilderGetObject[*gtk.ComboBoxText](builder, "geo-selector")
    if err != nil {
        log.Fatal(err)
    }

    addGeo, err := gtk_utils.BuilderGetObject[*gtk.Button](builder, "add-geo")
    if err != nil {
        log.Fatal(err)
    }

    addGeo.Connect("clicked", func() {
        name := geoSelector.GetActiveText()
        if name == "" {
            log.Print("No geometry selected")
            return
        }

        propNum, err := AddProp(name)
        if err != nil {
            log.Print(err)
            return
        }

        newRow := tempView.model.Append(nil)
        tempView.model.SetValue(newRow, NAME, name)
        tempView.model.SetValue(newRow, PROP_NUM, propNum)
    })

    removeGeo, err := gtk_utils.BuilderGetObject[*gtk.Button](builder, "remove-geo")
    if err != nil {
        log.Fatal(err)
    }

    removeGeo.Connect("clicked", func() {
        selection, err := tempView.view.GetSelection()
        if err != nil {
            log.Printf("Error getting selected")
            return
        }

        _, iter, ok := selection.GetSelected()
        if !ok {
            log.Printf("No geometry selected")
            return
        }

        model := &tempView.model.TreeModel
        propID, err := gtk_utils.ModelGetValue[int](model, iter, PROP_NUM)
        if err != nil {
            log.Printf("Error getting prop id (%s)", err)
            return
        }

        RemoveProp(propID)
        tempView.model.Remove(iter)
    })


    geoScroll, err := gtk_utils.BuilderGetObject[*gtk.ScrolledWindow](builder, "geo-win")
    if err != nil {
        log.Fatal(err)
    }

    geoScroll.Add(tempView.view)

    editBox, err := gtk_utils.BuilderGetObject[*gtk.Box](builder, "edit")
    if err != nil {
        log.Fatal(err)
    }

    editBox.PackStart(editView.Box, true, true, 0)
    
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

    for name, render := range conn {
        eng := NewEngineWidget(name, render)
        lowerBox.PackStart(eng.button)
    }

    win.ShowAll()
}

var visible = map[string]bool {
    "x": true,
    "y": true, 
    "width": true,
    "height": true,
    "inner_radius": true,
    "outer_radius": true,
    "start_angle": true,
    "end_angle": true,
    "color": true,
    "string": true,
    "node": true,
    "text": true,
    "clock": true,
    "scale": true,
    "parent": true,
}

func ArtistPage() *templates.Template {
    page := &templates.Template{
        Layer: 0,
        TempID: 0,
    }

    page.Geometry = make(map[int]*props.Property)

    return page
}

var geo_type = map[string]int {
    "Rectangle": props.RECT_PROP,
    "Circle": props.CIRCLE_PROP,
    "Text": props.TEXT_PROP,
    "Graph": props.GRAPH_PROP,
    "Ticker": props.TICKER_PROP,
    "Clock": props.CLOCK_PROP,
    "Image": props.IMAGE_PROP,
}

func AddProp(label string) (id int, err error) {
    cont := func() { SendPreview(template, tcp.CONTINUE) }

    geo_typed, ok := geo_type[label]
    if !ok {
        return 0, fmt.Errorf("Unknown label %s", label)
    }

    geom, ok := geoms[geo_typed]
    if !ok {
        return 0, fmt.Errorf("Unknown geom %s", label)
    }

    id, err = geom.allocGeom()
    if err != nil {
        return 
    }

    template.Geometry[id] = props.NewProperty(geo_typed, label, visible, cont)
    template.Geometry[id].Attr["parent"] = attribute.NewIntAttribute("parent")
    return
}

func RemoveProp(propID int) {
    prop := template.Geometry[propID]
    if prop == nil {
        log.Printf("No prop with prop id %d", propID)
        return
    }

    geom, ok := geoms[prop.PropType]
    if !ok {
        log.Printf("No geom with prop type %d", prop.PropType)
        return
    }

    geom.freeGeom(propID)
    template.Geometry[propID] = nil
}

func guiImportPage(win *gtk.ApplicationWindow, temp *TempTree) error {
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

        buf, err := os.ReadFile(filename)
        if err != nil {
            return err 
        }

        var newTemp templates.Template
        err = json.Unmarshal(buf, &newTemp)
        if err != nil {
            return err
        }

        // reset temp view geometry
        temp.model, err = gtk.TreeStoreNew(glib.TYPE_OBJECT, glib.TYPE_STRING, glib.TYPE_OBJECT, glib.TYPE_INT)
        temp.view.SetModel(temp.model)

        // reset geometry allocs

        importTemplate(template, &newTemp)

        // add geometry to temp view
        for id, geo := range template.Geometry {
            newRow := temp.model.Append(nil)
            temp.model.SetValue(newRow, NAME, geo.Name)
            temp.model.SetValue(newRow, PROP_NUM, id) 
        }
    }

    return nil
}

func guiExportPage(win *gtk.ApplicationWindow, temp *TempTree) error {
    dialog, err := gtk.FileChooserDialogNewWith2Buttons(
        "Save Template", win, gtk.FILE_CHOOSER_ACTION_SAVE, 
        "_Cancel", gtk.RESPONSE_CANCEL, "_Save", gtk.RESPONSE_ACCEPT)
    if err != nil {
        return err
    }
    defer dialog.Destroy()

    dialog.SetCurrentName(template.Title + ".json")
    res := dialog.Run()
    if res == gtk.RESPONSE_ACCEPT {
        filename := dialog.GetFilename()

        newTemp := templates.NewTemplate(
            template.Title, 
            template.TempID, 
            template.Layer, 
            template.NumGeo,
        )

        // sync parent attrs
        model := &temp.model.TreeModel
        if iter, ok := model.GetIterFirst(); ok {
            updateParentGeometry(model, iter, 0)
        }

        compressGeometry(template, newTemp)

        // TODO: sync visible attrs to template

        err := templates.ExportTemplate(newTemp, filename)
        if err != nil {
            return err
        }
    }

    return nil
}

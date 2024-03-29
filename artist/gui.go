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
	"time"

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

    tempView, err := NewTempTree(
        func(propID int) { 
            prop := template.Geometry[propID]
            editView.SetProperty(prop) 
        },
    )

    editView.AddAction("Save", true, func() { 
        // sync parent attrs
        model := tempView.model.ToTreeModel()
        if iter, ok := model.GetIterFirst(); ok {
            updateParentGeometry(model, iter, 0)
        }

        tempid := template.TempID
        template.TempID = 0

        editView.UpdateProps()
        SendPreview(editView.Page, tcp.ANIMATE_ON) 
        time.Sleep(50 * time.Millisecond)
        template.TempID = tempid
    })

    editView.PropertyEditor()
    editView.Page = template

    preview := setup_preview_window(port)

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

    importPage := glib.SimpleActionNew("import_page", nil)
    app.AddAction(importPage)

    exportPage := glib.SimpleActionNew("export_page", nil)
    app.AddAction(exportPage)

    app.SetMenubar(menu.(*glib.MenuModel))

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

    tempid, err := gtk_utils.BuilderGetObject[*gtk.Entry](builder, "tempid")
    if err != nil {
        log.Fatal(err)
    }

    layer, err := gtk_utils.BuilderGetObject[*gtk.Entry](builder, "layer")
    if err != nil {
        log.Fatal(err)
    }

    animSelector, err := gtk_utils.BuilderGetObject[*gtk.ComboBoxText](builder, "anim-selector")
    if err != nil {
        log.Fatal(err)
    }

    geoSelector, err := gtk_utils.BuilderGetObject[*gtk.ComboBoxText](builder, "geo-selector")
    if err != nil {
        log.Fatal(err)
    }

    addGeo, err := gtk_utils.BuilderGetObject[*gtk.Button](builder, "add-geo")
    if err != nil {
        log.Fatal(err)
    }

    removeGeo, err := gtk_utils.BuilderGetObject[*gtk.Button](builder, "remove-geo")
    if err != nil {
        log.Fatal(err)
    }

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

    /* actions */
    importPage.Connect("activate", func() { 
        err := guiImportPage(win, tempView)
        if err != nil {
            log.Print(err)
        }

        switch template.AnimateOn {
        case "":
            animSelector.SetActive(0)
        case "left_to_right":
            animSelector.SetActive(1)
        case "up":
            animSelector.SetActive(2)
        case "right_to_left":
            animSelector.SetActive(3)
        default:
        log.Printf("Unknown animation %s", template.AnimateOn)
        }

        title.SetText(template.Title)
        tempid.SetText(strconv.Itoa(template.TempID))
        layer.SetText(strconv.Itoa(template.Layer))
    })

    exportPage.Connect("activate", func() { 
        switch animSelector.GetActiveText() {
        case "None", "":
            template.AnimateOn = ""
            template.AnimateOff = ""
        case "Left":
            template.AnimateOn = "left_to_right"
            template.AnimateOff = "left_to_right"
        case "Up":
            template.AnimateOn = "up"
        case "Right":
            template.AnimateOn = "right_to_left"
            template.AnimateOff = "right_to_left"
        default:
            log.Printf("Unknown animation %s", animSelector.GetActiveText())
        }

        err := guiExportPage(win, tempView)
        if err != nil {
            log.Print(err)
        }
    })

    title.Connect("changed", func(entry *gtk.Entry) {
        text, err := entry.GetText()
        if err != nil {
            log.Print(err)
            return
        }

        template.Title = text
    })

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

    layer.Connect("changed", func(entry *gtk.Entry) {
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

        template.Layer = id
    })

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

        iter := tempView.model.Append(nil)
        tempView.AddRow(iter, name, name, propNum)
    })

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

var geo_name = map[int]string {
    props.RECT_PROP: "Rectangle",
    props.CIRCLE_PROP: "Circle",
    props.TEXT_PROP: "Text",
    props.GRAPH_PROP: "Graph",
    props.TICKER_PROP: "Ticker",
    props.CLOCK_PROP: "Clock",
    props.IMAGE_PROP: "Image",
}

func AddProp(label string) (id int, err error) {
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

    template.Geometry[id] = props.NewProperty(geo_typed, label, true, nil)
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
        temp.Clean()

        template.Title = newTemp.Title
        template.TempID = newTemp.TempID
        template.Layer = newTemp.Layer
        template.NumGeo = len(newTemp.Geometry)

        decompressGeometry(template, &newTemp)
        geometryToTreeView(temp, nil, 0)

        // set temp switch to true to send all props to chroma engine
        for _, geo := range template.Geometry {
            geo.SetTemp(true)
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
            len(template.Geometry),
        )

        newTemp.AnimateOn = template.AnimateOn
        newTemp.Continue = template.Continue
        newTemp.AnimateOff = template.AnimateOff

        // sync parent attrs
        model := temp.model.ToTreeModel()
        if iter, ok := model.GetIterFirst(); ok {
            updateParentGeometry(model, iter, 0)
        }

        compressGeometry(template, newTemp, temp.model.ToTreeModel())

        // TODO: sync visible attrs to template

        err := templates.ExportTemplate(newTemp, filename)
        if err != nil {
            return err
        }
    }

    return nil
}

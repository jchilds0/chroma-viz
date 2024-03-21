package artist

import (
	"chroma-viz/hub"
	"chroma-viz/library/attribute"
	"chroma-viz/library/editor"
	"chroma-viz/library/props"
	"chroma-viz/library/tcp"
	"chroma-viz/library/templates"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
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

    index := 0
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

    temp, err := NewTempTree(
        func(propID int) { 
            prop := template.GetPropMap()[propID]
            editView.SetProperty(prop) 
        },
        func(propID, parentID int) {
            prop := template.GetPropMap()[propID]
            if prop == nil {
                return
            }

            parentAttr := prop.Attr["parent"]
            if parentAttr == nil {
                return
            }

            parentAttr.(*attribute.IntAttribute).Value = parentID
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
        dialog, err := gtk.FileChooserDialogNewWith2Buttons(
            "Import Page", win, gtk.FILE_CHOOSER_ACTION_OPEN, 
            "_Cancel", gtk.RESPONSE_CANCEL, "_Open", gtk.RESPONSE_ACCEPT)
        if err != nil {
            log.Print(err)
            return 
        }
        defer dialog.Destroy()

        res := dialog.Run()
        if res == gtk.RESPONSE_ACCEPT {
            filename := dialog.GetFilename()

            buf, err := os.ReadFile(filename)
            if err != nil {
                log.Print(err)
                return 
            }

            err = json.Unmarshal(buf, template)
            if err != nil {
                log.Print(err)
                return 
            }

            temp.model, err = gtk.TreeStoreNew(glib.TYPE_OBJECT, glib.TYPE_STRING, glib.TYPE_OBJECT, glib.TYPE_INT)
            temp.view.SetModel(temp.model)

            // build a map of json geo id's to new geo id's
            geoRename := make(map[int]int, len(template.Geometry))
            newGeoMap := make(map[int]*props.Property, len(template.Geometry))

            for id, geo := range template.Geometry {
                geom, ok := geoms[geo.PropType]
                if !ok {
                    log.Printf("Missing Geom %s", props.PropType(geo.PropType))
                    continue
                }

                newID, err := geom.allocGeom()
                if err != nil {
                    log.Print(err)
                    continue
                }

                newGeoMap[newID] = geo
                geoRename[id] = newID

                newRow := temp.model.Append(nil)
                temp.model.SetValue(newRow, NAME, geo.Name)
                temp.model.SetValue(newRow, PROP_NUM, newID)
            }

            template.Geometry = newGeoMap

            // update parent geo id's using geoRename
            for _, geo := range template.Geometry {
                parentAttr := geo.Attr["parent"]
                if parentAttr == nil {
                    log.Printf("Missing parent attr for geo %s", geo.Name)
                    continue
                }

                parent := parentAttr.(*attribute.IntAttribute).Value
                parentAttr.(*attribute.IntAttribute).Value = geoRename[parent]
            }

            // set props to visible
            for _, geo := range template.Geometry {
                geo.Visible = visible
            }

            return
        }
    })
    app.AddAction(importPage)

    exportPage := glib.SimpleActionNew("export_page", nil)
    exportPage.Connect("activate", func() { 
        dialog, err := gtk.FileChooserDialogNewWith2Buttons(
        "Save Template", win, gtk.FILE_CHOOSER_ACTION_SAVE, 
        "_Cancel", gtk.RESPONSE_CANCEL, "_Save", gtk.RESPONSE_ACCEPT)
        if err != nil {
            log.Print(err)
            return 
        }
        defer dialog.Destroy()

        dialog.SetCurrentName(template.Title + ".json")
        res := dialog.Run()
        if res == gtk.RESPONSE_ACCEPT {
            filename := dialog.GetFilename()

            // build geo id map 
            geoRename := make(map[int]int, len(template.Geometry))
            newGeoMap := make(map[int]*props.Property, len(template.Geometry))

            i := 0
            for id, geo := range template.Geometry {
                newGeoMap[i] = geo
                geoRename[id] = i 
                i++
            }

            template.Geometry = newGeoMap

            // update parent geo id's
            for id, geo := range template.Geometry {
                parentAttr := geo.Attr["parent"]
                if parentAttr == nil {
                    continue
                }

                attr := parentAttr.(*attribute.IntAttribute)
                if attr.Value != index {
                    continue
                }

                attr.Value = geoRename[id]
            }

            // TODO: sync visible attrs to template

            err := templates.ExportTemplate(template, filename)
            if err != nil {
                log.Print(err)
                return
            }
        }
    })
    app.AddAction(exportPage)


    /* Body layout */
    builder, err = gtk.BuilderNew()
    if err := builder.AddFromFile("./gtk/artist-gui.ui"); err != nil {
        log.Fatal(err)
    }

    body, err := gtkGetObject[*gtk.Paned](builder, "body")
    if err != nil {
        log.Fatal(err)
    }

    box.PackStart(body, true, true, 0)

    title, err := gtkGetObject[*gtk.Entry](builder, "title")
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

    tempid, err := gtkGetObject[*gtk.Entry](builder, "tempid")
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

    geoSelector, err := gtkGetObject[*gtk.ComboBoxText](builder, "geo-selector")
    if err != nil {
        log.Fatal(err)
    }

    addGeo, err := gtkGetObject[*gtk.Button](builder, "add-geo")
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

        newRow := temp.model.Append(nil)
        temp.model.SetValue(newRow, NAME, name)
        temp.model.SetValue(newRow, PROP_NUM, propNum)
    })

    removeGeo, err := gtkGetObject[*gtk.Button](builder, "remove-geo")
    if err != nil {
        log.Fatal(err)
    }

    removeGeo.Connect("clicked", func() {
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
        id, err := temp.model.GetValue(iter, PROP_NUM)
        if err != nil {
            log.Printf("Error removing prop (%s)", err)
            return
        }

        val, err := id.GoValue()
        if err != nil {
            log.Printf("Error removing prop (%s)", err)
            return
        }

        propID, ok := val.(int)
        if !ok {
            log.Printf("Error removing prop (%s)", err)
            return
        }

        RemoveProp(propID)
    })


    geoScroll, err := gtkGetObject[*gtk.ScrolledWindow](builder, "geo-win")
    if err != nil {
        log.Fatal(err)
    }

    geoScroll.Add(temp.view)

    editBox, err := gtkGetObject[*gtk.Box](builder, "edit")
    if err != nil {
        log.Fatal(err)
    }

    editBox.PackStart(editView.Box, true, true, 0)
    
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

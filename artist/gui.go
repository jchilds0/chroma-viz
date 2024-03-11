package artist

import (
	"chroma-viz/attribute"
	"chroma-viz/editor"
	"chroma-viz/props"
	"chroma-viz/shows"
	"chroma-viz/tcp"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jchilds0/chroma-hub/chroma_hub"
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

func SendPreview(page *shows.Page, action int) {
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

var page = ArtistPage()
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
    geo := []string{"rect", "text", "circle", "graph", "image"}
    geo_count = []int{10, 10, 10, 10, 10}
    geoms = make(map[int]*geom, len(geo))

    index := 0
    for i, name := range geo {
        geoms[props.StringToProp[name]] = newGeom(index, geo_count[i])
        index += geo_count[i]
    }

    chroma_hub.GenerateTemplateHub(geo, geo_count, "artist/artist.json")
    go chroma_hub.StartHub(port, -1, "artist/artist.json")

    editView := editor.NewEditor(func(page *shows.Page, action int) {}, SendPreview)
    editView.AddAction("Save", true, func() { 
        editView.UpdateProps()
        SendPreview(editView.Page, tcp.ANIMATE_ON) 
    })
    editView.PropertyEditor()
    editView.Page = page

    temp, err := NewTempTree(
        func(propID int) { 
            editView.SetProperty(page.PropMap[propID]) 
        },
        func(propID, parentID int) {
            page.PropMap[propID].Attr["parent"].(*attribute.IntAttribute).Value = parentID
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

        propNum, err := AddProp(name)
        if err != nil {
            log.Print(err)
            return
        }

        newRow := temp.model.Append(nil)
        temp.model.SetValue(newRow, NAME, name)
        temp.model.SetValue(newRow, PROP_NUM, propNum)
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

    scroll1, err := gtk.ScrolledWindowNew(nil, nil)
    if err != nil {
        log.Fatalf("Error starting artist gui (%s)", err)
    }

    templates.PackStart(scroll1, true, true, 0)
    scroll1.Add(temp.view)

    /* right */
    rightBox.PackStart(editView.Box, true, true, 0)

    preview := setup_preview_window(port)
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

func ArtistPage() *shows.Page {
    page := &shows.Page{
        Layer: 0,
        TemplateID: 0,
        PageNum: 0,
    }

    page.PropMap = make(map[int]*props.Property)

    return page
}

var visible = map[string]bool{
    "x": true,
    "y": true,
    "parent": true,
    "width": true,
    "height": true,
    "inner_radius": true,
    "outer_radius": true,
    "start_angle": true,
    "end_angle": true,
    "text": true,
    "image": true,
    "graph": true,
    "string": true,
    "color": true,
}

var geo_type = map[string]int {
    "Rectangle": props.RECT_PROP,
    "Circle": props.CIRCLE_PROP,
    "Text": props.TEXT_PROP,
}

func AddProp(label string) (id int, err error) {
    cont := func() { SendPreview(page, tcp.CONTINUE) }

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

    page.PropMap[id] = props.NewProperty(geo_typed, label, visible, cont)
    page.PropMap[id].Attr["parent"] = attribute.NewIntAttribute("parent", "parent")
    return
}

func RemoveProp(propID int) {
    prop := page.PropMap[propID]
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
    page.PropMap[propID] = nil
}


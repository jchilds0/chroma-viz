package artist

import (
	"chroma-viz/hub"
	"chroma-viz/library/attribute"
	"chroma-viz/library/config"
	"chroma-viz/library/editor"
	"chroma-viz/library/gtk_utils"
	"chroma-viz/library/pages"
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
var conf *config.Config
var chromaHub *hub.DataBase

func SendPreview(page tcp.Animator, action int) {
	if page == nil {
		log.Println("SendPreview recieved nil page")
		return
	}

	for _, c := range conn {
		if c == nil {
			continue
		}

		c.SetPage <- page
		c.SetAction <- action
	}
}

var page = pages.NewPage(0, 0, 0, 10, "")
var geoms map[int]*geom

func ArtistGui(app *gtk.Application) {
	var tempIDEntry, titleEntry, layerEntry *gtk.Entry

	win, err := gtk.ApplicationWindowNew(app)
	if err != nil {
		log.Fatalf("Error starting artist gui (%s)", err)
	}

	win.SetDefaultSize(800, 600)
	win.SetTitle("Chroma Artist")

	editView, err := editor.NewEditor(func(page tcp.Animator, action int) {}, SendPreview)
	if err != nil {
		log.Fatal(err)
	}

	tempView, err := NewTempTree(
		func(propID int) {
			prop := page.PropMap[propID]
			editView.SetProperty(prop)
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	editView.AddAction("Save", true, func() {
		title, err := titleEntry.GetText()
		if err != nil {
			log.Print(err)
		}

		layer, err := layerEntry.GetText()
		if err != nil {
			log.Print(err)
		}

		tempID, err := tempIDEntry.GetText()
		if err != nil {
			log.Print(err)
		}

		template, err := artistPageToTemplate(*page, tempView, tempID, title, layer)
		if err != nil {
			log.Printf("Error creating template (%s)", err)
			return
		}

		err = chromaHub.ImportTemplate(*template)
		if err != nil {
			log.Printf("Error sending template to chroma hub (%s)", err)
			return
		}

		page.TemplateID = int(template.TempID)
		editView.UpdateProps()
		SendPreview(editView.Page, tcp.UPDATE)
		SendPreview(editView.Page, tcp.ANIMATE_ON)
	})

	editView.PropertyEditor()
	editView.Page = page

	preview := setupPreviewWindow(conf.HubPort, conf.PreviewDirectory, conf.PreviewName)

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

	newTemplate := glib.SimpleActionNew("new_template", nil)
	app.AddAction(newTemplate)

	importTemplate := glib.SimpleActionNew("import_template", nil)
	app.AddAction(importTemplate)

	exportTemplate := glib.SimpleActionNew("export_template", nil)
	app.AddAction(exportTemplate)

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

	titleEntry, err = gtk_utils.BuilderGetObject[*gtk.Entry](builder, "title")
	if err != nil {
		log.Fatal(err)
	}

	tempIDEntry, err = gtk_utils.BuilderGetObject[*gtk.Entry](builder, "tempid")
	if err != nil {
		log.Fatal(err)
	}

	layerEntry, err = gtk_utils.BuilderGetObject[*gtk.Entry](builder, "layer")
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

	geoScroll.Add(tempView.geoView)

	keyScroll, err := gtk_utils.BuilderGetObject[*gtk.ScrolledWindow](builder, "key-win")
	if err != nil {
		log.Fatal(err)
	}

	keyScroll.Add(tempView.keyView)

	keyGeo, err := gtk_utils.BuilderGetObject[*gtk.ComboBox](builder, "key-geo")
	if err != nil {
		log.Fatal(err)
	}

	keyAttr, err := gtk_utils.BuilderGetObject[*gtk.ComboBoxText](builder, "key-attr")
	if err != nil {
		log.Fatal(err)
	}

	addKey, err := gtk_utils.BuilderGetObject[*gtk.Button](builder, "add-key")
	if err != nil {
		log.Fatal(err)
	}

	removeKey, err := gtk_utils.BuilderGetObject[*gtk.Button](builder, "remove-key")
	if err != nil {
		log.Fatal(err)
	}

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
	newTemplate.Connect("activate", func() {
		tempView.geoModel.Clear()
		tempView.keyModel.Clear()

		page.PropMap = make(map[int]*props.Property)

		titleEntry.SetText("")
		tempIDEntry.SetText("")
		layerEntry.SetText("")

		SendPreview(editView.Page, tcp.UPDATE)
	})

	importTemplate.Connect("activate", func() {
		title, tempID, layer, err := guiImportPage(win, tempView)
		if err != nil {
			log.Print(err)
		}

		titleEntry.SetText(title)
		tempIDEntry.SetText(tempID)
		layerEntry.SetText(layer)
	})

	exportTemplate.Connect("activate", func() {
		title, err := titleEntry.GetText()
		if err != nil {
			log.Print(err)
		}

		layer, err := layerEntry.GetText()
		if err != nil {
			log.Print(err)
		}

		tempID, err := tempIDEntry.GetText()
		if err != nil {
			log.Print(err)
		}

		err = guiExportPage(win, tempView, title, tempID, layer)
		if err != nil {
			log.Print(err)
		}
	})

	titleEntry.Connect("changed", func(entry *gtk.Entry) {
		text, err := entry.GetText()
		if err != nil {
			log.Print(err)
			return
		}

		page.Title = text
	})

	tempIDEntry.Connect("changed", func(entry *gtk.Entry) {
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

		page.TemplateID = id
	})

	layerEntry.Connect("changed", func(entry *gtk.Entry) {
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

		page.Layer = id
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

		iter := tempView.geoModel.Append(nil)
		tempView.AddGeoRow(iter, name, name, propNum)
	})

	removeGeo.Connect("clicked", func() {
		selection, err := tempView.geoView.GetSelection()
		if err != nil {
			log.Printf("Error getting selected")
			return
		}

		_, iter, ok := selection.GetSelected()
		if !ok {
			log.Printf("No geometry selected")
			return
		}

		model := tempView.geoModel.ToTreeModel()
		geoID, err := gtk_utils.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Printf("Error getting prop id (%s)", err)
			return
		}

		RemoveProp(geoID)
		tempView.geoModel.Remove(iter)
		tempView.removeKeys(geoID)
		tempView.removeGeometry(geoID)
	})

	geoCell, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal(err)
	}

	keyGeo.PackStart(geoCell, true)
	keyGeo.CellLayout.AddAttribute(geoCell, "text", GEO_NAME)
	keyGeo.SetActive(GEO_NAME)
	keyGeo.SetModel(tempView.geoList)

	keyGeo.Connect("changed", func() {
		iter, err := keyGeo.GetActiveIter()
		if err != nil {
			log.Printf("No geometry selected")
			return
		}

		geoID, err := gtk_utils.ModelGetValue[int](tempView.geoList.ToTreeModel(), iter, GEO_NUM)
		if err != nil {
			log.Printf("Error getting geo id (%s)", err)
			return
		}

		keyAttr.RemoveAll()

		for name := range page.PropMap[geoID].Attr {
			keyAttr.AppendText(name)
		}
	})

	addKey.Connect("clicked", func() {
		iter, err := keyGeo.GetActiveIter()
		if err != nil {
			log.Printf("No geometry selected")
			return
		}

		model := tempView.geoList.ToTreeModel()

		geo, err := gtk_utils.ModelGetValue[string](model, iter, GEO_NAME)
		if err != nil {
			log.Printf("Error getting geo name (%s)", err)
			return
		}

		geoID, err := gtk_utils.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Printf("Error getting geo id (%s)", err)
			return
		}

		attr := keyAttr.GetActiveText()

		tempView.AddKeyRow(tempView.keyModel.Append(nil), geo, geoID, attr)
	})

	removeKey.Connect("clicked", func() {
		selection, err := tempView.keyView.GetSelection()
		if err != nil {
			log.Printf("Error getting selected")
			return
		}

		_, iter, ok := selection.GetSelected()
		if !ok {
			log.Printf("No geometry selected")
			return
		}

		tempView.keyModel.Remove(iter)
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

var geo_type = map[string]int{
	"Rectangle": props.RECT_PROP,
	"Circle":    props.CIRCLE_PROP,
	"Text":      props.TEXT_PROP,
	"Graph":     props.GRAPH_PROP,
	"Ticker":    props.TICKER_PROP,
	"Clock":     props.CLOCK_PROP,
	"Image":     props.IMAGE_PROP,
}

var geo_name = map[int]string{
	props.RECT_PROP:   "Rectangle",
	props.CIRCLE_PROP: "Circle",
	props.TEXT_PROP:   "Text",
	props.GRAPH_PROP:  "Graph",
	props.TICKER_PROP: "Ticker",
	props.CLOCK_PROP:  "Clock",
	props.IMAGE_PROP:  "Image",
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

	page.PropMap[id] = props.NewProperty(geo_typed, label, true, nil)
	page.PropMap[id].Attr["parent"] = attribute.NewIntAttribute("parent")
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

func guiImportPage(win *gtk.ApplicationWindow, tempView *TempTree) (title, tempID, layer string, err error) {
	dialog, err := gtk.FileChooserDialogNewWith2Buttons(
		"Import Page", win, gtk.FILE_CHOOSER_ACTION_OPEN,
		"_Cancel", gtk.RESPONSE_CANCEL, "_Open", gtk.RESPONSE_ACCEPT)
	if err != nil {
		return
	}
	defer dialog.Destroy()

	res := dialog.Run()
	if res == gtk.RESPONSE_ACCEPT {
		filename := dialog.GetFilename()

		var buf []byte
		buf, err = os.ReadFile(filename)
		if err != nil {
			return
		}

		var temp templates.Template
		err = json.Unmarshal(buf, &temp)
		if err != nil {
			return
		}

		// reset temp view geometry
		tempView.Clean()

		title = temp.Title
		tempID = strconv.FormatInt(temp.TempID, 10)
		layer = strconv.Itoa(temp.Layer)

		page = pages.NewPageFromTemplate(&temp)
		geometryToTreeView(page, tempView, nil, 0)

		tempView.addKeyframes(&temp)

		// set temp switch to true to send all props to chroma engine
		for _, geo := range page.PropMap {
			geo.SetTemp(true)
		}
	}

	return
}

func guiExportPage(win *gtk.ApplicationWindow, tempView *TempTree, title, tempID, layer string) error {
	dialog, err := gtk.FileChooserDialogNewWith2Buttons(
		"Save Template", win, gtk.FILE_CHOOSER_ACTION_SAVE,
		"_Cancel", gtk.RESPONSE_CANCEL, "_Save", gtk.RESPONSE_ACCEPT)
	if err != nil {
		return err
	}
	defer dialog.Destroy()

	dialog.SetCurrentName(page.Title + ".json")
	res := dialog.Run()
	if res == gtk.RESPONSE_ACCEPT {
		filename := dialog.GetFilename()

		template, err := artistPageToTemplate(*page, tempView, tempID, title, layer)
		if err != nil {
			return fmt.Errorf("Error creating template (%s)", err)
		}

		err = template.ExportTemplate(filename)
		if err != nil {
			return err
		}
	}

	return nil
}

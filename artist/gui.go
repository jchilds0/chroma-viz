package artist

import (
	"chroma-viz/hub"
	"chroma-viz/library/attribute"
	"chroma-viz/library/config"
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

var template = ArtistPage()
var geoms map[int]*geom

func ArtistGui(app *gtk.Application) {
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
			prop := template.Geometry[propID]
			editView.SetProperty(prop)
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	editView.AddAction("Save", true, func() {
		// sync parent attrs
		geoModel := tempView.geoModel.ToTreeModel()
		if iter, ok := geoModel.GetIterFirst(); ok {
			updateParentGeometry(geoModel, iter, 0)
		}

		tempid := template.TempID
		template.TempID = 0
		template.Keyframe = tempView.keyframes()
		template.NumKeyframe = len(template.Keyframe)

		editView.UpdateProps()
		SendPreview(editView.Page, tcp.ANIMATE_ON)
		time.Sleep(50 * time.Millisecond)
		template.TempID = tempid
	})

	editView.PropertyEditor()
	editView.Page = template

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

		template.Title = ""
		template.TempID = 0
		template.Layer = 0
		template.Geometry = make(map[int]*props.Property)

		title.SetText(template.Title)
		tempid.SetText("")
		layer.SetText("")

		SendPreview(editView.Page, tcp.CLEAN)
	})

	importTemplate.Connect("activate", func() {
		err := guiImportPage(win, tempView)
		if err != nil {
			log.Print(err)
		}

		title.SetText(template.Title)
		tempid.SetText(strconv.Itoa(template.TempID))
		layer.SetText(strconv.Itoa(template.Layer))
	})

	exportTemplate.Connect("activate", func() {
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

		for name := range template.Geometry[geoID].Attr {
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

func ArtistPage() *templates.Template {
	page := &templates.Template{
		Layer:  0,
		TempID: 0,
	}

	page.Geometry = make(map[int]*props.Property)

	return page
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
		template.NumKeyframe = len(newTemp.Keyframe)

		decompressGeometry(template, &newTemp)
		geometryToTreeView(temp, nil, 0)

		temp.addKeyframes(template)

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

		// get keyframes
		template.Keyframe = temp.keyframes()
		template.NumKeyframe = len(template.Keyframe)

		newTemp := templates.NewTemplate(
			template.Title,
			template.TempID,
			template.Layer,
			len(template.Geometry),
			len(template.Keyframe),
		)

		// sync parent attrs
		model := temp.geoModel.ToTreeModel()
		if iter, ok := model.GetIterFirst(); ok {
			updateParentGeometry(model, iter, 0)
		}

		compressGeometry(template, newTemp, temp.geoModel.ToTreeModel())

		// TODO: sync visible attrs to template

		err := templates.ExportTemplate(newTemp, filename)
		if err != nil {
			return err
		}
	}

	return nil
}

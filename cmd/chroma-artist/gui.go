package main

import (
	"chroma-viz/library"
	"chroma-viz/library/attribute"
	"chroma-viz/library/hub"
	"chroma-viz/library/pages"
	"chroma-viz/library/templates"
	"chroma-viz/library/util"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var conn map[string]*library.Connection
var conf *library.Config

func SendPreview(page library.Animator, action int) {
	if page == nil {
		log.Println("SendPreview recieved nil page")
		return
	}

	for _, c := range conn {
		if c == nil {
			continue
		}

		if !c.IsConnected() {
			continue
		}

		c.SetPage <- page
		c.SetAction <- action
	}
}

func ArtistGui(app *gtk.Application) {
	var tempIDEntry, titleEntry, layerEntry *gtk.Entry
	template := templates.NewTemplate("", 1, 0, 10, 10)

	win, err := gtk.ApplicationWindowNew(app)
	if err != nil {
		log.Fatalf("Error starting artist gui (%s)", err)
	}

	win.SetDefaultSize(800, 600)
	win.SetTitle("Chroma Artist")

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatalf("Error starting artist gui (%s)", err)
	}

	win.Add(box)

	/* Menu layout */
	builder, err := gtk.BuilderNew()
	if err := builder.AddFromFile("artist/menu.ui"); err != nil {
		log.Fatalf("Error starting artist gui (%s)", err)
	}

	menu, err := builder.GetObject("menubar")
	if err != nil {
		log.Fatalf("Error starting artist gui (%s)", err)
	}

	newTemplate := glib.SimpleActionNew("new_template", nil)
	app.AddAction(newTemplate)

	importTemplateJSON := glib.SimpleActionNew("import_template_json", nil)
	app.AddAction(importTemplateJSON)

	importTemplateHub := glib.SimpleActionNew("import_template_hub", nil)
	app.AddAction(importTemplateHub)

	exportTemplate := glib.SimpleActionNew("export_template", nil)
	app.AddAction(exportTemplate)

	app.SetMenubar(menu.(*glib.MenuModel))

	/* Body layout */
	builder, err = gtk.BuilderNew()
	if err := builder.AddFromFile("./artist/gui.ui"); err != nil {
		log.Fatal(err)
	}

	body, err := util.BuilderGetObject[*gtk.Paned](builder, "body")
	if err != nil {
		log.Fatal(err)
	}

	box.PackStart(body, true, true, 0)

	titleEntry, err = util.BuilderGetObject[*gtk.Entry](builder, "title")
	if err != nil {
		log.Fatal(err)
	}

	tempIDEntry, err = util.BuilderGetObject[*gtk.Entry](builder, "tempid")
	if err != nil {
		log.Fatal(err)
	}

	layerEntry, err = util.BuilderGetObject[*gtk.Entry](builder, "layer")
	if err != nil {
		log.Fatal(err)
	}

	geoSelector, err := util.BuilderGetObject[*gtk.ComboBox](builder, "geo-selector")
	if err != nil {
		log.Fatal(err)
	}

	addGeoButton, err := util.BuilderGetObject[*gtk.Button](builder, "add-geo")
	if err != nil {
		log.Fatal(err)
	}

	removeGeoButton, err := util.BuilderGetObject[*gtk.Button](builder, "remove-geo")
	if err != nil {
		log.Fatal(err)
	}

	geoScroll, err := util.BuilderGetObject[*gtk.ScrolledWindow](builder, "geo-win")
	if err != nil {
		log.Fatal(err)
	}

	framePane, err := util.BuilderGetObject[*gtk.Paned](builder, "keyframe-win")
	if err != nil {
		log.Fatal(err)
	}

	frameSideBar, err := util.BuilderGetObject[*gtk.StackSidebar](builder, "frame-sidebar")
	if err != nil {
		log.Fatal(err)
	}

	frameStack, err := util.BuilderGetObject[*gtk.Stack](builder, "frame-stack")
	if err != nil {
		log.Fatal(err)
	}

	keyGeo, err := util.BuilderGetObject[*gtk.ComboBox](builder, "key-geo")
	if err != nil {
		log.Fatal(err)
	}

	keyAttr, err := util.BuilderGetObject[*gtk.ComboBox](builder, "key-attr")
	if err != nil {
		log.Fatal(err)
	}

	addFrameButton, err := util.BuilderGetObject[*gtk.Button](builder, "add-frame")
	if err != nil {
		log.Fatal(err)
	}

	addKeyframeButton, err := util.BuilderGetObject[*gtk.Button](builder, "add-key")
	if err != nil {
		log.Fatal(err)
	}

	removeKeyframeButton, err := util.BuilderGetObject[*gtk.Button](builder, "remove-key")
	if err != nil {
		log.Fatal(err)
	}

	editBox, err := util.BuilderGetObject[*gtk.Box](builder, "edit")
	if err != nil {
		log.Fatal(err)
	}

	prevBox, err := util.BuilderGetObject[*gtk.Box](builder, "preview")
	if err != nil {
		log.Fatal(err)
	}

	/* create objects */
	chromaHub, err := hub.NewDataBase(10, "", "")
	if err != nil {
		log.Fatal(err)
	}

	hubConn := library.NewConnection("Hub", conf.HubAddr, conf.HubPort)
	hubConn.Connect()

	start := time.Now()
	err = attribute.ImportAssets(hubConn.Conn)
	if err != nil {
		log.Print(err)
	}

	end := time.Now()
	elapsed := end.Sub(start)
	log.Printf("Imported Assets in %s", elapsed)

	conn = make(map[string]*library.Connection)
	for _, c := range conf.Connections {
		conn[c.Name] = library.NewConnection(c.Name, c.Address, c.Port)
	}

	editView, err := templates.NewEditor()
	if err != nil {
		log.Fatal(err)
	}
	editView.CurrentTemp = template

	preview, err := library.SetupPreviewWindow(*conf,
		func() { SendPreview(editView.CurrentTemp, library.ANIMATE_ON) },
		func() { SendPreview(editView.CurrentTemp, library.CONTINUE) },
		func() { SendPreview(editView.CurrentTemp, library.ANIMATE_OFF) },
	)
	if err != nil {
		log.Fatalf("Error setting up preview window: %s", err)
	}

	geometryToEditor := func(geoID int) {
		editView.CurrentGeoID = geoID

		err := editView.UpdateEditor()
		if err != nil {
			log.Printf("Error sending prop %d to editor: %s", geoID, err)
		}
	}

	keyTree := NewKeyframeTree(keyGeo, keyAttr)

	geoTree, err := NewGeoTree(geoSelector, geometryToEditor, keyTree.UpdateGeometryName)
	if err != nil {
		log.Fatal(err)
	}

	editView.AddAction("Save", true, func() {
		err := updateTemplateFromUI(template, geoTree, keyTree, titleEntry, tempIDEntry, layerEntry)
		if err != nil {
			log.Printf("Error creating template (%s)", err)
			return
		}

		err = chromaHub.ImportTemplate(*template)
		if err != nil {
			log.Printf("Error sending template to chroma hub (%s)", err)
		}

		editView.UpdateGeometry()
		SendPreview(editView.CurrentTemp, library.UPDATE)
		time.Sleep(10 * time.Millisecond)
		SendPreview(editView.CurrentTemp, library.ANIMATE_ON)
	})

	/* actions */
	newTemplate.Connect("activate", func() {
		template.Clean()
		updateUIFromTemplate(
			template, geoTree, keyTree, frameSideBar, framePane,
			titleEntry, tempIDEntry, layerEntry,
		)

		SendPreview(editView.CurrentTemp, library.UPDATE)
	})

	importTemplateJSON.Connect("activate", func() {
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

			*template, err = templates.NewTemplateFromFile(filename)
			if err != nil {
				log.Print(err)
				return
			}

			updateUIFromTemplate(
				template, geoTree, keyTree, frameSideBar, framePane,
				titleEntry, tempIDEntry, layerEntry,
			)
		}
	})

	importTemplateHub.Connect("activate", func() {
		dialog := NewTemplateChooserDialog(&win.Window)
		defer dialog.Destroy()

		dialog.ImportTemplates(hubConn.Conn)
		res := dialog.Run()
		if res == gtk.RESPONSE_ACCEPT {
			selection, err := dialog.treeView.GetSelection()
			if err != nil {
				log.Print("No template selected")
				return
			}

			_, iter, ok := selection.GetSelected()
			if !ok {
				log.Print("No template selected")
				return
			}

			tempID, err := util.ModelGetValue[int](dialog.treeList.ToTreeModel(), iter, 1)

			template, err := templates.GetTemplate(hubConn.Conn, tempID)
			if err != nil {
				log.Printf("Error importing template: %s", err)
				return
			}

			editView.CurrentTemp = &template
			updateUIFromTemplate(
				&template, geoTree, keyTree, frameSideBar, framePane,
				titleEntry, tempIDEntry, layerEntry,
			)
		}
	})

	exportTemplate.Connect("activate", func() {
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

			err := updateTemplateFromUI(template, geoTree, keyTree, titleEntry, tempIDEntry, layerEntry)
			if err != nil {
				log.Printf("Error exporting template (%s)", err)
			}

			err = template.ExportTemplate(filename)
			if err != nil {
				log.Printf("Error exporting template (%s)", err)
			}
		}
	})

	titleEntry.Connect("changed", func(entry *gtk.Entry) {
		text, err := entry.GetText()
		if err != nil {
			log.Print(err)
			return
		}

		template.Title = text
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

		template.TempID = int64(id)
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

		template.Layer = id
	})

	addGeoButton.Connect("clicked", func() { addGeo(template, geoTree, keyTree) })
	removeGeoButton.Connect("clicked", func() { removeGeo(template, geoTree, keyTree) })
	geoScroll.Add(geoTree.geoView)

	keyGeo.Connect("changed", func() {
		geoID, _, err := keyTree.SelectedGeometry()
		if err != nil {
			log.Printf("Error getting selected geometry: %s", err)
			return
		}

		geo := template.Geos[geoID]
		if geo == nil {
			log.Printf("Missing geometry %d", geoID)
			return
		}

		keyTree.UpdateAttrList(geo.GeoType)
	})

	frameSideBar.SetStack(frameStack)

	addFrameButton.Connect("clicked", func() { addFrame(frameSideBar, keyTree) })
	addKeyframeButton.Connect("clicked", func() { addKeyframe(frameSideBar, keyTree) })
	removeKeyframeButton.Connect("clicked", func() { removeKeyframe(frameSideBar, keyTree) })

	editBox.PackStart(editView.Box, true, true, 0)

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

	for _, render := range conn {
		eng := library.NewEngineWidget(render)
		lowerBox.PackStart(eng.Button)
	}

	win.ShowAll()
}

func updateUIFromTemplate(temp *templates.Template, geoTree *GeoTree, keyTree *KeyTree, sidebar *gtk.StackSidebar, framePane *gtk.Paned,
	titleEntry, tempIDEntry, layerEntry *gtk.Entry) (page *pages.Page) {
	titleEntry.SetText(temp.Title)
	tempIDEntry.SetText(strconv.FormatInt(temp.TempID, 10))
	layerEntry.SetText(strconv.Itoa(temp.Layer))

	geoTree.ImportGeometry(temp)
	keyTree.ImportKeyframes(temp)

	stack := sidebar.GetStack()
	for frameNum := 1; frameNum < keyTree.nextFrame; frameNum++ {
		name := fmt.Sprintf("   Frame %d   ", frameNum)
		treeView := keyTree.keyframeView[frameNum]
		if treeView == nil {
			log.Printf("Missing keyframe %d view", frameNum)
			continue
		}

		stack.AddTitled(treeView, strconv.Itoa(frameNum), name)
	}

	return page
}

func updateTemplateFromUI(temp *templates.Template, geoTree *GeoTree, keyTree *KeyTree,
	titleEntry, tempIDEntry, layerEntry *gtk.Entry) (err error) {
	temp.Title, err = titleEntry.GetText()
	if err != nil {
		return
	}

	layer, err := layerEntry.GetText()
	if err != nil {
		return
	}

	temp.Layer, _ = strconv.Atoi(layer)

	tempID, err := tempIDEntry.GetText()
	if err != nil {
		return
	}

	temp.TempID, _ = strconv.ParseInt(tempID, 10, 64)

	geoTree.ExportGeometry(temp)

	err = keyTree.ExportKeyframes(temp)
	if err != nil {
		return
	}

	return
}

func addGeo(temp *templates.Template, geoTree *GeoTree, keyTree *KeyTree) {
	geoName, err := geoTree.GetSelectedGeoName()
	if err != nil {
		log.Printf("Error adding geometry: %s", err)
		return
	}

	geoType, err := geoTree.GetSelectedGeoType()
	if err != nil {
		log.Printf("Error adding geometry: %s", err)
		return
	}

	propNum, err := temp.AddGeometry(geoType, geoName)
	if err != nil {
		log.Printf("Error adding geometry: %s", err)
		return
	}

	geoTree.AddGeoRow(propNum, 0, geoName, geoName)
	keyTree.AddGeometry(geoName, propNum)
}

func removeGeo(temp *templates.Template, geoTree *GeoTree, keyTree *KeyTree) {
	geoID, err := geoTree.GetSelectedGeoID()
	if err != nil {
		log.Printf("Error removing geo: %s", err)
		return
	}

	temp.RemoveGeometry(geoID)

	geoTree.RemoveGeo(geoID)
	keyTree.RemoveGeo(geoID)
}

func addFrame(sidebar *gtk.StackSidebar, keyTree *KeyTree) (err error) {
	stack := sidebar.GetStack()
	frameNum, err := keyTree.AddFrame()
	if err != nil {
		return
	}

	name := fmt.Sprintf("   Frame %d   ", frameNum)
	treeView := keyTree.keyframeView[frameNum]

	stack.AddTitled(treeView, strconv.Itoa(frameNum), name)
	return
}

func addKeyframe(sidebar *gtk.StackSidebar, keyTree *KeyTree) (err error) {
	stack := sidebar.GetStack()

	geoID, geoName, err := keyTree.SelectedGeometry()
	if err != nil {
		err = fmt.Errorf("Error getting geo id: %s", err)
		return
	}

	attrType, attrName, err := keyTree.SelectedAttribute()
	if err != nil {
		err = fmt.Errorf("Error getting attribute: %s", err)
		return
	}

	frameString := stack.GetVisibleChildName()
	frameNum, err := strconv.Atoi(frameString)
	if err != nil {
		err = fmt.Errorf("Error getting frame num: %s", err)
		return
	}

	keyTree.AddKeyframe(frameNum, geoID, geoName, attrType, attrName)
	return
}

func removeKeyframe(sidebar *gtk.StackSidebar, keyTree *KeyTree) (err error) {
	stack := sidebar.GetStack()
	frameString := stack.GetVisibleChildName()
	frameNum, err := strconv.Atoi(frameString)
	if err != nil {
		err = fmt.Errorf("Error getting frame num: %s", err)
		return
	}

	model := keyTree.keyframeModel[frameNum]
	if model == nil {
		err = fmt.Errorf("Error getting selected keyframe model")
		return
	}

	view := keyTree.keyframeView[frameNum]
	if view == nil {
		err = fmt.Errorf("Error getting selected keyframe model")
		return
	}

	selection, err := view.GetSelection()
	if err != nil {
		return
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		return
	}

	model.Remove(iter)
	return
}

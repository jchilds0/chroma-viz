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

func SendPreview(page *templates.Template, action int) {
	if page == nil {
		log.Println("SendPreview recieved nil page")
		return
	}

	for _, c := range conn {
		c.SendPage(action, page)
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
	if err := builder.AddFromFile("cmd/chroma-artist/menu.ui"); err != nil {
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

	exportTemplate := glib.SimpleActionNew("export_template", nil)
	app.AddAction(exportTemplate)

	importTemplateHub := glib.SimpleActionNew("import_template_hub", nil)
	app.AddAction(importTemplateHub)

	importAsset := glib.SimpleActionNew("import_asset", nil)
	app.AddAction(importAsset)

	exportAsset := glib.SimpleActionNew("export_asset", nil)
	app.AddAction(exportAsset)

	fetchAssets := glib.SimpleActionNew("fetch_assets", nil)
	app.AddAction(fetchAssets)

	generateHub := glib.SimpleActionNew("gen", nil)
	app.AddAction(generateHub)

	cleanHub := glib.SimpleActionNew("clean", nil)
	app.AddAction(cleanHub)

	app.SetMenubar(menu.(*glib.MenuModel))

	/* Body layout */
	builder, err = gtk.BuilderNew()
	if err := builder.AddFromFile("cmd/chroma-artist/gui.ui"); err != nil {
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

	duplicateGeoButton, err := util.BuilderGetObject[*gtk.Button](builder, "duplicate-geo")
	if err != nil {
		log.Fatal(err)
	}

	geoScroll, err := util.BuilderGetObject[*gtk.ScrolledWindow](builder, "geo-win")
	if err != nil {
		log.Fatal(err)
	}

	keyframes, err := util.BuilderGetObject[*gtk.Box](builder, "keyframes")
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
	conn = make(map[string]*library.Connection)
	for _, c := range conf.Connections {
		conn[c.Name] = library.NewConnection(c.Name, c.Address, c.Port)
	}

	geoModel, err := gtk.ListStoreNew(
		glib.TYPE_STRING,  // GEO TYPE
		glib.TYPE_STRING,  // GEO NAME
		glib.TYPE_INT,     // GEO NUM
		glib.TYPE_BOOLEAN, // GEO VISIBLE
	)
	if err != nil {
		log.Fatal(err)
	}

	frameModel, err := gtk.ListStoreNew(
		glib.TYPE_INT,    // Frame Num
		glib.TYPE_STRING, // Frame Text
	)
	if err != nil {
		log.Fatal(err)
	}

	editView, err := templates.NewEditor(frameModel, geoModel)
	if err != nil {
		log.Fatal(err)
	}

	preview, err := library.SetupPreviewWindow(*conf,
		func() { SendPreview(template, library.ANIMATE_ON) },
		func() { SendPreview(template, library.CONTINUE) },
		func() { SendPreview(template, library.ANIMATE_OFF) },
	)
	if err != nil {
		log.Fatalf("Error setting up preview window: %s", err)
	}

	geometryToEditor := func(geoID int) {
		editView.CurrentGeoID = geoID

		err := editView.UpdateEditor(template)
		if err != nil {
			log.Printf("Error sending prop %d to editor: %s", geoID, err)
		}
	}

	keyTree, err := NewFrames(editView, geoModel, frameModel)
	if err != nil {
		log.Fatal(err)
	}

	keyframes.PackStart(keyTree.actions, false, false, 10)
	keyframes.PackStart(keyTree.keyframes, true, true, 0)

	geoTree, err := NewGeoTree(geoSelector, geoModel, geometryToEditor, keyTree.UpdateGeometryName)
	if err != nil {
		log.Fatal(err)
	}

	editView.AddAction("Save", true, func() {
		err := updateTemplateFromUI(template, geoTree, keyTree, titleEntry, tempIDEntry, layerEntry)
		if err != nil {
			log.Printf("Error creating template (%s)", err)
			return
		}

		if template.TempID == 0 {
			log.Print("Temp ID is 0, not exporting template to chroma hub")
			return
		}

		editView.UpdateGeometry(template)
		frame, ok := keyTree.keyFrames[editView.CurrentFrameID]
		if ok && frame != nil {
			frame.UpdateBindFrame(editView.BindFrameEdit, editView.CurrentKeyID)
			frame.UpdateSetFrame(editView.SetFrameEdit, editView.CurrentKeyID)
		}

		err = conf.ChromaHub.PutJSON("/template", template)
		if err != nil {
			log.Printf("Error sending template to chroma hub (%s)", err)
		}

		SendPreview(template, library.UPDATE)
		time.Sleep(10 * time.Millisecond)
		SendPreview(template, library.ANIMATE_ON)
	})

	/* actions */
	newTemplate.Connect("activate", func() {
		template.Clean()

		updateUIFromTemplate(
			template, geoTree, keyTree, titleEntry, tempIDEntry, layerEntry,
		)

		editView.Clear()

		SendPreview(template, library.UPDATE)
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

			template.Clean()
			*template, err = templates.NewTemplateFromFile(filename, false)
			if err != nil {
				log.Print(err)
				return
			}

			updateUIFromTemplate(
				template, geoTree, keyTree, titleEntry, tempIDEntry, layerEntry,
			)

			editView.Clear()
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

	importTemplateHub.Connect("activate", func() {
		dialog, err := NewTemplateChooserDialog(&win.Window)
		if err != nil {
			log.Printf("Error creating dialog: %s", err)
			return
		}
		defer dialog.Destroy()

		err = dialog.ImportTemplates(conf.ChromaHub)
		if err != nil {
			log.Printf("Error fetching template IDs: %s", err)
		}

		res := dialog.Run()
		if res == gtk.RESPONSE_ACCEPT {
			tempID, err := dialog.SelectedTemplateID()
			if err != nil {
				log.Println("Error importing template", err)
				return
			}

			path := fmt.Sprintf("/template/%d", tempID)

			template.Clean()
			err = conf.ChromaHub.GetJSON(path, template)
			if err != nil {
				log.Printf("Error importing template: %s", err)
				return
			}

			err = template.Init(false)
			if err != nil {
				log.Printf("Error importing template: %s", err)
				return
			}

			updateUIFromTemplate(
				template, geoTree, keyTree, titleEntry, tempIDEntry, layerEntry,
			)

			editView.Clear()
		}

		fetchAssetsHub()
	})

	importAsset.Connect("activate", func() {
		dialog, err := gtk.FileChooserDialogNewWith2Buttons(
			"Import Assets", win, gtk.FILE_CHOOSER_ACTION_OPEN,
			"_Cancel", gtk.RESPONSE_CANCEL, "_Open", gtk.RESPONSE_ACCEPT)
		if err != nil {
			return
		}
		defer dialog.Destroy()

		res := dialog.Run()
		if res == gtk.RESPONSE_ACCEPT {
			filename := dialog.GetFilename()

			assets, err := hub.AssetsFromFile(filename)
			if err != nil {
				log.Printf("Error importing assets: %s", err)
				return
			}

			err = conf.ChromaHub.PutJSON("/assets", assets)
			if err != nil {
				log.Printf("Error importing assets: %s", err)
				return
			}
		}
	})

	exportAsset.Connect("activate", func() {
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

	fetchAssets.Connect("activate", func() {
		fetchAssetsHub()
	})

	generateHub.Connect("activate", func() {
		go func() {
			err = conf.ChromaHub.Generate()
			if err != nil {
				log.Printf("Error generating hub: %s", err)
			}

			log.Print("Generated Hub")
		}()
	})

	cleanHub.Connect("activate", func() {
		err := conf.ChromaHub.Clean()
		if err != nil {
			log.Println("Error cleaning chroma hub:", err)
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

	addGeoButton.Connect("clicked", func() { addGeo(template, geoTree) })
	removeGeoButton.Connect("clicked", func() { removeGeo(template, geoTree, keyTree) })
	duplicateGeoButton.Connect("clicked", func() { duplicateGeo(template, geoTree, keyTree) })

	geoScroll.Add(geoTree.geoView)

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

func updateUIFromTemplate(temp *templates.Template, geoTree *GeoTree, keyTree *Frames,
	titleEntry, tempIDEntry, layerEntry *gtk.Entry) (page *pages.Page) {
	titleEntry.SetText(temp.Title)
	tempIDEntry.SetText(strconv.FormatInt(temp.TempID, 10))
	layerEntry.SetText(strconv.Itoa(temp.Layer))

	geoTree.Clear()
	geoTree.ImportGeometry(temp)

	keyTree.Clear()
	keyTree.ImportKeyframes(temp)

	return page
}

func updateTemplateFromUI(temp *templates.Template, geoTree *GeoTree, keyTree *Frames,
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

func addGeo(temp *templates.Template, geoTree *GeoTree) {
	geoName, err := geoTree.GetSelectedGeoName()
	if err != nil {
		log.Printf("Error adding geometry: %s", err)
		return
	}

	geoNum, err := temp.AddGeometry(geoName, geoName)
	if err != nil {
		log.Printf("Error adding geometry: %s", err)
		return
	}

	geoTree.AddGeoRow(geoNum, 0, geoName, geoName, true)
}

func removeGeo(temp *templates.Template, geoTree *GeoTree, keyTree *Frames) {
	iter, err := geoTree.GetSelectedGeometry()
	if err != nil {
		log.Printf("Error removing geometry: %s", err)
		return
	}

	geoID, err := util.ModelGetValue[int](geoTree.geoModel.ToTreeModel(), iter, GEO_NUM)
	if err != nil {
		log.Printf("Error removing geo: %s", err)
		return
	}

	temp.RemoveGeometry(geoID)
	geoTree.RemoveGeo(geoID)
	keyTree.RemoveGeo(geoID)
}

func duplicateGeo(temp *templates.Template, geoTree *GeoTree, keyTree *Frames) {
	iter, err := geoTree.GetSelectedGeometry()
	if err != nil {
		log.Printf("Error duplicating geometry: %s", err)
		return
	}

	geoID, err := util.ModelGetValue[int](geoTree.geoModel.ToTreeModel(), iter, GEO_NUM)
	if err != nil {
		log.Printf("Error duplicating geometry: %s", err)
		return
	}

	geoType, err := util.ModelGetValue[string](geoTree.geoModel.ToTreeModel(), iter, GEO_TYPE)
	if err != nil {
		log.Printf("Error duplicating geometry: %s", err)
		return
	}

	geoName, err := util.ModelGetValue[string](geoTree.geoModel.ToTreeModel(), iter, GEO_NAME)
	if err != nil {
		log.Printf("Error duplicating geometry: %s", err)
		return
	}

	geoVisible, err := util.ModelGetValue[bool](geoTree.geoModel.ToTreeModel(), iter, GEO_VISIBLE)
	if err != nil {
		log.Printf("Error duplicating geometry: %s", err)
		return
	}

	newGeoName := geoName + " copy"
	newGeoID, err := temp.AddGeometry(geoType, newGeoName)
	if err != nil {
		log.Printf("Error adding geometry: %s", err)
		return
	}

	err = temp.CopyGeometry(geoID, newGeoID)
	if err != nil {
		log.Printf("Error adding geometry: %s", err)
		return
	}

	geoTree.AddGeoRow(newGeoID, 0, newGeoName, geoType, geoVisible)
}

func fetchAssetsHub() {
	start := time.Now()
	assets := make([]hub.Asset, 0, 10)

	err := conf.ChromaHub.GetJSON("/assets", &assets)
	if err != nil {
		log.Printf("Error importing assets: %s", err)
		return
	}

	for _, a := range assets {
		attribute.InsertAsset(a.Directory, a.Name, int(a.AssetID))
	}

	end := time.Now()
	elapsed := end.Sub(start)
	log.Printf("Imported Assets in %s", elapsed)
}

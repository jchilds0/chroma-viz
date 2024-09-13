package main

import (
	"chroma-viz/library"
	"chroma-viz/library/geometry"
	"chroma-viz/library/hub"
	"chroma-viz/library/pages"
	"chroma-viz/library/templates"
	"chroma-viz/library/util"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const iconSize = 24
const addFrameIcon = "artist/add-file-icon.svg"
const removeFrameIcon = "artist/remove-file-icon.svg"

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

	duplicateGeoButton, err := util.BuilderGetObject[*gtk.Button](builder, "duplicate-geo")
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

	removeFrameButton, err := util.BuilderGetObject[*gtk.Button](builder, "remove-frame")
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

	err = chromaHub.SelectDatabase("chroma_hub", "", "")
	if err != nil {
		log.Fatal(err)
	}

	/*
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
	*/

	conn = make(map[string]*library.Connection)
	for _, c := range conf.Connections {
		conn[c.Name] = library.NewConnection(c.Name, c.Address, c.Port)
	}

	editView, err := templates.NewEditor()
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

	keyTree := NewKeyframeTree(keyGeo, keyAttr, frameSideBar)

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

		if template.TempID == 0 {
			log.Print("Temp ID is 0, not exporting template to chroma hub")
			return
		}

		editView.UpdateGeometry(template)

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
			template, geoTree, keyTree, framePane,
			titleEntry, tempIDEntry, layerEntry,
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

			*template, err = templates.NewTemplateFromFile(filename)
			if err != nil {
				log.Print(err)
				return
			}

			updateUIFromTemplate(
				template, geoTree, keyTree, framePane,
				titleEntry, tempIDEntry, layerEntry,
			)

			editView.Clear()
		}
	})

	importTemplateHub.Connect("activate", func() {
		dialog := NewTemplateChooserDialog(&win.Window)
		defer dialog.Destroy()

		dialog.ImportTemplates(conf.ChromaHub)
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
			path := fmt.Sprintf("/template/%d", tempID)

			err = conf.ChromaHub.GetJSON(path, template)
			if err != nil {
				log.Printf("Error importing template: %s", err)
				return
			}

			err = template.Init()
			if err != nil {
				log.Printf("Error importing template: %s", err)
				return
			}

			updateUIFromTemplate(
				template, geoTree, keyTree, framePane,
				titleEntry, tempIDEntry, layerEntry,
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
	duplicateGeoButton.Connect("clicked", func() { duplicateGeo(template, geoTree, keyTree) })

	geoScroll.Add(geoTree.geoView)

	keyGeo.Connect("changed", func() {
		geoID, _, err := keyTree.SelectedGeometry()
		if err != nil {
			return
		}

		geo := template.Geos[geoID]
		if geo == nil {
			log.Printf("Missing geometry %d", geoID)
			return
		}

		geometry.UpdateAttrList(keyTree.keyAttrList, geo.GeoType)
	})

	frameSideBar.SetStack(frameStack)

	{
		buf, err := gdk.PixbufNewFromFileAtSize(addFrameIcon, iconSize, iconSize)
		if err != nil {
			log.Fatal(err)
		}

		img, err := gtk.ImageNewFromPixbuf(buf)
		if err != nil {
			log.Print(err)
		}
		addFrameButton.SetImage(img)

		addFrameButton.Connect("clicked", func() {
			err := keyTree.AddFrame()
			if err != nil {
				log.Printf("Error adding frame: %s", err)
			}
		})
	}

	{
		buf, err := gdk.PixbufNewFromFileAtSize(removeFrameIcon, iconSize, iconSize)
		if err != nil {
			log.Fatal(err)
		}

		img, err := gtk.ImageNewFromPixbuf(buf)
		if err != nil {
			log.Print(err)
		}

		removeFrameButton.SetImage(img)

		removeFrameButton.Connect("clicked", func() {
			err := keyTree.RemoveFrame()
			if err != nil {
				log.Printf("Error removing frame: %s", err)
			}
		})
	}

	{
		img, err := gtk.ImageNewFromIconName("list-add", 3)
		if err != nil {
			log.Print(err)
		}
		addKeyframeButton.SetImage(img)

		addKeyframeButton.Connect("clicked", func() {
			err := keyTree.AddKeyframe()
			if err != nil {
				log.Printf("Error adding keyframe: %s", err)
			}
		})
	}

	{
		img, err := gtk.ImageNewFromIconName("list-remove", 3)
		if err != nil {
			log.Print(err)
		}
		removeKeyframeButton.SetImage(img)

		removeKeyframeButton.Connect("clicked", func() {
			err := keyTree.RemoveKeyframe()
			if err != nil {
				log.Printf("Error removing keyframe: %s", err)
			}
		})
	}

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

func updateUIFromTemplate(temp *templates.Template, geoTree *GeoTree, keyTree *KeyTree,
	framePane *gtk.Paned, titleEntry, tempIDEntry, layerEntry *gtk.Entry) (page *pages.Page) {
	titleEntry.SetText(temp.Title)
	tempIDEntry.SetText(strconv.FormatInt(temp.TempID, 10))
	layerEntry.SetText(strconv.Itoa(temp.Layer))

	geoTree.Clear()
	geoTree.ImportGeometry(temp)

	keyTree.Clear()
	keyTree.ImportKeyframes(temp)

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

	geoNum, err := temp.AddGeometry(geoName, geoName)
	if err != nil {
		log.Printf("Error adding geometry: %s", err)
		return
	}

	geoTree.AddGeoRow(geoNum, 0, geoName, geoName)
	keyTree.AddGeometry(geoName, geoNum)
}

func removeGeo(temp *templates.Template, geoTree *GeoTree, keyTree *KeyTree) {
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
	geoTree.RemoveGeo(iter, geoID)
	keyTree.RemoveGeo(geoID)
}

func duplicateGeo(temp *templates.Template, geoTree *GeoTree, keyTree *KeyTree) {
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

	geoTree.AddGeoRow(newGeoID, 0, newGeoName, geoType)
	keyTree.AddGeometry(newGeoName, newGeoID)
}

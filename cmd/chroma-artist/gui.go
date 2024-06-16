package main

import (
	"chroma-viz/library"
	"chroma-viz/library/attribute"
	"chroma-viz/library/hub"
	"chroma-viz/library/pages"
	"chroma-viz/library/props"
	"chroma-viz/library/templates"
	"chroma-viz/library/util"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var conn map[string]*library.Connection
var conf *library.Config
var chromaHub *hub.DataBase
var hubConn *library.Connection

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
	page := pages.NewPage(0, 0, 0, 10, "")

	win, err := gtk.ApplicationWindowNew(app)
	if err != nil {
		log.Fatalf("Error starting artist gui (%s)", err)
	}

	win.SetDefaultSize(800, 600)
	win.SetTitle("Chroma Artist")

	editView, err := library.NewEditor()
	if err != nil {
		log.Fatal(err)
	}

	propToEditor := func(propID int) {
		prop := page.PropMap[propID]
		err := editView.SetProperty(prop)
		if err != nil {
			log.Printf("Error sending prop %d to editor: %s", propID, err)
		}
	}

	geoTree, err := NewGeoTree(propToEditor)
	if err != nil {
		log.Fatal(err)
	}

	keyTree := NewKeyframeTree()

	editView.AddAction("Save", true, func() {
		template, err := exportPage(page, geoTree, keyTree, titleEntry, tempIDEntry, layerEntry)
		if err != nil {
			log.Printf("Error creating template (%s)", err)
			return
		}

		err = chromaHub.ImportTemplate(*template)
		if err != nil {
			log.Printf("Error sending template to chroma hub (%s)", err)
			return
		}

		editView.CurrentPage = page
		editView.UpdateProps()

		SendPreview(editView.CurrentPage, library.UPDATE)
		time.Sleep(10 * time.Millisecond)
		SendPreview(editView.CurrentPage, library.ANIMATE_ON)
	})

	editView.PropertyEditor()
	editView.CurrentPage = page

	preview, err := library.SetupPreviewWindow(*conf,
		func() { SendPreview(editView.CurrentPage, library.ANIMATE_ON) },
		func() { SendPreview(editView.CurrentPage, library.CONTINUE) },
		func() { SendPreview(editView.CurrentPage, library.ANIMATE_OFF) },
	)
	if err != nil {
		log.Fatalf("Error setting up preview window: %s", err)
	}

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

	geoSelectorModel(geoSelector)

	addGeo, err := util.BuilderGetObject[*gtk.Button](builder, "add-geo")
	if err != nil {
		log.Fatal(err)
	}

	removeGeo, err := util.BuilderGetObject[*gtk.Button](builder, "remove-geo")
	if err != nil {
		log.Fatal(err)
	}

	geoScroll, err := util.BuilderGetObject[*gtk.ScrolledWindow](builder, "geo-win")
	if err != nil {
		log.Fatal(err)
	}

	geoScroll.Add(geoTree.geoView)

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

	keyAttr, err := util.BuilderGetObject[*gtk.ComboBoxText](builder, "key-attr")
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

	editBox.PackStart(editView.Box, true, true, 0)

	prevBox, err := util.BuilderGetObject[*gtk.Box](builder, "preview")
	if err != nil {
		log.Fatal(err)
	}

	prevBox.PackStart(preview, true, true, 0)

	/* actions */
	newTemplate.Connect("activate", func() {
		page.PropMap = make(map[int]*props.Property)

		newPage(geoTree, keyTree, frameSideBar, framePane, titleEntry, tempIDEntry, layerEntry)
		SendPreview(editView.CurrentPage, library.UPDATE)
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

			page = importPage(&temp, geoTree, keyTree, frameSideBar, framePane, titleEntry, tempIDEntry, layerEntry)
			editView.CurrentPage = page
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

			page = importPage(&template, geoTree, keyTree, frameSideBar, framePane, titleEntry, tempIDEntry, layerEntry)
			editView.CurrentPage = page
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

		dialog.SetCurrentName(page.Title + ".json")
		res := dialog.Run()
		if res == gtk.RESPONSE_ACCEPT {
			filename := dialog.GetFilename()

			template, err := exportPage(page, geoTree, keyTree, titleEntry, tempIDEntry, layerEntry)
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
		iter, err := geoSelector.GetActiveIter()
		if err != nil {
			log.Print("No geometry selected")
			return
		}

		model, err := geoSelector.GetModel()
		if err != nil {
			log.Printf("Error getting geometry model: %s", err)
			return
		}

		propType, err := util.ModelGetValue[string](model.ToTreeModel(), iter, 0)
		if err != nil {
			log.Printf("Error getting geometry model: %s", err)
			return
		}

		propNum, err := AddProp(page, propType)
		if err != nil {
			log.Print(err)
			return
		}

		iter = geoTree.geoModel.Append(nil)
		prop := page.PropMap[propNum]
		geoTree.AddGeoRow(iter, prop.Name, prop.Name, propNum)
	})

	removeGeo.Connect("clicked", func() {
		selection, err := geoTree.geoView.GetSelection()
		if err != nil {
			log.Printf("Error getting selected")
			return
		}

		_, iter, ok := selection.GetSelected()
		if !ok {
			log.Printf("No geometry selected")
			return
		}

		model := geoTree.geoModel.ToTreeModel()
		geoID, err := util.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Printf("Error getting prop id (%s)", err)
			return
		}

		RemoveProp(page, geoID)
		geoTree.geoModel.Remove(iter)
		geoTree.removeGeometry(geoID)
	})

	geoCell, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal(err)
	}

	keyGeo.PackStart(geoCell, true)
	keyGeo.CellLayout.AddAttribute(geoCell, "text", GEO_NAME)
	keyGeo.SetActive(GEO_NAME)
	keyGeo.SetModel(geoTree.geoList)

	keyGeo.Connect("changed", func() {
		iter, err := keyGeo.GetActiveIter()
		if err != nil {
			log.Printf("No geometry selected")
			return
		}

		geoID, err := util.ModelGetValue[int](geoTree.geoList.ToTreeModel(), iter, GEO_NUM)
		if err != nil {
			log.Printf("Error getting geo id (%s)", err)
			return
		}

		keyAttr.RemoveAll()

		for name := range page.PropMap[geoID].Attr {
			keyAttr.AppendText(name)
		}
	})

	frameSideBar.SetStack(frameStack)

	addFrameButton.Connect("clicked", func() { addFrame(frameSideBar, keyTree) })
	addKeyframeButton.Connect("clicked", func() { addKeyframe(frameSideBar, keyGeo, keyAttr, geoTree, keyTree) })
	removeKeyframeButton.Connect("clicked", func() { removeKeyframe(frameSideBar, keyTree) })

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

var geoNames = map[string]string{
	props.RECT_PROP:   "Rectangle",
	props.CIRCLE_PROP: "Circle",
	props.TEXT_PROP:   "Text",
	props.TICKER_PROP: "Ticker",
	props.CLOCK_PROP:  "Clock",
	props.IMAGE_PROP:  "Image",
}

func geoSelectorModel(combo *gtk.ComboBox) (err error) {
	model, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		return
	}

	for propName, name := range geoNames {
		iter := model.Append()

		model.SetValue(iter, 0, propName)
		model.SetValue(iter, 1, name)
	}

	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	combo.SetModel(model)
	combo.CellLayout.PackStart(cell, true)
	combo.CellLayout.AddAttribute(cell, "text", 1)
	combo.SetActive(1)
	return
}

func AddProp(page *pages.Page, propType string) (id int, err error) {
	ok := true
	for id = 1; ok; id++ {
		prop, ok := page.PropMap[id]
		if !ok {
			break
		}

		if prop == nil {
			break
		}
	}

	page.PropMap[id] = props.NewProperty(propType, geoNames[propType], true, nil)

	if propType != props.CLOCK_PROP {
		return
	}

	/*
	   Clock requires a way to send updates to viz
	   to animate the clock. We manually add this
	   after parsing the page.
	*/
	attr, ok := page.PropMap[id].Attr["string"]
	if !ok {
		err = fmt.Errorf("Clock Prop missing string attr")
		return
	}

	clockAttr, ok := attr.(*attribute.ClockAttribute)
	if !ok {
		err = fmt.Errorf("String attr is not a clock attribute")
		return
	}

	clockAttr.SetClock(func() { SendPreview(page, library.CONTINUE) })
	return
}

func RemoveProp(page *pages.Page, propID int) {
	prop := page.PropMap[propID]
	if prop == nil {
		log.Printf("No prop with prop id %d", propID)
		return
	}

	page.PropMap[propID] = nil
}

func newPage(geoTree *GeoTree, keyTree *KeyTree, frameSideBar *gtk.StackSidebar, framePane *gtk.Paned,
	titleEntry, tempIDEntry, layerEntry *gtk.Entry) {
	geoTree.Clear()
	keyTree.Clear()

	framePane.Remove(frameSideBar.GetStack())

	stack, _ := gtk.StackNew()
	stack.SetVExpand(true)
	stack.SetVisible(true)

	frameSideBar.SetStack(stack)
	framePane.Add2(stack)

	titleEntry.SetText("")
	tempIDEntry.SetText("")
	layerEntry.SetText("")

}

func importPage(temp *templates.Template, geoTree *GeoTree, keyTree *KeyTree, sidebar *gtk.StackSidebar, framePane *gtk.Paned,
	titleEntry, tempIDEntry, layerEntry *gtk.Entry) (page *pages.Page) {
	newPage(geoTree, keyTree, sidebar, framePane, titleEntry, tempIDEntry, layerEntry)
	page = pages.NewPageFromTemplate(temp)

	titleEntry.SetText(temp.Title)
	tempIDEntry.SetText(strconv.FormatInt(temp.TempID, 10))
	layerEntry.SetText(strconv.Itoa(temp.Layer))

	geoTree.ImportGeometry(page)
	keyTree.ImportKeyframes(temp)

	// set temp switch to true to send all props to chroma engine
	for geoID, geo := range page.PropMap {
		geo.SetTemp(true)
		keyTree.UpdateKeys(geoID, geo.Name)

		if geo.PropType != props.CLOCK_PROP {
			continue
		}

		/*
		   Clock requires a way to send updates to viz
		   to animate the clock. We manually add this
		   after parsing the page.
		*/
		attr, ok := geo.Attr["string"]
		if !ok {
			log.Printf("Clock Prop missing string attr")
			continue
		}

		clockAttr, ok := attr.(*attribute.ClockAttribute)
		if !ok {
			log.Printf("String attr is not a clock attribute")
			continue
		}

		clockAttr.SetClock(func() { SendPreview(page, library.CONTINUE) })
	}

	stack := sidebar.GetStack()

	for frameNum := range keyTree.keyframeView {
		name := fmt.Sprintf("Frame %d", frameNum)
		treeView := keyTree.keyframeView[frameNum]
		if treeView == nil {
			log.Printf("Missing keyframe %d view", frameNum)
			continue
		}

		stack.AddTitled(treeView, strconv.Itoa(frameNum), name)
	}

	geoTree.currentPage = page
	return page
}

func exportPage(page *pages.Page, geoTree *GeoTree, keyTree *KeyTree, titleEntry, tempIDEntry, layerEntry *gtk.Entry) (temp *templates.Template, err error) {
	// update geometry parents
	model := geoTree.geoModel.ToTreeModel()
	if iter, ok := model.GetIterFirst(); ok {
		updateParentGeometry(page, model, iter, 0)
	}

	temp = page.CreateTemplate()

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

	err = keyTree.ExportKeyframes(temp)
	if err != nil {
		return
	}

	return
}

func addFrame(sidebar *gtk.StackSidebar, keyTree *KeyTree) (err error) {
	stack := sidebar.GetStack()
	frameNum, err := keyTree.AddFrame()
	if err != nil {
		return
	}

	name := fmt.Sprintf("Frame %d", frameNum)
	treeView := keyTree.keyframeView[frameNum]

	stack.AddTitled(treeView, strconv.Itoa(frameNum), name)
	return
}

func addKeyframe(sidebar *gtk.StackSidebar, keyGeo *gtk.ComboBox, keyAttr *gtk.ComboBoxText,
	geoTree *GeoTree, keyTree *KeyTree) (err error) {
	stack := sidebar.GetStack()
	iter, err := keyGeo.GetActiveIter()
	if err != nil {
		err = fmt.Errorf("No geometry selected")
		return
	}

	model := geoTree.geoList.ToTreeModel()

	geo, err := util.ModelGetValue[string](model, iter, GEO_NAME)
	if err != nil {
		err = fmt.Errorf("Error getting geo name: %s", err)
		return
	}

	geoID, err := util.ModelGetValue[int](model, iter, GEO_NUM)
	if err != nil {
		err = fmt.Errorf("Error getting geo id: %s", err)
		return
	}

	attr := keyAttr.GetActiveText()

	frameString := stack.GetVisibleChildName()
	frameNum, err := strconv.Atoi(frameString)
	if err != nil {
		err = fmt.Errorf("Error getting frame num: %s", err)
		return
	}

	keyTree.AddKeyframe(frameNum, geoID, geo, attr)
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

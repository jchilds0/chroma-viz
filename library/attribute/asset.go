package attribute

import (
	"chroma-viz/library/util"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	IMAGE_ID = iota
	NAME
	NUM_COL
)

var cols = []glib.Type{glib.TYPE_INT, glib.TYPE_STRING}

var defaultNodeSize = 10
var rootNode = newAssetNode(defaultNodeSize)

type assetNode struct {
	assetNames map[int]string
	childNodes map[string]*assetNode
}

func newAssetNode(numAssets int) *assetNode {
	node := &assetNode{}
	node.assetNames = make(map[int]string, numAssets)
	node.childNodes = make(map[string]*assetNode)

	return node
}

func InsertAsset(path, name string, id int) {
	dirs := strings.Split(path, "/")

	currentNode := rootNode
	for _, dir := range dirs {
		nextNode, ok := currentNode.childNodes[dir]
		if !ok {
			nextNode = newAssetNode(defaultNodeSize)
			currentNode.childNodes[dir] = nextNode
		}

		currentNode = nextNode
	}

	currentNode.assetNames[id] = name
}

type AssetAttribute struct {
	Name  string
	Value int
	dir   *gtk.TreePath
	asset *gtk.TreePath
}

func NewAssetAttribute(name string) *AssetAttribute {
	asset := &AssetAttribute{
		Name: name,
	}

	return asset
}

func (assetAttr *AssetAttribute) Directory() string {
	return ""
}

func (asset *AssetAttribute) UpdateAttribute(assetEdit *AssetEditor) (err error) {
	selection, err := assetEdit.dirs.GetSelection()
	if err == nil {
		_, iter, ok := selection.GetSelected()

		if ok {
			asset.dir, _ = assetEdit.dirsStore.GetPath(iter)
		}
	}

	selection, err = assetEdit.assets.GetSelection()
	if err != nil {
		return
	}
	_, selected, ok := selection.GetSelected()
	if !ok {
		return
	}

	asset.Value, err = util.ModelGetValue[int](assetEdit.assetsStore.ToTreeModel(), selected, IMAGE_ID)
	return
}

type AssetEditor struct {
	Name        string
	Box         *gtk.Box
	dirs        *gtk.TreeView
	dirsStore   *gtk.TreeStore
	assets      *gtk.TreeView
	assetsStore *gtk.ListStore
}

func NewAssetEditor(name string) *AssetEditor {
	assetEdit := &AssetEditor{Name: name}

	assetEdit.Box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	assetEdit.Box.SetVisible(true)
	assetEdit.Box.SetVExpand(true)

	paned, _ := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	assetEdit.Box.PackStart(paned, true, true, 0)
	paned.SetPosition(200)
	paned.SetVisible(true)

	dirsScroll, _ := gtk.ScrolledWindowNew(nil, nil)
	paned.Pack1(dirsScroll, true, true)
	dirsScroll.SetVisible(true)

	assetEdit.dirs, _ = gtk.TreeViewNew()
	dirsScroll.Add(assetEdit.dirs)
	assetEdit.dirs.SetVisible(true)

	cell, _ := gtk.CellRendererTextNew()
	col, _ := gtk.TreeViewColumnNewWithAttribute("File", cell, "text", 0)
	assetEdit.dirs.AppendColumn(col)

	assetEdit.dirs.Connect("row-activated",
		func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
			iter, _ := assetEdit.dirsStore.GetIter(path)
			assets := assetEdit.GetAssets(iter)

			assetEdit.assetsStore.Clear()

			for id, name := range assets.assetNames {
				row := assetEdit.assetsStore.Append()
				assetEdit.assetsStore.SetValue(row, NAME, name)
				assetEdit.assetsStore.SetValue(row, IMAGE_ID, id)
			}
		})

	assetsScroll, _ := gtk.ScrolledWindowNew(nil, nil)
	paned.Pack2(assetsScroll, true, true)
	assetsScroll.SetVisible(true)

	assetEdit.assets, _ = gtk.TreeViewNew()
	assetsScroll.Add(assetEdit.assets)
	assetEdit.assets.SetVisible(true)

	assetEdit.assetsStore, _ = gtk.ListStoreNew(cols...)
	assetEdit.assets.SetModel(assetEdit.assetsStore)

	cell, _ = gtk.CellRendererTextNew()
	col, _ = gtk.TreeViewColumnNewWithAttribute("Image ID", cell, "text", IMAGE_ID)
	assetEdit.assets.AppendColumn(col)

	cell, _ = gtk.CellRendererTextNew()
	col, _ = gtk.TreeViewColumnNewWithAttribute("Name", cell, "text", NAME)
	assetEdit.assets.AppendColumn(col)

	assetEdit.RefreshDirs()

	return assetEdit
}

func (asset *AssetEditor) RefreshDirs() {
	model, _ := gtk.TreeStoreNew(glib.TYPE_STRING)
	addChildren(model, nil, rootNode)

	asset.dirsStore = model
	asset.dirs.SetModel(model)
}

func addChildren(model *gtk.TreeStore, iter *gtk.TreeIter, node *assetNode) {
	for name, child := range node.childNodes {
		newDir := model.Append(iter)
		model.SetValue(newDir, 0, name)

		addChildren(model, newDir, child)
	}
}

func (asset *AssetEditor) GetAssets(iter *gtk.TreeIter) *assetNode {
	var parent gtk.TreeIter
	var parentNode *assetNode

	ok := asset.dirsStore.IterParent(&parent, iter)
	if ok {
		parentNode = asset.GetAssets(&parent)
	} else {
		parentNode = rootNode
	}

	name, _ := util.ModelGetValue[string](asset.dirsStore.ToTreeModel(), iter, 0)
	return parentNode.childNodes[name]
}

func (asset *AssetEditor) UpdateEditor(assetAttr *AssetAttribute) error {
	asset.RefreshDirs()

	dirSelection, err := asset.dirs.GetSelection()
	if err != nil {
		return nil
	}

	if assetAttr.dir == nil {
		return nil
	}

	dirSelection.SelectPath(assetAttr.dir)
	_, iter, _ := dirSelection.GetSelected()
	assets := asset.GetAssets(iter)
	asset.assetsStore.Clear()

	for id, name := range assets.assetNames {
		row := asset.assetsStore.Append()
		asset.assetsStore.SetValue(row, NAME, name)
		asset.assetsStore.SetValue(row, IMAGE_ID, id)
	}

	return nil
}

package attribute

import (
	"bufio"
	"chroma-viz/library/gtk_utils"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
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
	assetIDs   []int
	assetNames []string
	childNodes map[string]*assetNode
}

func newAssetNode(numAssets int) *assetNode {
	node := &assetNode{}
	node.assetIDs = make([]int, 0, numAssets)
	node.assetNames = make([]string, 0, numAssets)
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

	currentNode.assetIDs = append(currentNode.assetIDs, id)
	currentNode.assetNames = append(currentNode.assetNames, name)
}

func ImportAssets(hub net.Conn) (err error) {
	hub.Write([]byte("ver 0 1 assets;"))
	buf := bufio.NewReader(hub)

	assetBytes, err := buf.ReadBytes(0)
	if err != nil {
		return
	}

	var asset struct {
		Dirs  map[int]string
		Names map[int]string
	}

	bytes := assetBytes[:len(assetBytes)-1]
	err = json.Unmarshal(bytes, &asset)
	if err != nil {
		return
	}

	for id, path := range asset.Dirs {
		InsertAsset(path, asset.Names[id], id)
	}

	return
}

type AssetAttribute struct {
	Name  string
	Type  int
	Value int
	dir   *gtk.TreePath
	asset *gtk.TreePath
}

func NewAssetAttribute(name string) *AssetAttribute {
	asset := &AssetAttribute{
		Name: name,
		Type: ASSET,
	}

	return asset
}

func (asset *AssetAttribute) String() string {
	return fmt.Sprintf("%s=%d#", asset.Name, asset.Value)
}

func (asset *AssetAttribute) Update(edit Editor) (err error) {
	assetEdit, ok := edit.(*AssetEditor)
	if !ok {
		return fmt.Errorf("AssetAttribute.Update requires AssetEditor")
	}

	selection, err := assetEdit.dirs.GetSelection()
	if err == nil {
		_, iter, _ := selection.GetSelected()
		asset.dir, _ = assetEdit.dirsStore.GetPath(iter)
	}

	selection, err = assetEdit.assets.GetSelection()
	_, selected, _ := selection.GetSelected()

	asset.Value, err = gtk_utils.ModelGetValue[int](assetEdit.assetsStore.ToTreeModel(), selected, IMAGE_ID)
	return err
}

func (asset *AssetAttribute) Copy(attr Attribute) {
	assetAttrCopy, ok := attr.(*AssetAttribute)
	if !ok {
		log.Printf("Attribute not an AssetAttribute")
		return
	}

	asset.Value = assetAttrCopy.Value
}

func (asset *AssetAttribute) Encode() string {
	return fmt.Sprintf("{'name': '%s', 'value': '%d'}",
		asset.Name, asset.Value)
}

func (asset *AssetAttribute) Decode(s string) {
	var err error

	asset.Value, err = strconv.Atoi(s)
	if err != nil {
		log.Printf("Error decoding int attr (%s)", err)
	}
}

type AssetEditor struct {
	name        string
	box         *gtk.Box
	dirs        *gtk.TreeView
	dirsStore   *gtk.TreeStore
	assets      *gtk.TreeView
	assetsStore *gtk.ListStore
}

func NewAssetEditor(name string) *AssetEditor {
	assetEdit := &AssetEditor{name: name}

	assetEdit.box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	assetEdit.box.SetVisible(true)
	assetEdit.box.SetVExpand(true)

	paned, _ := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	assetEdit.box.PackStart(paned, true, true, 0)
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

			for i, name := range assets.assetNames {
				row := assetEdit.assetsStore.Append()
				assetEdit.assetsStore.SetValue(row, NAME, name)
				assetEdit.assetsStore.SetValue(row, IMAGE_ID, assets.assetIDs[i])
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

	name, _ := gtk_utils.ModelGetValue[string](asset.dirsStore.ToTreeModel(), iter, 0)
	return parentNode.childNodes[name]
}

func (asset *AssetEditor) Box() *gtk.Box {
	return asset.box
}

func (asset *AssetEditor) Expand() bool {
	return true
}

func (asset *AssetEditor) Update(attr Attribute) error {
	assetAttr, ok := attr.(*AssetAttribute)
	if !ok {
		return fmt.Errorf("AssetEditor.Update requires AssetAttribute")
	}

	dirSelection, err := asset.dirs.GetSelection()
	if err == nil && assetAttr.dir != nil {
		dirSelection.SelectPath(assetAttr.dir)
		_, iter, _ := dirSelection.GetSelected()
		assets := asset.GetAssets(iter)
		asset.assetsStore.Clear()

		for i, name := range assets.assetNames {
			row := asset.assetsStore.Append()
			asset.assetsStore.SetValue(row, NAME, name)
			asset.assetsStore.SetValue(row, IMAGE_ID, assets.assetIDs[i])
		}
	}

	return nil
}

func (asset *AssetEditor) Name() string {
	return asset.name
}

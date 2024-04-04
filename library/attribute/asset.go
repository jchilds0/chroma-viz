package attribute

import (
	"chroma-viz/library/gtk_utils"
	"fmt"
	"log"
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

var cols = []glib.Type{glib.TYPE_STRING, glib.TYPE_INT}

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

type AssetAttribute struct {
	Name  string
	Type  int
	Value int
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

func (asset *AssetAttribute) Update(edit Editor) error {
	assetEdit, ok := edit.(*AssetEditor)
	if !ok {
		return fmt.Errorf("AssetAttribute.Update requires AssetEditor")
	}

	selection, err := assetEdit.assets.GetSelection()
	if err != nil {
		return err
	}

	model, selected, ok := selection.GetSelected()
	if !ok {
		return fmt.Errorf("No asset selected")
	}

	asset.Value, err = gtk_utils.ModelGetValue[int](model.ToTreeModel(), selected, IMAGE_ID)
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
	name   string
	box    *gtk.Box
	dirs   *gtk.TreeView
	assets *gtk.TreeView
}

func NewAssetEditor(name string) *AssetEditor {
	assetEdit := &AssetEditor{name: name}

	assetEdit.box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	assetEdit.box.SetVisible(true)
	assetEdit.box.SetHExpand(true)

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

	assetsScroll, _ := gtk.ScrolledWindowNew(nil, nil)
	paned.Pack2(assetsScroll, true, true)
	assetsScroll.SetVisible(true)

	assetEdit.assets, _ = gtk.TreeViewNew()
	assetsScroll.Add(assetEdit.assets)
	assetEdit.assets.SetVisible(true)

	assetEdit.RefreshDirs()

	return assetEdit
}

func (asset *AssetEditor) RefreshDirs() {
	model, _ := gtk.TreeStoreNew(glib.TYPE_STRING)
	addChildren(model, nil, rootNode)

	asset.dirs.SetModel(model)
}

func addChildren(model *gtk.TreeStore, iter *gtk.TreeIter, node *assetNode) {
	for name, child := range node.childNodes {
		newDir := model.Append(iter)
		model.SetValue(newDir, 0, name)

		addChildren(model, newDir, child)
	}
}

func (asset *AssetEditor) Box() *gtk.Box {
	return asset.box
}

func (asset *AssetEditor) Update(attr Attribute) error {
	_, ok := attr.(*AssetAttribute)
	if !ok {
		return fmt.Errorf("AssetEditor.Update requires AssetAttribute")
	}

	return nil
}

func (asset *AssetEditor) Name() string {
	return asset.name
}

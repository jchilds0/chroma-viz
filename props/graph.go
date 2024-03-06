package props

import (
	"chroma-viz/attribute"

	"github.com/gotk3/gotk3/gtk"
)

type GraphEditor struct {
    box         *gtk.Box
    treeView    *gtk.TreeView
    listStore   *gtk.ListStore
    edit        map[string]attribute.Editor
}

func NewGraphEditor(width, height int, animate func()) (g *GraphEditor, err error) {
    g = &GraphEditor{}
    g.edit = make(map[string]attribute.Editor, 5)

    g.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        return
    }

    g.edit["x"], err = attribute.NewIntEditor("x", 0, float64(width), animate)
    if err != nil {
        return
    }

    g.edit["y"], err = attribute.NewIntEditor("y", 0, float64(height), animate)
    if err != nil {
        return
    }

    columns := []string{"x Pos", "y Pos"}
    g.edit["node"], err = attribute.NewListEditor("Graph", columns, animate)
    if err != nil {
        return
    }

    g.box.SetVisible(true)

    g.box.PackStart(g.edit["x"].Box(), false, false, padding)
    g.box.PackStart(g.edit["y"].Box(), false, false, padding)
    g.box.PackStart(g.edit["node"].Box(), false, false, padding)

    return
}

func (g *GraphEditor) Box() *gtk.Box {
    return g.box
}

func (g *GraphEditor) Editors() map[string]attribute.Editor {
    return g.edit
}

type GraphProp struct {
    name      string
    attrs     map[string]attribute.Attribute
    visible   map[string]bool
}

func NewGraphProp(name string, visible map[string]bool) Property {
    g := &GraphProp{name: name, visible: visible}
    g.attrs = make(map[string]attribute.Attribute, 5)

    g.attrs["x"] = attribute.NewIntAttribute("rel_x")
    g.attrs["y"] = attribute.NewIntAttribute("rel_y")
    g.attrs["node"] = attribute.NewListAttribute("graph", "graph_node", 2, false)

    return g
}

func (g *GraphProp) Type() int {
    return GRAPH_PROP 
}

func (g *GraphProp) Name() string {
    return g.name
}

func (g *GraphProp) Visible() map[string]bool {
    return g.visible
}

func (g *GraphProp) Attributes() map[string]attribute.Attribute {
    return g.attrs
}

package props

import (
	"chroma-viz/library/attribute"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

/*

   Creating PropertyEditor's has a non negligible cost due to
   the gtk c calls, so to speed up initialisation we use a
   limited number of PropertyEditor's (enough to show a Page).

   PropertyEditor's consist of gtk ui elements which are used
   to edit Properties. When a user selects a Page in the Show,
   the Page is sent to the editor, which calls UpdateEditor
   (or UpdateEditorAllProp). This updates the PropertyEditor
   with the data from the Property.

   The changes made to the PropertyEditor are synced back to
   the Property using UpdateProp (see props.go).

   The Properties are built up from a collection of Attributes,
   and in a similar way, PropertyEditor's are built up from a
   collection of AttributeEditors.

*/

type PropertyEditor struct {
	PropType int
	Box      *gtk.Box
	editor   map[string]attribute.Editor
	visible  map[string]*gtk.CheckButton
}

func NewPropertyEditor(typed int) (propEdit *PropertyEditor, err error) {
	propEdit = &PropertyEditor{PropType: typed}
	width := 1920
	height := 1080

	propEdit.editor = make(map[string]attribute.Editor, 10)
	propEdit.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	propEdit.Box.SetVisible(true)

	switch typed {
	case RECT_PROP:
		propEdit.editor["rel_x"], _ = attribute.NewIntEditor("x", -float64(width), float64(width))
		propEdit.editor["rel_y"], _ = attribute.NewIntEditor("y", -float64(width), float64(height))
		propEdit.editor["width"], _ = attribute.NewIntEditor("Width", 0, float64(width))
		propEdit.editor["height"], _ = attribute.NewIntEditor("Height", 0, float64(height))
		propEdit.editor["rounding"], _ = attribute.NewIntEditor("Rounding", 0, float64(width))
		propEdit.editor["color"], _ = attribute.NewColorEditor("Color")

	case TEXT_PROP:
		propEdit.editor["rel_x"], _ = attribute.NewIntEditor("x", 0, float64(width))
		propEdit.editor["rel_y"], _ = attribute.NewIntEditor("y", 0, float64(height))
		propEdit.editor["string"], _ = attribute.NewStringEditor("Text")
		propEdit.editor["color"], _ = attribute.NewColorEditor("Color")

	case CIRCLE_PROP:
		propEdit.editor["rel_x"], _ = attribute.NewIntEditor("Center x", 0, float64(width))
		propEdit.editor["rel_y"], _ = attribute.NewIntEditor("Center y", 0, float64(height))
		propEdit.editor["inner_radius"], _ = attribute.NewIntEditor("Inner Radius", 0, float64(width))
		propEdit.editor["outer_radius"], _ = attribute.NewIntEditor("Outer Radius", 0, float64(width))
		propEdit.editor["start_angle"], _ = attribute.NewIntEditor("Start Angle", 0, 360)
		propEdit.editor["end_angle"], _ = attribute.NewIntEditor("End Angle", 0, 360)
		propEdit.editor["color"], _ = attribute.NewColorEditor("Color")

	case GRAPH_PROP:
		propEdit.editor["rel_x"], _ = attribute.NewIntEditor("x", 0, float64(width))
		propEdit.editor["rel_y"], _ = attribute.NewIntEditor("y", 0, float64(height))
		propEdit.editor["graph_node"], _ = attribute.NewListEditor("Graph", []string{"x Pos", "y Pos"})
		propEdit.editor["color"], _ = attribute.NewColorEditor("Color")

	case TICKER_PROP:
		propEdit.editor["rel_x"], _ = attribute.NewIntEditor("x", 0, float64(width))
		propEdit.editor["rel_y"], _ = attribute.NewIntEditor("y", 0, float64(height))
		propEdit.editor["string"], _ = attribute.NewListEditor("Ticker", []string{"Text"})
		propEdit.editor["color"], _ = attribute.NewColorEditor("Color")

	case CLOCK_PROP:
		propEdit.editor["rel_x"], _ = attribute.NewIntEditor("x", 0, float64(width))
		propEdit.editor["rel_y"], _ = attribute.NewIntEditor("y", 0, float64(height))
		propEdit.editor["string"], _ = attribute.NewClockEditor("Time")
		propEdit.editor["color"], _ = attribute.NewColorEditor("Color")

	case IMAGE_PROP:
		propEdit.editor["rel_x"], _ = attribute.NewIntEditor("x", 0, float64(width))
		propEdit.editor["rel_y"], _ = attribute.NewIntEditor("y", 0, float64(height))
		propEdit.editor["scale"], _ = attribute.NewFloatEditor("Scale", 0.01, 10, 0.01)
		propEdit.editor["image_id"] = attribute.NewAssetEditor("Image")

	default:
		return nil, fmt.Errorf("Unknown Prop %d", typed)
	}

	propEdit.AddEditors()
	return
}

var PropAttrs = map[int][]string{
	RECT_PROP:   {"rel_x", "rel_y", "width", "height", "rounding", "color"},
	TEXT_PROP:   {"rel_x", "rel_y", "color", "string"},
	CIRCLE_PROP: {"rel_x", "rel_y", "inner_radius", "outer_radius", "start_angle", "end_angle", "color"},
	GRAPH_PROP:  {"rel_x", "rel_y", "color", "graph_node"},
	TICKER_PROP: {"rel_x", "rel_y", "color", "string"},
	CLOCK_PROP:  {"rel_x", "rel_y", "color", "string"},
	IMAGE_PROP:  {"rel_x", "rel_y", "scale", "image_id"},
}

func (propEdit *PropertyEditor) AddEditors() (err error) {
	order := PropAttrs[propEdit.PropType]

	if len(order) == 0 {
		err = fmt.Errorf("Prop order for %d has length 0", propEdit.PropType)
		return
	}

	for _, name := range order {
		edit := propEdit.editor[name]
		if edit == nil {
			continue
		}

		propEdit.Box.PackStart(edit.Box(), edit.Expand(), edit.Expand(), padding)
	}

	return
}

func (propEdit *PropertyEditor) CreateVisibleEditor() (box *gtk.Box, err error) {
	widthChars := 12
	order := PropAttrs[propEdit.PropType]
	propEdit.visible = make(map[string]*gtk.CheckButton)

	box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	box.SetVisible(true)

	row, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	attrLabel, err := gtk.LabelNew("Attribute")
	if err != nil {
		return
	}

	visibleLabel, err := gtk.LabelNew("Visible")
	if err != nil {
		return
	}

	row.SetVisible(true)
	attrLabel.SetVisible(true)
	visibleLabel.SetVisible(true)

	attrLabel.SetWidthChars(widthChars)
	visibleLabel.SetWidthChars(widthChars)

	row.PackStart(attrLabel, false, false, padding)
	row.PackStart(visibleLabel, false, false, padding)
	box.PackStart(row, false, false, padding)

	for _, name := range order {
		attr := propEdit.editor[name]

		row, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
		if err != nil {
			return
		}

		attrLabel, err = gtk.LabelNew(attr.Name())
		if err != nil {
			return
		}

		propEdit.visible[name], err = gtk.CheckButtonNew()
		if err != nil {
			return
		}

		row.SetVisible(true)
		attrLabel.SetVisible(true)
		attrLabel.SetWidthChars(widthChars)

		propEdit.visible[name].SetVisible(true)
		propEdit.visible[name].SetMarginStart(40)

		row.PackStart(attrLabel, false, false, padding)
		row.PackStart(propEdit.visible[name], false, false, padding)
		box.PackStart(row, false, false, padding)
	}

	return
}

/*
Update PropertyEditor with the data in Property
*/
func (propEdit *PropertyEditor) UpdateEditor(prop *Property) (err error) {
	for name, edit := range propEdit.editor {
		if _, ok := prop.Attr[name]; !ok {
			continue
		}

		edit.Box().SetVisible(prop.Visible[name])
		err = edit.Update(prop.Attr[name])
		if err != nil {
			return
		}
	}

	return
}

func (propEdit *PropertyEditor) UpdateEditorAllProp(prop *Property) (err error) {
	for name, edit := range propEdit.editor {
		if _, ok := prop.Attr[name]; !ok {
			continue
		}

		edit.Box().SetVisible(true)
		if check := propEdit.visible[name]; check != nil {
			check.SetActive(prop.Visible[name])
		}

		err = edit.Update(prop.Attr[name])
		if err != nil {
			return
		}
	}

	return
}

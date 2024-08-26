package attribute

import (
	"chroma-viz/library/util"
	"fmt"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type PolyAttribute struct {
	Name       string
	Type       int
	numPoints  int
	PointStore *gtk.ListStore
}

func NewPolygonAttribute(name string) (poly *PolyAttribute, err error) {
	poly = &PolyAttribute{
		Name: name,
		Type: POLY,
	}

	poly.PointStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_INT, glib.TYPE_INT)
	return
}

func (poly *PolyAttribute) String() (s string) {
	var b strings.Builder
	b.WriteString("num_node=")
	b.WriteString(strconv.Itoa(poly.numPoints))
	b.WriteRune('#')

	iter, ok := poly.PointStore.GetIterFirst()
	model := poly.PointStore.ToTreeModel()
	for ok {
		index, err := util.ModelGetValue[int](model, iter, 0)
		if err != nil {
			continue
		}

		posX, err := util.ModelGetValue[int](model, iter, 1)
		if err != nil {
			continue
		}

		posY, err := util.ModelGetValue[int](model, iter, 2)
		if err != nil {
			continue
		}

		b.WriteString("point=")
		b.WriteString(strconv.Itoa(index))
		b.WriteRune(' ')
		b.WriteString(strconv.Itoa(posX))
		b.WriteRune(' ')
		b.WriteString(strconv.Itoa(posY))
		b.WriteRune('#')
	}

	return b.String()
}

func (poly *PolyAttribute) Encode() (s string) {
	return fmt.Sprintf("%d", poly.numPoints)
}

func (poly *PolyAttribute) Decode(s string) error {
	return nil
}

func (poly *PolyAttribute) Update(editor Editor) error {
	return nil
}

func (poly *PolyAttribute) Copy(attr Attribute) error {
	return nil
}

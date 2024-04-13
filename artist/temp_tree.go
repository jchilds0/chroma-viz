package artist

import (
	"chroma-viz/library/gtk_utils"
	"log"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	GEO_TYPE = iota
	GEO_NAME
	GEO_NUM
	GEO_NUM_COLS
)

const (
    FRAME_NUM = iota
    FRAME_GEOMETRY
    FRAME_GEOMETRY_ID
    FRAME_VALUE
    FRAME_VALUE_FRAME
    FRAME_VALUE_GEO
    FRAME_VALUE_ATTR
    FRAME_USER_VALUE
    FRAME_NUM_COLS
)

type TempTree struct {
	geoModel *gtk.TreeStore
	keyModel *gtk.TreeStore
	geoView  *gtk.TreeView
    keyView  *gtk.TreeView
}

func NewTempTree(propToEditor func(propID int)) (*TempTree, error) {
	temp := &TempTree{}

    err := temp.createGeometryTree(propToEditor)
    if err != nil {
        return nil, err
    }

    err = temp.createKeyTree()
    if err != nil {
        return nil, err
    }

	return temp, nil
}

func (temp *TempTree) createGeometryTree(propToEditor func(propID int)) (err error) {
    temp.geoView, err = gtk.TreeViewNew()
	if err != nil {
		return 
	}

	temp.geoView.Set("reorderable", true)

	typeCell, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	column, err := gtk.TreeViewColumnNewWithAttribute("Geometry", typeCell, "text", GEO_TYPE)
	if err != nil {
		return
	}
	temp.geoView.AppendColumn(column)

	nameCell, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	nameCell.SetProperty("editable", true)
	nameCell.Connect("edited", func(cell *gtk.CellRendererText, path, text string) {
		iter, err := temp.geoModel.GetIterFromString(path)
		if err != nil {
			log.Printf("Error editing geometry (%s)", err)
			return
		}

		model := temp.geoModel.ToTreeModel()
		geoID, err := gtk_utils.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Printf("Error editing geometry (%s)", err)
			return
		}

		geo := template.Geometry[geoID]
		if geo == nil {
			log.Printf("Error getting geometry %d", geoID)
			return
		}

		geo.Name = text
		temp.geoModel.SetValue(iter, GEO_NAME, text)
        temp.updateGeometry(geoID, text)
	})

	column, err = gtk.TreeViewColumnNewWithAttribute("Name", nameCell, "text", GEO_NAME)
	if err != nil {
		return
	}
	temp.geoView.AppendColumn(column)

	temp.geoModel, err = gtk.TreeStoreNew(
        glib.TYPE_STRING,       // GEO TYPE
        glib.TYPE_STRING,       // GEO NAME
        glib.TYPE_INT,          // GEO NUM
    )
	if err != nil {
		return
	}

	temp.geoView.SetModel(temp.geoModel)

	temp.geoView.Connect("row-activated",
		func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
			iter, err := temp.geoModel.GetIter(path)
			if err != nil {
				log.Printf("Error sending prop to editor (%s)", err)
				return
			}

			model := &temp.geoModel.TreeModel
			propID, err := gtk_utils.ModelGetValue[int](model, iter, GEO_NUM)
			if err != nil {
				log.Printf("Error sending prop to editor (%s)", err)
				return
			}

			propToEditor(propID)
		})

    return nil
}

func (temp *TempTree) createKeyTree() (err error) {
    temp.keyView, err = gtk.TreeViewNew()
	if err != nil {
		return
	}

	temp.keyView.Set("reorderable", true)

	temp.keyModel, err = gtk.TreeStoreNew(
        glib.TYPE_INT,          // Frame Num 
        glib.TYPE_STRING,       // Geometry Name
        glib.TYPE_INT,          // Geoemtry Num
        glib.TYPE_STRING,       // Value Entry
        glib.TYPE_STRING,       // Derived Value Frame
        glib.TYPE_STRING,       // Derived Value Geo
        glib.TYPE_STRING,       // Derived Value Attr 
        glib.TYPE_BOOLEAN,       // User Value Selector
    )
    if err != nil {
        return
    }

    temp.keyView.SetModel(temp.keyModel)

    // Frame Number
    {

        var frameNumCell *gtk.CellRendererText
        var column *gtk.TreeViewColumn

        frameNumCell, err = gtk.CellRendererTextNew()
        if err != nil {
            return
        }

        frameNumCell.SetProperty("editable", true)
        frameNumCell.Connect("edited", func(cell *gtk.CellRendererText, path, text string) {
            iter, err := temp.keyModel.GetIterFromString(path)
            if err != nil {
                log.Printf("Error editing keyframe (%s)", err)
                return
            }

            num, err := strconv.Atoi(text)
            if err != nil {
                log.Printf("Error editing keyframe (%s)", err)
                return
            }

            temp.keyModel.SetValue(iter, FRAME_NUM, num)
        })

        column, err = gtk.TreeViewColumnNewWithAttribute("Frame Number", frameNumCell, "text", FRAME_NUM)
        if err != nil {
            return
        }
        temp.keyView.AppendColumn(column)

    }

    // Geometry Name
    {

        var geoCell *gtk.CellRendererText
        var column *gtk.TreeViewColumn

        geoCell, err = gtk.CellRendererTextNew()
        if err != nil {
            return
        }

        column, err = gtk.TreeViewColumnNewWithAttribute("Geometry", geoCell, "text", FRAME_GEOMETRY)
        if err != nil {
            return
        }

        temp.keyView.AppendColumn(column)

    }

    // Set Value 
    {

        var valueCell *gtk.CellRendererText
        var column *gtk.TreeViewColumn

        valueCell, err = gtk.CellRendererTextNew()
        if err != nil {
            return
        }

        valueCell.SetProperty("editable", true)
        valueCell.Connect("edited", func(cell *gtk.CellRendererText, path, text string) {
            iter, err := temp.keyModel.GetIterFromString(path)
            if err != nil {
                log.Printf("Error editing geometry (%s)", err)
                return
            }

            temp.keyModel.SetValue(iter, FRAME_VALUE, text)
        })

        column, err = gtk.TreeViewColumnNewWithAttribute("Set Value", valueCell, "text", FRAME_VALUE)
        if err != nil {
            return
        }

        temp.keyView.AppendColumn(column)

    }

    // Derived Value
    {

        var valueText, valueCell *gtk.CellRendererText
        var column *gtk.TreeViewColumn

        column, err = gtk.TreeViewColumnNew()
        if err != nil {
            return
        }

        column.SetTitle("Value From KeyFrame")

        names := []string{"Frame", "Geometry", "Attr"}
        cols := []int{FRAME_VALUE_FRAME, FRAME_VALUE_GEO, FRAME_VALUE_ATTR}

        for i, name := range names {
            valueText, err = gtk.CellRendererTextNew()
            if err != nil {
                return
            }

            valueText.SetProperty("text", name + ": ")

            valueCell, err = gtk.CellRendererTextNew()
            if err != nil {
                return 
            }

            valueCell.SetProperty("editable", true)

            column.PackStart(valueText, false)
            column.PackStart(valueCell, true)

            column.AddAttribute(valueCell, "text", cols[i])

            valueCell.Connect("edited", func(cell *gtk.CellRendererText, path, text string) {
                iter, err := temp.keyModel.GetIterFromString(path)
                if err != nil {
                    log.Printf("Error editing geometry (%s)", err)
                    return
                }

                temp.keyModel.SetValue(iter, cols[i], text)
            })
        } 

        temp.keyView.AppendColumn(column)

    }

    // User Value 
    {

        var toggleCell *gtk.CellRendererToggle
        var column *gtk.TreeViewColumn

        toggleCell, err = gtk.CellRendererToggleNew()
        if err != nil {
            return
        }

        toggleCell.SetProperty("activatable", true)
        toggleCell.Connect("toggled", 
            func(cell *gtk.CellRendererToggle, path string) {
                iter, err := temp.keyModel.GetIterFromString(path)
                if err != nil {
                    log.Printf("Error toggling toggle (%s)", err)
                    return
                }

                state, err := gtk_utils.ModelGetValue[bool](temp.keyModel.ToTreeModel(), iter, FRAME_USER_VALUE)
                if err != nil {
                    log.Printf("Error toggling toggle (%s)", err)
                    return
                }

                temp.keyModel.SetValue(iter, FRAME_USER_VALUE, !state)
            })

        column, err = gtk.TreeViewColumnNewWithAttribute("User Value", toggleCell, "active", FRAME_USER_VALUE)
        if err != nil {
            return
        }

        temp.keyView.AppendColumn(column)

    }

    return nil
}

func (tempView *TempTree) updateGeometry(geoID int, name string) {
    iter, ok := tempView.keyModel.GetIterFirst()
    model := tempView.keyModel.ToTreeModel()

    for ok {
        currentID, err := gtk_utils.ModelGetValue[int](model, iter, FRAME_GEOMETRY_ID)
        if err != nil {
            log.Printf("Error getting keyframe geo id (%s)", err)
            ok = model.IterNext(iter)
            continue
        }

        if (currentID == geoID) {
            tempView.keyModel.SetValue(iter, FRAME_GEOMETRY, name)
        }

        ok = model.IterNext(iter)
    }
}

func (tempView *TempTree) removeGeometry(propID int) {
    iter, ok := tempView.keyModel.GetIterFirst()
    for ok {
        currentID, err := gtk_utils.ModelGetValue[int](tempView.keyModel.ToTreeModel(), iter, FRAME_GEOMETRY_ID)
        if err != nil {
            log.Printf("Error getting keyframe geo id (%s)", err)
            ok = tempView.keyModel.IterNext(iter)
            continue
        }

        if (currentID == propID) {
            tempView.keyModel.Remove(iter)
            iter, ok = tempView.keyModel.GetIterFirst()
        } else {
            ok = tempView.keyModel.IterNext(iter)
        }
    }
}

func (tempView *TempTree) AddGeoRow(iter *gtk.TreeIter, name, propName string, propNum int) {
	tempView.geoModel.SetValue(iter, GEO_TYPE, propName)
	tempView.geoModel.SetValue(iter, GEO_NAME, name)
	tempView.geoModel.SetValue(iter, GEO_NUM, propNum)
}

func (tempView *TempTree) AddKeyRow(iter *gtk.TreeIter, name string, propNum int) {
    tempView.keyModel.SetValue(iter, FRAME_NUM, 0)
    tempView.keyModel.SetValue(iter, FRAME_GEOMETRY, name)
    tempView.keyModel.SetValue(iter, FRAME_GEOMETRY_ID, propNum)
    tempView.keyModel.SetValue(iter, FRAME_USER_VALUE, true)
}

func (tempView *TempTree) Clean() {
	var err error
	tempView.geoModel, err = gtk.TreeStoreNew(
		glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_OBJECT, glib.TYPE_INT)
	if err != nil {
		log.Print(err)
		return
	}

	tempView.geoView.SetModel(tempView.geoModel)
}

package props

import (
	"chroma-viz/attribute"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type ImageEditor struct {
    box *gtk.Box 
    edit map[string]attribute.Editor
}

func NewImageEditor(width, height int, animate func()) (imageEdit *ImageEditor, err error) {
    imageEdit = &ImageEditor{}
    imageEdit.edit = make(map[string]attribute.Editor, 5)

    imageEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Fatalf("Error creating image prop editor (%s)", err) 
    }

    // TODO: replace with file explorer
    imageEdit.edit["string"], err = attribute.NewStringEditor("Image", animate)
    if err != nil {
        return
    }

    imageEdit.edit["x"], err = attribute.NewIntEditor("x", 0, float64(width), animate)
    if err != nil { 
        return
    }

    imageEdit.edit["y"], err = attribute.NewIntEditor("y", 0, float64(height), animate)
    if err != nil { 
        return 
    }

    imageEdit.edit["scale"], err = attribute.NewFloatEditor("Scale", 0.01, 10, 0.01, animate)
    if err != nil { 
        return
    }

    imageEdit.box.SetVisible(true)
    imageEdit.box.PackStart(imageEdit.edit["x"].Box(), false, false, padding)
    imageEdit.box.PackStart(imageEdit.edit["y"].Box(), false, false, padding)
    imageEdit.box.PackStart(imageEdit.edit["string"].Box(), false, false, padding)
    imageEdit.box.PackStart(imageEdit.edit["scale"].Box(), false, false, padding)

    return 
}

func (img *ImageEditor) Box() *gtk.Box {
    return img.box
}

func (img *ImageEditor) Editors() map[string]attribute.Editor {
    return img.edit
}

type ImageProp struct {
    name    string
    attrs   map[string]attribute.Attribute
    visible map[string]bool
}

func NewImageProp(name string, visible map[string]bool) *ImageProp {
    image := &ImageProp{name: name, visible: visible}
    image.attrs = make(map[string]attribute.Attribute, 5)

    image.attrs["x"] = attribute.NewIntAttribute("x", "rel_x")
    image.attrs["y"] = attribute.NewIntAttribute("y", "rel_y")
    image.attrs["scale"] = attribute.NewFloatAttribute("scale", "scale")
    image.attrs["string"] = attribute.NewStringAttribute("string", "string")

    return image 
}

func (image *ImageProp) Name() string {
    return image.name
}

func (image *ImageProp) Type() int {
    return IMAGE_PROP 
}

func (image *ImageProp) Visible() map[string]bool {
    return image.visible
}

func (image *ImageProp) Attributes() map[string]attribute.Attribute {
    return image.attrs
}

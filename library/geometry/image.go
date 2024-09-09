package geometry

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/util"
	"strings"
)

type Image struct {
	Geometry

	Scale attribute.FloatAttribute
	Image attribute.AssetAttribute
}

func NewImage(geo Geometry) *Image {
	image := &Image{
		Geometry: geo,
	}

	image.Scale.Name = ATTR_SCALE
	image.Image.Name = "image_id"
	return image
}

func (i *Image) UpdateGeometry(iEdit *ImageEditor) (err error) {
	err = i.Geometry.UpdateGeometry(&iEdit.GeometryEditor)
	if err != nil {
		return
	}

	err = i.Scale.UpdateAttribute(iEdit.Scale)
	if err != nil {
		return
	}

	err = i.Image.UpdateAttribute(iEdit.Image)
	return
}

func (i *Image) Encode(b *strings.Builder) {
	i.Geometry.Encode(b)

	util.EngineAddKeyValue(b, i.Image.Name, i.Image.Value)
	util.EngineAddKeyValue(b, i.Scale.Name, i.Scale.Value)
}

type ImageEditor struct {
	GeometryEditor

	Scale *attribute.FloatEditor
	Image *attribute.AssetEditor
}

func NewImageEditor() (imgEdit *ImageEditor, err error) {
	geoEdit, err := NewGeometryEditor()
	if err != nil {
		return
	}

	imgEdit = &ImageEditor{
		GeometryEditor: *geoEdit,
	}

	imgEdit.Scale, err = attribute.NewFloatEditor(Attrs[ATTR_SCALE], 0.0, 10.0, 0.01)
	if err != nil {
		return
	}

	imgEdit.Image = attribute.NewAssetEditor("Image")

	imgEdit.ScrollBox.PackStart(imgEdit.Scale.Box, false, false, padding)
	imgEdit.ScrollBox.PackStart(imgEdit.Image.Box, true, true, padding)
	return
}

func (iEdit *ImageEditor) UpdateEditor(i *Image) (err error) {
	err = iEdit.GeometryEditor.UpdateEditor(&i.Geometry)
	if err != nil {
		return
	}

	err = iEdit.Scale.UpdateEditor(&i.Scale)
	if err != nil {
		return
	}

	err = iEdit.Image.UpdateEditor(&i.Image)
	return
}

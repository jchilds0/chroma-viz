package geometry

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/parser"
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
	return
}

func (i *Image) Encode(b *strings.Builder) {
	i.Geometry.Encode(b)

	parser.EngineAddKeyValue(b, i.Image.Name, i.Image.Value)
	parser.EngineAddKeyValue(b, i.Scale.Name, i.Scale.Value)
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
	return
}

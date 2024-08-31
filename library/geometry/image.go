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

	image.Scale.Name = "scale"
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

	Scale attribute.FloatEditor
	Image attribute.AssetEditor
}

func NewImageEditor() (*ImageEditor, error) {
	return nil, nil
}

func (iEdit *ImageEditor) UpdateEditor(i *Image) (err error) {
	return
}

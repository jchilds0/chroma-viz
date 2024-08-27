package geometry

import (
	"chroma-viz/library/attribute"
	"strings"
)

type Image struct {
	Geometry

	Scale attribute.FloatAttribute
	Image attribute.AssetAttribute
}

func NewImage(geo Geometry) Image {
	image := Image{
		Geometry: geo,
	}

	image.Scale.Name = "scale"
	image.Image.Name = "image_id"
	return image
}

func (i *Image) UpdateGeometry(iEdit *ImageEditor) (err error) {
	return
}

func (i *Image) EncodeEngine(b strings.Builder) {

}

func (i *Image) EncodeJSON(b strings.Builder) {

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

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

func (i *Image) UpdateGeometry(iEdit *ImageEditor) (err error) {
	return
}

func (i *Image) EncodeEngine(b strings.Builder) {

}

type ImageEditor struct {
	GeometryEditor

	Scale attribute.FloatEditor
	Image attribute.AssetEditor
}

func NewImageEditor() *ImageEditor {
	return nil
}

func (iEdit *ImageEditor) UpdateEditor(i *Image) (err error) {
	return
}

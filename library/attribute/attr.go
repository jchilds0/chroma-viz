package attribute

import (
	"encoding/json"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

const padding = 10

/*

   Attributes and Editors come in pairs, for example
   in int.go we have IntAttribute which stores an
   integer and IntEditor which creates a gtk box with
   a label for the integer, a text entry and spin button
   for incrementing and decrementing the integer.

   Attributes are the building blocks of Properties.
   Each attribute has a String method for sending to
   Chroma Engine, Encode and Decode for import and
   export of shows respectively and Update for
   transferring the values of the corresponding Editor
   to the Attribute.

*/

const (
	INT = iota + 1
	STRING
	CLOCK
	FLOAT
	LIST
	COLOR
	ASSET
)

type Attribute interface {
	String() string
	Update(Editor) error
	Copy(Attribute) error

	// A -> {'name': string, 'value': string} | A, A
	Encode() string
	Decode(string) error
}

type AttributeJSON struct {
	Name string
	Type int
	attr Attribute
}

func (attrJSON *AttributeJSON) UnmarshalJSON(b []byte) error {
	var tempAttrJSON struct {
		AttributeJSON
		UnmarshalJSON struct{}
	}

	err := json.Unmarshal(b, &tempAttrJSON)
	if err != nil {
		return err
	}

	*attrJSON = tempAttrJSON.AttributeJSON

	switch attrJSON.Type {
	case INT:
		intAttr := &IntAttribute{}
		err = json.Unmarshal(b, intAttr)
		attrJSON.attr = intAttr

	case STRING:
		stringAttr := &StringAttribute{}
		err = json.Unmarshal(b, stringAttr)
		attrJSON.attr = stringAttr

	case FLOAT:
		floatAttr := &FloatAttribute{}
		err = json.Unmarshal(b, floatAttr)
		attrJSON.attr = floatAttr

	case LIST:
		listAttr := &ListAttribute{}
		err = json.Unmarshal(b, listAttr)
		attrJSON.attr = listAttr

	case CLOCK:
		clockAttr := &ClockAttribute{}
		err = json.Unmarshal(b, clockAttr)
		attrJSON.attr = clockAttr

	case COLOR:
		colorAttr := &ColorAttribute{}
		err = json.Unmarshal(b, colorAttr)
		attrJSON.attr = colorAttr

	case ASSET:
		assetAttr := &AssetAttribute{}
		err = json.Unmarshal(b, assetAttr)
		attrJSON.attr = assetAttr

	default:
		return fmt.Errorf("Error unknown attribute type %d", attrJSON.Type)
	}

	if err != nil {
		attrJSON.attr = nil
		return err
	}

	return nil
}

func (attrJSON *AttributeJSON) Attr() Attribute {
	return attrJSON.attr
}

/*

   Editors are the building block of PropertyEditors.
   Each Editor creates the gtk ui elements to edit
   the associated Attribute.

*/

type Editor interface {
	Box() *gtk.Box
	Update(Attribute) error
	Name() string
	Expand() bool
}

func SetIntValue(attr Attribute, val int) error {
	intAttr, ok := attr.(*IntAttribute)
	if !ok {
		return fmt.Errorf("Attribute is not an IntAttribute")
	}

	intAttr.Value = val
	return nil
}

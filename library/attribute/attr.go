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
	POLY
)

/*

   Editors are the building block of PropertyEditors.
   Each Editor creates the gtk ui elements to edit
   the associated Attribute.

*/

package attribute

import "github.com/gotk3/gotk3/gtk"

const padding = 10

type Attribute interface {
    String() string
    Encode() string
    Decode(string) error 
    Update(Editor) error
}

type Editor interface {
    Box() *gtk.Box
    Update(Attribute) error
}

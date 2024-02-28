package props

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type ImageEditor struct {
    box *gtk.Box 
    value [3]*gtk.SpinButton
    entry *gtk.Entry
    input *gtk.Box
}

func NewImageEditor(width, height int, animate func()) PropertyEditor {
    var err error
    image := &ImageEditor{}

    image.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Fatalf("Error creating image prop editor (%s)", err) 
    }

    image.box.SetVisible(true)

    // TODO: replace with file explorer
    image.input, image.entry = StringEditor("Image: ", animate)
    image.box.PackStart(image.input, false, false, padding)

    image.value[0], err = gtk.SpinButtonNewWithRange(-float64(width), float64(width), 1)
    if err != nil { 
        log.Fatalf("Error creating image prop editor (%s)", err) 
    }

    image.value[1], err = gtk.SpinButtonNewWithRange(-float64(width), float64(height), 1)
    if err != nil { 
        log.Fatalf("Error creating image prop editor (%s)", err) 
    }

    image.value[2], err = gtk.SpinButtonNewWithRange(0.01, 10, 0.01)
    if err != nil { 
        log.Fatalf("Error creating image prop editor (%s)", err) 
    }

    image.box.PackStart(IntEditor("x Pos", image.value[0], animate), false, false, padding)
    image.box.PackStart(IntEditor("y Pos", image.value[1], animate), false, false, padding)
    image.box.PackStart(IntEditor("Scale", image.value[2], animate), false, false, padding)

    return image
}

func (img *ImageEditor) Box() *gtk.Box {
    return img.box
}

func (imgEdit *ImageEditor) Update(img Property) {
    imgProp, ok := img.(*ImageProp)
    if !ok {
        log.Printf("ImageEditor.Update requires ImageProp")
        return
    }

    imgEdit.value[0].SetValue(float64(imgProp.Value[0]))
    imgEdit.value[1].SetValue(float64(imgProp.Value[1]))
    imgEdit.value[2].SetValue(imgProp.Scale)
    imgEdit.entry.SetText(imgProp.path)
}

type ImageProp struct {
    name    string
    path    string 
    Scale   float64
    Value   [2]int
}

func NewImageProp(name string) *ImageProp {
    image := &ImageProp{name: name, Scale: 1}
    return image 
}

func (image *ImageProp) Type() int {
    return IMAGE_PROP 
}

func (image *ImageProp) Name() string {
    return image.name
}

func (image *ImageProp) String() string {
    return fmt.Sprintf("string=%s#rel_x=%d#rel_y=%d#scale=%f#", 
        image.path, image.Value[0], image.Value[1], image.Scale)
}
 
func (image *ImageProp) Encode() string {
    return fmt.Sprintf("string %s;x %d;y %d;scale %f", 
        image.path, image.Value[0], image.Value[1], image.Scale)
}

func (image *ImageProp) Decode(input string) {
    attrs := strings.Split(input, ";")

    for _, attr := range attrs[1:] {
        line := strings.Split(attr, " ")
        name := line[0]

        switch (name) {
        case "x":
            value, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Fatalf("Error decoding image prop (%s)", err) 
            }

            image.Value[0] = value
        case "y":
            value, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Printf("Error decoding image prop (%s)", err) 
            }

            image.Value[1] = value
        case "string":
            image.path = strings.TrimPrefix(attr, "string ")
        case "scale":
            value, err := strconv.ParseFloat(line[1], 64)
            if err != nil {
                log.Printf("Error decoding image prop (%s)", err) 
            }

            image.Scale = value
        case "":
        default:
            log.Printf("Unknown ImageProp attr name (%s)\n", name)
        }
    }
}

func (imgProp *ImageProp) Update(image PropertyEditor, action int) {
    var err error
    imgEdit, ok := image.(*ImageEditor)
    if !ok {
        log.Printf("ImageProp.Update requires ImageEditor")
        return
    }

    imgProp.Value[0] = imgEdit.value[0].GetValueAsInt()
    imgProp.Value[1] = imgEdit.value[1].GetValueAsInt()
    imgProp.Scale = imgEdit.value[2].GetValue()
    imgProp.path, err = imgEdit.entry.GetText()

    if err != nil {
        log.Printf("Error getting path from editor entry (%s)", err)
        imgProp.path = ""
    }
}

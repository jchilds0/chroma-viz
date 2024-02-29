package props

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type VideoEditor struct {
    box *gtk.Box 
    value [3]*gtk.SpinButton
    entry *gtk.Entry
    input *gtk.Box
}

func NewVideoEditor(width, height int, animate func()) PropertyEditor {
    var err error
    video := &VideoEditor{}

    video.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Fatalf("Error creating video prop editor (%s)", err) 
    }

    video.box.SetVisible(true)

    // TODO: replace with file explorer
    video.input, video.entry = StringEditor("Video: ", animate)
    video.box.PackStart(video.input, false, false, padding)

    video.value[0], err = gtk.SpinButtonNewWithRange(-float64(width), float64(width), 1)
    if err != nil { 
        log.Fatalf("Error creating video prop editor (%s)", err) 
    }

    video.value[1], err = gtk.SpinButtonNewWithRange(-float64(width), float64(height), 1)
    if err != nil { 
        log.Fatalf("Error creating video prop editor (%s)", err) 
    }

    video.value[2], err = gtk.SpinButtonNewWithRange(0.01, 10, 0.01)
    if err != nil { 
        log.Fatalf("Error creating video prop editor (%s)", err) 
    }

    video.box.PackStart(IntEditor("x Pos", video.value[0], animate), false, false, padding)
    video.box.PackStart(IntEditor("y Pos", video.value[1], animate), false, false, padding)
    video.box.PackStart(IntEditor("Scale", video.value[2], animate), false, false, padding)

    return video
}

func (vid *VideoEditor) Box() *gtk.Box {
    return vid.box
}

func (vidEdit *VideoEditor) Update(vid Property) {
    vidProp, ok := vid.(*VideoProp)
    if !ok {
        log.Printf("VideoEditor.Update requires VideoProp")
        return
    }

    vidEdit.value[0].SetValue(float64(vidProp.Value[0]))
    vidEdit.value[1].SetValue(float64(vidProp.Value[1]))
    vidEdit.value[2].SetValue(vidProp.Scale)
    vidEdit.entry.SetText(vidProp.path)
}

type VideoProp struct {
    name    string
    path    string 
    Scale   float64
    Value   [2]int
}

func NewVideoProp(name string) *VideoProp {
    video := &VideoProp{name: name, Scale: 1}
    return video 
}

func (video *VideoProp) Type() int {
    return VIDEO_PROP 
}

func (video *VideoProp) Name() string {
    return video.name
}

func (video *VideoProp) String() string {
    return fmt.Sprintf("string=%s#rel_x=%d#rel_y=%d#scale=%f#", 
        video.path, video.Value[0], video.Value[1], video.Scale)
}
 
func (video *VideoProp) Encode() string {
    return fmt.Sprintf("string %s;x %d;y %d;scale %f", 
        video.path, video.Value[0], video.Value[1], video.Scale)
}

func (video *VideoProp) Decode(input string) {
    attrs := strings.Split(input, ";")

    for _, attr := range attrs[1:] {
        line := strings.Split(attr, " ")
        name := line[0]

        switch (name) {
        case "x":
            value, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Fatalf("Error decoding video prop (%s)", err) 
            }

            video.Value[0] = value
        case "y":
            value, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Printf("Error decoding video prop (%s)", err) 
            }

            video.Value[1] = value
        case "string":
            video.path = strings.TrimPrefix(attr, "string ")
        case "scale":
            value, err := strconv.ParseFloat(line[1], 64)
            if err != nil {
                log.Printf("Error decoding video prop (%s)", err) 
            }

            video.Scale = value
        case "":
        default:
            log.Printf("Unknown VideoProp attr name (%s)\n", name)
        }
    }
}

func (vidProp *VideoProp) Update(video PropertyEditor, action int) {
    var err error
    vidEdit, ok := video.(*VideoEditor)
    if !ok {
        log.Printf("VideoProp.Update requires VideoEditor")
        return
    }

    vidProp.Value[0] = vidEdit.value[0].GetValueAsInt()
    vidProp.Value[1] = vidEdit.value[1].GetValueAsInt()
    vidProp.Scale = vidEdit.value[2].GetValue()
    vidProp.path, err = vidEdit.entry.GetText()

    if err != nil {
        log.Printf("Error getting path from editor entry (%s)", err)
        vidProp.path = ""
    }
}

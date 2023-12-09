package gui

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

const padding = 10

const (
    START = iota
    PAUSE
    STOP
)

func IntEditor(name string, lowerBound, upperBound int, 
    spin *gtk.SpinButton, animate interface{}) *gtk.Box {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    box.SetVisible(true)

    label, _ := gtk.LabelNew(name)
    label.SetVisible(true)
    label.SetWidthChars(7)
    box.PackStart(label, false, false, uint(padding))

    spin.SetVisible(true)
    spin.SetValue(0)
    box.PackStart(spin, false, false, 0)

    spin.Connect("value-changed", animate)
    return box
}

func TextEditor(name string, animate interface{}) (*gtk.Box, *gtk.Entry) {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    box.SetVisible(true)

    label, _ := gtk.LabelNew(name)
    label.SetVisible(true)
    label.SetWidthChars(7)
    box.PackStart(label, false, false, uint(padding))

    buf, _ := gtk.EntryBufferNew("", 0)
    text, _ := gtk.EntryNewWithBuffer(buf)
    text.SetVisible(true)
    box.PackStart(text, false, false, 0)

    text.Connect("changed", animate)

    return box, text
}

type Property interface {
    Tab() *gtk.Box
    String() string
    Encode() string
    Decode(string)
}

type RectProp struct {
    value [4]*gtk.SpinButton
    input [4]*gtk.Box
    box *gtk.Box
}

func NewRectProp(width, height int, animate func()) Property {
    rect := &RectProp{}

    rect.value[0], _ = gtk.SpinButtonNewWithRange(float64(0), float64(width), 1)
    rect.value[1], _ = gtk.SpinButtonNewWithRange(float64(0), float64(height), 1)
    rect.value[2], _ = gtk.SpinButtonNewWithRange(float64(0), float64(width), 1)
    rect.value[3], _ = gtk.SpinButtonNewWithRange(float64(0), float64(height), 1)

    rect.input[0] = IntEditor("x Pos", 0, width, rect.value[0], animate)
    rect.input[1] = IntEditor("y Pos", 0, height, rect.value[1], animate)
    rect.input[2] = IntEditor("width", 0, width, rect.value[2], animate)
    rect.input[3] = IntEditor("height", 0, height, rect.value[3], animate)

    rect.box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

    for _, in := range rect.input {
        rect.box.PackStart(in, false, false, padding)
    }

    rect.box.SetVisible(true)

    return rect
}

func (rect *RectProp) Tab() *gtk.Box {
    return rect.box
}

func (rect *RectProp) String() string {
    return fmt.Sprintf("pos_x#%d#pos_y#%d#width#%d#height#%d#", 
        rect.value[0].GetValueAsInt(),
        rect.value[1].GetValueAsInt(),
        rect.value[2].GetValueAsInt(),
        rect.value[3].GetValueAsInt())
}

func (rect *RectProp) Encode() string {
    return fmt.Sprintf("x %d;y %d;width %d;height %d;", 
        rect.value[0].GetValueAsInt(),
        rect.value[1].GetValueAsInt(),
        rect.value[2].GetValueAsInt(),
        rect.value[3].GetValueAsInt())
}

func (rect *RectProp) Decode(input string) {
    attrs := strings.Split(input, ";")

    for _, attr := range attrs[1:] {
        line := strings.Split(attr, " ")

        if len(line) != 2 {
            continue
        }

        name := line[0]
        value, _ := strconv.Atoi(line[1])

        switch (name) {
        case "x":
            rect.value[0].SetValue(float64(value))
        case "y":
            rect.value[1].SetValue(float64(value))
        case "width":
            rect.value[2].SetValue(float64(value))
        case "height":
            rect.value[3].SetValue(float64(value))
        default:
            log.Printf("Unknown RectProp attr name (%s)\n", name)
        }
    }
}

type TextProp struct {
    text []*gtk.Entry
    input []*gtk.Box
    numLines int
    box *gtk.Box
}

func NewTextProp(numLines int, f func()) *TextProp {
    text := &TextProp{numLines: numLines}
    text.text = make([]*gtk.Entry, numLines)
    text.input = make([]*gtk.Box, numLines)

    for i := range text.text {
        text.input[i], text.text[i] = TextEditor("Line " + strconv.Itoa(i + 1), f)
    }

    text.box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    text.box.SetVisible(true)

    for _, in := range text.input {
        text.box.PackStart(in, false, false, padding)
    }

    return text
}

func (text *TextProp) Tab() *gtk.Box {
    return text.box 
}

func (text *TextProp) String() string {
    str := ""
    for i, entry := range text.text {
        entryText, _ := entry.GetText()
        str = str + fmt.Sprintf("text%d#%s#", i, entryText)
    }

    return str
}
 
func (text *TextProp) Encode() string {
    str := ""
    for _, entry := range text.text {
        entryText, _ := entry.GetText()
        str = str + fmt.Sprintf("%s;", entryText)
    }

    return str
}

func (text *TextProp) Decode(input string) {
    strings := strings.Split(input, ";")

    for i := range text.text {
        text.text[i].SetText(strings[i + 1])
    }
}

type ClockProp struct {
    box *gtk.Box
    page *Page
    active bool
    c chan int
    currentTime time.Time
    editTime *time.Time
    timeFormat string
}

func NewClockProp(page *Page, animate func()) *ClockProp {
    clock := &ClockProp{page: page}
    clock.box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    clock.box.SetVisible(true)
    clock.c = make(chan int, 1)

    actions, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    clock.box.PackStart(actions, false, false, padding)
    actions.SetVisible(true)

    startButton, _ := gtk.ButtonNewWithLabel("Start")
    actions.PackStart(startButton, false, false, padding)
    startButton.SetVisible(true)
    startButton.Connect("clicked", func() {
        clock.c <- START
    })

    pauseButton, _ := gtk.ButtonNewWithLabel("Pause")
    actions.PackStart(pauseButton, false, false, padding)
    pauseButton.SetVisible(true)
    pauseButton.Connect("clicked", func() {
        clock.c <- PAUSE 
    })

    stopButton, _ := gtk.ButtonNewWithLabel("Stop")
    actions.PackStart(stopButton, false, false, padding)
    stopButton.SetVisible(true)
    stopButton.Connect("clicked", func() {
        clock.c <- STOP
    })

    clock.timeFormat = "04:05"
    edit, _ := time.Parse(clock.timeFormat, "00:00")
    clock.editTime = &edit

    inputBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    inputBox.SetVisible(true)
    clock.box.PackStart(inputBox, false, false, padding)

    labelTime, _ := gtk.LabelNew("Start Time: ")
    labelTime.SetVisible(true)
    labelTime.SetWidthChars(7)
    inputBox.PackStart(labelTime, false, false, padding)

    editTime := clock.editTime.Format(clock.timeFormat)
    bufTime, _ := gtk.EntryBufferNew(editTime, len(editTime))
    textTime, _ := gtk.EntryNewWithBuffer(bufTime)
    textTime.SetVisible(true)
    inputBox.PackStart(textTime, false, false, 0)

    textTime.Connect("changed", 
        func(e *gtk.Entry) { 
            text, err := e.GetText()
            newTime, err := time.Parse(clock.timeFormat, text)
            if err != nil {
                log.Printf("Incorrect time format entered (%s)\n", text)
                return
            }
            *clock.editTime = newTime
            animate()
        })

    go clock.RunClock(animate)
    return clock
}

func (clock *ClockProp) Tab() *gtk.Box {
    return clock.box
}

func (clock *ClockProp) String() string {
    return fmt.Sprintf("text0#%s#", 
        clock.currentTime.Format(clock.timeFormat))
}

func (clock *ClockProp) Encode() string {
    return ""
}

func (clock *ClockProp) Decode(input string) {
}

func (clock *ClockProp) RunClock(animate func()) {
    state := PAUSE
    clock.currentTime = *clock.editTime

    for {
        select {
        case state = <-clock.c:
        default:
        }    

        switch state {
        case START:
            animate()
            clock.currentTime = clock.currentTime.Add(time.Second)
            time.Sleep(1 * time.Second)
        case PAUSE:
            // block until we recieve an instruction
            state = <-clock.c
        case STOP:
            // reset the time and block
            clock.currentTime = *clock.editTime
            animate()
            state = <-clock.c
        default:
            log.Printf("Clock recieved unknown value through channel %d\n", state)
        }
    }
}

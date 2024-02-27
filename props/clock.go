package props

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

// TODO: untangle clock updating

type ClockEditor struct {
    box             *gtk.Box
    value           [2]*gtk.SpinButton
    engine          bool
    preview         bool
    c               chan int
    timeFormat      string 
    editTime        *time.Time
    currentTime     time.Time
    entry           *gtk.Entry
}

func NewClockEditor(width, height int, animate, cont func()) PropertyEditor {
    var err error
    clock := &ClockEditor{
        timeFormat: "04:05",
    }

    clock.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    clock.box.SetVisible(true)
    clock.c = make(chan int, 1)

    actions, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    clock.box.PackStart(actions, false, false, padding)
    actions.SetVisible(true)

    startButton, err := gtk.ButtonNewWithLabel("Start")
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    actions.PackStart(startButton, false, false, padding)
    startButton.SetVisible(true)
    startButton.Connect("clicked", func() {
        clock.c <- START
    })

    pauseButton, err := gtk.ButtonNewWithLabel("Pause")
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    actions.PackStart(pauseButton, false, false, padding)
    pauseButton.SetVisible(true)
    pauseButton.Connect("clicked", func() {
        clock.c <- PAUSE 
    })

    stopButton, err := gtk.ButtonNewWithLabel("Stop")
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    actions.PackStart(stopButton, false, false, padding)
    stopButton.SetVisible(true)
    stopButton.Connect("clicked", func() {
        clock.c <- STOP
    })

    var input *gtk.Box
    input, clock.entry = StringEditor("Time: ", animate)
    clock.box.PackStart(input, false, false, padding)

    clock.value[0], err = gtk.SpinButtonNewWithRange(-float64(width), float64(width), 1)
    if err != nil { 
        log.Fatalf("Error creating text prop (%s)", err) 
    }

    clock.value[1], err = gtk.SpinButtonNewWithRange(-float64(height), float64(height), 1)
    if err != nil { 
        log.Fatalf("Error creating text prop (%s)", err) 
    }

    clock.box.PackStart(IntEditor("x Pos", clock.value[0], animate), false, false, padding)
    clock.box.PackStart(IntEditor("y Pos", clock.value[1], animate), false, false, padding)

    go clock.RunClock(cont)
    return clock
}

func (clock *ClockEditor) RunClock(cont func()) {
    var err error
    state := PAUSE
    tick := time.NewTicker(time.Second)

    if err != nil {
        log.Printf("Error parsing edit time (%s)", err)
        return
    }

    for {
        select {
        case state = <-clock.c:
        case <-tick.C:
        }    

        switch state {
        case START:
            // update time and animate
            cont()
            clock.currentTime = clock.currentTime.Add(time.Second)
        case PAUSE:
            // block until we recieve an instruction
            state = <-clock.c
        case STOP:
            // reset the time and block
            clock.currentTime = *clock.editTime
            cont()
            state = <-clock.c
        default:
            log.Printf("Clock recieved unknown value through channel %d\n", state)
        }
    }
}

func (clock *ClockEditor) Box() *gtk.Box {
    return clock.box
}

func (clockEdit *ClockEditor) Update(clock Property) {
    clockProp, ok := clock.(*ClockProp)
    if !ok {
        log.Printf("ClockEditor.Update requires ClockProp")
        return
    }

    clockEdit.editTime = clockProp.editTime
    text := clockEdit.editTime.Format(clockEdit.timeFormat)
    clockEdit.entry.SetText(text)
    clockEdit.value[0].SetValue(float64(clockProp.Value[0]))
    clockEdit.value[1].SetValue(float64(clockProp.Value[1]))
}

type ClockProp struct {
    name            string
    Value           [2]int
    editTime        *time.Time
    CurrentTime     string
    timeFormat      string
}

func NewClockProp(name string) *ClockProp {
    clock := &ClockProp{
        name: name,
        timeFormat: "04:05",
    }

    edit, err := time.Parse(clock.timeFormat, "00:00")
    if err != nil { 
        log.Printf("Error creating clock prop (%s)", err) 
    }

    clock.editTime = &edit
    return clock
}

func (clock *ClockProp) Type() int {
    return CLOCK_PROP 
}

func (clock *ClockProp) Name() string {
    return clock.name
}

func (clock *ClockProp) String() string {
    currentString := fmt.Sprintf("string=%s#rel_x=%d#rel_y=%d#", 
        clock.CurrentTime,
        clock.Value[0],
        clock.Value[1],
    )
    
    return currentString
}

func (clock *ClockProp) Encode() string {
    return fmt.Sprintf("string %s;x %d;y %d;", 
        clock.editTime.Format(clock.timeFormat), clock.Value[0], clock.Value[1])
}

func (clock *ClockProp) Decode(input string) {
    attrs := strings.Split(input, ";")

    for _, attr := range attrs[1:] {
        line := strings.Split(attr, " ")
        name := line[0]

        switch (name) {
        case "x":
            value, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Fatalf("Error decoding text prop (%s)", err) 
            }

            clock.Value[0] = value
        case "y":
            value, err := strconv.Atoi(line[1])
            if err != nil { 
                log.Printf("Error decoding text prop (%s)", err) 
            }

            clock.Value[1] = value
        case "string":
            textTime := strings.TrimPrefix(attr, "string ")
            edit, err := time.Parse(clock.timeFormat, textTime)
            if err != nil { 
                log.Printf("Error decoding clock (%s)", err) 
            }

            clock.editTime = &edit
        case "":
        default:
            log.Printf("Unknown TextProp attr name (%s)\n", name)
        }
    }
}

func (clockProp *ClockProp) Update(clock PropertyEditor, action int) {
    clockEdit, ok := clock.(*ClockEditor)

    if !ok {
        log.Printf("ClockProp.Update requires ClockEditor") 
        return
    }

    switch action {
    case ANIMATE_ON:
        clockProp.Value[0] = clockEdit.value[0].GetValueAsInt()
        clockProp.Value[1] = clockEdit.value[1].GetValueAsInt()

        editText := clockEdit.editTime.Format(clockEdit.timeFormat)
        editTime, err := time.Parse(clockProp.timeFormat, editText)
        if err != nil {
            log.Printf("Error parsing time in ClockProp.Update (%s)", err)
            return
        }
        clockProp.editTime = &editTime
    case CONTINUE:
        clockProp.CurrentTime = clockEdit.currentTime.Format(clockProp.timeFormat)
    case ANIMATE_OFF:
    default:
        log.Printf("Unknown action")
    }
}

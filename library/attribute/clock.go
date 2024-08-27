package attribute

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

const padding = 10

const (
	START = iota
	PAUSE
	STOP
)

type ClockAttribute struct {
	Name        string
	Type        int
	c           chan int
	CurrentTime string
	TimeFormat  string
	m           sync.Mutex
}

func NewClockAttribute(name string) *ClockAttribute {
	clockAttr := &ClockAttribute{
		Name:       name,
		TimeFormat: "04:05",
		c:          make(chan int),
	}

	return clockAttr
}

func (clockAttr *ClockAttribute) UnmarshalJSON(b []byte) error {
	var clockAttrJSON struct {
		ClockAttribute
		UnmarshalJSON struct{}
	}

	err := json.Unmarshal(b, &clockAttrJSON)
	if err != nil {
		return err
	}

	clockAttr.Name = clockAttrJSON.Name
	clockAttr.TimeFormat = clockAttrJSON.TimeFormat
	clockAttr.CurrentTime = clockAttrJSON.CurrentTime

	clockAttr.c = make(chan int)
	return nil
}

func (clockAttr *ClockAttribute) SetClock(cont func()) {
	go clockAttr.RunClock(cont)
}

func (clockAttr *ClockAttribute) EncodeEngine() string {
	return fmt.Sprintf("%s=%s#", clockAttr.Name, clockAttr.CurrentTime)
}

func (clockAttr *ClockAttribute) Update(clockEdit *ClockEditor) error {
	var err error
	clockAttr.CurrentTime, err = clockEdit.entry.GetText()
	return err
}

func (clock *ClockAttribute) RunClock(cont func()) {
	state := PAUSE
	tick := time.NewTicker(time.Second)
	run := false

	for {
		if run {
			select {
			case state = <-clock.c:
			case <-tick.C:
				if run {
					cont()
					clock.tickTime()
				}

				continue
			}
		} else {
			state = <-clock.c
		}

		switch state {
		case START:
			tick.Reset(time.Second)
			run = true
		case PAUSE:
			run = false
		case STOP:
			run = false

			clock.m.Lock()
			clock.CurrentTime = "00:00"
			clock.m.Unlock()

			cont()
		default:
			log.Printf("Clock recieved unknown value through channel %d\n", state)
		}
	}
}

func (clock *ClockAttribute) tickTime() {
	clock.m.Lock()
	defer clock.m.Unlock()

	currentTime, err := time.Parse(clock.TimeFormat, clock.CurrentTime)
	if err != nil {
		log.Println(err)
		return
	}

	currentTime = currentTime.Add(time.Second)
	clock.CurrentTime = currentTime.Format(clock.TimeFormat)
}

type ClockEditor struct {
	box   *gtk.Box
	entry *gtk.Entry
	c     chan int
	name  string
}

func NewClockEditor(name string) (clockEdit *ClockEditor, err error) {
	clockEdit = &ClockEditor{name: name}

	clockEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	clockEdit.box.SetVisible(true)

	actions, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	clockEdit.box.PackStart(actions, false, false, padding)
	actions.SetVisible(true)

	startButton, err := gtk.ButtonNewWithLabel("Start")
	if err != nil {
		return
	}

	actions.PackStart(startButton, false, false, padding)
	startButton.SetVisible(true)
	startButton.Connect("clicked", func() {
		clockEdit.c <- START
	})

	pauseButton, err := gtk.ButtonNewWithLabel("Pause")
	if err != nil {
		return
	}

	actions.PackStart(pauseButton, false, false, padding)
	pauseButton.SetVisible(true)
	pauseButton.Connect("clicked", func() {
		clockEdit.c <- PAUSE
	})

	stopButton, err := gtk.ButtonNewWithLabel("Stop")
	if err != nil {
		return
	}

	actions.PackStart(stopButton, false, false, padding)
	stopButton.SetVisible(true)
	stopButton.Connect("clicked", func() {
		clockEdit.c <- STOP
	})

	timeBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	timeBox.SetVisible(true)
	clockEdit.box.PackStart(timeBox, false, false, padding)

	label, err := gtk.LabelNew(name)
	if err != nil {
		return
	}

	label.SetVisible(true)
	label.SetWidthChars(12)
	timeBox.PackStart(label, false, false, padding)

	buf, err := gtk.EntryBufferNew("", 0)
	if err != nil {
		return
	}

	clockEdit.entry, err = gtk.EntryNewWithBuffer(buf)
	if err != nil {
		return
	}

	clockEdit.entry.SetVisible(true)
	timeBox.PackStart(clockEdit.entry, false, false, 0)

	return
}

func (clockEdit *ClockEditor) Name() string {
	return clockEdit.name
}

func (clockEdit *ClockEditor) Update(clockAttr *ClockAttribute) error {
	clockEdit.c = clockAttr.c
	clockEdit.entry.SetText(clockAttr.CurrentTime)
	return nil
}

func (clockEdit *ClockEditor) Box() *gtk.Box {
	return clockEdit.box
}

func (clockEdit *ClockEditor) Expand() bool {
	return false
}

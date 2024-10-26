package attribute

import (
	"log"
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
	CurrentTime string
	TimeFormat  string

	cont   func()
	c      chan int
	active bool
}

func NewClockAttribute(name string) *ClockAttribute {
	clockAttr := &ClockAttribute{
		Name:       name,
		TimeFormat: "04:05",
	}

	return clockAttr
}

func (clockAttr *ClockAttribute) SetClock(cont func()) {
	clockAttr.cont = cont
}

func (clockAttr *ClockAttribute) UpdateAttribute(clockEdit *ClockEditor) error {
	var err error
	clockAttr.CurrentTime, err = clockEdit.entry.GetText()
	return err
}

func (clock *ClockAttribute) InitClock() {
	clock.c = make(chan int, 1)
	clock.active = true
	state := PAUSE
	tick := time.NewTicker(time.Second)
	run := false
	startTime := clock.CurrentTime
	currentTime := clock.CurrentTime

	go func() {
		for {
			var ok bool
			if run {
				select {
				case state, ok = <-clock.c:
				case <-tick.C:
					if run {
						clock.cont()
						currentTime = clock.tickTime(currentTime)
					}

					continue
				}
			} else {
				state, ok = <-clock.c
			}

			if !ok {
				return
			}

			switch state {
			case START:
				run = true
				tick.Reset(time.Second)
			case PAUSE:
				run = false
			case STOP:
				run = false
				currentTime = startTime
			default:
				log.Printf("Clock recieved unknown value through channel %d\n", state)
			}
		}
	}()
}

func (clock *ClockAttribute) tickTime(currentTime string) string {
	t, err := time.Parse(clock.TimeFormat, currentTime)
	if err != nil {
		log.Println(err)
		return currentTime
	}

	tick := t.Add(time.Second)
	return tick.Format(clock.TimeFormat)
}

type ClockEditor struct {
	Box   *gtk.Box
	entry *gtk.Entry
	c     chan int
	name  string
}

func NewClockEditor(name string) (clockEdit *ClockEditor, err error) {
	clockEdit = &ClockEditor{name: name}

	clockEdit.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	clockEdit.Box.SetVisible(true)

	actions, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	clockEdit.Box.PackStart(actions, false, false, padding)
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
		close(clockEdit.c)
	})

	timeBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	timeBox.SetVisible(true)
	clockEdit.Box.PackStart(timeBox, false, false, padding)

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

func (clockEdit *ClockEditor) UpdateEditor(clockAttr *ClockAttribute) error {
	if !clockAttr.active {
		clockAttr.InitClock()
	}

	clockEdit.c = clockAttr.c
	clockEdit.entry.SetText(clockAttr.CurrentTime)
	return nil
}

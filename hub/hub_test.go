package hub

import (
	"chroma-viz/library/props"
	"chroma-viz/library/templates"
	"log"
	"net"
	"testing"
	"time"
)

func TestImportTemplates(t *testing.T) {
	hub := NewDataBase(10)

	err := hub.ImportArchive("test_archive.json")
	if err != nil {
		t.Errorf("Error importing test archive (%s)", err)
	}
	go StartHub(hub, 9000)

	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		t.Fatalf("Error connecting to graphics hub (%s)", err)
	}

	for i := 0; i < 5; i++ {
		template, err := templates.GetTemplate(conn, i)
		if err != nil {
			t.Fatalf("Error receiving template (%s)", err)
		}

		log.Printf("Received Template %d", template.TempID)

		switch template.Title {
		case "Teal Box":
			propTest(t, template, 0, "Background", "rect")
			propTest(t, template, 1, "Circle", "circle")
			propTest(t, template, 2, "Title", "text")
			propTest(t, template, 3, "Subtitle", "text")
		case "Clock":
			propTest(t, template, 0, "Background", "rect")
			propTest(t, template, 2, "Circle", "circle")
			propTest(t, template, 3, "Left Split", "rect")
			propTest(t, template, 4, "Team 1", "text")
			propTest(t, template, 5, "Score 1", "text")
			propTest(t, template, 6, "Mid Split", "rect")
			propTest(t, template, 1, "Team 2", "text")
			propTest(t, template, 7, "Score 2", "text")
			propTest(t, template, 8, "Right Split", "rect")
			propTest(t, template, 9, "Clock", "clock")
			propTest(t, template, 10, "Period", "text")
		case "White Circle":
			propTest(t, template, 0, "Circle", "circle")
		case "Graph":
			propTest(t, template, 0, "Background", "rect")
			propTest(t, template, 1, "Graph", "graph")
			propTest(t, template, 2, "Title", "text")
		case "Ticker":
			propTest(t, template, 0, "Background", "rect")
			propTest(t, template, 1, "Box", "rect")
			propTest(t, template, 2, "Ticker", "ticker")
		default:
			t.Errorf("Unknown template %s", template.Title)
		}
	}
}

func propTest(t *testing.T, template *templates.Template, i int, name, typed string) {
	geo := template.Geometry[i]
	if geo.Name != name {
		t.Errorf("(%s) Incorrect prop name, expected %s, recieved %s",
			template.Title, name, geo.Name)
	}

	if props.PropType(geo.PropType) != typed {
		t.Errorf("(%s) Incorrect prop type, expected %s, recieved %s",
			template.Title, typed, props.PropType(geo.PropType))
	}
}

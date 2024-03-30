package shows

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/props"
	"testing"
)

func TestImportShow(t *testing.T) {
	show := NewShow()
	show.ImportShow("test_show.json")

	if len(show.Pages) != 4 {
		t.Errorf("Incorrect number of pages (len(show.Pages) = %d)", len(show.Pages))
	}

	for _, page := range show.Pages {
		switch page.Title {
		case "Teal Box":
			rectPropTest(t, page.PropMap[0], 50, 100, 850, 180)
			circlePropTest(t, page.PropMap[1], 99, 90, 30, 75, 45, 315)
			textPropTest(t, page.PropMap[2], 190, 100, "Lower Frame")
			textPropTest(t, page.PropMap[3], 190, 30, "A lower frame subtitle")
		case "Clock Box":
			rectPropTest(t, page.PropMap[0], 50, 965, 810, 65)
			circlePropTest(t, page.PropMap[1], 35, 32, 0, 25, 0, 360)
			rectPropTest(t, page.PropMap[2], 76, 0, 10, 65)
			textPropTest(t, page.PropMap[3], 20, 15, "TEAM")
			textPropTest(t, page.PropMap[4], 180, 15, "0")
			rectPropTest(t, page.PropMap[5], 230, 0, 10, 65)
			textPropTest(t, page.PropMap[6], 20, 15, "TEAM")
			textPropTest(t, page.PropMap[7], 180, 15, "0")
			rectPropTest(t, page.PropMap[8], 230, 0, 10, 65)
			clockPropTest(t, page.PropMap[9], 110, 15)
			textPropTest(t, page.PropMap[10], 30, 15, "Q1")
		case "Ticker":
			rectPropTest(t, page.PropMap[0], 0, 0, 1920, 75)
			rectPropTest(t, page.PropMap[1], 1700, 25, 400, 100)
			tickerPropTest(t, page.PropMap[2], 25, 20, "Hello there", "world")
		case "Graph":
			rectPropTest(t, page.PropMap[0], 50, 100, 570, 240)
			graphPropTest(t, page.PropMap[1], 25, 25)
			textPropTest(t, page.PropMap[2], 50, 175, "Step Graph")
		default:
			t.Errorf("Unknown page %s", page.Title)
		}
	}
}

func rectPropTest(t *testing.T, prop *props.Property, x, y, w, h int) {
	if prop.PropType != props.RECT_PROP {
		t.Errorf("Prop %s is not a rect prop", prop.Name)
		return
	}

	intAttrTest(t, prop.Attr["x"], x)
	intAttrTest(t, prop.Attr["y"], y)
	intAttrTest(t, prop.Attr["width"], w)
	intAttrTest(t, prop.Attr["height"], h)
}

func textPropTest(t *testing.T, prop *props.Property, x, y int, text string) {
	if prop.PropType != props.TEXT_PROP {
		t.Errorf("Prop %s is not a text prop", prop.Name)
		return
	}

	intAttrTest(t, prop.Attr["x"], x)
	intAttrTest(t, prop.Attr["y"], y)
	stringAttrTest(t, prop.Attr["string"], text)
}

func circlePropTest(t *testing.T, prop *props.Property, x, y, ir, or, sa, ea int) {
	if prop.PropType != props.CIRCLE_PROP {
		t.Errorf("Prop %s is not a circle prop", prop.Name)
		return
	}

	intAttrTest(t, prop.Attr["x"], x)
	intAttrTest(t, prop.Attr["y"], y)
	intAttrTest(t, prop.Attr["inner_radius"], ir)
	intAttrTest(t, prop.Attr["outer_radius"], or)
	intAttrTest(t, prop.Attr["start_angle"], sa)
	intAttrTest(t, prop.Attr["end_angle"], ea)
}

func clockPropTest(t *testing.T, prop *props.Property, x, y int) {
	if prop.PropType != props.CLOCK_PROP {
		t.Errorf("Prop %s is not a clock prop", prop.Name)
		return
	}

	intAttrTest(t, prop.Attr["x"], x)
	intAttrTest(t, prop.Attr["y"], y)
}

func tickerPropTest(t *testing.T, prop *props.Property, x, y int, s ...string) {
	if prop.PropType != props.TICKER_PROP {
		t.Errorf("Prop %s is not a ticker prop", prop.Name)
		return
	}

	intAttrTest(t, prop.Attr["x"], x)
	intAttrTest(t, prop.Attr["y"], y)
	// check list store values
}

func graphPropTest(t *testing.T, prop *props.Property, x, y int) {
	if prop.PropType != props.GRAPH_PROP {
		t.Errorf("Prop %s is not a graph prop", prop.Name)
		return
	}

	intAttrTest(t, prop.Attr["x"], x)
	intAttrTest(t, prop.Attr["y"], y)
	// check list store values
}

func intAttrTest(t *testing.T, attr attribute.Attribute, val int) {
	intAttr, ok := attr.(*attribute.IntAttribute)
	if !ok {
		t.Errorf("Attr %v is not an int attr", attr)
		return
	}

	if intAttr.Value != val {
		t.Errorf("Int attr incorrect value (intAttr.Value = %d), expected %d", intAttr.Value, val)
	}
}

func stringAttrTest(t *testing.T, attr attribute.Attribute, val string) {
	stringAttr, ok := attr.(*attribute.StringAttribute)
	if !ok {
		t.Errorf("Attr %v is not an string attr", attr)
		return
	}

	if stringAttr.Value != val {
		t.Errorf("String attr incorrect value (stringAttr.Value = %s), expected %s", stringAttr.Value, val)
	}
}

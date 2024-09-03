package pages

import (
	"chroma-viz/library/geometry"
	"testing"
)

func TestImportShow(t *testing.T) {
	show := NewShow()
	show.ImportShow("test_show.json")

	if len(show.Pages) != 3 {
		t.Errorf("Incorrect number of pages (len(show.Pages) = %d)", len(show.Pages))
	}

	for _, page := range show.Pages {
		switch page.Title {
		case "Teal Box":
			rectPropTest(t, page.Rect[0], 50, 100, 850, 180)
			circlePropTest(t, page.Circle[0], 99, 90, 30, 75, 45, 315)
			textPropTest(t, page.Text[0], 190, 100, "Lower Frame")
			textPropTest(t, page.Text[1], 190, 30, "A lower frame subtitle")
		case "Clock Box":
			rectPropTest(t, page.Rect[0], 50, 965, 810, 65)
			circlePropTest(t, page.Circle[0], 35, 32, 0, 25, 0, 360)
			rectPropTest(t, page.Rect[1], 310, 0, 10, 65)
			textPropTest(t, page.Text[0], 180, 15, "0")
			textPropTest(t, page.Text[1], 20, 15, "TEAM")
			rectPropTest(t, page.Rect[2], 76, 0, 10, 65)
			textPropTest(t, page.Text[2], 180, 15, "0")
			textPropTest(t, page.Text[3], 20, 15, "TEAM")
			rectPropTest(t, page.Rect[3], 540, 0, 10, 65)
			//clockPropTest(t, page.Clock[0], 110, 15)
			textPropTest(t, page.Text[4], 30, 15, "Q1")
		case "Ticker":
			rectPropTest(t, page.Rect[0], 0, 0, 1920, 75)
			rectPropTest(t, page.Rect[1], 1700, 25, 400, 100)
			//tickerPropTest(t, page.Ticker[0], 25, 20, "Hello there", "world")
		default:
			t.Errorf("Unknown page %s", page.Title)
		}
	}
}

func rectPropTest(t *testing.T, rect *geometry.Rectangle, x, y, w, h int) {
	assert(t, rect.RelX.Value, x, "Rectangle rel_x")
	assert(t, rect.RelY.Value, y, "Rectangle rel_y")
	assert(t, rect.Width.Value, w, "Rectangle width")
	assert(t, rect.Height.Value, h, "Rectangle height")
}

func textPropTest(t *testing.T, string *geometry.Text, x, y int, text string) {
	assert(t, string.RelX.Value, x, "Text rel_x")
	assert(t, string.RelY.Value, y, "Text rel_y")
	assert(t, string.String.Value, text, "Text string")
}

func circlePropTest(t *testing.T, c *geometry.Circle, x, y, ir, or, sa, ea int) {
	assert(t, c.RelX.Value, x, "Circle rel_x")
	assert(t, c.RelY.Value, y, "Circle rel_y")
	assert(t, c.InnerRadius.Value, ir, "Circle inner_radius")
	assert(t, c.OuterRadius.Value, or, "Circle outer_radius")
	assert(t, c.StartAngle.Value, sa, "Circle start_angle")
	assert(t, c.EndAngle.Value, ea, "Circle end_angle")
}

func clockPropTest(t *testing.T, c *geometry.Clock, x, y int) {
	assert(t, c.RelX.Value, x, "Clock rel_x")
	assert(t, c.RelY.Value, x, "Clock rel_y")
}

func tickerPropTest(t *testing.T, ticker *geometry.List, x, y int, s ...string) {
	assert(t, ticker.RelX.Value, x, "Ticker rel_x")
	assert(t, ticker.RelY.Value, y, "Ticker rel_y")
	// check list store values

}

func assert[T comparable](t *testing.T, v1 T, v2 T, s string) {
	if v1 == v2 {
		return
	}

	t.Errorf("%s: expected %v, received %v", s, v2, v1)
}

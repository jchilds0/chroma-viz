package hub

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/templates"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type DataBase struct {
	db        *sql.DB
	templates map[int64]*templates.Template
	assets    map[int64]Asset
	stmt      map[string]*sql.Stmt
	lock      sync.Mutex
}

type Templates struct {
	NumTemplates int
	Templates    []*templates.Template
}

func NewDataBase(numTemp int, username, password string) (hub *DataBase, err error) {
	hub = &DataBase{}
	hub.assets = make(map[int64]Asset, 128)
	hub.templates = make(map[int64]*templates.Template, 128)

	hub.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/?multiStatements=true", username, password))
	if err != nil {
		err = fmt.Errorf("Error opening database: %s", err)
		return
	}

	err = hub.db.Ping()
	if err != nil {
		err = fmt.Errorf("Error connecting to database: %s", err)
		return
	}

	hub.stmt = make(map[string]*sql.Stmt, 20)
	return
}

func (hub *DataBase) ImportSchema(filename string) (err error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	s := string(buf)

	_, err = hub.db.Exec(s)
	return
}

func (hub *DataBase) SelectDatabase(name, username, password string) (err error) {
	hub.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", username, password, name))

	for geo, s := range stmts {
		hub.stmt[geo], err = hub.db.Prepare(s)
		if err != nil {
			err = fmt.Errorf("geometry %s: %s", geo, err)
			return
		}
	}

	return
}

func (hub *DataBase) CleanDB() {
	_, err := hub.db.Exec("DELETE FROM template")
	if err != nil {
		Logger("Error clearing db: %s", err)
	}

	return
}

type TemplateHeader struct {
	TemplateID int
	Title      string
	Layer      int
}

func (hub *DataBase) TempIDs() (temps []TemplateHeader, err error) {
	q := `
        SELECT t.templateID, t.Name, t.Layer
        FROM template t;
    `

	rows, err := hub.db.Query(q)
	if err != nil {
		return
	}

	temps = make([]TemplateHeader, 0, 1024)

	var (
		tempID int
		title  string
		layer  int
	)
	for rows.Next() {
		err = rows.Scan(&tempID, &title, &layer)
		if err != nil {
			return
		}

		temps = append(temps, TemplateHeader{TemplateID: tempID, Title: title, Layer: layer})
	}

	return
}

func (hub *DataBase) randomTemplate(tempID int64, numGeo int) (err error) {
	err = hub.AddTemplate(tempID, "Template "+strconv.FormatInt(tempID, 10), 0)
	if err != nil {
		err = fmt.Errorf("Adding template %d: %s", tempID, err)
		return
	}

	geos := []string{geometry.GEO_RECT, geometry.GEO_CIRCLE, geometry.GEO_TEXT}

	for j := 1; j < numGeo; j++ {
		geoIndex := rand.Int() % len(geos)

		geo := geometry.Geometry{
			GeometryID: j,
			Name:       geos[geoIndex],
			GeoType:    geos[geoIndex],
			Visible:    j < 10,
		}

		geo.RelX.Value = rand.Int() % 2000
		geo.RelY.Value = rand.Int() % 2000

		switch geos[geoIndex] {
		case geometry.GEO_RECT:
			rect := geometry.NewRectangle(geo)
			rect.Width.Value = rand.Int() % 1000
			rect.Height.Value = rand.Int() % 1000
			rect.Rounding.Value = rand.Int() % 300
			rect.Color.Red = float64(rand.Int()%255) / 255
			rect.Color.Green = float64(rand.Int()%255) / 255
			rect.Color.Blue = float64(rand.Int()%255) / 255
			rect.Color.Alpha = float64(rand.Int()%255) / 255
			err = hub.AddRectangle(tempID, *rect)

		case geometry.GEO_CIRCLE:
			circle := geometry.NewCircle(geo)
			circle.InnerRadius.Value = rand.Int() % 200
			circle.OuterRadius.Value = rand.Int() % 200
			circle.StartAngle.Value = rand.Int() % 10
			circle.EndAngle.Value = rand.Int() % 200
			circle.Color.Red = float64(rand.Int()%255) / 255
			circle.Color.Green = float64(rand.Int()%255) / 255
			circle.Color.Blue = float64(rand.Int()%255) / 255
			circle.Color.Alpha = float64(rand.Int()%255) / 255
			err = hub.AddCircle(tempID, *circle)

		case geometry.GEO_TEXT:
			text := geometry.NewText(geo)
			text.String.Value = "some text"
			text.Scale.Value = 1.0
			text.Color.Red = float64(rand.Int()%255) / 255
			text.Color.Green = float64(rand.Int()%255) / 255
			text.Color.Blue = float64(rand.Int()%255) / 255
			text.Color.Alpha = float64(rand.Int()%255) / 255
			err = hub.AddText(tempID, *text)

		}

		if j%3 == 0 {
			var tempFrame *templates.Keyframe
			if j%2 == 0 {
				tempFrame = templates.NewKeyFrame(1, j, "rel_x", false)
			} else {
				tempFrame = templates.NewKeyFrame(1, j, "rel_y", false)
			}

			startFrame := templates.NewSetFrame(*tempFrame, rand.Float64()*2000)
			hub.AddSetFrame(tempID, *startFrame)

			tempFrame.FrameNum = 2

			endFrame := templates.NewUserFrame(*tempFrame)
			hub.AddUserFrame(tempID, *endFrame)
		}

		if err != nil {
			return
		}
	}

	return
}

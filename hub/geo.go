package hub

import (
	"chroma-viz/library/templates"
	"fmt"
)

func (hub *DataBase) addGeometry(tempID int64, geo templates.Geometry) (geoID int64, err error) {
	q := `
        INSERT INTO geometry VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	result, err := hub.db.Exec(q, tempID, geo.GeoNum, geo.Name, geo.GeoType,
		geo.PropType, geo.RelX, geo.RelY, geo.Parent)
	if err != nil {
		return
	}

	geoID, err = result.LastInsertId()
	return
}

func (hub *DataBase) AddRectangle(tempID int64, rect templates.Rectangle) (err error) {
	q := `
        INSERT INTO rectangle VALUES (?, ?, ?, ?, ?);
    `

	geoID, err := hub.addGeometry(tempID, rect.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, geoID, rect.Width, rect.Height, rect.Rounding, rect.Color)
	return
}

func (hub *DataBase) AddText(tempID int64, text templates.Text) (err error) {
	q := `
        INSERT INTO text VALUES (?, ?, ?);
    `

	geoID, err := hub.addGeometry(tempID, text.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, geoID, text.Text, text.Color)
	return
}

func (hub *DataBase) AddCircle(tempID int64, circle templates.Circle) (err error) {
	q := `
        INSERT INTO circle VALUES (?, ?, ?, ?, ?, ?);
    `

	geoID, err := hub.addGeometry(tempID, circle.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, geoID, circle.InnerRadius, circle.OuterRadius, circle.StartAngle, circle.EndAngle, circle.Color)
	return
}

func (hub *DataBase) AddAsset(tempID int64, a templates.Asset) (err error) {
	q := `
        INSERT INTO asset VALUES (?, ?, ?, ?, ?);
    `

	geoID, err := hub.addGeometry(tempID, a.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, geoID, a.Dir, a.Name, a.ID, a.Scale)
	return
}

func (hub *DataBase) GetGeometry(geoID int64) (geo templates.Geometry, err error) {
	q := `
        SELECT g.geoNum, g.name, g.propType, g.geoType, g.rel_x, g.rel_y, g.parent
        FROM geometry g 
        WHERE g.geometryID = ?;
    `

	row := hub.db.QueryRow(q, geoID)

	err = row.Scan(&geo.GeoNum, &geo.Name, &geo.PropType, &geo.GeoType, &geo.RelX, &geo.RelY, &geo.Parent)
	return
}

func (hub *DataBase) GetRectangles(temp *templates.Template) (err error) {
	q := `
        SELECT r.geometryID, r.width, r.height, r.rounding, r.color
        FROM rectangle r 
        INNER JOIN geometry g
        ON r.geometryID = g.geometryID
        WHERE g.templateID = ?;
    `

	rows, err := hub.db.Query(q, temp.TempID)
	if err != nil {
		return
	}

	var geoID int64
	var geo templates.Geometry
	var width, height, rounding int
	var color string

	for rows.Next() {
		err = rows.Scan(&geoID, &width, &height, &rounding, &color)
		if err != nil {
			return
		}

		geo, err = hub.GetGeometry(geoID)
		if err != nil {
			return
		}

		rect := templates.NewRectangle(geo, width, height, rounding, color)
		temp.Rectangle = append(temp.Rectangle, *rect)
	}

	return
}

func (hub *DataBase) GetCircles(temp *templates.Template) (err error) {
	q := `
        SELECT c.geometryID, c.inner_radius, c.outer_radius, c.start_angle, c.end_angle, c.color
        FROM circle c
        INNER JOIN geometry g
        ON c.geometryID = g.geometryID
        WHERE g.templateID = ?;
    `

	rows, err := hub.db.Query(q, temp.TempID)

	var inner, outer, start, end int
	var color string
	var geoID int64
	var geo templates.Geometry

	for rows.Next() {
		err = rows.Scan(&geoID, &inner, &outer, &start, &end, &color)
		if err != nil {
			return
		}

		geo, err = hub.GetGeometry(geoID)
		if err != nil {
			return
		}

		c := templates.NewCircle(geo, inner, outer, start, end, color)
		temp.Circle = append(temp.Circle, *c)
	}

	return
}

func (hub *DataBase) GetTexts(temp *templates.Template) (err error) {
	q := `
        SELECT t.geometryID, t.text, t.color
        FROM text t
        INNER JOIN geometry g
        ON g.geometryID = t.geometryID
        WHERE g.templateID = ?;
    `

	rows, err := hub.db.Query(q, temp.TempID)
	if err != nil {
		return
	}

	var text, color string
	var geoID int64
	var geo templates.Geometry
	for rows.Next() {
		err = rows.Scan(&geoID, &text, &color)
		if err != nil {
			err = fmt.Errorf("Text: %s", err)
		}

		geo, err = hub.GetGeometry(geoID)
		if err != nil {
			return
		}

		t := templates.NewText(geo, text, color)
		temp.Text = append(temp.Text, *t)
	}

	return
}

func (hub *DataBase) GetAssets(temp *templates.Template) (err error) {
	q := `
        SELECT a.geometryID, a.directory, a.name, a.assetID, a.scale
        FROM asset a 
        INNER JOIN geometry g 
        ON a.geometryID = g.geometryID 
        WHERE g.templateID = ?;
    `

	rows, err := hub.db.Query(q, temp.TempID)
	if err != nil {
		return
	}

	var (
		geoID, assetID int64
		geo            templates.Geometry
		dir, name      string
		scale          float64
	)

	for rows.Next() {
		err = rows.Scan(&geoID, &dir, &name, &assetID, &scale)
		if err != nil {
			return
		}

		geo, err = hub.GetGeometry(geoID)
		if err != nil {
			return
		}

		a := templates.NewAsset(geo, name, dir, int(assetID), scale)
		temp.Asset = append(temp.Asset, *a)
	}

	return
}

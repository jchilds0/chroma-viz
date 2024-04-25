package hub

import (
	"chroma-viz/library/templates"
	"fmt"
)

func (hub *DataBase) addGeometry(tempID int64, geo templates.Geometry) (err error) {
	q := `
        INSERT INTO geometry VALUES (?, ?, ?, ?, ?, ?, ?, ?);
    `

	_, err = hub.db.Exec(q, geo.GeoID, tempID, geo.Name, geo.GeoType,
		geo.PropType, geo.RelX, geo.RelY, geo.Parent)
	return
}

func (hub *DataBase) AddRectangle(tempID int64, rect templates.Rectangle) (err error) {
	q := `
        INSERT INTO rectangle VALUES (?, ?, ?, ?, ?, ?);
    `

	err = hub.addGeometry(tempID, rect.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, rect.GeoID, tempID, rect.Width, rect.Height, rect.Rounding, rect.Color)
	return
}

func (hub *DataBase) AddText(tempID int64, text templates.Text) (err error) {
	q := `
        INSERT INTO text VALUES (?, ?, ?, ?);
    `

	err = hub.addGeometry(tempID, text.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, text.GeoID, tempID, text.Text, text.Color)
	return
}

func (hub *DataBase) AddCircle(tempID int64, circle templates.Circle) (err error) {
	q := `
        INSERT INTO circle VALUES (?, ?, ?, ?, ?, ?, ?);
    `

	err = hub.addGeometry(tempID, circle.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, circle.GeoID, tempID, circle.InnerRadius, circle.OuterRadius, circle.StartAngle, circle.EndAngle, circle.Color)
	return
}

func (hub *DataBase) AddAsset(geoID int64, dir, name string, id int) (err error) {
	q := `
        INSERT INTO asset VALUES (?, ?, ?, ?);
    `

	_, err = hub.db.Exec(q, geoID, dir, name, id)
	return
}

func (hub *DataBase) GetGeometry(temp *templates.Template) (err error) {
	q := `
        SELECT g.geometryID, g.Name, g.prop_type, g.geo_type, g.rel_x, g.rel_y, g.parent
        FROM geometry g 
        WHERE g.templateID = ?;
    `

	rows, err := hub.db.Query(q, temp.TempID)
	if err != nil {
		return
	}

	var geo templates.Geometry
	for rows.Next() {
		err = rows.Scan(&geo.GeoID, &geo.Name, &geo.PropType, &geo.GeoType, &geo.RelX, &geo.RelY, &geo.Parent)
		if err != nil {
			return
		}

		switch geo.GeoType {
		case templates.GEO_RECT:
			var rect templates.Rectangle
			rect, err = hub.GetRectangle(temp.TempID, geo)

			temp.Rectangle = append(temp.Rectangle, rect)

		case templates.GEO_CIRCLE:
			var circle templates.Circle
			circle, err = hub.GetCircle(temp.TempID, geo)

			temp.Circle = append(temp.Circle, circle)

		case templates.GEO_TEXT:
			var text templates.Text
			text, err = hub.GetText(temp.TempID, geo)

			temp.Text = append(temp.Text, text)

		default:
			Logger("Geo type %s not implemented in chroma hub", templates.GeoName[geo.GeoType])
			continue
		}

		if err != nil {
			return
		}
	}

	return
}

func (hub *DataBase) GetRectangle(tempID int64, geo templates.Geometry) (rect templates.Rectangle, err error) {
	q := `
        SELECT g.width, g.height, g.rounding, g.color
        FROM rectangle g 
        WHERE g.geometryID = ? AND g.templateID = ?;
    `

	row := hub.db.QueryRow(q, geo.GeoID, tempID)

	var width, height, rounding int
	var color string

	err = row.Scan(&width, &height, &rounding, &color)
	if err != nil {
		err = fmt.Errorf("Rectangle: %s", err)
	}

	rect = *templates.NewRectangle(geo, width, height, rounding, color)

	return
}

func (hub *DataBase) GetCircle(tempID int64, geo templates.Geometry) (c templates.Circle, err error) {
	q := `
        SELECT g.inner_radius, g.outer_radius, g.start_angle, g.end_angle, g.color
        FROM circle g 
        WHERE g.geometryID = ? AND g.templateID = ?;
    `

	row := hub.db.QueryRow(q, geo.GeoID, tempID)

	var inner, outer, start, end int
	var color string

	err = row.Scan(&inner, &outer, &start, &end, &color)
	if err != nil {
		err = fmt.Errorf("Circle: %s", err)
	}

	c = *templates.NewCircle(geo, inner, outer, start, end, color)

	return
}

func (hub *DataBase) GetText(tempID int64, geo templates.Geometry) (t templates.Text, err error) {
	q := `
        SELECT g.text, g.color
        FROM text g 
        WHERE g.geometryID = ? AND g.templateID = ?;
    `

	row := hub.db.QueryRow(q, geo.GeoID, tempID)

	var text, color string

	err = row.Scan(&text, &color)
	if err != nil {
		err = fmt.Errorf("Text: %s", err)
	}

	t = *templates.NewText(geo, text, color)

	return
}

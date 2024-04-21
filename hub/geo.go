package hub

import (
	"chroma-viz/library/props"
	"chroma-viz/library/templates"
	"log"
)

func (hub *DataBase) addGeometry(tempID int64, geo templates.Geometry) (geo_id int64, err error) {
	q := `
        INSERT INTO geometry VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	result, err := hub.db.Exec(q,
		tempID, geo.Name, geo.GeoType, geo.RelX, geo.RelY,
		geo.Color[0], geo.Color[1], geo.Color[2], geo.Color[3], geo.Parent)
	if err != nil {
		return
	}

	geo_id, err = result.LastInsertId()
	return
}

func (hub *DataBase) AddRectangle(tempID int64, rect templates.Rectangle) (err error) {
	q := `
        INSERT INTO rectangle VALUES (?, ?, ?, ?);
    `

	geoID, err := hub.addGeometry(tempID, rect.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, geoID, rect.Width, rect.Height, rect.Rounding)
	return
}

func (hub *DataBase) AddText(tempID int64, text templates.Text) (err error) {
	q := `
        INSERT INTO text VALUES (?, ?);
    `

	geoID, err := hub.addGeometry(tempID, text.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, geoID, text.Text)
	return
}

func (hub *DataBase) AddCircle(tempID int64, circle templates.Circle) (err error) {
	q := `
        INSERT INTO circle VALUES (?, ?, ?, ?, ?);
    `

	geoID, err := hub.addGeometry(tempID, circle.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, geoID, circle.InnerRadius, circle.OuterRadius, circle.StartAngle, circle.EndAngle)
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
        SELECT g.geometryID, g.Name, g.Type, g.rel_x, g.rel_y, g.color_r, g.color_g, g.color_b, g.color_a, g.parent
        FROM geometry g 
        WHERE g.templateID = ?;
    `

	rows, err := hub.db.Query(q, temp.TempID)
	if err != nil {
		return
	}

	var (
		id    int64
		geoID int
		geo   templates.Geometry
	)
	for rows.Next() {
		err = rows.Scan(&id, &geo.Name, &geo.GeoType, &geo.RelX, &geo.RelY,
			&geo.Color[0], &geo.Color[1], &geo.Color[2], &geo.Color[3], &geo.Parent)
		if err != nil {
			return
		}

		temp.AddGeometry(geo.Name, geoID, geo.GeoType, nil)
		temp.Geometry[geoID].Visible["rel_x"] = true
		temp.Geometry[geoID].Visible["rel_y"] = true
		temp.Geometry[geoID].Visible["color"] = true

		var attr templates.GeometryEncoder
		switch geo.GeoType {
		case props.RECT_PROP:
			attr, err = hub.GetRectangle(id, geo)
		case props.CIRCLE_PROP:
			attr, err = hub.GetCircle(id, geo)
		case props.TEXT_PROP:
			attr, err = hub.GetText(id, geo)
		default:
			log.Printf("Prop type %s not implemented in chroma hub", props.PropType(geo.GeoType))
			geoID++
			continue
		}
		if err != nil {
			return
		}

		for name, a := range attr.Encode() {
			temp.Geometry[geoID].Attr[name] = a
			temp.Geometry[geoID].Visible[name] = true
		}

		geoID++
	}

	return
}

func (hub *DataBase) GetRectangle(geoID int64, geo templates.Geometry) (rect *templates.Rectangle, err error) {
	q := `
        SELECT g.width, g.height, g.rounding
        FROM rectangle g 
        WHERE g.geometryID = ?;
    `

	row := hub.db.QueryRow(q, geoID)
	rect = templates.NewRectangle(geo, 0, 0, 0)

	err = row.Scan(&rect.Width, &rect.Height, &rect.Rounding)
	return
}

func (hub *DataBase) GetCircle(geoID int64, geo templates.Geometry) (c *templates.Circle, err error) {
	q := `
        SELECT g.inner_radius, g.outer_radius, g.start_angle, g.end_angle
        FROM circle g 
        WHERE g.geometryID = ?;
    `

	row := hub.db.QueryRow(q, geoID)
	c = templates.NewCircle(geo, 0, 0, 0, 0)

	err = row.Scan(&c.InnerRadius, &c.OuterRadius, &c.StartAngle, &c.EndAngle)
	return
}

func (hub *DataBase) GetText(geoID int64, geo templates.Geometry) (t *templates.Text, err error) {
	q := `
        SELECT g.text 
        FROM text g 
        WHERE g.geometryID = ?;
    `

	row := hub.db.QueryRow(q, geoID)
	t = templates.NewText(geo, "")

	err = row.Scan(&t.Text)
	return
}

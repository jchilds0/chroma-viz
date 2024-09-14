package hub

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/geometry"
	"chroma-viz/library/templates"
	"database/sql"
	"fmt"
)

func (hub *DataBase) addGeometry(tempID int64, geo geometry.Geometry) (geoID int64, err error) {
	q := `
        INSERT INTO geometry VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	result, err := hub.db.Exec(q, tempID, geo.GeometryID, geo.Name, geo.GeoType,
		geo.RelX.Value, geo.RelY.Value, geo.Parent.Value, geo.Mask.Value)
	if err != nil {
		return
	}

	geoID, err = result.LastInsertId()
	return
}

func (hub *DataBase) AddRectangle(tempID int64, rect geometry.Rectangle) (err error) {
	q := `
        INSERT INTO rectangle VALUES (?, ?, ?, ?, ?);
    `

	geoID, err := hub.addGeometry(tempID, rect.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, geoID, rect.Width.Value, rect.Height.Value,
		rect.Rounding.Value, rect.Color.ToString())
	return
}

func (hub *DataBase) AddText(tempID int64, text geometry.Text) (err error) {
	q := `
        INSERT INTO text VALUES (?, ?, ?, ?, ?);
    `

	geoID, err := hub.addGeometry(tempID, text.Geometry)
	if err != nil {
		return
	}

	fontFace := ""
	_, err = hub.db.Exec(q, geoID, text.String.Value, text.Scale.Value, fontFace, text.Color.ToString())
	return
}

func (hub *DataBase) AddCircle(tempID int64, circle geometry.Circle) (err error) {
	q := `
        INSERT INTO circle VALUES (?, ?, ?, ?, ?, ?);
    `

	geoID, err := hub.addGeometry(tempID, circle.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, geoID, circle.InnerRadius.Value, circle.OuterRadius.Value,
		circle.StartAngle.Value, circle.EndAngle.Value, circle.Color.ToString())
	return
}

func (hub *DataBase) AddAsset(tempID int64, a geometry.Image) (err error) {
	q := `
        INSERT INTO asset VALUES (?, ?, ?, ?, ?);
    `

	geoID, err := hub.addGeometry(tempID, a.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, geoID, a.Image.Directory(), a.Image.Name, a.Image.Value, a.Scale.Value)
	return
}

func (hub *DataBase) AddClock(tempID int64, c *geometry.Clock) (err error) {
	q := `
        INSERT INTO clock VALUES (?, ?, ?);
    `

	geoID, err := hub.addGeometry(tempID, c.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, geoID, c.Scale.Value, c.Color.ToString())
	return
}

func (hub *DataBase) AddPolygon(tempID int64, p geometry.Polygon) (err error) {
	qPoly := `
	       INSERT INTO polygon VALUES (?, ?);
	   `

	qPoint := `
	       INSERT INTO point VALUES (?, ?, ?, ?);
	   `

	geoID, err := hub.addGeometry(tempID, p.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(qPoly, geoID, p.Color.ToString())
	if err != nil {
		return
	}

	for i := range p.Polygon.PosX {
		posX := p.Polygon.PosX[i]
		posY := p.Polygon.PosY[i]

		_, err = hub.db.Exec(qPoint, geoID, i, posX, posY)
		if err != nil {
			return
		}
	}

	return
}

func (hub *DataBase) AddList(tempID int64, l geometry.List) (err error) {
	qList := `
        INSERT INTO list VALUES (?, ?, ?, ?);
    `

	qRows := `
        INSERT INTO row VALUES (?, ?, ?);
    `

	geoID, err := hub.addGeometry(tempID, l.Geometry)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(qList, geoID, l.Color.ToString(), l.String.Selected, l.Scale.Value)
	if err != nil {
		return
	}

	for index, row := range l.String.Rows {
		_, err = hub.db.Exec(qRows, geoID, index, row.ToString())
		if err != nil {
			return
		}
	}

	return
}

func (hub *DataBase) GetGeometry(geoID int64) (geo geometry.Geometry, err error) {
	q := `
        SELECT g.geoNum, g.name, g.geoType, g.rel_x, g.rel_y, g.parent, g.mask
        FROM geometry g 
        WHERE g.geometryID = ?;
    `

	row := hub.db.QueryRow(q, geoID)

	var geoNum int
	var name, geoType string
	var relX, relY, parent, mask int
	err = row.Scan(&geoNum, &name, &geoType, &relX, &relY, &parent, &mask)
	if err != nil {
		return
	}

	geo = geometry.NewGeometry(geoNum, name, geoType)
	geo.RelX.Value = relX
	geo.RelY.Value = relY
	geo.Parent.Value = parent
	geo.Mask.Value = mask

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

	var geo geometry.Geometry
	var geoID int64
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

		rect := geometry.NewRectangle(geo)
		rect.Width.Value = width
		rect.Height.Value = height
		rect.Rounding.Value = rounding

		err = rect.Color.FromString(color)
		if err != nil {
			return
		}

		temp.Rectangle = append(temp.Rectangle, rect)
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

	var geo geometry.Geometry
	var inner, outer, start, end int
	var color string
	var geoID int64

	for rows.Next() {
		err = rows.Scan(&geoID, &inner, &outer, &start, &end, &color)
		if err != nil {
			return
		}

		geo, err = hub.GetGeometry(geoID)
		if err != nil {
			return
		}

		c := geometry.NewCircle(geo)
		c.InnerRadius.Value = inner
		c.OuterRadius.Value = outer
		c.StartAngle.Value = start
		c.EndAngle.Value = end
		err = c.Color.FromString(color)
		if err != nil {
			return
		}

		temp.Circle = append(temp.Circle, c)
	}

	return
}

func (hub *DataBase) GetTexts(temp *templates.Template) (err error) {
	q := `
        SELECT t.geometryID, t.text, t.fontSize, t.color
        FROM text t
        INNER JOIN geometry g
        ON g.geometryID = t.geometryID
        WHERE g.templateID = ?;
    `

	rows, err := hub.db.Query(q, temp.TempID)
	if err != nil {
		return
	}

	var geo geometry.Geometry
	var text, color string
	var geoID int64
	var scale float64

	for rows.Next() {
		err = rows.Scan(&geoID, &text, &scale, &color)
		if err != nil {
			return
		}

		geo, err = hub.GetGeometry(geoID)
		if err != nil {
			return
		}

		t := geometry.NewText(geo)
		t.Scale.Value = scale
		t.String.Value = text
		t.Color.FromString(color)
		if err != nil {
			return
		}

		temp.Text = append(temp.Text, t)
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
		geo            geometry.Geometry
		geoID, assetID int64
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

		a := geometry.NewImage(geo)
		a.Image.Value = int(assetID)
		a.Scale.Value = scale

		temp.Image = append(temp.Image, a)
	}

	return
}

func (hub *DataBase) GetPolygons(temp *templates.Template) (err error) {
	q := `
        SELECT p.geometryID, p.color
        FROM polygon p
        INNER JOIN geometry g 
        ON p.geometryID = g.geometryID 
        WHERE g.templateID = ?;
    `

	rows, err := hub.db.Query(q, temp.TempID)
	if err != nil {
		return
	}

	var (
		geo   geometry.Geometry
		geoID int64
		color string
	)

	for rows.Next() {
		err = rows.Scan(&geoID, &color)
		if err != nil {
			return
		}

		geo, err = hub.GetGeometry(geoID)
		if err != nil {
			return
		}

		poly := geometry.NewPolygon(geo)
		poly.Color.FromString(color)
		temp.Polygon = append(temp.Polygon, poly)
	}

	err = hub.GetPolyPoints(temp)
	return
}

func (hub *DataBase) GetPolyPoints(temp *templates.Template) (err error) {
	q := `
        SELECT point.pointID, point.pos_x, point.pos_y
        FROM point
        INNER JOIN polygon
        INNER JOIN geometry g
        ON point.geometryID = polygon.geometryID
        AND polygon.geometryID = g.geometryID
        WHERE g.geoNum = ?
        AND g.templateID = ?;
    `

	var (
		rows       *sql.Rows
		pointIndex int
		posX, posY int
	)
	for _, poly := range temp.Polygon {
		rows, err = hub.db.Query(q, poly.GeometryID, temp.TempID)
		if err != nil {
			return
		}

		for rows.Next() {
			err = rows.Scan(&pointIndex, &posX, &posY)
			if err != nil {
				return
			}

			poly.Polygon.AddPoint(pointIndex, posX, posY)
		}
	}

	return
}

func (hub *DataBase) GetClocks(temp *templates.Template) (err error) {
	q := `
        SELECT c.geometryID, c.scale, c.color
        FROM clock c
        INNER JOIN geometry g
        ON g.geometryID = c.geometryID
        WHERE g.templateID = ?;
    `

	rows, err := hub.db.Query(q, temp.TempID)
	if err != nil {
		return
	}

	var geo geometry.Geometry
	var color string
	var geoID int64
	var scale float64

	for rows.Next() {
		err = rows.Scan(&geoID, &scale, &color)
		if err != nil {
			return
		}

		geo, err = hub.GetGeometry(geoID)
		if err != nil {
			return
		}

		c := geometry.NewClock(geo)
		temp.Clock = append(temp.Clock, c)
	}

	return
}

func (hub *DataBase) GetLists(temp *templates.Template) (err error) {
	q := `
        SELECT l.geometryID, l.color, l.single_row, l.scale
        FROM list l
        INNER JOIN geometry g
        ON g.geometryID = l.geometryID
        WHERE g.templateID = ?;
    `

	rows, err := hub.db.Query(q, temp.TempID)
	if err != nil {
		return
	}

	var geoID int64
	var color string
	var singleRow bool
	var scale float64
	var geo geometry.Geometry

	for rows.Next() {
		err = rows.Scan(&geoID, &color, &singleRow, &scale)
		if err != nil {
			return
		}

		geo, err = hub.GetGeometry(geoID)
		if err != nil {
			return
		}

		l := geometry.NewList(geo)
		l.Color.FromString(color)
		l.String.Selected = singleRow
		l.Scale.Value = scale

		temp.List = append(temp.List, l)
	}

	err = hub.GetListRows(temp)
	return
}

func (hub *DataBase) GetListRows(temp *templates.Template) (err error) {
	q := `
        SELECT r.rowID, r.row
        FROM row r
        INNER JOIN list l 
        INNER JOIN geometry g
        ON r.geometryID = l.geometryID
        AND l.geometryID = g.geometryID
        WHERE g.geoNum = ?
        AND g.templateID = ?;
    `

	var (
		rows  *sql.Rows
		row   string
		index int
	)

	for _, list := range temp.List {
		rows, err = hub.db.Query(q, list.GeometryID, temp.TempID)
		if err != nil {
			fmt.Println(err)
			return
		}

		for rows.Next() {
			err = rows.Scan(&index, &row)
			if err != nil {
				return
			}

			rowAttr := attribute.NewListRow(row)
			list.AddRow(index, rowAttr)
		}
	}

	return
}

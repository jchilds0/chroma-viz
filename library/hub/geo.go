package hub

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/templates"
	"fmt"
)

func (hub *DataBase) addGeometry(tempID int64, geo geometry.Geometry) (geoID int64, err error) {
	q := `
        INSERT INTO geometry VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?);
    `

	result, err := hub.db.Exec(q, tempID, geo.GeometryID, geo.Name, geo.GeoType,
		geo.RelX, geo.RelY, geo.Parent, geo.Mask)
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

	_, err = hub.db.Exec(q, geoID, a.Image.Directory(), a.Image.Name, a.Image.AssetID(), a.Scale.Value)
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

		temp.Rect = append(temp.Rect, rect)
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
			err = fmt.Errorf("Text: %s", err)
		}

		geo, err = hub.GetGeometry(geoID)
		if err != nil {
			return
		}

		t := geometry.NewText(geo)
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
        SELECT p.geometryID, p.point_index, p.pos_x, p.pos_y
        FROM polygon p
        INNER JOIN geometry g 
        ON p.geometryID = g.geometryID 
        WHERE g.templateID = ?;
    `

	rows, err := hub.db.Query(q, temp.TempID)
	if err != nil {
		return
	}

	pointsX := make(map[int64]map[int]int, 128)
	pointsY := make(map[int64]map[int]int, 128)

	var (
		geo        geometry.Geometry
		geoID      int64
		pointIndex int
		posX, posY int
	)

	for rows.Next() {
		err = rows.Scan(&geoID, &pointIndex, &posX, &posY)
		if err != nil {
			return
		}

		if _, ok := pointsX[geoID]; !ok {
			pointsX[geoID] = make(map[int]int, 128)
			pointsY[geoID] = make(map[int]int, 128)
		}

		pointsX[geoID][pointIndex] = posX
		pointsY[geoID][pointIndex] = posY
	}

	for geoID := range pointsX {
		geo, err = hub.GetGeometry(geoID)
		if err != nil {
			return
		}

		poly := geometry.NewPolygon(geo, len(pointsX[geoID])+10)
		for i := range len(pointsX[geoID]) {
			if _, ok := pointsX[geoID][i]; !ok {
				err = fmt.Errorf("Missing point %d for geometry %d", i, geoID)
				return
			}

			if _, ok := pointsY[geoID][i]; !ok {
				err = fmt.Errorf("Missing point %d for geometry %d", i, geoID)
				return
			}

			poly.Polygon.PosX = append(poly.Polygon.PosX, pointsX[geoID][i])
			poly.Polygon.PosY = append(poly.Polygon.PosY, pointsY[geoID][i])
		}

		temp.Poly = append(temp.Poly, poly)
	}

	return
}

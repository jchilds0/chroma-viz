package hub

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/geometry"
	"chroma-viz/library/templates"
	"database/sql"
	"fmt"
)

func (hub *DataBase) addGeometry(tempID int64, geo geometry.Geometry) (geoID int64, err error) {
	result, err := hub.stmt[GEOMETRY_INSERT].Exec(tempID, geo.GeometryID, geo.Name, geo.GeoType,
		geo.RelX.Value, geo.RelY.Value, geo.Parent.Value, geo.Mask.Value, geo.Visible)
	if err != nil {
		return
	}

	geoID, err = result.LastInsertId()
	return
}

func (hub *DataBase) AddRectangle(tempID int64, rect geometry.Rectangle) (err error) {
	geoID, err := hub.addGeometry(tempID, rect.Geometry)
	if err != nil {
		return
	}

	_, err = hub.stmt[RECTANGLE_INSERT].Exec(geoID,
		rect.Width.Value, rect.Height.Value, rect.Rounding.Value,
		rect.Color.Red, rect.Color.Green, rect.Color.Blue, rect.Color.Alpha)
	return
}

func (hub *DataBase) AddText(tempID int64, text geometry.Text) (err error) {
	geoID, err := hub.addGeometry(tempID, text.Geometry)
	if err != nil {
		return
	}

	fontFace := ""
	_, err = hub.stmt[TEXT_INSERT].Exec(geoID,
		text.String.Value, text.Scale.Value, fontFace,
		text.Color.Red, text.Color.Blue, text.Color.Green, text.Color.Alpha)
	return
}

func (hub *DataBase) AddCircle(tempID int64, circle geometry.Circle) (err error) {
	geoID, err := hub.addGeometry(tempID, circle.Geometry)
	if err != nil {
		return
	}

	_, err = hub.stmt[CIRCLE_INSERT].Exec(geoID, circle.InnerRadius.Value, circle.OuterRadius.Value,
		circle.StartAngle.Value, circle.EndAngle.Value,
		circle.Color.Red, circle.Color.Green, circle.Color.Blue, circle.Color.Alpha)
	return
}

func (hub *DataBase) AddAssetGeo(tempID int64, a geometry.Image) (err error) {
	geoID, err := hub.addGeometry(tempID, a.Geometry)
	if err != nil {
		return
	}

	_, err = hub.stmt[ASSET_GEO_INSERT].Exec(geoID, a.Image.Value, a.Scale.Value)
	return
}

func (hub *DataBase) AddClock(tempID int64, c *geometry.Clock) (err error) {
	geoID, err := hub.addGeometry(tempID, c.Geometry)
	if err != nil {
		return
	}

	_, err = hub.stmt[CLOCK_INSERT].Exec(geoID, c.Scale.Value,
		c.Color.Red, c.Color.Green, c.Color.Blue, c.Color.Alpha)
	return
}

func (hub *DataBase) AddPolygon(tempID int64, p geometry.Polygon) (err error) {
	geoID, err := hub.addGeometry(tempID, p.Geometry)
	if err != nil {
		return
	}

	_, err = hub.stmt[POLYGON_INSERT].Exec(geoID, p.Color.Red, p.Color.Green, p.Color.Blue, p.Color.Alpha)
	if err != nil {
		return
	}

	for i := range p.Polygon.PosX {
		posX := p.Polygon.PosX[i]
		posY := p.Polygon.PosY[i]

		_, err = hub.stmt[POINT_INSERT].Exec(geoID, i, posX, posY)
		if err != nil {
			return
		}
	}

	return
}

func (hub *DataBase) AddList(tempID int64, l geometry.List) (err error) {
	geoID, err := hub.addGeometry(tempID, l.Geometry)
	if err != nil {
		return
	}

	_, err = hub.stmt[LIST_INSERT].Exec(geoID,
		l.Color.Red, l.Color.Green, l.Color.Blue, l.Color.Alpha,
		l.String.Selected, l.Scale.Value)
	if err != nil {
		return
	}

	for index, row := range l.String.Rows {
		_, err = hub.stmt[ROW_INSERT].Exec(geoID, index, row.ToString())
		if err != nil {
			return
		}
	}

	return
}

func (hub *DataBase) GetGeometry(tempID int64) (geos map[int64]geometry.Geometry, err error) {
	geos = make(map[int64]geometry.Geometry, 128)
	rows, err := hub.stmt[GEOMETRY_SELECT].Query(tempID)
	if err != nil {
		return
	}

	var (
		geoID                    int64
		geoNum                   int
		name, geoType            string
		relX, relY, parent, mask int
		visible                  bool
	)
	for rows.Next() {
		err = rows.Scan(&geoID, &geoNum, &name, &geoType, &relX, &relY, &parent, &mask, &visible)
		if err != nil {
			return
		}

		geo := geometry.NewGeometry(geoNum, name, geoType, visible)
		geo.RelX.Value = relX
		geo.RelY.Value = relY
		geo.Parent.Value = parent
		geo.Mask.Value = mask

		geos[geoID] = geo
	}

	return
}

func (hub *DataBase) GetRectangles(temp *templates.Template, geos map[int64]geometry.Geometry) (err error) {
	rows, err := hub.stmt[RECTANGLE_SELECT].Query(temp.TempID)
	if err != nil {
		return
	}

	var (
		geoID                   int64
		width, height, rounding int
		red, green, blue, alpha float64
	)

	for rows.Next() {
		err = rows.Scan(&geoID, &width, &height, &rounding, &red, &green, &blue, &alpha)
		if err != nil {
			return
		}

		geo, ok := geos[geoID]
		if !ok {
			return fmt.Errorf("Missing geometry %d", geoID)
		}

		rect := geometry.NewRectangle(geo)
		rect.Width.Value = width
		rect.Height.Value = height
		rect.Rounding.Value = rounding
		rect.Color.Red = red
		rect.Color.Green = green
		rect.Color.Blue = blue
		rect.Color.Alpha = alpha

		temp.Rectangle = append(temp.Rectangle, rect)
	}

	return
}

func (hub *DataBase) GetCircles(temp *templates.Template, geos map[int64]geometry.Geometry) (err error) {
	rows, err := hub.stmt[CIRCLE_SELECT].Query(temp.TempID)
	if err != nil {
		return
	}

	var (
		geoID                    int64
		inner, outer, start, end int
		red, green, blue, alpha  float64
	)

	for rows.Next() {
		err = rows.Scan(&geoID, &inner, &outer, &start, &end, &red, &green, &blue, &alpha)
		if err != nil {
			return
		}

		geo, ok := geos[geoID]
		if !ok {
			return fmt.Errorf("Missing geometry %d", geoID)
		}

		c := geometry.NewCircle(geo)
		c.InnerRadius.Value = inner
		c.OuterRadius.Value = outer
		c.StartAngle.Value = start
		c.EndAngle.Value = end
		c.Color.Red = red
		c.Color.Green = green
		c.Color.Blue = blue
		c.Color.Alpha = alpha

		temp.Circle = append(temp.Circle, c)
	}

	return
}

func (hub *DataBase) GetTexts(temp *templates.Template, geos map[int64]geometry.Geometry) (err error) {
	rows, err := hub.stmt[TEXT_SELECT].Query(temp.TempID)
	if err != nil {
		return
	}

	var (
		geoID                   int64
		text                    string
		scale                   float64
		red, green, blue, alpha float64
	)

	for rows.Next() {
		err = rows.Scan(&geoID, &text, &scale, &red, &green, &blue, &alpha)
		if err != nil {
			return
		}

		geo, ok := geos[geoID]
		if !ok {
			return fmt.Errorf("Missing geometry %d", geoID)
		}

		t := geometry.NewText(geo)
		t.Scale.Value = scale
		t.String.Value = text
		t.Color.Red = red
		t.Color.Green = green
		t.Color.Blue = blue
		t.Color.Alpha = alpha

		temp.Text = append(temp.Text, t)
	}

	return
}

func (hub *DataBase) GetAssetGeos(temp *templates.Template, geos map[int64]geometry.Geometry) (err error) {
	rows, err := hub.stmt[ASSET_GEO_SELECT].Query(temp.TempID)
	if err != nil {
		return
	}

	var (
		geoID, assetID int64
		scale          float64
	)

	for rows.Next() {
		err = rows.Scan(&geoID, &assetID, &scale)
		if err != nil {
			return
		}

		geo, ok := geos[geoID]
		if !ok {
			return fmt.Errorf("Missing geometry %d", geoID)
		}

		a := geometry.NewImage(geo)
		a.Image.Value = int(assetID)
		a.Scale.Value = scale

		temp.Image = append(temp.Image, a)
	}

	return
}

func (hub *DataBase) GetPolygons(temp *templates.Template, geos map[int64]geometry.Geometry) (err error) {
	rows, err := hub.stmt[POLYGON_SELECT].Query(temp.TempID)
	if err != nil {
		return
	}

	var (
		geoID                   int64
		red, green, blue, alpha float64
	)

	for rows.Next() {
		err = rows.Scan(&geoID, &red, &green, &blue, &alpha)
		if err != nil {
			return
		}

		geo, ok := geos[geoID]
		if !ok {
			return fmt.Errorf("Missing geometry %d", geoID)
		}

		poly := geometry.NewPolygon(geo)
		poly.Color.Red = red
		poly.Color.Green = green
		poly.Color.Blue = blue
		poly.Color.Alpha = alpha

		temp.Polygon = append(temp.Polygon, poly)
	}

	err = hub.GetPolyPoints(temp)
	return
}

func (hub *DataBase) GetPolyPoints(temp *templates.Template) (err error) {
	var (
		rows       *sql.Rows
		pointIndex int
		posX, posY int
	)
	for _, poly := range temp.Polygon {
		rows, err = hub.stmt[POINT_SELECT].Query(poly.GeometryID, temp.TempID)
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

func (hub *DataBase) GetClocks(temp *templates.Template, geos map[int64]geometry.Geometry) (err error) {
	rows, err := hub.stmt[CLOCK_SELECT].Query(temp.TempID)
	if err != nil {
		return
	}

	var (
		geoID                   int64
		scale                   float64
		red, green, blue, alpha float64
	)

	for rows.Next() {
		err = rows.Scan(&geoID, &scale, &red, &green, &blue, &alpha)
		if err != nil {
			return
		}

		geo, ok := geos[geoID]
		if !ok {
			return fmt.Errorf("Missing geometry %d", geoID)
		}

		c := geometry.NewClock(geo)
		c.Scale.Value = scale
		c.Color.Red = red
		c.Color.Green = green
		c.Color.Blue = blue
		c.Color.Alpha = alpha
		temp.Clock = append(temp.Clock, c)
	}

	return
}

func (hub *DataBase) GetLists(temp *templates.Template, geos map[int64]geometry.Geometry) (err error) {
	rows, err := hub.stmt[LIST_SELECT].Query(temp.TempID)
	if err != nil {
		return
	}

	var (
		geoID                   int64
		singleRow               bool
		scale                   float64
		red, green, blue, alpha float64
	)

	for rows.Next() {
		err = rows.Scan(&geoID, &red, &green, &blue, &alpha, &singleRow, &scale)
		if err != nil {
			return
		}

		geo, ok := geos[geoID]
		if !ok {
			return fmt.Errorf("Missing geometry %d", geoID)
		}

		l := geometry.NewList(geo)
		l.Color.Red = red
		l.Color.Green = green
		l.Color.Blue = blue
		l.Color.Alpha = alpha
		l.String.Selected = singleRow
		l.Scale.Value = scale

		temp.List = append(temp.List, l)
	}

	err = hub.GetListRows(temp)
	return
}

func (hub *DataBase) GetListRows(temp *templates.Template) (err error) {
	var (
		rows  *sql.Rows
		row   string
		index int
	)

	for _, list := range temp.List {
		rows, err = hub.stmt[ROW_SELECT].Query(list.GeometryID, temp.TempID)
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

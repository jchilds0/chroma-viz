package hub

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/props"
	"chroma-viz/library/templates"
	"fmt"
	"log"
)

type Geometry struct {
	GeoID   int64
	Name    string
	GeoType int
	RelX    int
	RelY    int
	Color   [4]byte
	Parent  int
}

func NewGeometry(name string, geoType, rel_x, rel_y int, r, g, b, a byte, parent int) *Geometry {
	geo := &Geometry{
		Name:    name,
		GeoType: geoType,
		RelX:    rel_x,
		RelY:    rel_y,
		Color:   [4]byte{r, g, b, a},
		Parent:  parent,
	}

	return geo
}

func (hub *DataBase) AddGeometry(tempID int64, geo Geometry) (geo_id int64, err error) {
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

func (hub *DataBase) AddRectangle(geoID int64, width int, height int, rounding int) (err error) {
	q := `
        INSERT INTO rectangle VALUES (?, ?, ?, ?);
    `

	_, err = hub.db.Exec(q, geoID, width, height, rounding)
	return
}

func (hub *DataBase) AddText(geoID int64, text string) (err error) {
	q := `
        INSERT INTO text VALUES (?, ?);
    `

	_, err = hub.db.Exec(q, geoID, text)
	return
}

func (hub *DataBase) AddCircle(geoID int64, innerRadius, outerRadius, startAngle, endAngle int) (err error) {
	q := `
        INSERT INTO circle VALUES (?, ?, ?, ?, ?);
    `

	_, err = hub.db.Exec(q, geoID, innerRadius, outerRadius, startAngle, endAngle)
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
		id     int64
		name   string
		typed  int
		rel_x  int
		rel_y  int
		r      byte
		g      byte
		b      byte
		a      byte
		parent int
	)
	geo_id := 0
	for rows.Next() {
		err = rows.Scan(&id, &name, &typed, &rel_x, &rel_y, &r, &g, &b, &a, &parent)
		if err != nil {
			return
		}

		temp.AddGeometry(name, geo_id, typed, nil)
		temp.Geometry[geo_id].Visible["rel_x"] = true
		temp.Geometry[geo_id].Visible["rel_y"] = true
		temp.Geometry[geo_id].Visible["color"] = true

		attribute.SetIntValue(temp.Geometry[geo_id].Attr["rel_x"], rel_x)
		attribute.SetIntValue(temp.Geometry[geo_id].Attr["rel_y"], rel_y)

		color := temp.Geometry[geo_id].Attr["color"].(*attribute.ColorAttribute)
		color.Red = float64(r) / 255
		color.Green = float64(g) / 255
		color.Blue = float64(b) / 255
		color.Alpha = float64(a) / 255

		switch typed {
		case props.RECT_PROP:
			err = hub.GetRectangle(id, temp.Geometry[geo_id])
		case props.CIRCLE_PROP:
			err = hub.GetCircle(id, temp.Geometry[geo_id])
		case props.TEXT_PROP:
			err = hub.GetText(id, temp.Geometry[geo_id])
		default:
			log.Printf("Prop type %s not implemented in chroma hub", props.PropType(typed))
			geo_id++
			continue
		}

		if err != nil {
			return
		}

		geo_id++
	}

	return
}

func (hub *DataBase) GetRectangle(geoID int64, prop *props.Property) (err error) {
	q := `
        SELECT g.width, g.height, g.rounding
        FROM rectangle g 
        WHERE g.geometryID = ?;
    `

	row := hub.db.QueryRow(q, geoID)

	var (
		width    int
		height   int
		rounding int
	)
	err = row.Scan(&width, &height, &rounding)
	if err != nil {
		return
	}

	err = attribute.SetIntValue(prop.Attr["width"], width)
	if err != nil {
		return
	}

	err = attribute.SetIntValue(prop.Attr["height"], height)
	if err != nil {
		return
	}

	err = attribute.SetIntValue(prop.Attr["rounding"], width)
	if err != nil {
		return
	}

	prop.Visible["width"] = true
	prop.Visible["height"] = true
	prop.Visible["rounding"] = true
	return
}

func (hub *DataBase) GetCircle(geoID int64, prop *props.Property) (err error) {
	q := `
        SELECT g.inner_radius, g.outer_radius, g.start_angle, g.end_angle
        FROM circle g 
        WHERE g.geometryID = ?;
    `

	row := hub.db.QueryRow(q, geoID)

	var (
		inner_radius int
		outer_radius int
		start_angle  int
		end_angle    int
	)
	err = row.Scan(&inner_radius, &outer_radius, &start_angle, &end_angle)
	if err != nil {
		return
	}

	err = attribute.SetIntValue(prop.Attr["inner_radius"], inner_radius)
	if err != nil {
		return
	}

	err = attribute.SetIntValue(prop.Attr["outer_radius"], inner_radius)
	if err != nil {
		return
	}

	err = attribute.SetIntValue(prop.Attr["start_angle"], inner_radius)
	if err != nil {
		return
	}

	err = attribute.SetIntValue(prop.Attr["end_angle"], inner_radius)
	if err != nil {
		return
	}

	prop.Visible["inner_radius"] = true
	prop.Visible["outer_radius"] = true
	prop.Visible["start_angle"] = true
	prop.Visible["end_angle"] = true

	return
}

func (hub *DataBase) GetText(geoID int64, prop *props.Property) (err error) {
	q := `
        SELECT g.text 
        FROM text g 
        WHERE g.geometryID = ?;
    `

	row := hub.db.QueryRow(q, geoID)

	var text string
	err = row.Scan(&text)
	if err != nil {
		return
	}

	if prop.Attr["string"] == nil {
		err = fmt.Errorf("String attr missing")
		return
	}

	attr, ok := prop.Attr["string"].(*attribute.StringAttribute)
	if !ok {
		err = fmt.Errorf("Attribute is not a StringAttribute")
		return
	}

	attr.Value = text
	prop.Visible["string"] = true
	return
}

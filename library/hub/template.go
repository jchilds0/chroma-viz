package hub

import (
	"chroma-viz/library/templates"
	"fmt"
	"reflect"
)

func (hub *DataBase) ImportTemplate(temp templates.Template) (err error) {
	err = hub.AddTemplate(temp.TempID, temp.Title, temp.Layer)
	if err != nil {
		return
	}

	err = importPointer(temp.TempID, temp.Rectangle, hub.AddRectangle)
	if err != nil {
		return
	}

	err = importPointer(temp.TempID, temp.Circle, hub.AddCircle)
	if err != nil {
		return
	}

	err = importPointer(temp.TempID, temp.Text, hub.AddText)
	if err != nil {
		return
	}

	err = importPointer(temp.TempID, temp.Image, hub.AddAsset)
	if err != nil {
		return
	}

	err = importPointer(temp.TempID, temp.Polygon, hub.AddPolygon)
	if err != nil {
		return
	}

	err = importPointer(temp.TempID, temp.List, hub.AddList)
	if err != nil {
		return
	}

	err = importStruct(temp.TempID, temp.Clock, hub.AddClock)
	if err != nil {
		return
	}

	err = importStruct(temp.TempID, temp.SetFrame, hub.AddSetFrame)
	if err != nil {
		return
	}

	err = importStruct(temp.TempID, temp.BindFrame, hub.AddBindFrame)
	if err != nil {
		return
	}

	err = importStruct(temp.TempID, temp.UserFrame, hub.AddUserFrame)
	if err != nil {
		return
	}

	hub.Templates[temp.TempID] = &temp

	return
}

func importPointer[T any](tempID int64, geos []*T, f func(tempID int64, geo T) error) (err error) {
	for _, geo := range geos {
		if geo == nil {
			continue
		}

		err = f(tempID, *geo)
		if err != nil {
			err = fmt.Errorf("Error adding %s: %s", reflect.TypeOf(geo).String(), err)
			return
		}
	}

	return
}

func importStruct[T any](tempID int64, keys []T, f func(tempID int64, geo T) error) (err error) {
	for _, key := range keys {
		err = f(tempID, key)
		if err != nil {
			err = fmt.Errorf("Error adding %s: %s", reflect.TypeOf(key).String(), err)
			return
		}
	}

	return
}

func (hub *DataBase) AddTemplate(tempID int64, name string, layer int) (err error) {
	// TODO: run as a transaction
	deleteTemp := `
        DELETE FROM template WHERE templateID = ?;
    `
	_, err = hub.db.Exec(deleteTemp, tempID)
	if err != nil {
		Logger(err.Error())
	}

	addTemp := `
        INSERT INTO template VALUES (?, ?, ?);
    `

	_, err = hub.db.Exec(addTemp, tempID, name, layer)
	return
}

func (hub *DataBase) GetTemplate(tempID int64) (temp *templates.Template, err error) {
	// Needs template versioning to update correctly
	// temp, ok := hub.Templates[tempID]
	// if ok {
	// 	return
	// }

	tempQuery := `
        SELECT t.Name, t.Layer, COUNT(*)
        FROM template t
        INNER JOIN geometry g 
        ON g.templateID = t.templateID
        WHERE t.templateID = ?;
    `
	var (
		name    string
		layer   int
		num_geo int
	)

	row := hub.db.QueryRow(tempQuery, tempID)
	if err = row.Scan(&name, &layer, &num_geo); err != nil {
		return
	}

	temp = templates.NewTemplate(name, tempID, layer, num_geo, 0)
	err = hub.GetRectangles(temp)
	if err != nil {
		err = fmt.Errorf("Rectangle: %s", err)
		return
	}

	err = hub.GetCircles(temp)
	if err != nil {
		err = fmt.Errorf("Circle: %s", err)
		return
	}

	err = hub.GetTexts(temp)
	if err != nil {
		err = fmt.Errorf("Text: %s", err)
		return
	}

	err = hub.GetAssets(temp)
	if err != nil {
		err = fmt.Errorf("Assets: %s", err)
		return
	}

	err = hub.GetClocks(temp)
	if err != nil {
		err = fmt.Errorf("Clock: %s", err)
		return
	}

	err = hub.GetLists(temp)
	if err != nil {
		err = fmt.Errorf("Lists: %s", err)
		return
	}

	err = hub.GetSetFrame(temp)
	if err != nil {
		err = fmt.Errorf("Set Frame: %s", err)
		return
	}

	err = hub.GetBindFrames(temp)
	if err != nil {
		err = fmt.Errorf("Bind Frame: %s", err)
		return
	}

	err = hub.GetUserFrames(temp)
	if err != nil {
		err = fmt.Errorf("User Frame: %s", err)
		return
	}

	hub.Templates[tempID] = temp
	return
}

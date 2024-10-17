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

	hub.lock.Lock()
	hub.templates[int(temp.TempID)] = &temp
	hub.lock.Unlock()

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
	_, err = hub.stmt[TEMPLATE_DELETE].Exec(tempID)
	if err != nil {
		Logger(err.Error())
	}

	_, err = hub.stmt[TEMPLATE_INSERT].Exec(tempID, name, layer)
	return
}

func (hub *DataBase) GetTemplate(tempID int64) (temp *templates.Template, err error) {
	hub.lock.Lock()
	temp, ok := hub.templates[int(tempID)]
	hub.lock.Unlock()
	if ok {
		return
	}

	var (
		name    string
		layer   int
		num_geo int
	)

	row := hub.stmt[TEMPLATE_SELECT].QueryRow(tempID)
	if err = row.Scan(&name, &layer, &num_geo); err != nil {
		err = fmt.Errorf("Template %d: %s", tempID, err)
		return
	}

	temp = templates.NewTemplate(name, tempID, layer, num_geo, 0)

	geos, err := hub.GetGeometry(temp.TempID)
	if err != nil {
		err = fmt.Errorf("Geometry: %s", err)
		return
	}

	err = hub.GetRectangles(temp, geos)
	if err != nil {
		err = fmt.Errorf("Rectangle: %s", err)
		return
	}

	err = hub.GetCircles(temp, geos)
	if err != nil {
		err = fmt.Errorf("Circle: %s", err)
		return
	}

	err = hub.GetTexts(temp, geos)
	if err != nil {
		err = fmt.Errorf("Text: %s", err)
		return
	}

	err = hub.GetAssets(temp, geos)
	if err != nil {
		err = fmt.Errorf("Assets: %s", err)
		return
	}

	err = hub.GetPolygons(temp, geos)
	if err != nil {
		err = fmt.Errorf("Polygons: %s", err)
		return
	}

	err = hub.GetClocks(temp, geos)
	if err != nil {
		err = fmt.Errorf("Clock: %s", err)
		return
	}

	err = hub.GetLists(temp, geos)
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

	hub.lock.Lock()
	hub.templates[int(tempID)] = temp
	hub.lock.Unlock()

	return
}

package hub

import (
	"chroma-viz/library/templates"
	"fmt"
)

func (hub *DataBase) ImportTemplate(temp templates.Template) (err error) {
	err = hub.AddTemplate(temp.TempID, temp.Title, temp.Layer)
	if err != nil {
		return
	}

	for _, rect := range temp.Rectangle {
		err = hub.AddRectangle(temp.TempID, rect)

		if err != nil {
			err = fmt.Errorf("Error adding rect: %s", err)
			return
		}
	}

	for _, text := range temp.Text {
		err = hub.AddText(temp.TempID, text)

		if err != nil {
			err = fmt.Errorf("Error adding text: %s", err)
			return
		}
	}

	for _, circle := range temp.Circle {
		err = hub.AddCircle(temp.TempID, circle)

		if err != nil {
			err = fmt.Errorf("Error adding circle: %s", err)
			return
		}
	}

	for _, asset := range temp.Asset {
		err = hub.AddAsset(temp.TempID, asset)

		if err != nil {
			err = fmt.Errorf("Error adding asset: %s", err)
			return
		}
	}

	for _, bind := range temp.BindFrame {
		err = hub.AddBindFrame(temp.TempID, bind)
		if err != nil {
			err = fmt.Errorf("Error adding bind frame: %s", err)
			return
		}
	}

	for _, set := range temp.SetFrame {
		err = hub.AddSetFrame(temp.TempID, set)

		if err != nil {
			err = fmt.Errorf("Error adding set frame: %s", err)
			return
		}
	}

	for _, user := range temp.UserFrame {
		err = hub.AddUserFrame(temp.TempID, user)
		if err != nil {
			err = fmt.Errorf("Error adding user frame: %s", err)
			return
		}
	}

	hub.Templates[temp.TempID] = &temp

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

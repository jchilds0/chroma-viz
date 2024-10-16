package hub

import (
	"chroma-viz/library/templates"
)

func (hub *DataBase) addKeyframe(tempID int64, frame templates.Keyframe) (frameID int64, err error) {
	q := `
        INSERT INTO keyframe VALUES (NULL, ?, ?, ?, ?, ?, ?);
    `

	result, err := hub.db.Exec(q, tempID, frame.FrameNum, frame.GeoID, frame.GeoAttr, frame.Type, frame.Expand)
	if err != nil {
		return
	}

	frameID, err = result.LastInsertId()
	return
}

func (hub *DataBase) AddBindFrame(tempID int64, frame templates.BindFrame) (err error) {
	q := `
        INSERT INTO bindFrame VALUES (?, ?);
    `

	frameID, err := hub.addKeyframe(tempID, frame.Keyframe)
	if err != nil {
		return
	}

	bindFrameID, err := hub.addKeyframe(tempID, frame.Bind)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, frameID, bindFrameID)
	return
}

func (hub *DataBase) AddUserFrame(tempID int64, frame templates.UserFrame) (err error) {
	q := `
        INSERT INTO userFrame VALUES (?);
    `

	frameID, err := hub.addKeyframe(tempID, frame.Keyframe)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, frameID)
	return
}

func (hub *DataBase) AddSetFrame(tempID int64, frame templates.SetFrame) (err error) {
	q := `
        INSERT INTO setFrame VALUES (?, ?);
    `

	frameID, err := hub.addKeyframe(tempID, frame.Keyframe)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, frameID, frame.Value)
	return
}

func (hub *DataBase) GetKeyframe(frameID int64) (frame templates.Keyframe, err error) {
	q := `
        SELECT k.frameNum, k.geoNum, k.attr, k.type, k.expand 
        FROM keyframe k 
        WHERE k.frameID = ?;
    `

	row := hub.db.QueryRow(q, frameID)

	err = row.Scan(&frame.FrameNum, &frame.GeoID, &frame.GeoAttr, &frame.Type, &frame.Expand)
	return
}

func (hub *DataBase) GetUserFrames(temp *templates.Template) (err error) {
	q := `
        SELECT u.frameID 
        FROM userFrame u 
        INNER JOIN keyframe k 
        ON k.frameID = u.frameID 
        WHERE k.templateID = ?;
    `

	rows, err := hub.db.Query(q, temp.TempID)
	if err != nil {
		return
	}

	var frameID int64
	var frame templates.Keyframe
	for rows.Next() {
		err = rows.Scan(&frameID)
		if err != nil {
			return
		}

		frame, err = hub.GetKeyframe(frameID)
		if err != nil {
			return
		}

		userFrame := templates.NewUserFrame(frame)
		temp.UserFrame = append(temp.UserFrame, *userFrame)
	}

	return
}

func (hub *DataBase) GetSetFrame(temp *templates.Template) (err error) {
	q := `
        SELECT s.frameID, s.value
        FROM setFrame s 
        INNER JOIN keyframe k 
        ON k.frameID = s.frameID 
        WHERE k.templateID = ?;
    `

	rows, err := hub.db.Query(q, temp.TempID)
	if err != nil {
		return
	}

	var (
		frame   templates.Keyframe
		frameID int64
		value   float64
	)
	for rows.Next() {
		err = rows.Scan(&frameID, &value)
		if err != nil {
			return
		}

		frame, err = hub.GetKeyframe(frameID)
		if err != nil {
			return
		}

		setFrame := templates.NewSetFrame(frame, value)
		temp.SetFrame = append(temp.SetFrame, *setFrame)
	}

	return
}
func (hub *DataBase) GetBindFrames(temp *templates.Template) (err error) {
	q := `
        SELECT b.frameID, b.bindFrameID
        FROM bindFrame b
        INNER JOIN keyframe k 
        ON k.frameID = b.frameID 
        WHERE k.templateID = ?;
    `

	rows, err := hub.db.Query(q, temp.TempID)
	if err != nil {
		return
	}

	var frameID, bindFrameID int64
	var frame, bindFrame templates.Keyframe
	for rows.Next() {
		err = rows.Scan(&frameID, &bindFrameID)
		if err != nil {
			return
		}

		frame, err = hub.GetKeyframe(frameID)
		if err != nil {
			return
		}

		bindFrame, err = hub.GetKeyframe(bindFrameID)
		if err != nil {
			return
		}

		bind := templates.NewBindFrame(frame, bindFrame)
		temp.BindFrame = append(temp.BindFrame, *bind)
	}

	return
}

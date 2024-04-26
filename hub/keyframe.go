package hub

import "chroma-viz/library/templates"

func (hub *DataBase) addKeyframe(tempID int64, frame templates.Keyframe) (frameID int64, err error) {
	q := `
        INSERT INTO keyframe VALUES (NULL, ?, ?, ?, ?, ?, ?, ?);
    `

	result, err := hub.db.Exec(q, tempID, frame.FrameNum, frame.GeoID, frame.GeoAttr, frame.Type, frame.Mask, frame.Expand)
	if err != nil {
		return
	}

	frameID, err = result.LastInsertId()
	return
}

func (hub *DataBase) AddBindFrame(tempID int64, frame templates.BindFrame) (err error) {
	q := `
        INSERT INTO bindFrame VALUES (?, ?, ?);
    `

	frameID, err := hub.addKeyframe(tempID, frame.Keyframe)
	if err != nil {
		return
	}

	bindFrameID, err := hub.addKeyframe(tempID, frame.Bind)
	if err != nil {
		return
	}

	_, err = hub.db.Exec(q, frameID, bindFrameID, tempID)
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

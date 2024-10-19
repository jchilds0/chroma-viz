package hub

import (
	"chroma-viz/library/templates"
	"fmt"
)

func (hub *DataBase) addKeyframe(tempID int64, frame templates.Keyframe) (frameID int64, err error) {
	result, err := hub.stmt[KEYFRAME_INSERT].Exec(tempID, frame.FrameNum, frame.GeoID, frame.GeoAttr, frame.Type, frame.Expand)
	if err != nil {
		return
	}

	frameID, err = result.LastInsertId()
	return
}

func (hub *DataBase) AddBindFrame(tempID int64, frame templates.BindFrame) (err error) {
	frameID, err := hub.addKeyframe(tempID, frame.Keyframe)
	if err != nil {
		return
	}

	bindFrameID, err := hub.addKeyframe(tempID, frame.Bind)
	if err != nil {
		return
	}

	_, err = hub.stmt[BIND_INSERT].Exec(frameID, bindFrameID)
	return
}

func (hub *DataBase) AddUserFrame(tempID int64, frame templates.UserFrame) (err error) {
	frameID, err := hub.addKeyframe(tempID, frame.Keyframe)
	if err != nil {
		return
	}

	_, err = hub.stmt[USER_INSERT].Exec(frameID)
	return
}

func (hub *DataBase) AddSetFrame(tempID int64, frame templates.SetFrame) (err error) {

	frameID, err := hub.addKeyframe(tempID, frame.Keyframe)
	if err != nil {
		return
	}

	_, err = hub.stmt[SET_INSERT].Exec(frameID, frame.Value)
	return
}

func (hub *DataBase) GetKeyframes(tempID int64) (frames map[int64]templates.Keyframe, err error) {
	rows, err := hub.stmt[KEYFRAME_SELECT].Query(tempID)
	if err != nil {
		return
	}

	var (
		frame   templates.Keyframe
		frameID int64
	)

	frames = make(map[int64]templates.Keyframe, 128)
	for rows.Next() {
		err = rows.Scan(&frameID, &frame.FrameNum, &frame.GeoID, &frame.GeoAttr, &frame.Type, &frame.Expand)
		if err != nil {
			return
		}

		frames[frameID] = frame
	}

	return
}

func (hub *DataBase) GetUserFrames(temp *templates.Template, frames map[int64]templates.Keyframe) (err error) {
	rows, err := hub.stmt[USER_SELECT].Query(temp.TempID)
	if err != nil {
		return
	}

	var frameID int64
	for rows.Next() {
		err = rows.Scan(&frameID)
		if err != nil {
			return
		}

		frame, ok := frames[frameID]
		if !ok {
			return fmt.Errorf("Missing Keyframe %d for user frame", frameID)
		}

		userFrame := templates.NewUserFrame(frame)
		temp.UserFrame = append(temp.UserFrame, *userFrame)
	}

	return
}

func (hub *DataBase) GetSetFrame(temp *templates.Template, frames map[int64]templates.Keyframe) (err error) {
	rows, err := hub.stmt[SET_SELECT].Query(temp.TempID)
	if err != nil {
		return
	}

	var (
		frameID int64
		value   float64
	)
	for rows.Next() {
		err = rows.Scan(&frameID, &value)
		if err != nil {
			return
		}

		frame, ok := frames[frameID]
		if !ok {
			return fmt.Errorf("Missing Keyframe %d for user frame", frameID)
		}

		setFrame := templates.NewSetFrame(frame, value)
		temp.SetFrame = append(temp.SetFrame, *setFrame)
	}

	return
}
func (hub *DataBase) GetBindFrames(temp *templates.Template, frames map[int64]templates.Keyframe) (err error) {
	rows, err := hub.stmt[BIND_SELECT].Query(temp.TempID)
	if err != nil {
		return
	}

	var frameID, bindFrameID int64
	for rows.Next() {
		err = rows.Scan(&frameID, &bindFrameID)
		if err != nil {
			return
		}

		frame, ok := frames[frameID]
		if !ok {
			return fmt.Errorf("Missing Keyframe %d for bind frame", frameID)
		}

		bindFrame, ok := frames[bindFrameID]
		if !ok {
			return fmt.Errorf("Missing Bind Keyframe %d for bind frame", frameID)
		}

		bind := templates.NewBindFrame(frame, bindFrame)
		temp.BindFrame = append(temp.BindFrame, *bind)
	}

	return
}

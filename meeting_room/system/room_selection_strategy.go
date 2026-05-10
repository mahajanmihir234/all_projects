package system

import "errors"

type RoomSelectionStrategy interface {
	SelectRoom(rooms []MeetingRoom, requiredCapacity int) (*MeetingRoom, error)
}

type FirstAvailableStrategy struct{}

func (s FirstAvailableStrategy) SelectRoom(rooms []MeetingRoom, requiredCapacity int) (*MeetingRoom, error) {
	if len(rooms) == 0 {
		return nil, errors.New("no room available")
	}
	for _, room := range rooms {
		if room.capacity >= requiredCapacity {
			return &room, nil
		}
	}

	return nil, errors.New("no room available for the given capacity")
}

package system

import "errors"

type BookingSystem struct {
	rooms                 map[string]*MeetingRoom
	meetings              map[string]*Meeting
	roomSelectionStrategy RoomSelectionStrategy
	roomMeetings          map[string][]*Meeting
}

func NewBookingSystem(strategy RoomSelectionStrategy) BookingSystem {
	return BookingSystem{
		rooms:                 map[string]*MeetingRoom{},
		meetings:              map[string]*Meeting{},
		roomSelectionStrategy: strategy,
		roomMeetings:          map[string][]*Meeting{},
	}
}

func (b *BookingSystem) AddRoom(r MeetingRoom) {
	b.rooms[r.roomId] = &r
	b.roomMeetings[r.roomId] = []*Meeting{}
}

func (b *BookingSystem) SetRoomSelectionStrategy(strategy RoomSelectionStrategy) {
	b.roomSelectionStrategy = strategy
}

func (b *BookingSystem) getAvailableRooms(t TimeSlot, capacity int) []MeetingRoom {
	availableRooms := []MeetingRoom{}
	for _, room := range b.rooms {
		if room.capacity < capacity {
			continue
		}
		isAvailable := true
		for _, meeting := range b.roomMeetings[room.roomId] {
			if meeting.timeSlot.Overlap(t) {
				isAvailable = false
			}
		}
		if isAvailable {
			availableRooms = append(availableRooms, *room)
		}
	}
	return availableRooms
}

func (b *BookingSystem) CancelMeeting(meetingId string) error {
	meeting, ok := b.meetings[meetingId]
	if !ok {
		return errors.New("meeting ID does not exist")
	}
	return meeting.Cancel()
}

func (b *BookingSystem) ScheduleMeeting(meetingId string, t TimeSlot, capacity int) (*Meeting, error) {
	availableRooms := b.getAvailableRooms(t, capacity)
	if len(availableRooms) == 0 {
		return nil, errors.New("no rooms available")
	}
	room, err := b.roomSelectionStrategy.SelectRoom(availableRooms, capacity)
	if err != nil {
		return nil, err
	}

	meeting := NewMeeting(meetingId, *room, t)
	return &meeting, nil
}

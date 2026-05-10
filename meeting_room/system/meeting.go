package system

import "errors"

type MeetingStatus string

const (
	SCHEDULED MeetingStatus = "SCHEDULED"
	CANCELLED MeetingStatus = "CANCELLED"
	COMPLETED MeetingStatus = "COMPLETED"
)

type TimeSlot struct {
	startTime int
	endTime   int
}

func (t TimeSlot) Overlap(other TimeSlot) bool {
	return t.endTime > other.startTime && t.startTime < other.endTime
}

type Meeting struct {
	id       string
	room     MeetingRoom
	timeSlot TimeSlot
	status   MeetingStatus
}

func NewMeeting(meetingId string, room MeetingRoom, timeSlot TimeSlot) Meeting {
	return Meeting{id: meetingId, room: room, timeSlot: timeSlot, status: SCHEDULED}
}

func (m *Meeting) Complete() error {
	switch m.status {
	case SCHEDULED:
		m.status = COMPLETED
		return nil
	case CANCELLED:
		return errors.New("cannot complete a cancelled meeting")
	case COMPLETED:
		return errors.New("cannot complete an already completed meeting")
	default:
		return errors.New("Unknown status")
	}
}

func (m *Meeting) Cancel() error {
	switch m.status {
	case SCHEDULED:
		m.status = CANCELLED
		return nil
	case CANCELLED:
		return errors.New("cannot cancel a cancelled meeting")
	case COMPLETED:
		return errors.New("cannot cancel an already completed meeting")
	default:
		return errors.New("Unknown status")
	}
}

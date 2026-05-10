package system

type RoomType string

const (
	CONFERENCE   RoomType = "CONFERENCE"
	BOARD_ROOM   RoomType = "BOARD_ROOM"
	HUDDLE_SPACE RoomType = "HUDDLE_SPACE"
)

type MeetingRoom struct {
	roomId   string
	roomType RoomType
	capacity int
}

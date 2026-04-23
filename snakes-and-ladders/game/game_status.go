package game

type GameStatus string

const (
	IN_PROGRESS GameStatus = "IN_PROGRESS"
	FINISHED    GameStatus = "FINISHED"
	NOT_STARTED GameStatus = "NOT_STARTED"
)

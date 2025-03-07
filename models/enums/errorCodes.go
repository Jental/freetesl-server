package enums

type ErrorCode int

const (
	ErrorCodePlayerHasMatch ErrorCode = 1
)

var ErrorCodeMessages = map[ErrorCode]string{
	ErrorCodePlayerHasMatch: "Player already have a match started",
}

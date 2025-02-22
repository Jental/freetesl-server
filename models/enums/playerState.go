package enums

type PlayerState byte

const (
	PlayerStateOffline            PlayerState = 1
	PlayerStateOnline             PlayerState = 2
	PlayerStateLookingForOpponent PlayerState = 3
	PlayerStateInMatch            PlayerState = 4
)

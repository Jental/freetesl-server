package dtos

type LoginResponseDTO struct {
	Valid    bool    `json:"valid"`
	Token    *string `json:"token"`
	PlayerID int     `json:"playerID"`
}

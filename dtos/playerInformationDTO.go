package dtos

type PlayerInformationDTO struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	AvatarName string `json:"avatarName"`
	State      byte   `json:"state"`
}

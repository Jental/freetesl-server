package dtos

type DeckDTO struct {
	ID         int                 `json:"id"`
	Name       string              `json:"name"`
	AvatarName string              `json:"avatarName"`
	Attributes []string            `json:"attributes"`
	Cards      []*CardWithCountDTO `json:"cards"`
}

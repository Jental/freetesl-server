package dtos

type LoginDTO struct {
	Login          string `json:"login"`
	PasswordSha512 string `json:"passwordSha512"`
}

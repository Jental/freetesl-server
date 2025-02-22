package dtos

type ListDTO[T interface{}] struct {
	Items []T `json:"items"`
}

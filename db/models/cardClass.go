package models

type CardClass struct {
	ID         byte
	Name       string
	Attributes []*Attribute
}

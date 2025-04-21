package models

type CardEffect struct {
	Name       string
	EffectID   byte
	Parameter0 *string // nullable
	Parameter1 *string // nullable
}

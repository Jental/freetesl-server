package dtos

type CardDTO struct {
	ID       int   `json:"id"`
	Power    int   `json:"power"`
	Health   int   `json:"health"`
	Cost     int   `json:"cost"`
	Type     byte  `json:"type"`
	Keywords []int `json:"keywords"` // we don't use []byte here because it's serialized as Base64
}

package models

type Neon struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Sizes  string `json:"sizes"`
	Colors string `json:"colors"`
	Cost   int    `json:"cost"`
	Theme  string `json:"theme"`
}

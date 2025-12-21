package model

type MediaImage struct {
	Src   string `json:"src"`
	Sizes string `json:"sizes,omitempty"`
	Type  string `json:"type,omitempty"`
}

type MediaMetadata struct {
	Title   string       `json:"title"`
	Artist  string       `json:"artist"`
	Album   string       `json:"album"`
	Artwork []MediaImage `json:"artwork,omitempty"`
}

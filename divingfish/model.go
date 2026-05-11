package divingfish

type MaiMaiMusicData struct {
	ID        string               `json:"id"`
	Title     string               `json:"title"`
	Type      string               `json:"type"`
	DS        []float64            `json:"ds"`
	Level     []string             `json:"level"`
	Cids      []int                `json:"cids"`
	BasicInfo MaiMaiMusicBasicInfo `json:"basic_info"`
}

type MaiMaiMusicBasicInfo struct {
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Genre       string `json:"genre"`
	BPM         int    `json:"bpm"`
	ReleaseDate string `json:"release_date"`
	From        string `json:"from"`
	IsNew       bool   `json:"is_new"`
}

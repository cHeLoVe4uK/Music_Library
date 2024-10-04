package models

// Модель песни, представляющая собой способ хранения сущности, используемой в нашей БД
type Song struct {
	Group       string   `json:"group"`
	Song        string   `json:"song"`
	ReleaseDate string   `json:"releaseDate,omitempty"`
	Text        []string `json:"text"`
	Link        string   `json:"link,omitempty"`
}

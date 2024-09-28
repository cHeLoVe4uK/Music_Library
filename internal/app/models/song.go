package models

// Модель, представляющая собой способ хранения сущности, используемой в нашей БД
type Song struct {
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}

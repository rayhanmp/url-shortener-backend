package models

type url struct {
	Original string `json:"original" gorm:"text; not null"`
	Short string `json:"short" gorm:"text; not null, unique"`
}


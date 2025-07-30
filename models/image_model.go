package models

type ImageModel struct {
	Model
	FileName string `gorm:"size:256" json:"file_name"`
	Path     string `gorm:"size:256" json:"path"`
	Size     int64  `json:"size"`
	Hash     string `gorm:"size:256" json:"hash"`
}

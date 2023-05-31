package model

type Video struct {
	ID        uint   `gorm:"primary_key"`
	Title     string `gorm:"type:varchar(255);not null"`
	Cover     string `gorm:"type:varchar(255);column:cover"`
	PlayUrl   string `gorm:"type:varchar(255);column:play_url"`
	CreatedAt string `gorm:"type:timestamp;column:created_at"`
}

package model

type Book struct {
	Id        uint64 `gorm:"primaryKey;auto_increment" json:"id"`
	Name      string `gorm:"" json:"name"`
	Author    string `gorm:"" json:"author"`
	CreatedAt int64  `gorm:"autoCreateTime:milli" json:"createdAt"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli" json:"-"`
}

package data

import "time"

type CommonConfig struct {
	ID       int64     `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

func (m *CommonConfig) TableName() string {
	return "common_config"
}

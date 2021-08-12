package model

type Index struct {
	ID       int64  `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	CreateTs int64  `json:"create_ts" gorm:"autoCreateTime"`
	UpdateTs int64  `json:"update_tx" gorm:"autoUpdateTime"`
}

func (m *Index) TableName() string {
	return "index"
}

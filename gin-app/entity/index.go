package entity

type ReqIndex struct {
	Name string `json:"name" binding:"required"`
}

type ResIndex struct {
	Words string `json:"words"`
}

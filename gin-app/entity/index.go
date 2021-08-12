package entity

type ReqIndex struct {
	Name string `form:"name" binding:"required"`
}

type ResIndex struct {
	Words string `json:"words"`
}

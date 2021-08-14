package entitydemo

type ReqGitUser struct {
	Name string `form:"name" binding:"required"`
}

type ResGitUser struct {
	Info string `json:"words"`
}

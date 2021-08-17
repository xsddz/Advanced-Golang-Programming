package entitydemo

// ReqGitUser 请求参数定义结构体
type ReqGitUser struct {
	Name string `form:"name" binding:"required"`
}

// ResGitUser 请求响应定义结构体
type ResGitUser struct {
	Info string `json:"words"`
}

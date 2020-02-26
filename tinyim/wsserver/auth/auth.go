package auth

import (
	"net/http"
)

// ClientInfo 客户端信息
type ClientInfo struct {
	Channel string
	UserID  string
}

// Login 认证请求
func Login(r *http.Request) (*ClientInfo, bool) {
	// fmt.Println("============ begin authorization request ================")

	// fmt.Println("method:\t", r.Method)
	// fmt.Println("url:\t", r.URL, r.URL.Path)
	// fmt.Println("proto:\t", r.Proto)
	// fmt.Println("header:\t", r.Header)
	// fmt.Println("param[auth]:\t", r.FormValue("auth"))
	// fmt.Println("remote addr:\t", r.RemoteAddr)
	// fmt.Println("request url:\t", r.RequestURI)

	// fmt.Println("room:", r.URL.Query().Get("room"))

	// fmt.Println("============ end authorization request ================")

	return &ClientInfo{
		Channel: r.URL.Query().Get("room"),
		UserID:  r.RemoteAddr,
	}, true
}

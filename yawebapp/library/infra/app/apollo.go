package app

import "github.com/shima-park/agollo"

var (
	apolloTable = make(map[string]*apollo)
)

type apollo struct {
	agollo.Agollo
}

// Apollo -
func Apollo() *apollo {
	if a, ok := apolloTable["Default"]; ok {
		return a
	}

	apolloTable["Default"] = &apollo{initAgollo()}
	return apolloTable["Default"]
}

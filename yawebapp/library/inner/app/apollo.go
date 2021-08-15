package app

import "github.com/shima-park/agollo"

var (
	apolloDriver = "agollo"
	apolloTable  = make(map[string]*apollo)
)

type apollo struct {
	agollo.Agollo
}

func Apollo() *apollo {
	if a, ok := apolloTable[apolloDriver]; ok {
		return a
	}

	apolloTable[apolloDriver] = &apollo{initAgollo()}

	return apolloTable[apolloDriver]
}

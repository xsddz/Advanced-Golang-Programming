package app

import "yawebapp/library/inner/trace"

var (
	Logger *trace.Logger
)

func loadLogger() *trace.Logger {
	if Logger == nil {
		Logger = initLogger()
	}
	return Logger
}

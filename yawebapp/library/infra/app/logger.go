package app

import "yawebapp/library/infra/trace"

var (
	Logger *trace.Logger
)

func initGlobalLogger() *trace.Logger {
	if Logger == nil {
		Logger = initLogger()
	}
	return Logger
}

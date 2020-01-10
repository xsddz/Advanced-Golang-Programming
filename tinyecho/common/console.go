package common

const (
	colorNone   = "\033[0m"
	colorBlack  = "\033[0;30m"
	colorGreen  = "\033[0;32m"
	colorYellow = "\033[0;33m"
)

// GreenString  GreenString
func GreenString(s string) string {
	return colorGreen + s + colorNone
}

// YellowString YellowString
func YellowString(s string) string {
	return colorYellow + s + colorNone
}

package common

import "fmt"

const (
	colorNone    = "\033[0m"
	colorBlack   = "\033[0;30m"
	colorRed     = "\033[0;31m"
	colorGreen   = "\033[0;32m"
	colorYellow  = "\033[0;33m"
	colorBule    = "\033[0;34m"
	colorMagenta = "\033[0;35m"
	colorCyan    = "\033[0;36m"
	colorGray    = "\033[0;37m"

	fontBold      = "\033[1m"
	fontUnderline = "\033[4m"
	fontBlink     = "\033[5m"
	fontReverse   = "\033[7m"
	fontHide      = "\033[8m"

	screenClear     = "\033[2J"
	screenClearLine = "\033[1K\r"
)

// BlackString Black String
func BlackString(s string) string {
	return colorBlack + s + colorNone
}

// RedString Red String
func RedString(s string) string {
	return colorRed + s + colorNone
}

// GreenString  Green String
func GreenString(s string) string {
	return colorGreen + s + colorNone
}

// YellowString Yellow String
func YellowString(s string) string {
	return colorYellow + s + colorNone
}

// BuleString Bule String
func BuleString(s string) string {
	return colorBule + s + colorNone
}

// MagentaString Magenta String
func MagentaString(s string) string {
	return colorMagenta + s + colorNone
}

// CyanString Cyan String
func CyanString(s string) string {
	return colorCyan + s + colorNone
}

// GrayString Gray String
func GrayString(s string) string {
	return colorGray + s + colorNone
}

// BoldString Bold String
func BoldString(s string) string {
	return fontBold + s + colorNone
}

// UnderlineString Underline String
func UnderlineString(s string) string {
	return fontUnderline + s + colorNone
}

// BlinkString Blink String
func BlinkString(s string) string {
	return fontBlink + s + colorNone
}

// ReverseString Reverse String
func ReverseString(s string) string {
	return fontReverse + s + colorNone
}

// HideString Hide String
func HideString(s string) string {
	return fontHide + s + colorNone
}

// ScreenClear Clear Screen
func ScreenClear() {
	fmt.Printf(screenClear)
}

// ScreenClearLine Clear Current Line
func ScreenClearLine() {
	fmt.Printf(screenClearLine)
}

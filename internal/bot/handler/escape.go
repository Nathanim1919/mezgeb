package handler

import "strings"

// escMD escapes Telegram Markdown v1 special characters in user-provided text.
var mdReplacer = strings.NewReplacer(
	"*", "\\*",
	"_", "\\_",
	"`", "\\`",
	"[", "\\[",
)

func escMD(s string) string {
	return mdReplacer.Replace(s)
}

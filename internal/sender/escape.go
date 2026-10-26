package sender

import "strings"

// EscapeText escapes Telegram's MarkdownV2 reserved characters so text renders
// literally inside a message sent with the package's parse mode. The reserved
// set is the MarkdownV2 syntax: _ * [ ] ( ) ~ ` > # + - = | { } . !
func EscapeText(text string) string {
	return strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	).Replace(text)
}

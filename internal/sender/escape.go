package sender

import "strings"

// EscapeText escapes Telegram's MarkdownV2 reserved characters so text renders
// literally inside a message sent with the package's parse mode. The reserved
// set is the MarkdownV2 syntax: _ * [ ] ( ) ~ ` > # + - = | { } . ! — plus a
// literal backslash, the escape character itself, which would otherwise
// consume the escape of the character that follows it.
func EscapeText(text string) string {
	return strings.NewReplacer(
		"\\", "\\\\",
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

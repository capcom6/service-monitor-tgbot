package templates

type Config struct {
	Templates map[string]string // key is template name
	EscapeFn  func(string) string
}

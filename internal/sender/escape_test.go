package sender_test

import (
	"testing"

	"github.com/capcom6/service-monitor-tgbot/internal/sender"
)

func TestEscapeText(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "plain text is untouched", in: "Google is online", want: "Google is online"},
		{name: "parentheses are escaped", in: "3 (2 offline)", want: `3 \(2 offline\)`},
		{name: "dots are escaped", in: "api.v2.example.com", want: `api\.v2\.example\.com`},
		{name: "emphasis and brackets are escaped", in: "_*[]", want: `\_\*\[\]`},
		{name: "operators are escaped", in: "#+-=|{}.!", want: `\#\+\-\=\|\{\}\.\!`},
		{name: "tilde is escaped", in: "a~b", want: `a\~b`},
		{name: "backtick is escaped", in: "a`b", want: "a\\`b"},
		{name: "greater-than is escaped", in: "a>b", want: `a\>b`},
		{name: "backslash is untouched", in: `a\b`, want: `a\b`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sender.EscapeText(tt.in); got != tt.want {
				t.Errorf("EscapeText(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

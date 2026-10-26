package sender_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/capcom6/service-monitor-tgbot/internal/sender"
	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	"go.uber.org/zap"
)

// testToken is well formed but not real: telego validates the token format at
// construction and never dials it, so the graph only reaches the fake API when
// a test actually sends.
const testToken = "123456:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

// sentRequest is the decoded sendMessage payload the fake API records.
type sentRequest struct {
	Path      string `json:"-"`
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// TestService_Send pins the outbound contract of the message sender: the chat
// ID, the verbatim text and - the detail a missing field would break silently
// - MarkdownV2 as the parse mode the templates are escaped for.
func TestService_Send(t *testing.T) {
	received := make(chan sentRequest, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			t.Errorf("read request: %v", readErr)
			return
		}

		var req sentRequest
		if unmarshalErr := json.Unmarshal(body, &req); unmarshalErr != nil {
			t.Errorf("decode request %q: %v", body, unmarshalErr)
			return
		}
		req.Path = r.URL.Path
		received <- req

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"ok":true,"result":{"message_id":7,"date":0,"chat":{"id":-1001,"type":"private"},"text":"reply"}}`,
		))
	}))
	t.Cleanup(server.Close)

	bot, err := telegofx.New(
		telegofx.Config{Token: testToken, Mode: telegofx.ModePolling},
		&http.Client{Timeout: time.Second},
		[]telego.BotOption{telego.WithAPIServer(server.URL)},
		zap.NewNop(),
	)
	if err != nil {
		t.Fatalf("create bot: %v", err)
	}

	svc := sender.New(bot)
	if sendErr := svc.Send(context.Background(), -1001, `✅ api.v2 is *offline*`); sendErr != nil {
		t.Fatalf("send message: %v", sendErr)
	}

	got := <-received
	if got.ChatID != -1001 {
		t.Errorf("chat_id = %d, want -1001", got.ChatID)
	}
	if got.Text != `✅ api.v2 is *offline*` {
		t.Errorf("text = %q, want the verbatim message", got.Text)
	}
	if got.ParseMode != telego.ModeMarkdownV2 {
		t.Errorf("parse_mode = %q, want %q", got.ParseMode, telego.ModeMarkdownV2)
	}
	if !strings.HasSuffix(got.Path, "/sendMessage") {
		t.Errorf("path = %q, want a sendMessage call", got.Path)
	}
}

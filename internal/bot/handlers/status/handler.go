// Package status implements the /status command handler.
package status

import (
	"fmt"

	"github.com/capcom6/service-monitor-tgbot/internal/bot/handler"
	"github.com/capcom6/service-monitor-tgbot/internal/messages"
	"github.com/capcom6/service-monitor-tgbot/internal/monitor"
	"github.com/capcom6/service-monitor-tgbot/internal/sender"
	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

// Handler is the /status command: it replies with the current services list.
type Handler struct {
	monitorSvc  *monitor.Service
	messagesSvc *messages.Service
	senderSvc   *sender.Service
}

// New constructs the status handler. It returns handler.Handler - not the
// concrete type - because dig keys value groups by the exact result type: a
// concrete return type would make the []handler.Handler slice resolve empty
// with no error.
func New(
	monitorSvc *monitor.Service,
	messagesSvc *messages.Service,
	senderSvc *sender.Service,
) handler.Handler {
	return &Handler{
		monitorSvc:  monitorSvc,
		messagesSvc: messagesSvc,
		senderSvc:   senderSvc,
	}
}

// Register mounts the /status route on router; the predicate guarantees the
// update carries a command message, so the chat ID is always usable.
func (h *Handler) Register(router *telegofx.Router) {
	router.Handle(h.handleStatus, th.CommandEqual("status"))
}

// handleStatus renders the services list and replies in the chat the command
// arrived in. Errors return to the caller: telegofx logs them through the bot
// logger together with the update ID.
func (h *Handler) handleStatus(ctx *th.Context, update telego.Update) error {
	chatID := update.Message.Chat.ID

	services := h.monitorSvc.GetCurrentStatuses()
	if len(services) == 0 {
		if err := h.senderSvc.Send(ctx, chatID, "No services configured."); err != nil {
			return fmt.Errorf("send empty status response: %w", err)
		}

		return nil
	}

	data := make(messages.ServicesListContext, len(services))
	for i, service := range services {
		data[i] = messages.NewServiceState(service.Name, string(service.State), service.Error, service.ChangedAt)
	}

	msg, err := h.messagesSvc.ServicesList(data)
	if err != nil {
		return fmt.Errorf("render services list: %w", err)
	}

	if sendErr := h.senderSvc.Send(ctx, chatID, msg); sendErr != nil {
		return fmt.Errorf("send services list: %w", sendErr)
	}

	return nil
}

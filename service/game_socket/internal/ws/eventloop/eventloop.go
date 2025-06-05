package eventloop

import (
	"errors"
	"log/slog"
	"ms4me/game_socket/internal/models"
	dto_ws "ms4me/game_socket/internal/ws/dto"
	ws "ms4me/game_socket/internal/ws/server"
)

var (
	ErrUserNotFound  = errors.New("пользователь не найден")
	ErrInternalError = errors.New("внутренняя ошибка")
	ErrInvalidChatID = errors.New("неверный chat_id")
)

type EventLoop struct {
	log            *slog.Logger
	eventQueue     *chan models.Event
	ws             *ws.Server
	eventsShutdown chan struct{}
}

func New(log *slog.Logger, eventQueue *chan models.Event, ws *ws.Server) *EventLoop {
	return &EventLoop{
		log:            log,
		eventQueue:     eventQueue,
		ws:             ws,
		eventsShutdown: make(chan struct{}),
	}
}

func (s *EventLoop) EventLoop() {
	const op = "ws.processEvents"
	log := s.log.With(slog.String("op", op))

	queue := *s.eventQueue

	for {
		select {
		case event := <-queue:
			log.Info("received event", slog.Any("event", event))

			switch event.Type {
			case models.TypeCreateGame:
				resp := &dto_ws.Response{
					Status:    dto_ws.StatusOK,
					EventType: dto_ws.CreateRoomEventType,
					Payload:   event.Payload,
				}
				log.Info("broadcast event", slog.Any("event", resp))
				go s.ws.BroadcastEvent(resp)

			default:
				log.Warn("unknown event type", slog.String("type", string(event.Type)))
			}

		case <-s.eventsShutdown:
			log.Info("event loop shutting down")
			return
		}
	}
}

func (s *EventLoop) Stop() {
	s.eventsShutdown <- struct{}{}
}

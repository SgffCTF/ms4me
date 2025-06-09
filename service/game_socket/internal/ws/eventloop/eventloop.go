package eventloop

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"ms4me/game_socket/internal/models"
	storage "ms4me/game_socket/internal/redis"
	dto_ws "ms4me/game_socket/internal/ws/dto"
	ws "ms4me/game_socket/internal/ws/server"

	"github.com/jacute/prettylogger"
	"github.com/redis/go-redis/v9"
)

var (
	ErrUserNotFound  = errors.New("пользователь не найден")
	ErrInternalError = errors.New("внутренняя ошибка")
	ErrInvalidChatID = errors.New("неверный chat_id")
)

type EventLoop struct {
	log    *slog.Logger
	ws     *ws.Server
	redis  *storage.Redis
	pubsub *redis.PubSub
}

func New(log *slog.Logger, ws *ws.Server, redis *storage.Redis) *EventLoop {
	return &EventLoop{
		log:    log,
		ws:     ws,
		redis:  redis,
		pubsub: redis.DB.Subscribe(context.Background(), storage.PUBLIC_QUEUE),
	}
}

func (s *EventLoop) EventLoop() {
	const op = "eventloop.processEvents"
	log := s.log.With(slog.String("op", op))

	defer s.pubsub.Close()

	ch := s.pubsub.Channel()

	for msg := range ch {
		log.Info("received msg", slog.Any("msg", msg))

		var event models.Event
		err := json.Unmarshal([]byte(msg.Payload), &event)
		if err != nil {
			log.Error("error unmarshaling event", slog.String("msg", msg.Payload), prettylogger.Err(err))
			continue
		}

		var resp *dto_ws.Response
		switch event.Type {
		case models.TypeCreateGame:
			resp = &dto_ws.Response{
				Status:    dto_ws.StatusOK,
				EventType: dto_ws.CreateRoomEventType,
				Payload:   event.Payload,
			}
			var eventUnmarshalled models.CreateEvent
			err := json.Unmarshal(event.Payload, &eventUnmarshalled)
			if err != nil {
				log.Error("error unmarshalling event", slog.Any("event", event), prettylogger.Err(err))
				continue
			}
			err = s.redis.AddClientToChannel(context.Background(), event.GameID, event.UserID, &models.RoomParticipant{
				ID:       event.UserID,
				Username: event.Username,
				IsOwner:  true,
			})
			if err != nil {
				log.Error("error adding event to channel", slog.Any("event", event), prettylogger.Err(err))
				continue
			}
			log.Info("broadcast event", slog.Any("event", resp))
			go s.ws.BroadcastEvent(resp)
		case models.TypeUpdateGame:
			resp = &dto_ws.Response{
				Status:    dto_ws.StatusOK,
				EventType: dto_ws.UpdateRoomEventType,
				Payload:   event.Payload,
			}
			users, err := s.redis.GetUsersInChannel(context.Background(), event.GameID)
			if err != nil {
				log.Error("error reading channel clients from redis", slog.Any("event", resp), prettylogger.Err(err))
				return
			}
			if event.IsPublic {
				go s.ws.BroadcastEvent(resp)
			}
			go s.ws.MulticastEvent(event.GameID, users, resp)
		case models.TypeDeleteGame:
			payloadMarshalled, err := json.Marshal(map[string]any{"id": event.GameID, "user_id": event.UserID})
			if err != nil {
				log.Error("error marshalling event", slog.Any("event", event))
				continue
			}
			resp = &dto_ws.Response{
				Status:    dto_ws.StatusOK,
				EventType: dto_ws.DeleteRoomEventType,
				Payload:   payloadMarshalled,
			}
			users, err := s.redis.GetUsersInChannel(context.Background(), event.GameID)
			if err != nil {
				log.Error("error reading channel clients from redis", slog.Any("event", resp), prettylogger.Err(err))
				return
			}
			err = s.redis.DeleteChannel(context.Background(), event.GameID)
			if err != nil {
				log.Error("error deleting channel", slog.Any("event", event), prettylogger.Err(err))
				continue
			}
			if event.IsPublic {
				go s.ws.BroadcastEvent(resp)
			}
			go s.ws.MulticastEvent(event.GameID, users, resp)
		case models.TypeJoinGame:
			payloadMarshalled, err := json.Marshal(map[string]any{
				"id":       event.GameID,
				"user_id":  event.UserID,
				"username": event.Username,
			})
			if err != nil {
				log.Error("error marshalling event", slog.Any("event", event))
				continue
			}
			resp = &dto_ws.Response{
				Status:    dto_ws.StatusOK,
				EventType: dto_ws.JoinRoomEventType,
				Payload:   payloadMarshalled,
			}
			users, err := s.redis.GetUsersInChannel(context.Background(), event.GameID)
			if err != nil {
				log.Error("error reading channel clients from redis", slog.Any("event", resp), prettylogger.Err(err))
				return
			}
			err = s.redis.AddClientToChannel(context.Background(), event.GameID, event.UserID, &models.RoomParticipant{
				ID:       event.UserID,
				Username: event.Username,
				IsOwner:  false,
			})
			if err != nil {
				log.Error("error adding client to channel", slog.Any("event", event), prettylogger.Err(err))
				continue
			}
			if event.IsPublic {
				go s.ws.BroadcastEvent(resp)
			}
			go s.ws.MulticastEvent(event.GameID, users, resp)
		case models.TypeExitGame:
			payloadMarshalled, err := json.Marshal(map[string]any{
				"id":       event.GameID,
				"user_id":  event.UserID,
				"username": event.Username,
			})
			if err != nil {
				log.Error("error marshalling event", slog.Any("event", event))
				continue
			}
			resp = &dto_ws.Response{
				Status:    dto_ws.StatusOK,
				EventType: dto_ws.ExitRoomEventType,
				Payload:   payloadMarshalled,
			}
			err = s.redis.RemoveClientFromChannel(context.Background(), event.GameID, event.UserID)
			if err != nil {
				log.Error("error removing client from channel", slog.Any("event", event), prettylogger.Err(err))
				continue
			}
			users, err := s.redis.GetUsersInChannel(context.Background(), event.GameID)
			if err != nil {
				log.Error("error reading channel clients from redis", slog.Any("event", resp), prettylogger.Err(err))
				return
			}
			if event.IsPublic {
				go s.ws.BroadcastEvent(resp)
			}
			go s.ws.MulticastEvent(event.GameID, users, resp)
		case models.TypeStartGame:
			payloadMarshalled, err := json.Marshal(map[string]any{
				"id": event.GameID,
			})
			if err != nil {
				log.Error("error marshalling event", slog.Any("event", event))
				continue
			}
			resp = &dto_ws.Response{
				Status:    dto_ws.StatusOK,
				EventType: dto_ws.StartGameEventType,
				Payload:   payloadMarshalled,
			}
			s.redis.AddGameInfoIntoRoom(context.Background(), event.GameID, &models.GameInfo{})
			users, err := s.redis.GetUsersInChannel(context.Background(), event.GameID)
			if err != nil {
				log.Error("error reading channel clients from redis", slog.Any("event", resp), prettylogger.Err(err))
				return
			}
			if event.IsPublic {
				go s.ws.BroadcastEvent(resp)
			}
			go s.ws.MulticastEvent(event.GameID, users, resp)
		default:
			log.Warn("unknown event type", slog.String("type", string(event.Type)))
			continue
		}
	}
}

func (s *EventLoop) Stop() {
	s.pubsub.Close()
}

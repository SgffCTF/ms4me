package events

import (
	"errors"
	"fmt"
	"log/slog"
	eventsdto "ms4me/game/internal/http/dto/events"
	"ms4me/game/internal/storage/channel"

	"github.com/jacute/prettylogger"
)

type EventsService struct {
	log      *slog.Logger
	channels *channel.Channels
}

func New(log *slog.Logger, channels *channel.Channels) *EventsService {
	return &EventsService{
		log:      log,
		channels: channels,
	}
}

func (s *EventsService) ProcessEvent(event *eventsdto.Event) {
	const op = "service.events.ProcessEvent"
	log := s.log.With(slog.Int64("user_id", event.UserID), slog.String("game_id", event.GameID))

	switch event.Type {
	case eventsdto.StartGame:
		log.Info("Received start game event")

		// channelName := redis.GetChannel()

		// s.CreateChannel(fmt.Sprintf("game_%s", event.GameID))
		// if err := s.AddToChannel(fmt.Sprintf("game_%s", event.GameID), event.UserID); err != nil {
		// 	log.Error("Failed to add user to channel", prettylogger.Err(err), slog.String("game_id", event.GameID), slog.Int64("user_id", event.UserID))
		// }
	case eventsdto.CreateGame:
		log.Info("Received create game event")

		s.CreateChannel(event.UserID)
		if err := s.AddToChannel(event.UserID, event.GameID); err != nil && !errors.Is(err, ErrUserAlreadyInChannel) {
			log.Error("Failed to add user to channel", prettylogger.Err(err))
		}
	case eventsdto.JoinGame:
		log.Info("Received join game event")

		s.CreateChannel(event.UserID)
		if err := s.AddToChannel(event.UserID, event.GameID); err != nil && !errors.Is(err, ErrUserAlreadyInChannel) {
			log.Error("Failed to add user to channel", prettylogger.Err(err))
		}
	default:
		s.log.Warn("Unknown event type", slog.Int("type", int(event.Type)))
	}
}

func (s *EventsService) CreateChannel(userID int64) {
	const op = "service.events.CreateChannel"

	s.channels.Mu.Lock()
	if _, ok := s.channels.Channels[userID]; !ok {
		s.channels.Channels[userID] = ""
		s.log.Info("Created channel", slog.String("op", op), slog.Int64("user_id", userID))
	}
	s.channels.Mu.Unlock()
}

func (s *EventsService) AddToChannel(userID int64, gameID string) error {
	const op = "service.events.AddToChannel"

	s.channels.Mu.Lock()
	defer s.channels.Mu.Unlock()
	if channel, ok := s.channels.Channels[userID]; ok {
		if channel != gameID {
			s.channels.Channels[userID] = gameID
		} else {
			return fmt.Errorf("%s: %w", op, ErrUserAlreadyInChannel)
		}
	} else {
		return fmt.Errorf("%s: %w", op, ErrChannelNotFound)
	}

	return nil
}

package ws

import (
	"errors"
	"io"
	"log/slog"
	"time"

	dto_ws "ms4me/game_socket/internal/ws/dto"

	"github.com/google/uuid"
	"github.com/jacute/prettylogger"
	"golang.org/x/net/websocket"
)

func (s *Server) Handle(conn *websocket.Conn) {
	const op = "ws.Handle"
	ctx := conn.Request().Context()
	requestID := uuid.NewString()
	log := s.log.With(slog.String("request_id", requestID), slog.String("op", op))

	log.Info("new connection", slog.String("origin", conn.RemoteAddr().String()))

	user, err := s.auth(ctx, conn)
	if err != nil {
		_, err = conn.Write(dto_ws.Error(err, dto_ws.AuthEventType).Serialize())
		if err != nil {
			log.Error("failed to send error message", prettylogger.Err(err))
		}
		log.Info("failed to auth client. close connection")
		err = conn.Close()
		if err != nil {
			log.Error("failed to close connect", prettylogger.Err(err))
		}
		return
	}
	client := &Client{ctx: ctx, conn: conn, user: user, requestID: requestID}

	s.clientsMu.Lock()
	err = s.write(conn, dto_ws.OK("Authenticated successfully", dto_ws.AuthEventType).Serialize())
	if err != nil {
		log.Error("failed to send message", prettylogger.Err(err))
		err = conn.Close()
		if err != nil {
			log.Error("failed to close connect", prettylogger.Err(err))
		}
		s.clientsMu.Unlock()
		return
	}
	s.clients[user.ID] = client
	s.clientsMu.Unlock()

	s.pingLoop(client)
}

func (s *Server) Close() error {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	for _, client := range s.clients {
		err := client.conn.Close()
		if err != nil {
			s.log.Error("failed to close client connection", prettylogger.Err(err))
		}
		delete(s.clients, client.user.ID)
	}

	return nil
}

func (s *Server) read(conn *websocket.Conn) ([]byte, error) {
	buf := make([]byte, BUF_SIZE)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, errors.New("empty message received")
	}
	msg := buf[:n]

	return msg, nil
}

func (s *Server) write(conn *websocket.Conn, msg []byte) error {
	const op = "ws.write"
	log := s.log.With(slog.String("op", op))

	if conn == nil {
		log.Warn("attempted to write to nil connection")
		return errors.New("attempted to write to nil connection")
	}

	if _, err := conn.Write(msg); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, websocket.ErrBadFrame) {
			log.Info("client connection closed, skipping write")
			return err
		}
		log.Warn("error writing to conn", prettylogger.Err(err))
		return err
	}

	return nil
}

func (s *Server) pingLoop(client *Client) {
	const op = "ws.pingLoop"
	log := s.log.With(slog.String("op", op), slog.String("request_id", client.requestID), slog.Int64("user_id", client.user.ID))

	ticker := time.NewTicker(s.cfg.WSPingTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-client.ctx.Done():
			return
		case <-ticker.C:
			err := websocket.Message.Send(client.conn, "")
			if err != nil {
				log.Debug("ping failed, closing connection", prettylogger.Err(err))
				s.disconnect(client)
				return
			}
			log.Debug("ping succeeded")
		}
	}
}

func (s *Server) disconnect(client *Client) error {
	const op = "ws.disconnect"
	log := s.log.With(slog.String("op", op), slog.String("request_id", client.requestID), slog.Int64("user_id", client.user.ID))

	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	err := client.conn.Close()
	if err != nil {
		return err
	}
	delete(s.clients, client.user.ID)

	log.Info("client disconnected")

	return nil
}

func (s *Server) BroadcastEvent(res *dto_ws.Response) {
	const op = "ws.BroadcatEvent"
	log := s.log.With(slog.String("op", op))

	log.Debug("start broadcast", slog.Any("clients", s.clients))
	for _, client := range s.clients {
		err := s.write(client.conn, res.Serialize())
		if err != nil {
			log.Error("error writing event to client", slog.Any("event", res))
			continue
		}
		log.Debug("event sent to client", slog.Any("event", res), slog.Int64("user_id", client.user.ID))
	}
}

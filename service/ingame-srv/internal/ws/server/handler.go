package ws

import (
	"errors"
	"io"
	"log/slog"
	"sync"
	"time"

	dto_ws "ms4me/game_socket/internal/ws/dto"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jacute/prettylogger"
	"golang.org/x/net/websocket"
)

const multicastWriteRetriesCount = 3
const multicastWriteRetriesTimeoutMLS = 200

func (s *Server) Handle(conn *websocket.Conn) {
	const op = "ws.Handle"
	r := conn.Request()
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	requestID := uuid.NewString()
	log := s.log.With(slog.String("request_id", requestID), slog.String("op", op))

	log.Debug("new connection", slog.String("origin", conn.RemoteAddr().String()))

	user, err := s.auth(ctx, conn)
	if err != nil {
		_, err = conn.Write(dto_ws.Error(err, dto_ws.AuthEventType).Serialize())
		if err != nil {
			log.Error("failed to send error message", prettylogger.Err(err))
		}
		log.Debug("failed to auth client. close connection")
		err = conn.Close()
		if err != nil {
			log.Error("failed to close connect", prettylogger.Err(err))
		}
		return
	}
	if id != "" {
		_, err = s.redis.GetClientInChannel(ctx, id, user.ID)
		if err != nil {
			_, err = conn.Write(dto_ws.ErrPlayerNotInGame.Serialize())
			if err != nil {
				log.Error("failed to send error message", prettylogger.Err(err))
			}
			log.Error("failed to get info in room about user", slog.Any("user", user))
			conn.Close()
			return
		}
		log.Info("user is participant of room", slog.Any("user", user), slog.String("room_id", id))
	}
	client := &Client{ctx: ctx, conn: conn, user: user, requestID: requestID, room: id}

	s.usersMu.Lock()
	err = s.write(conn, dto_ws.OK("Authenticated successfully", dto_ws.AuthEventType).Serialize())
	if err != nil {
		log.Error("failed to send message", prettylogger.Err(err))
		err = conn.Close()
		if err != nil {
			log.Error("failed to close connect", prettylogger.Err(err))
		}
		s.usersMu.Unlock()
		return
	}
	s.users[user.ID] = append(s.users[user.ID], client)
	s.usersMu.Unlock()

	go s.readLoop(client)
	s.pingLoop(client)
}

func (s *Server) Close() error {
	s.usersMu.Lock()
	defer s.usersMu.Unlock()

	for _, clients := range s.users {
		for _, client := range clients {
			err := client.conn.Close()
			if err != nil {
				s.log.Error("failed to close client connection", prettylogger.Err(err))
			}
		}
		if len(clients) > 0 {
			delete(s.users, clients[0].user.ID)
		}
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
			err := websocket.Message.Send(client.conn, "ping")
			if err != nil {
				log.Debug("ping failed, closing connection", prettylogger.Err(err))
				s.disconnect(client)
				return
			}
			log.Debug("ping succeeded")
		}
	}
}

// CheckConn проверяет, что есть хотя бы одно активное соединение по вебсокету с юзером
func (s *Server) CheckRoomConn(userID int64, room string) bool {
	const op = "ws.CheckConn"
	log := s.log.With(slog.String("op", op))

	clients := s.users[userID]

	connExists := false
	for _, client := range clients {
		if client.room == room {
			err := websocket.Message.Send(client.conn, "ping")
			if err != nil {
				log.Debug("ping failed, closing connection", prettylogger.Err(err))
				s.disconnect(client)
				continue
			}
			log.Debug("ping succeeded")
			connExists = true
			break
		}
	}

	if !connExists {
		return false
	}
	return true
}

func (s *Server) disconnect(client *Client) error {
	const op = "ws.disconnect"
	log := s.log.With(slog.String("op", op), slog.String("request_id", client.requestID), slog.Int64("user_id", client.user.ID))

	s.usersMu.Lock()
	defer s.usersMu.Unlock()

	clients := s.users[client.user.ID]
	for i, c := range clients {
		if c == client {
			s.users[client.user.ID] = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	err := client.conn.Close()
	if err != nil {
		return err
	}
	if len(s.users[client.user.ID]) == 0 {
		delete(s.users, client.user.ID)
	}

	log.Info("client disconnected")

	return nil
}

func (s *Server) MulticastEvent(roomID string, users []int, res *dto_ws.Response) {
	const op = "ws.MulticastEvent"
	log := s.log.With(slog.String("op", op), slog.String("room_id", roomID))
	startTime := time.Now().UTC()

	var wg sync.WaitGroup
	log.Debug("start multicast")
	for _, userID := range users {
		clients, ok := s.users[int64(userID)]
		if !ok {
			log.Warn("user with this id not found in ws clients", slog.Int("user_id", userID))
			continue
		}
		for _, client := range clients {
			if client.room == roomID {
				wg.Add(1)
				go func() {
					for i := 0; i < multicastWriteRetriesCount; i++ {
						err := s.write(client.conn, res.Serialize())
						if err != nil {
							log.Error("error writing event to client", slog.Any("event", res))
							time.Sleep(multicastWriteRetriesTimeoutMLS * time.Millisecond)
							continue
						}
						break
					}
					wg.Done()
					// log.Debug("event sent to client", slog.Any("event", res), slog.Int64("user_id", client.user.ID))
				}()
			}
		}
	}
	wg.Wait()
	log.Debug("end multicast", slog.Duration("time", time.Now().UTC().Sub(startTime)))
}

func (s *Server) DisconnectRoom(roomID string, users []int) {
	const op = "ws.DisconnectRoom"
	log := s.log.With(slog.String("op", op), slog.String("room_id", roomID))

	for _, userID := range users {
		clients, ok := s.users[int64(userID)]
		if !ok {
			log.Warn("user with this id not found in ws clients", slog.Int("user_id", userID))
			continue
		}
		for _, client := range clients {
			if client.room == roomID {
				err := s.disconnect(client)
				if err != nil {
					log.Error("error disconnecting", slog.Int64("user_id", client.user.ID), prettylogger.Err(err))
				}
			}
		}
	}
}

func (s *Server) BroadcastEvent(res *dto_ws.Response) {
	const op = "ws.BroadcastEvent"
	log := s.log.With(slog.String("op", op))

	log.Debug("start broadcast")
	for _, clients := range s.users {
		for _, client := range clients {
			if client.room == "" {
				err := s.write(client.conn, res.Serialize())
				if err != nil {
					log.Error("error writing event to client", slog.Any("event", res))
					continue
				}
				// log.Debug("event sent to client", slog.Any("event", res), slog.Int64("user_id", client.user.ID))
			}
		}
	}
	log.Debug("end broadcast")
}

func (s *Server) readLoop(client *Client) {
	const op = "ws.readLoop"
	log := s.log.With(slog.String("op", op), slog.String("request_id", client.requestID), slog.Int64("user_id", client.user.ID))

	for {
		select {
		case <-client.ctx.Done():
			log.Debug("context canceled, stop reading")
			return
		default:
			var msg string
			err := websocket.Message.Receive(client.conn, &msg)
			if err != nil {
				if err == io.EOF {
					log.Info("client closed connection")
				} else {
					log.Warn("error reading from websocket", prettylogger.Err(err))
				}
				s.disconnect(client)
				return
			}
			log.Debug("received message", slog.String("msg", msg))
		}
	}
}

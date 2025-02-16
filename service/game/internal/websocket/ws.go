package ws

import (
	"context"
	"io"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/jacute/prettylogger"
	"golang.org/x/net/websocket"
	"google.golang.org/grpc"

	ssov1 "ms4me/game/pkg/grpc/sso"
)

const BUF_SIZE = 4096

type AuthService interface {
	VerifyToken(ctx context.Context, in *ssov1.VerifyTokenRequest, opts ...grpc.CallOption) (*ssov1.VerifyTokenResponse, error)
}

type Server struct {
	log        *slog.Logger
	authClient AuthService
	db         UserInfo
	clients    map[*Client]bool
	clientsMu  sync.Mutex
}

type Client struct {
	ctx  context.Context
	conn *websocket.Conn
}

var requestIDContextKey contextKey = "request_id"

func New(log *slog.Logger, authClient AuthService, db UserInfo) *Server {
	return &Server{
		log:        log,
		authClient: authClient,
		clients:    make(map[*Client]bool),
		clientsMu:  sync.Mutex{},
		db:         db,
	}
}

func (s *Server) Handle(conn *websocket.Conn) {
	const op = "ws.Handle"
	requestID := uuid.NewString()
	log := s.log.With(slog.String("request_id", requestID), slog.String("op", op))

	log.Info("new connection", slog.String("addr", conn.RemoteAddr().String()))

	client := &Client{conn: conn, ctx: context.Background()}

	s.clientsMu.Lock()
	s.clients[client] = true
	s.clientsMu.Unlock()

	client.ctx = context.WithValue(context.Background(), requestIDContextKey, requestID)

	err := s.auth(client)
	if err != nil {
		jsonErr := NewError(err)
		conn.Write([]byte(jsonErr.Error()))

		log.Debug("error, closing connection", prettylogger.Err(err))
		conn.Close()
		return
	}

	s.readLoop(client)
}

func (s *Server) readLoop(client *Client) {
	const op = "ws.readLoop"
	request_id := client.ctx.Value(requestIDContextKey).(string)
	log := s.log.With(slog.String("op", op), slog.String("request_id", request_id))

	for {
		msg, err := s.read(client.conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Warn("Error reading from socket", prettylogger.Err(err))
			continue
		}
		log.Debug("readed successfully", slog.String("msg", string(msg)))

		// s.broadcast(msg)
	}
}

func (s *Server) read(conn *websocket.Conn) ([]byte, error) {
	buf := make([]byte, BUF_SIZE)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	msg := buf[:n]

	return msg, nil
}

func (s *Server) broadcast(msg []byte) {
	const op = "wb.broadcast"
	log := s.log.With(slog.String("op", op))

	for client, _ := range s.clients {
		go func(conn *websocket.Conn) {
			if _, err := conn.Write(msg); err != nil {
				log.Warn("error writing to conn", prettylogger.Err(err))
			}
		}(client.conn)
	}
}

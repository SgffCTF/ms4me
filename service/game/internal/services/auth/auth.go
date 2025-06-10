package auth

import (
	"context"
	"errors"
	"log/slog"
	"ms4me/game/internal/models"
	"ms4me/game/internal/storage"
	"ms4me/game/pkg/lib/jwt"
	"time"

	"github.com/jacute/prettylogger"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrIncorrectPassword = errors.New("Неправильный пароль")
)

type AuthStorage interface {
	CreateUser(ctx context.Context, username string, password string) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type AuthService struct {
	log       *slog.Logger
	db        AuthStorage
	jwtSecret []byte
	jwtTTL    time.Duration
}

func New(log *slog.Logger, db AuthStorage, jwtSecret []byte, jwtTTL time.Duration) *AuthService {
	return &AuthService{log, db, jwtSecret, jwtTTL}
}

func (as *AuthService) Register(ctx context.Context, username, password string) (int64, error) {
	const op = "auth.Register"
	log := as.log.With(slog.String("op", op), slog.String("username", username))

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error hashing password", prettylogger.Err(err))
		return 0, err
	}

	id, err := as.db.CreateUser(ctx, username, string(hash))
	if err != nil {
		log.Error("error creating user", prettylogger.Err(err))
		return 0, err
	}

	return id, nil
}

func (as *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	const op = "auth.Login"
	log := as.log.With(slog.String("op", op), slog.String("username", username))

	user, err := as.db.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found")
			return "", storage.ErrUserNotFound
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrIncorrectPassword
		}
		log.Error("error comparing hash with password", prettylogger.Err(err))
		return "", err
	}

	token, err := jwt.NewToken(user, as.jwtSecret, as.jwtTTL)
	if err != nil {
		log.Error("error creating token", prettylogger.Err(err))
		return "", err
	}

	return token, nil
}

package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"ms4me/sso/internal/database"
	"ms4me/sso/internal/lib/jwt"
	"ms4me/sso/internal/models"
	"time"

	"github.com/jacute/prettylogger"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	jwtSecret    []byte
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		username string,
		passwordHash []byte,
	) (id int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, username string) (*models.User, error)
}

// New creates a new Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	jwtSecret []byte,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		jwtSecret:    jwtSecret,
		tokenTTL:     tokenTTL,
	}
}

// Login checks if the user with given credentials exists and give the token
func (a *Auth) Login(
	ctx context.Context,
	username string,
	password string,
) (string, error) {
	const op = "auth.Login"
	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	user, err := a.userProvider.User(ctx, username)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			log.Warn("user not found")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("failed to get user", prettylogger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Info("invalid password", prettylogger.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, a.jwtSecret, a.tokenTTL)
	if err != nil {
		log.Info("failed to create token", prettylogger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// Register add a new user
func (a *Auth) Register(
	ctx context.Context,
	username string,
	password string,
) (int64, error) {
	const op = "auth.Register"
	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", prettylogger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, username, passwordHash)
	if err != nil {
		if errors.Is(err, database.ErrUserExists) {
			log.Warn("user already exists", prettylogger.Err(err))
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save user", prettylogger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered successfully")
	return id, nil
}

func (a *Auth) VerifyToken(ctx context.Context, token string) (*models.User, error) {
	const op = "auth.VerifyToken"
	log := a.log.With(
		slog.String("op", op),
	)

	claims, err := jwt.VerifyToken(token, a.jwtSecret)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Warn("token expired")
			return nil, fmt.Errorf("%s: %w", op, jwt.ErrTokenExpired)
		}
		if errors.Is(err, jwt.ErrTokenInvalid) {
			log.Warn("invalid token")
			return nil, fmt.Errorf("%s: %w", op, jwt.ErrTokenInvalid)
		}
		log.Error("unexpected error verifying token", prettylogger.Err(err))
		return nil, err
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		log.Error("invalid user_id parameter in claims", slog.Any("claims", claims))
		return nil, fmt.Errorf("%s: invalid user_id parameter in claims", op)
	}
	username, ok := claims["username"].(string)
	if !ok {
		log.Error("invalid username parameter in claims", slog.Any("claims", claims))
		return nil, fmt.Errorf("%s: invalid username parameter in claims", op)
	}

	return &models.User{
		ID:       int64(userID),
		Username: username,
	}, nil
}

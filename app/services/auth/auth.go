package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog"
	"github.com/victor-felix/chat-service/app/errors"
	"github.com/victor-felix/chat-service/app/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService models.UserService
	jwtSecret string
	jwtTTL int
	log zerolog.Logger
}

func NewAuthService(userService models.UserService, jwtSecret string, jwtTTL int, log zerolog.Logger) models.AuthService {
	return &AuthService{
		userService: userService,
		jwtSecret: jwtSecret,
		jwtTTL: jwtTTL,
		log: log,
	}
}

func (as *AuthService) Login(ctx context.Context, payload *models.AuthPayload) (*models.AuthResponsePayload, error) {
	if err := payload.Validate(); err != nil {
		return nil, err
	}

	user, err := as.userService.FindByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	if  user == nil {
		return nil, errors.NewErrNotFound(models.ErrorUserNotFound).WithMessage("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return nil, errors.NewErrInvalidArgument(models.ErrorInvalidPassword).WithMessage("invalid password")
	}

	expiresAt := time.Now().Add(time.Hour * time.Duration(as.jwtTTL)).Unix()

	jwtPayload := &models.Token{
		ID: user.ID,
		UserID: user.ID,
		UserName: user.UserName,
		Email: user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtPayload)

	tokenString, err := token.SignedString([]byte(as.jwtSecret))
	if err != nil {
		as.log.Error().Msg(err.Error())
		return nil, err
	}

	return &models.AuthResponsePayload{ Token: tokenString }, nil
}

package user

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog"
	"github.com/victor-felix/chat-service/app/errors"
	"github.com/victor-felix/chat-service/app/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userStorage models.UserStorage
	jwtSecret string
	jwtTTL int
	log zerolog.Logger
}

func NewUserService(userStorage models.UserStorage, jwtSecret string, jwtTTL int, log zerolog.Logger) models.UserService {
	return &UserService{
		userStorage: userStorage,
		jwtSecret: jwtSecret,
		jwtTTL: jwtTTL,
		log: log,
	}
}

func (us *UserService) Create(ctx context.Context, user *models.User) (*models.AuthResponsePayload, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := us.UserValidate(ctx, user); err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		us.log.Error().Msg(err.Error())
		return nil, err
	}

	user.Password = string(passwordHash)
	us.userStorage.Insert(ctx, user)
	user.Password = ""

	expiresAt := time.Now().Add(time.Hour * time.Duration(us.jwtTTL)).Unix()

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

	tokenString, err := token.SignedString([]byte(us.jwtSecret))
	if err != nil {
		us.log.Error().Msg(err.Error())
		return nil, err
	}

	return &models.AuthResponsePayload{ Token: tokenString }, nil
}

func (us *UserService) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return us.userStorage.FindOneByEmail(ctx, email)
}

func (us *UserService) UserValidate(ctx context.Context, user *models.User) error {
	checkUser, err := us.userStorage.FindOneByEmail(ctx, user.Email)

	if checkUser != nil {
		return errors.NewErrInvalidArgument(models.ErrorEmailDuplicated).WithMessage("email already exists")
	}

	if err != nil && isUserNotFoundError(err) {
		return nil
	}

	if err != nil {
		us.log.Error().Msg(err.Error())
	}

	return err
}

func isUserNotFoundError(err error) bool {
	return strings.Contains(err.Error(), string(models.ErrorUserNotFound))
}

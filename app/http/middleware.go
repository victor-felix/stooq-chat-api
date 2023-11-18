package http

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/victor-felix/chat-service/app/errors"
	"github.com/victor-felix/chat-service/app/models"
)

const JwtClaimsAttribute = "jwtClaims"

func (h *handler) authGuard(c *gin.Context) {
	token := c.GetHeader("Authorization")
	access_token := c.Query("access_token")

	if len(strings.TrimSpace(token)) == 0 && len(strings.TrimSpace(access_token)) == 0 {
		h.responseProblem(c, errors.NewErrUnauthorized(models.ErrorAuthTokenMissing).WithMessage("auth token is missing"))
		c.Abort()
		return
	}

	if len(strings.TrimSpace(access_token)) > 0 {
		token = access_token
	}

	token = strings.Replace(token, "Bearer ", "", 1)

	jwtSecret := h.config.JwtSecret

	jwtToken, err := jwt.Parse(token, func (token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewErrInvalidArgument(models.ErrorInvalidToken).WithMessage("invalid token")
		}
	
		return []byte(jwtSecret), nil
	})

	if err != nil {
		h.responseProblem(c, errors.NewErrInvalidArgument(models.ErrorInvalidToken).WithMessage("invalid token"))
		c.Abort()
		return
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		h.responseProblem(c, errors.NewErrInvalidArgument(models.ErrorInvalidToken).WithMessage("invalid token"))
		c.Abort()
		return
	}

	c.Set(JwtClaimsAttribute, claims)
	c.Next()
}
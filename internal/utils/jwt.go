package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Krajiyah/nimble-interview-backend/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func ValidateJWT(db *gorm.DB, token string) (*models.User, error) {
	tokenObj, err := jwt.Parse(token, parseJWT(db))
	if err != nil || !tokenObj.Valid {
		return nil, errors.Wrap(err, "Invalid JWT")
	}
	s := tokenObj.Claims.(jwt.MapClaims)["jti"].(string)
	id, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return models.GetUserByID(db, uint(id))
}

func parseJWT(db *gorm.DB) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("Malformed Claims in JWT")
		}

		s := claims["jti"].(string)
		id, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		user, err := models.GetUserByID(db, uint(id))
		if err != nil {
			return nil, err
		}
		return []byte(user.PasswordHash), nil
	}
}

func NewJWT(user *models.User, sessionDuration time.Duration) (string, error) {
	claims := jwt.StandardClaims{
		Id:        fmt.Sprintf("%d", user.ID),
		ExpiresAt: time.Now().Add(sessionDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(user.PasswordHash))
}

package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/xenking/kitsu-media-server/pkg/config"
)

func GenerateJWT(id uint) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, _ := token.SignedString(config.Global.JWTSecret)
	return t
}

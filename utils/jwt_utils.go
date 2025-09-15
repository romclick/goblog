package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JwtConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireHour int    `mapstructure:"expire_hour"`
}

type Claims struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

func InitConfig() {
	viper.SetConfigFile("configs/apps.yaml")
	viper.ReadInConfig()

	var jwtConfig JwtConfig
	viper.UnmarshalKey("jwtConfig", &jwtConfig)
	//jwt.SetSecret(jwtConfig.Secret)
	jwtSecret = []byte(jwtConfig.Secret)
}

func GenerateToken(userId uint, expirSeconds int64) (string, error) {
	expirTime := time.Now().Add(time.Duration(expirSeconds) * time.Second)
	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "Goblog",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
	if len(jwtSecret) == 0 {
		return nil, errors.New("JWT密钥未初始化，请先调用InitConfig")
	}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("无效的令牌")
}

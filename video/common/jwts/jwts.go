package jwts

import (
	"errors"
	"fmt"
	"west2-video/common/logs"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"time"
)

type JwtToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessExp    int64  `json:"access_exp"`
	RefreshExp   int64  `json:"refresh_exp"`
}

func CreateToken(val, username string, accessExp, refreshExp time.Duration, secret, refreshSecret string, ip string) (*JwtToken, error) {
	aExp := time.Now().Add(accessExp * time.Second).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token":    val,
		"exp":      aExp,
		"username": username,
		"ip":       ip,
	})
	aToken, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	rExp := time.Now().Add(refreshExp * time.Second).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token":    val,
		"exp":      rExp,
		"username": username,
		"ip":       ip,
	})
	rToken, err := refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return nil, err
	}
	return &JwtToken{
		AccessExp:    aExp,
		AccessToken:  aToken,
		RefreshExp:   rExp,
		RefreshToken: rToken,
	}, nil
}
func ParseToken(tokenString string, secret string, refreshSecret, ip string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		logs.LG.Info("token parse error ", zap.Error(err))
		return "", "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		logs.LG.Info("token ", zap.Any("parse", claims))
		val := claims["token"].(string)
		exp := int64(claims["exp"].(float64))
		username := claims["username"].(string)
		if exp <= time.Now().Unix() {
			return "", "", errors.New("token expired")
		}
		if claims["ip"] != ip {
			return "", "", errors.New("ip不合法")
		}
		return val, username, nil
	} else {
		logs.LG.Info("token cannot parse ")
		return "", "", nil
	}

}
func ParseRefreshToken(refreshToken string, refreshSecret, ip string) (string, string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(refreshSecret), nil
	})
	if err != nil {
		logs.LG.Info("refresh token parse error ", zap.Error(err))
		return "", "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		logs.LG.Info("refresh token ", zap.Any("parse", claims))
		val := claims["token"].(string)
		username := claims["username"].(string)
		exp := int64(claims["exp"].(float64))
		if exp <= time.Now().Unix() {
			return "", "", errors.New("token expired")
		}
		if claims["ip"] != ip {
			return "", "", errors.New("ip不合法")
		}
		return val, username, nil
	} else {
		logs.LG.Info("refresh token cannot parse ")
		return "", "", nil
	}

}
func RefreshToken(refreshToken string, accessExp, refreshExp time.Duration, secret, refreshSecret string, ip string) (*JwtToken, error) {
	id, username, err := ParseRefreshToken(refreshToken, refreshSecret, ip)
	if err != nil {
		return nil, err
	}
	return CreateToken(id, username, accessExp, refreshExp, secret, refreshSecret, ip)
}

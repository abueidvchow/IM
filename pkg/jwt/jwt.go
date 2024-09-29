package jwt

import (
	"IM/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const TokenExpireDuration = time.Hour * 24 * 7
const RefreshTokenExpireDuration = time.Hour * 24 * 30

var mySecret = []byte("IMHAPPY")

// MyClaims自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
type MyClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken生成Access Token和Refresh Token
func GentToken(userId int64, username string) (aToken, rToken string, err error) {
	c := MyClaims{
		UserId:   userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    config.Conf.Name,
		},
	}
	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	aToken, err = token.SignedString(mySecret)
	if err != nil {
		return "", "", err
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RefreshTokenExpireDuration).Unix(),
		Issuer:    config.Conf.Name,
	})

	rToken, err = token.SignedString(mySecret)
	if err != nil {
		return "", "", err
	}

	// 使用指定的secret签名并获得完整的编码后的字符串token
	return aToken, rToken, nil
}

// 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return token.Claims.(*MyClaims), nil
}

func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	_, err = jwt.Parse(rToken, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})

	if err != nil {
		return "", "", err
	}

	c := &MyClaims{}
	_, err = jwt.ParseWithClaims(aToken, c, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})

	v, _ := err.(*jwt.ValidationError)
	if v.Errors == jwt.ValidationErrorExpired {
		return GentToken(c.UserId, c.Username)
	}
	return
}

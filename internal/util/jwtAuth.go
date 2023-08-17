package util

import (
	"doushen_by_liujun/internal/common"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// CustomClaims /**
// CustomClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type CustomClaims struct {
	UserID               int64  `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        // 内嵌标准的声明
}

// GenToken 生成JWT
func GenToken(userId int64, username string) (string, error) {
	// 创建一个我们自己声明的数据
	claims := CustomClaims{
		userId,
		username, // 自定义字段
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * common.JwtExpire)), // 定义过期时间
			Issuer:    common.OrganizationName,                                          // 签发人
			ID:        common.OrganizationName,                                          // 编号
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成签名字符串
	return token.SignedString([]byte(common.JwtSecret))
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	//var mc = new(CustomClaims)
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(common.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}

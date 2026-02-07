package myjwt

import (
	"GopherAI/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 用于生成JWT Token
func GenerateToken(id int64, username string) (string, error) {
	claims := Claims{
		ID:       id,       //用户唯一标识
		Username: username, //用户名
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.GetConfig().ExpireDuration) * time.Hour)), //过期时间
			Issuer:    config.GetConfig().Issuer,                                                                        //签发者
			Subject:   config.GetConfig().Subject,                                                                       //主题
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                                   //签发时间
			//该函数会把time.Time转成JWT规范的时间格式

		},
	}

	//生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//使用HS256签名算法创建一个新的Token
	return token.SignedString([]byte(config.GetConfig().Key))
	//使用对称密钥进行签名，生成最终的JWT字符串
	//最终的字符串格式为header.payload.signature
}

// ParseToken解析Token
func ParseToken(token string) (string, bool) {
	claims := new(Claims)
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		//ParseWithClaims会解析token的结构，校验签名，反序列化Payload到claims
		return []byte(config.GetConfig().Key), nil
		//返回校验签名时的密钥，必须和生成token时用的key完全一致
	})
	if !t.Valid || err != nil || claims.Username == "" {
		return "", false
	}
	//t.Valid:JWT结构，签名，exp等是否通过校验

	return claims.Username, true
	//校验通过，返回业务所需要的信息
}

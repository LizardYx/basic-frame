package middleware

import (
	"basic-frame/util/common"
	"basic-frame/util/ginx"
	"basic-frame/util/ginx/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type MyClaims struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint64, userName string) (string, error) {
	newClaims := MyClaims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(common.SysConfig.JWTAuth.Expired))), // 过期时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                                                    // 生效时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                                    // 签发时间
		},
	}
	tokenInfo := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return tokenInfo.SignedString([]byte(common.SysConfig.JWTAuth.SecretKey))
}

// ParseToken 解析 JWT token字符串
func ParseToken(token string) (*MyClaims, error) {
	tokenInfo, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(common.SysConfig.JWTAuth.SecretKey), nil
	})
	if err != nil {
		if validationError, ok := err.(*jwt.ValidationError); ok {
			if validationError.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("that's not even a token")
			} else if validationError.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else if validationError.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token not active yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		}
	}
	if claims, ok := tokenInfo.Claims.(*MyClaims); ok && tokenInfo.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}

// RefreshToken 刷新JWT token字符串
func RefreshToken(token string) (string, error) {
	myClaims, err := ParseToken(token)
	if err != nil {
		return "", errors.WithMessage(err, "销毁Token失败")
	}

	// 更新Claims过期时间
	NewExpiresTime := time.Now().Add(time.Second * time.Duration(common.SysConfig.JWTAuth.Expired))
	myClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(NewExpiresTime)
	tokenInfo := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
	return tokenInfo.SignedString([]byte(common.SysConfig.JWTAuth.SecretKey))
}

// UserJWTAuth 用户登陆验证
func UserJWTAuth(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		tokenInfo := ginx.GetToken(c)
		if tokenInfo == "" {
			ginx.ResError(c, "", errors.NewResponse(401, 401, "Please login in"))
			c.Abort()
			return
		}
		myClaims, err := ParseToken(tokenInfo)
		if err != nil {
			ginx.ResError(c, "", err)
			c.Abort()
			return
		}
		ginx.SetUserInfo(c, myClaims.UserID, myClaims.UserName)
		c.Next()
	}
}

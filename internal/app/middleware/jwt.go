package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTValidator сервис
type JWTValidator struct {
	Claims Claims
	domain string
}

// Claims структура
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// TokenExp время жизни токена
const TokenExp = 3 * time.Hour

// SecretKey ключ шифрования
const SecretKey = "0N#6Ke|+OR:(`G;"

// UserIDProperty в каком поле храниться ID пользователя
const UserIDProperty = "UserId"

// NewJWTValidator конструктор
func NewJWTValidator(domain string) *JWTValidator {
	return &JWTValidator{domain: domain}
}

// Handle обработка проверки токена
func (JWTValidator *JWTValidator) Handle(ctx *gin.Context) {
	var err error
	token, _ := ctx.Cookie("Authorization")

	if token == "" {
		token, err = JWTValidator.BuildJWTString(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, err)
		}
	} else {
		errorValidateToken := JWTValidator.ValidateJWT(ctx, token)
		if errorValidateToken != nil {
			token, err = JWTValidator.BuildJWTString(ctx)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusForbidden, err)
			}
		}
	}

	claims := &Claims{}
	// парсим из строки токена tokenString в структуру claims
	jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	ctx.Set(UserIDProperty, claims.UserID)
	ctx.Header("Authorization", token)

	ctx.SetCookie("Authorization", token, 100000, "*", JWTValidator.domain, false, true)

}

// ValidateJWT проверка токена
func (JWTValidator *JWTValidator) ValidateJWT(ctx *gin.Context, tokenString string) error {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("token is not valid")
	}

	ctx.Set(UserIDProperty, claims.UserID)

	return nil
}

// BuildJWTString формирование токена
func (JWTValidator *JWTValidator) BuildJWTString(ctx *gin.Context) (string, error) {
	userID := JWTValidator.getNewUserID()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	ctx.Set(UserIDProperty, userID)

	return tokenString, nil
}

func (JWTValidator *JWTValidator) getNewUserID() string {
	return uuid.New().String()
}

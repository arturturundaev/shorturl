package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type JWTValidator struct {
	Claims Claims
	domain string
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const TOKEN_EXP = 3 * time.Hour
const SECRET_KEY = "0N#6Ke|+OR:(`G;"
const USER_ID_PROPERTY = "UserId"

func NewJWTValidator(domain string) *JWTValidator {
	return &JWTValidator{domain: domain}
}

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
		return []byte(SECRET_KEY), nil
	})

	ctx.Set(USER_ID_PROPERTY, claims.UserID)
	ctx.Header("Authorization", token)

	ctx.SetCookie("Authorization", token, 100000, "*", JWTValidator.domain, false, true)

}

func (JWTValidator *JWTValidator) ValidateJWT(ctx *gin.Context, tokenString string) error {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
		})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("Token is not valid")
	}

	ctx.Set(USER_ID_PROPERTY, claims.UserID)

	return nil
}

func (JWTValidator *JWTValidator) BuildJWTString(ctx *gin.Context) (string, error) {
	userId := JWTValidator.getNewUserId()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: userId,
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	ctx.Set(USER_ID_PROPERTY, userId)

	return tokenString, nil
}

func (JWTValidator *JWTValidator) getNewUserId() string {
	return uuid.New().String()
}

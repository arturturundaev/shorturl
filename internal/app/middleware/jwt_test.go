package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJWTValidator_Handle(t *testing.T) {

	tests := []struct {
		name  string
		token string
		err   error
	}{
		{
			name:  "init token",
			token: "",
			err:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwt := NewJWTValidator("/test")
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
			ctx.SetCookie("Authorization", tt.token, 100000, "*", "", false, true)
			jwt.Handle(ctx)

			token, _ := ctx.Cookie("Authorization")
			err := jwt.ValidateJWT(ctx, token)
			assert.Error(t, err, tt.err)
		})
	}
}

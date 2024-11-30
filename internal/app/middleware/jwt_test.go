package middleware

import (
	"github.com/gin-gonic/gin"
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
		{
			name:  "bad token",
			token: "asdas",
			err:   nil,
		},
		{
			name:  "valid token",
			token: "",
			err:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwt := NewJWTValidator("/test")
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
			jwt.Handle(ctx)
			if tt.name == "valid token" {
				tt.token, _ = jwt.BuildJWTString(ctx)
			}
			ctx.Request.AddCookie(&http.Cookie{
				Name:   "Authorization",
				Value:  tt.token,
				Path:   "/test",
				Domain: "/test",
			})
			jwt.Handle(ctx)
		})
	}
}

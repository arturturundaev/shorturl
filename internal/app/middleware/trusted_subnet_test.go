package middleware

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_isIPInSubnet(t *testing.T) {
	type args struct {
		ip      string
		subnets string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "IP exists",
			args: args{
				ip:      "127.45.45.1",
				subnets: "192.0.2.1/24",
			},
			want: false,
		},
		{
			name: "IP not exists",
			args: args{
				ip:      "192.168.0.5",
				subnets: "192.168.0.1/24",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ipnetArr []*net.IPNet
			_, ipnet, _ := net.ParseCIDR(tt.args.subnets)
			ipnetArr = append(ipnetArr, ipnet)
			if got := isIPInSubnet(tt.args.ip, ipnetArr); got != tt.want {
				t.Errorf("isIPInSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrustedSubnetMiddleware_Handle(t *testing.T) {
	type args struct {
		ip      string
		subnets string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "IP exists",
			args: args{
				ip:      "127.45.45.1",
				subnets: "192.0.2.1/24",
			},
			want: http.StatusForbidden,
		},
		{
			name: "IP not exists",
			args: args{
				ip:      "192.168.0.5",
				subnets: "192.168.0.1/24",
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ipnetArr []*net.IPNet
			_, ipnet, _ := net.ParseCIDR(tt.args.subnets)
			ipnetArr = append(ipnetArr, ipnet)

			h := &TrustedSubnetMiddleware{subnets: ipnetArr}
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = &http.Request{RemoteAddr: tt.args.ip + ":80"}
			h.Handle(ctx)

			if ctx.Writer.Status() != tt.want {
				t.Error(tt.name + " Bad request")
			}
		})
	}
}

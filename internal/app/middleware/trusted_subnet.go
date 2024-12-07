package middleware

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

type TrustedSubnetMiddleware struct {
	subnets []*net.IPNet
}

// NewTrustedSubnetMiddleware
func NewTrustedSubnetMiddleware(subnets []*net.IPNet) *TrustedSubnetMiddleware {
	return &TrustedSubnetMiddleware{subnets: subnets}
}

// Handle
func (h *TrustedSubnetMiddleware) Handle(c *gin.Context) {
	ip := c.ClientIP()
	if !isIPInSubnet(ip, h.subnets) {
		c.AbortWithStatus(http.StatusForbidden)
	}
	c.Next()
}

// IsIPInSubnet
func isIPInSubnet(ip string, subnets []*net.IPNet) bool {
	ipNet := net.ParseIP(ip)
	for _, subnet := range subnets {
		if subnet.Contains(ipNet) {
			return true
		}
	}
	return false
}

package logging

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// GCPが発行するTraceIDを取得する
func GetRequestTraceID(c *gin.Context) string {
	header := c.Request.Header.Get("X-Cloud-Trace-Context")
	if header == "" {
		return ""
	}
	parts := strings.Split(header, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

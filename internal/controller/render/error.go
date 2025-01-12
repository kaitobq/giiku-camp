package render

import (
	"encoding/json"
	"giiku-camp/internal/infra/logging"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Message string `json:"message" example:"error message"`
	Code    int    `json:"code" example:"0"`
	Status  int    `json:"status" example:"400"`
}

func ErrorJSON(c *gin.Context, err string, status int) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)

	if err := json.NewEncoder(c.Writer).Encode(&Error{
		Message: err,
		Status:  status,
	}); err != nil {
		logging.Warnf(c, "ErrorJSON json.NewEncoder %v", err)
	}
}

// エラーコードを返すときに使う
func ErrorCodeJSON(c *gin.Context, err string, status, code int) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)

	if err := json.NewEncoder(c.Writer).Encode(&Error{
		Message: err,
		Code:    code,
		Status:  status,
	}); err != nil {
		logging.Warnf(c, "ErrorCodeJSON json.NewEncoder %v", err)
	}
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealCheck(c *gin.Context) {

	msg := `# TYPE health_info gauge
health_info{status="ok", name="adp-backend", namespace="devops"} 0
`
	c.String(http.StatusOK, msg)
}

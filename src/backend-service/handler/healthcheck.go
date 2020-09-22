package handler

import "github.com/gin-gonic/gin"

func HealCheck(c *gin.Context) {

	c.String(200, "# TYPE health_info gauge\n"+"health_info{status=\"ok\", name=\"service-adp-deploy\", namespace=\"devops\"} 0\n")
}

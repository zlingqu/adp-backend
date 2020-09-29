package server

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
)

func Run(router *gin.Engine, port string, debug bool) {

	fmt.Println("port:", port)

	runtime.GOMAXPROCS(runtime.NumCPU())

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Run(":" + port)

}

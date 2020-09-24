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

	// srv := &http.Server{
	// 	Addr:    ":" + port,
	// 	Handler: router,
	// }

	// go func() {
	// 	// service connections
	// 	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("listen: %s\n", err)
	// 	}
	// }()

	// // Wait for interrupt signal to gracefully shutdown the server with
	// // a timeout of 10 seconds.
	// quit := make(chan os.Signal)
	// signal.Notify(quit, os.Interrupt)
	// <-quit
	// log.Println("Shutdown Server ...")

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// if err := srv.Shutdown(ctx); err != nil {
	// 	log.Fatal("Server Shutdown:", err)
	// }
	// log.Println("Server exiting")

	//pid := fmt.Sprintf("%d", os.Getpid())
	//_, openErr := os.OpenFile("pid", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if openErr == nil {
	//	_ = ioutil.WriteFile("pid", []byte(pid), 0)
	//}
}

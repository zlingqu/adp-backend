package server

import (
	"context"
	"net/http"

	// "encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"gopkg.in/antage/eventsource.v1"

	// "gopkg.in/antage/eventsource.v1"
	"log"
	// "net/http"
	"os"
	"os/signal"
	"time"
)

func RunApi(router *gin.Engine, port string, debug bool) *http.Server {
	fmt.Println("port:", port)

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return srv
}

func RunEventSource(port string) {
	fmt.Println("event source port:", port)
	es := eventsource.New(&eventsource.Settings{
		Timeout:        10 * time.Second,
		CloseOnTimeout: true,
		IdleTimeout:    50 * time.Minute,
	}, func(req *http.Request) [][]byte {
		return [][]byte{
			[]byte("X-Accel-Buffering: no"),
			[]byte("Access-Control-Allow-Origin: *"),
		}
	})
	defer es.Close()

	http.Handle("/events", es)

	// deployEnv := config.GetEnv().DeployEnv

	// go func() {
	// 	var appId uint
	// 	log.Println(appId)
	// 	for {
	// 		// result, err := handler.SearchAdpResultInfo(deployEnv, "x4c-mgepap-tp-service")
	// 		result, err := handler.SearchAdpResultInfo(deployEnv, "cpm")
	// 		if err != nil {
	// 			log.Println("错误：", err)
	// 			time.Sleep(2 * time.Second)
	// 			continue
	// 		}
	// 		if result.ID > appId {
	// 			r, _ := json.Marshal(result)
	// 			es.SendEventMessage(string(r), "", "")
	// 		}
	// 		appId = result.ID

	// 		log.Printf("Hello has been sent (consumers: %d)", es.ConsumersCount())
	// 		log.Printf("当前appId为：%d\n", appId)
	// 		time.Sleep(2 * time.Second)
	// 	}
	// }()

	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal("listen event source exception: ", err)
		}
	}()

}

func WaitInterrupt(apiServer *http.Server) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

package server

import (
	"app-deploy-platform/backend-service/config"
	"app-deploy-platform/backend-service/handler"
	"context"
	"encoding/json"

	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/antage/eventsource.v1"
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

	deployEnv := config.GetEnv().DeployEnv

	go func() {
		for {
			result, err := handler.SearchAdpResultInfo(deployEnv, "x4c-mgepap-tp-service")
			if err != nil {
				log.Println("错误：", err)
				time.Sleep(10 * time.Second)
				continue
			}
			r, _ := json.Marshal(result)
			es.SendEventMessage(string(r), "", "")
			// log.Printf("respData: %v", string(r))
			log.Printf("Hello has been sent (consumers: %d)", es.ConsumersCount())
			time.Sleep(20 * time.Second)
		}
	}()
	log.Fatal(http.ListenAndServe(":"+port, nil))
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

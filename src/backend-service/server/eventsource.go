package server

import (
	"app-deploy-platform/backend-service/model"
	"app-deploy-platform/backend-service/service"
	"encoding/json"
	"fmt"
	"gopkg.in/antage/eventsource.v1"
	"log"
	"net/http"
	"time"
)

var (
	channelSize = 10
	cache       chan model.Result
)

func init() {
	cache = make(chan model.Result, channelSize)
}

func PushResult(result model.Result) {
	cache <- result
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

	go func() {
		sleepTime := 2 * time.Second
		for res := range cache {
			deploy := service.GetDeployByResult(res)

			if deploy.ID == 0 {
				log.Printf("错误：%v\n", deploy)
				time.Sleep(sleepTime)
				continue
			}

			r, _ := json.Marshal(deploy)
			es.SendEventMessage(string(r), "", "")
			log.Printf("Hello has been sent (consumers: %d)", es.ConsumersCount())
			time.Sleep(sleepTime)
		}
	}()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

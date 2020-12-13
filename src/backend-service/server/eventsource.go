package server

import (
	"app-deploy-platform/backend-service/model"
	"app-deploy-platform/backend-service/service"
	"container/list"
	"fmt"
	"gopkg.in/antage/eventsource.v1"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	cache      = list.New()
	cacheMutex sync.Mutex
)

func PushResult(result model.Result) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	cache.PushBack(result)
}

func getAllResult() (results []model.Result) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	var next *list.Element
	for e := cache.Front(); e != nil; e = next {
		next = e.Next()
		results = append(results, cache.Remove(e).(model.Result))
	}
	return
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
		timer := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-timer.C:
				consumersCount := es.ConsumersCount()
				if consumersCount > 0 {
					jenkinsBuildTokens := ""
					for _, t := range service.GetDeployByResult(getAllResult()...) {
						jenkinsBuildTokens += t.JenkinsBuildToken + ","
					}

					if jenkinsBuildTokens != "" {
						es.SendEventMessage(jenkinsBuildTokens, "", "")
						log.Printf("Hello has been sent (consumers: %d)", consumersCount)
					}
				}

			}
		}
	}()

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

package server

import (
	"app-deploy-platform/backend-service/service"
	"fmt"
	"gopkg.in/antage/eventsource.v1"
	"log"
	"net/http"
	"time"
)

var (
	scheduleTime                 = time.Second * 5
	lastSearchResultTimeInterval = time.Second * 60
)

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
		timer := time.NewTicker(scheduleTime)
		for {
			select {
			case t := <-timer.C:
				consumersCount := es.ConsumersCount()
				if consumersCount > 0 {
					jenkinsBuildTokens := ""

					results := service.GetResultByCreateTime(t.Add(-lastSearchResultTimeInterval))
					for _, t := range service.GetDeployByResult(results...) {
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

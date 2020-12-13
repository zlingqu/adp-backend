package service

import (
	m "app-deploy-platform/backend-service/model"
	"errors"
	"log"
	"time"
)

var (
	allEnvWithIdAndName = make(map[uint]string)
	allEnvWithNameAndId = make(map[string]uint)
)

func init() {
	cacheEnv()

	go func() {
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticker.C:
				cacheEnv()
			}
		}
	}()
}

func cacheEnv() {
	var envs []m.Env
	m.DB.Find(&envs)
	log.Println("缓存env")
	for _, env := range envs {
		allEnvWithIdAndName[env.ID] = env.Name
		allEnvWithNameAndId[env.Name] = env.ID
	}
}

func GetEnvByID(id string) (env m.Env, err error) {
	RowsAffected := m.DB.First(&env, id).RowsAffected

	if RowsAffected == 0 {
		err = errors.New("Not Found " + id)
	}

	return
}

func GetEnvIDByName(name string) uint {
	return allEnvWithNameAndId[name]
}

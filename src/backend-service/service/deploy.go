package service

import (
	"app-deploy-platform/backend-service/model"
	"strings"
)

func GetDeployByResult(results ...model.Result) (deploys []model.Deploy) {
	l := len(results)
	wheres, args := make([]string, l), make([]interface{}, 0)
	if l > 0 {
		for i, result := range results {
			wheres[i] = "(name = ? and env_id = ? and version = ?)"
			args = append(args, result.Name, GetEnvIDByName(result.DeployEnv), result.Version)
		}
		model.DB.Where(strings.Join(wheres, " and "), args...).Find(&deploys)
	}
	return
}

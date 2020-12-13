package service

import (
	"app-deploy-platform/backend-service/model"
	"strings"
)

func GetDeployByResult(results ...model.Result) (deploys []model.Deploy) {
	l := len(results)
	wheres, args := make([]string, l), make([]interface{}, 0)
	if l > 0 {
		//  result version传入的是0.0.0没有意义
		for i, result := range results {
			//wheres[i] = "(name = ? and env_id = ? and version = ?)"
			//args = append(args, result.Name, GetEnvIDByName(result.DeployEnv), result.Version)
			wheres[i] = "(name = ? and env_id = ?)"
			args = append(args, result.Name, GetEnvIDByName(result.DeployEnv))
		}
		model.DB.Where(strings.Join(wheres, " and "), args...).Find(&deploys)
	}
	return
}

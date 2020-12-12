package service

import "app-deploy-platform/backend-service/model"

func GetDeployByResult(result model.Result) (deploy model.Deploy) {
	model.DB.Where("name = ? and env_id = ? and version = ?", result.Name, GetEnvIDByName(result.Name), result.Version).Limit(1).Find(&deploy)
	return
}

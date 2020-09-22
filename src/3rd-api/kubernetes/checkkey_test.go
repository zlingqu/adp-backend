package main

import (
	"errors"
	"fmt"
	log "github.com/zuoshenglo/libs/logs/logrus"
	"github.com/zuoshenglo/tools/file"
	"os"
	"testing"
	"time"
)

const keyServer = "http://service-k8s-key-manager.dm-ai.cn"

func TestCheckKey(t *testing.T) {
	envList := []string{"prd", "dev", "test", "stage", "jenkins", "mlcloud-dev"}
	dir, _ := os.Getwd()
	dir = dir + "/temp/"

	for _, v := range envList {
		log.Info(fmt.Sprintf("开始下载环境：%s的证书文件。", v))
		url := fmt.Sprintf(keyServer+"/api/v1/get-k8s-key-file?env=%s", v)
		file.DownFile(url, dir+v)
		log.Info(fmt.Sprintf("下载环境：%s的证书文件成功。", v))
	}

	go func() {
		for {
			time.Sleep(30 * time.Second)
			for _, v := range envList {
				log.Info(fmt.Sprintf("下载环境%s的证书文件%s-new。", v, v))
				url := fmt.Sprintf("http://service-k8s-key-manager.dm-ai.cn/api/v1/get-k8s-key-file?env=%s", v)
				file.DownFile(url, dir+v+"-new")
				log.Info(fmt.Sprintf("下载环境%s的证书文件%s-new, 完成。", v, v))
				if file.GetFileMd5Sum(dir+v) != file.GetFileMd5Sum(dir+v+"-new") {
					panic(errors.New("证书文件变更，触发查询规则，应用退出！"))
				}
			}
		}
	}()

	select {}
}

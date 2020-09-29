package model

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type PodsStatus struct {
	KubernetesConfigFile string                `json:"kubernetes_config_file"`
	ClientSet            *kubernetes.Clientset `json:"client_set"`
	Res                  struct {
		Res    string `json:"res"`
		Msg    string `json:"msg"`
		Status string `json:"status"`
	} `json:"res"`
}

func NewPodsStatus() *PodsStatus {
	p := &PodsStatus{}
	p.Res.Res = "fail"
	p.Res.Msg = "未满足到查询条件，此为初始化信息，错误。"
	p.Res.Status = "fail"
	return p
}

func (p *PodsStatus) GetKubernetesClient() {
	config, err := clientcmd.BuildConfigFromFlags("", p.KubernetesConfigFile)
	if err != nil {
		log.Error("初始化k8s的客户端的配置文件出错：", err)
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Error("获得k8s的clientset出错：", err)
		panic(err.Error())
	}
	p.ClientSet = clientset
}

func (p *PodsStatus) GetPodsInfo(namespace string, appName string, imageSha string) {
	pods, err := p.ClientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", appName),
	})

	if err != nil {
		log.Error("List失败")
		panic(err.Error())
	}

	if len(pods.Items) == 0 {
		log.Error("查询pod状态失败！")
		p.Res.Res = "fail"
		p.Res.Msg = "查不到此pod，请检查appName、namespace是否匹配"
		p.Res.Status = "fail"
		return
	}

	var containerRestartCount int32 = 0

	for _, pod := range pods.Items { //根据标签查询的，所以可能查出多个pod，逐个判断

		// pod.Status.Conditions,正常启动后分为4个阶段，Initialized、Ready、ContainersReady、PodScheduled。
		if len(pod.Status.Conditions) == 1 { //异常情况
			podInfo := pod.Status.Conditions[0]
			if pod.Status.Phase == "Pending" && podInfo.Type == "PodScheduled" && podInfo.Status == "False" && podInfo.Reason == "Unschedulable" {
				p.Res.Res = "fail"
				p.Res.Status = "fail"
				p.Res.Msg = fmt.Sprintf("因集群资源限制，无法调度pod， %s", podInfo.Message)
				return
			}
		}
		if len(pod.Status.Conditions) == 4 && pod.Status.Phase == "Pending" {
			for _, container := range pod.Status.ContainerStatuses {
				if matchIstio, _ := regexp.MatchString(".*istio/.*", container.Image); matchIstio { //跳过istio容器
					continue
				}
				switch container.State.Waiting.Reason {
				case "ContainerCreating":
					p.Res.Res = "ok"
					p.Res.Status = "continue"
					p.Res.Msg = fmt.Sprintf("pod中的容器正在创建，可能是正在下载镜像,%s,请稍等", container.Image)
					return
				case "ImagePullBackOff":
					p.Res.Res = "ok"
					p.Res.Status = "continue"
					p.Res.Msg = fmt.Sprintf("下载镜像异常, %s, 请检查", container.State.Waiting.Message)
					return
				case "PodInitializing":
					p.Res.Res = "ok"
					p.Res.Status = "continue"
					p.Res.Msg = fmt.Sprintf("正在进行pod的初始化，PodInitializing, 请稍等")
					return

				default:
					p.Res.Res = "fail"
					p.Res.Status = "fail"
					p.Res.Msg = fmt.Sprintf("pod状态一直处于Pending状态，异常, 请手动检查原因")
					return

				}
			}
		}

		if len(pod.Status.Conditions) == 4 && pod.Status.Phase == "Running" {
			for _, container := range pod.Status.ContainerStatuses {
				if matchIstio, _ := regexp.MatchString(".*istio/.*", container.Image); matchIstio { //跳过istio容器
					continue
				}
				if container.State.Waiting != nil {

					p.Res.Res = "fail"
					p.Res.Status = "fail"
					p.Res.Msg = fmt.Sprintf("容器启动失败%s，%s", container.State.Waiting.Reason, container.State.Waiting.Message)
					return

				}

				// 针对 container.State 来就行判断，排除正在消亡的pod
				if container.State.Terminated != nil {
					p.Res.Res = "ok"
					p.Res.Status = "continue"
					p.Res.Msg = fmt.Sprintf("查询到正在Terminated的容器镜像：%s,请稍等", container.Image)
					return
				}

				if !strings.Contains(container.ImageID, imageSha) { //镜像id hash值不匹配
					fmt.Println(container.ImageID, imageSha)
					p.Res.Res = "ok"
					p.Res.Status = "continue"
					p.Res.Msg = fmt.Sprintf("镜像id不匹配, %s,请稍等", container.Image)
					return
				}

				if container.Ready == false && container.RestartCount > containerRestartCount && container.LastTerminationState.Terminated.Reason == "Error" {
					p.Res.Res = "fail"
					p.Res.Status = "fail"
					p.Res.Msg = fmt.Sprintf("部署在k8s中的pod启动次数大于%d，请检查应用的日志, 确定启动失败的原因, 镜像：%s,", containerRestartCount, container.Image)
					return
				}
				if container.Ready == false { //正常情况
					p.Res.Res = "ok"
					p.Res.Status = "continue"
					p.Res.Msg = fmt.Sprintf("pod中的容器创建完成，正在启动及进行启动过程中的端口检查, %s,请稍等", container.Image)
					return
				}
			}

			for _, podCondition := range pod.Status.Conditions {
				if podCondition.Status == "False" {
					p.Res.Res = "fail"
					p.Res.Msg = fmt.Sprintf("%s, 请确定端口设置正常，请检查应用日志, podName: %s", podCondition.Message, pod.Name)
					p.Res.Status = "fail"
					return
				}
				p.Res.Res = "ok"
				p.Res.Msg = "部署成功"
				p.Res.Status = "ok"

			}
		}
	}
}

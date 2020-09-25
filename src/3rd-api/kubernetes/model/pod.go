package model

import (
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
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
		log.Error("查询pod状态失败！")
		panic(err.Error())
	}

	if len(pods.Items) == 0 {
		p.Res.Res = "fail"
		p.Res.Msg = "未在集群中查询到此服务的pod，无法检查状态。"
		p.Res.Status = "fail"
		return
	}

	// 检查过程中的兼容问题处理，如果检查的pod中包含有istiode容器，则在running的过程中，启动次数可能为2次。
	var containerRestartCount int32 = 0

	// if pods.Items == 0  没有在k8s集群中找到对应的应用的pod的信息。
	for _, pod := range pods.Items {
		log.WithFields(logrus.Fields{
			"name":                             pod.ObjectMeta.Name,
			"create_time":                      pod.ObjectMeta.CreationTimestamp,
			"pod.Status.Phase":                 pod.Status.Phase,
			"pod.Status.Conditions":            pod.Status.Conditions,
			"pod.Status.Message":               pod.Status.Message,
			"pod.Status.Reason":                pod.Status.Reason,
			"pod.Status.NominatedNodeName":     pod.Status.NominatedNodeName,
			"pod.Status.HostIP":                pod.Status.HostIP,
			"pod.Status.PodIP":                 pod.Status.PodIP,
			"pod.Status.StartTime":             pod.Status.StartTime,
			"pod.Status.InitContainerStatuses": pod.Status.InitContainerStatuses,
			"pod.Status.ContainerStatuses":     pod.Status.ContainerStatuses,
			"pod.Status.QOSClass":              pod.Status.QOSClass,
		})

		// 处理因为资源不够，导致的pod无法调度的情况。
		if len(pod.Status.Conditions) == 1 {
			podInfo := pod.Status.Conditions[0]
			if pod.Status.Phase == "Pending" && podInfo.Type == "PodScheduled" && podInfo.Status == "False" && podInfo.Reason == "Unschedulable" {
				p.Res.Res = "fail"
				p.Res.Status = "fail"
				p.Res.Msg = fmt.Sprintf("因集群资源限制，无法调度pod， %s", podInfo.Message)
				return
			}
		}

		// 处理部署的时候，某些镜像很大的服务，正在下载镜像的情况
		if len(pod.Status.Conditions) == 4 {
			if pod.Status.Phase == "Pending" {
				for _, container := range pod.Status.ContainerStatuses {

					if container.State.Waiting.Reason == "ContainerCreating" && container.RestartCount == 0 {
						p.Res.Res = "ok"
						p.Res.Status = "continue"
						p.Res.Msg = fmt.Sprintf("pod中的容器正在创建，可能是正在下载镜像,%s,请稍等", container.Image)
						return
					}

					if container.State.Waiting.Reason == "ImagePullBackOff" && container.RestartCount == 0 {
						p.Res.Res = "ok"
						p.Res.Status = "continue"
						p.Res.Msg = fmt.Sprintf("下载镜像异常, %s, 请检查", container.State.Waiting.Message)
						return
					}

					if container.State.Waiting.Reason == "PodInitializing" && container.RestartCount == 0 {
						p.Res.Res = "ok"
						p.Res.Status = "continue"
						p.Res.Msg = fmt.Sprintf("正在进行pod的初始化，PodInitializing, 请稍等")
						return
					}
				}

				// 除了上面的2种情况还处于Pending的状态的话，就认为失败了。
				p.Res.Res = "fail"
				p.Res.Status = "fail"
				p.Res.Msg = fmt.Sprintf("pod状态一直处于Pending状态，异常, 请手动检查检查")
				return
			}

			if pod.Status.Phase == "Running" {
				// 如果循环检查容器的状态中，检查到istio的容器, 则运行容器的最大重启次数为2
				for _, container := range pod.Status.ContainerStatuses {
					matchIstio, _ := regexp.MatchString(".*istio/.*", container.Image)
					if matchIstio {
						containerRestartCount = 2
						continue
					}
				}

				for _, container := range pod.Status.ContainerStatuses {

					// 针对 container.State 来就行判断，排除正在消亡的pod
					if container.State.Terminated != nil {
						p.Res.Res = "ok"
						p.Res.Status = "continue"
						p.Res.Msg = fmt.Sprintf("查询到正在Terminated的容器镜像：%s,请稍等", container.Image)
						return
					}

					// 如果是不是匹配istio 或者 传入的镜像地址，则不检查，检查匹配istio 和传入的镜像启动的容器，排除不同的构建镜像
					matchIstio, _ := regexp.MatchString(".*istio/.*", container.Image)
					matchAppImage, _ := regexp.MatchString(".*"+imageSha+".*", container.ImageID)
					if !matchAppImage && !matchIstio && imageSha != "" {
						continue
						//p.Res.Res = "ok"
						//p.Res.Status = "continue"
						//p.Res.Msg = fmt.Sprintf("检查到的容器不为部署的容器镜像, %s,请稍等", container.Image)
						//return
					}

					// ##
					if container.Ready == false && container.RestartCount == 0 {
						p.Res.Res = "ok"
						p.Res.Status = "continue"
						p.Res.Msg = fmt.Sprintf("pod中的容器创建完成，正在启动及进行启动过程中的端口检查, %s,请稍等", container.Image)
						return
					}

					// 针对含有iustio的服务，如果重启次数大于0 <= containerRestartCount,则continue
					if container.Ready == false && container.RestartCount > 0 && container.RestartCount <= containerRestartCount && container.LastTerminationState.Terminated.Reason == "Error" {
						p.Res.Res = "ok"
						p.Res.Status = "continue"
						p.Res.Msg = fmt.Sprintf("服务可能使用了istio，当前的重启测试为：%d, 镜像：%s", container.RestartCount, container.Image)
						return
					}

					if container.Ready == false && container.RestartCount > containerRestartCount && container.LastTerminationState.Terminated.Reason == "Error" {
						p.Res.Res = "fail"
						p.Res.Status = "fail"
						p.Res.Msg = fmt.Sprintf("部署在k8s中的pod启动失败次数大于%d，请检查应用的日志, 确定启动失败的原因, 镜像：%s,", containerRestartCount, container.Image)
						return
					}
				}

				for _, podCondition := range pod.Status.Conditions {
					if podCondition.Status == "False" {
						log.WithFields(logrus.Fields{
							"name":                             pod.ObjectMeta.Name,
							"create_time":                      pod.ObjectMeta.CreationTimestamp,
							"pod.Status.Phase":                 pod.Status.Phase,
							"pod.Status.Conditions":            pod.Status.Conditions,
							"pod.Status.Message":               pod.Status.Message,
							"pod.Status.Reason":                pod.Status.Reason,
							"pod.Status.NominatedNodeName":     pod.Status.NominatedNodeName,
							"pod.Status.HostIP":                pod.Status.HostIP,
							"pod.Status.PodIP":                 pod.Status.PodIP,
							"pod.Status.StartTime":             pod.Status.StartTime,
							"pod.Status.InitContainerStatuses": pod.Status.InitContainerStatuses,
							"pod.Status.ContainerStatuses":     pod.Status.ContainerStatuses,
							"pod.Status.QOSClass":              pod.Status.QOSClass,
						})
						p.Res.Res = "fail"
						p.Res.Msg = fmt.Sprintf("%s, 请确定端口设置正常，请检查应用日志, podName: %s", podCondition.Message, pod.Name)
						p.Res.Status = "fail"
						return
					} else {
						p.Res.Res = "ok"
						p.Res.Msg = "部署成功"
						p.Res.Status = "ok"
					}
				}
			}
		}
	}
}

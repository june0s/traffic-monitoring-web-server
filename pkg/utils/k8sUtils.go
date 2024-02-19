package utils

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"
	"time"

	"log"
)

var (
	ctxValue  = "/Users/june0/.kube/dev3-kubeconfig"
	clientSet *kubernetes.Clientset
	mtxClient *metricsclient.Clientset
)

// 초기화 함수 - 패키지 로드 시 가장 먼저 호출됨.
func init() {
	GetKubernetesClientSet(ctxValue)

}

func GetKubernetesClientSet(filepath string) {

	if config, err := GetRestConfig(ctxValue); err != nil {
		log.Fatal("Failed to get K8s rest config! check context file!!")
	} else if client, err := kubernetes.NewForConfig(config); err != nil {
		log.Fatal("Failed to kubernetes NewForConfig!")
	} else if mClient, err := metricsclient.NewForConfig(config); err != nil {
		log.Fatal("Failed to metricsclient NewForConfig!")
	} else {
		clientSet = client
		mtxClient = mClient
		fmt.Println("Get K8s, metrics client set -> ok!")
	}
}

func GetRestConfig(filepath string) (*rest.Config, error) {

	var configLoadingRules clientcmd.ClientConfigLoader
	if filepath == "" {
		configLoadingRules = clientcmd.NewDefaultClientConfigLoadingRules()
	} else {
		configLoadingRules = &clientcmd.ClientConfigLoadingRules{ExplicitPath: filepath}
	}

	if apiConfig, err := configLoadingRules.Load(); err != nil {
		return nil, err
	} else if restConfig, err := clientcmd.NewNonInteractiveClientConfig(*apiConfig, apiConfig.CurrentContext, &clientcmd.ConfigOverrides{}, nil).ClientConfig(); err != nil {
		return nil, err
	} else {
		return restConfig, nil
	}
}

func GetNamespaces() ([]string, error) {
	names := []string{}
	if nsList, err := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{}); err != nil {
		return nil, err
	} else {
		for _, ns := range nsList.Items {
			nm := ns.GetName()
			names = append(names, nm)
		}
	}
	return names, nil
}

func GetServices(namespace string) ([]string, error) {
	names := []string{}
	if svcList, err := clientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{}); err != nil {
		fmt.Println("get service error = ", err)
		return nil, err
	} else {
		for _, svc := range svcList.Items {
			names = append(names, svc.GetName())
		}
	}
	return names, nil
}

func GetWorkloads(namespace string) ([]string, error) {
	names := []string{}
	if wlList, err := clientSet.AppsV1().Deployments(namespace).List(context.TODO(), v1.ListOptions{}); err != nil {
		fmt.Println("get pods error = ", err)
		return nil, err
	} else {
		for _, item := range wlList.Items {
			names = append(names, item.GetName())
		}
	}
	return names, nil
}

// node metric 수집
// name, cpu(m cores), memory(Mi bytes)
type Point struct {
	Name      string `json:"name"`
	Timestamp string `json:"timestamp"`
	CPU       uint64 `json:"cpu"`
	Memory    uint64 `json:"memory"`
}

func GetNodeMetric() ([]Point, error) {
	nodeMetrics := &v1beta1.NodeMetricsList{}
	var err error

	fmt.Println("= = GetNodeMetric = = ")
	nodeMetrics, err = mtxClient.MetricsV1beta1().NodeMetricses().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		fmt.Println("Error scraping node metrics: ", err)
		return nil, err
	}
	//fmt.Println(nodeMetrics)

	resultList := []Point{}

	for _, v := range nodeMetrics.Items {
		fmt.Printf("name = %s, cpu = %v, mem = %v(%v), time = %v (%v)",
			v.GetName(), v.Usage.Cpu().MilliValue(),
			v.Usage.Memory().MilliValue()/1000, v.Usage.Memory().MilliValue()/1000/1024/1024,
			v.Timestamp, time.Now())
		fmt.Println()

		resultList = append(resultList, Point{
			Name:      v.Name,
			Timestamp: time.Now().String(),
			CPU:       uint64(v.Usage.Cpu().MilliValue()),
			Memory:    uint64(v.Usage.Memory().MilliValue() / 1000 / 1024 / 1024),
		})
	}
	return resultList, nil
}

package utils

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

var (
	ctxValue  = "/Users/june0/.kube/dev3-kubeconfig"
	clientSet *kubernetes.Clientset
)

// 초기화 함수 - 패키지 로드 시 가장 먼저 호출됨.
func init() {
	if cs, err := GetKubernetesClientSet(ctxValue); err != nil {
		log.Fatal("Failed to get K8s client set! check context file!!")
	} else {
		clientSet = cs
		fmt.Println("Get K8s client set -> ok!")
	}
}

func GetKubernetesClientSet(filepath string) (*kubernetes.Clientset, error) {

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
	} else if clientSet, err := kubernetes.NewForConfig(restConfig); err != nil {
		return nil, err
	} else {
		return clientSet, nil
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

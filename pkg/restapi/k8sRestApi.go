package restapi

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"traffic-monitoring-web-server/pkg/utils"
)

type ackNamespaces struct {
	Namespaces []string `json:"namespaces" example: "["kube-system", "istio-system",...]"`
}

type ackServices struct {
	Namespace string   `json:"namespace" example: "kube-system"`
	Services  []string `json:"services" example: "["kube-dns", "metrics-server",...]"`
}

type ackWorkloads struct {
	Namespace string   `json:"namespace" example: "kube-system"`
	Workloads []string `json:"workloads" example: "["kube-dns", "metrics-server",...]"`
}

type ackNodeMetrics struct {
	Metrics []utils.Point `json:"metrics"`
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodPatch},
			MaxAge:       12 * time.Hour,
		}))
	router.GET("/api/namespaces", getNamespaces)
	router.GET("/api/:namespace/services", getServices)
	router.GET("/api/:namespace/workloads", getWorkloads)
	router.GET("/api/node/metrics", getNodeMetrics)

	return router
}

// K8s cluster 의 네임스페이스 이름 리스트를 가져와 반환
func getNamespaces(c *gin.Context) {
	if namespaces, err := utils.GetNamespaces(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, GetServerErrResp(err))
	} else {
		resp := GetOkResp(ackNamespaces{namespaces})
		c.IndentedJSON(http.StatusOK, resp)
	}
}

// K8s cluster 의 특정 네임스페이스의 서비스 이름 리스트를 가져와 반환
func getServices(c *gin.Context) {
	namespace := c.Param("namespace")

	if services, err := utils.GetServices(namespace); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, GetServerErrResp(err))
	} else {
		resp := GetOkResp(ackServices{namespace, services})
		c.IndentedJSON(http.StatusOK, resp)
	}
}

// K8s cluster 의 특정 네임스페이스의 워크로드 이름 리스트를 가져와 반환
func getWorkloads(c *gin.Context) {
	namespace := c.Param("namespace")

	if workloads, err := utils.GetWorkloads(namespace); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, GetServerErrResp(err))
	} else {
		resp := GetOkResp(ackWorkloads{namespace, workloads})
		c.IndentedJSON(http.StatusOK, resp)
	}
}

func getNodeMetrics(c *gin.Context) {

	if metrics, err := utils.GetNodeMetric(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, GetServerErrResp(err))
	} else {
		resp := GetOkResp(ackNodeMetrics{metrics})
		c.IndentedJSON(http.StatusOK, resp)
	}
}

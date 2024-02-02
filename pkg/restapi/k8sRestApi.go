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

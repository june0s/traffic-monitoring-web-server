package main

import (
	"fmt"
	"traffic-monitoring-web-server/pkg/restapi"
)

func main() {
	fmt.Println("Hello, Traffic monitoring web server :D")

	router := restapi.SetupRouter()
	router.Run("localhost:9010")
}

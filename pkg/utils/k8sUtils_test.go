package utils

import (
	. "github.com/onsi/ginkgo"
	"testing"
)

func TestGetNodeMetric(t *testing.T) {

	GetNodeMetric()
	RunSpecs(t, "Get node metric test")
}

func TestGetWorkloads(t *testing.T) {
	ns := "travel-agency"
	if pods, err := GetWorkloads(ns); err != nil {
		println("error occur! err = ", err)
	} else {
		for _, po := range pods {
			println("pod = ", po)
		}
	}
}

package scraper

import (
	"k8s.io/klog"
	"os"
)

type systemInfo struct {
	Hostname string
}

var (
	SystemInformation = &systemInfo{}
)

func init() {
	var err error
	SystemInformation.Hostname, err = os.Hostname()
	if err != nil {
		klog.Exitf("Getting System Info Error: %v", err)
	}
}

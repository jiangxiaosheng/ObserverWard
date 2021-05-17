package main

import (
	"flag"
	"fmt"
	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog"
	"net/http"
	"observerward/pkg/prometheus"
	"strconv"
)

var (
	addr               = flag.String("address", "0.0.0.0", "The ip address to listen on for HTTP requests.")
	port               = flag.Int64("port", 9909, "The port to listen on for HTTP requests.")
	scrapeTimeInterval = flag.Int("interval", 5, "The time interval for scrapping metrics.")
)

func init() {
	flag.Parse()
}

func main() {
	defer nvml.Shutdown()

	listenAddr := *addr + ":" + strconv.FormatInt(*port, 10)
	fmt.Printf("observerward is listening on " + listenAddr)

	prometheus.Run(*scrapeTimeInterval)

	http.Handle("/metrics", promhttp.Handler())
	klog.Fatal(http.ListenAndServe(listenAddr, nil))
}

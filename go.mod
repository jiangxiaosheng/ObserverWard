module github.com/jiangxiaosheng/observerward

go 1.15

require (
	github.com/NVIDIA/gpu-monitoring-tools v0.0.0-20210420192559-75e0a1138db5
	github.com/observerward v0.0.0-00010101000000-000000000000
	github.com/prometheus/client_golang v1.10.0
	k8s.io/klog v1.0.0
)

replace github.com/observerward => ./

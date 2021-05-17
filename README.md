# ObserverWard

# Introduction

ObserverWard is a monitoring tool which could be utilized to implement real-time GPU resources monitoring in kubernetes clusters.

This tool uses Nvidia/nvml go-bindings to scrape GPU metrics such as bandwidth, model, memory utilization, etc. It also utilizes Prometheus go SDK to work together with Prometheus, offering strong querying languages and visualization.

# Usage

Run the command below:

```
kubectl apply -f deploy/deploy.yaml
```

Then visit the URL: http://<your-prometheus-machine-ip>:30090
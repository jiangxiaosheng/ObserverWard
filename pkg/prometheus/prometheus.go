package prometheus

import (
	"github.com/observerward/pkg/scraper"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog"
	"os"
	"strconv"
	"time"
)

var (
	namespace      = "observerward"
	dynamicMetrics = "dynamic"
	staticAttr     = "static"
	labels         = []string{"id", "uuid", "model"}
	ticker         *time.Ticker
)

var (
	gpuUsedGlobalMemory = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: dynamicMetrics,
			Name:      "gpu_used_global_memory_MiB",
			Help:      "GPU used global memory (in MiB).",
		}, labels)

	gpuFreeGlobalMemory = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: dynamicMetrics,
			Name:      "gpu_free_global_memory_MiB",
			Help:      "GPU free global memory (in MiB).",
		}, labels)

	gpuPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: dynamicMetrics,
			Name:      "gpu_power_usage_W",
			Help:      "GPU power draw (in W).",
		}, labels)

	gpuEncoderUtilization = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: dynamicMetrics,
			Name:      "gpu_encoder_utilization",
			Help:      "GPU encoder utilization (in %).",
		}, labels)

	gpuDecoderUtilization = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: dynamicMetrics,
			Name:      "gpu_decoder_utilization",
			Help:      "GPU decoder utilization (in %).",
		}, labels)

	gpuMemoryUtilization = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: dynamicMetrics,
			Name:      "gpu_memory_utilization",
			Help:      "GPU memory utilization (in %).",
		}, labels)

	gpuAttrMemorySize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: staticAttr,
			Name:      "gpu_memory_size",
			Help:      "Total size of memory of this GPU (in MB).",
		}, labels)

	gpuAttrMultiProcessorCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: staticAttr,
			Name:      "gpu_multiprocessor_count",
			Help:      "Count of multiprocessor on this GPU.",
		}, labels)

	gpuAttrSharedDecoderCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: staticAttr,
			Name:      "gpu_shared_decoder_count",
			Help:      "Count of shared decoder on this GPU.",
		}, labels)

	gpuAttrSharedEncoderCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: staticAttr,
			Name:      "gpu_shared_encoder_count",
			Help:      "Count of shared encoder on this GPU.",
		}, labels)
)

func Run(interval int) {
	ticker = time.NewTicker(time.Second * time.Duration(interval))
	sigs := make(chan os.Signal, 1)

	gpuMetrics, err := scraper.NewGPUMetrics()
	if err != nil {
		klog.Exitf("Prometheus Init Failed, %v", err)
	}
	sp := scraper.NewGPUMetricsScraper()

	go func() {
		{
			for {
				select {
				case <-ticker.C:
					err = sp.ScrapeGPUMetrics(gpuMetrics)
					if err != nil {
						klog.Error(err)
						return
					}
					for _, gpu := range gpuMetrics.GPUs {
						gpuUsedGlobalMemory.WithLabelValues(id2Str(gpu.StaticAttr.ID), gpu.StaticAttr.UUID, gpu.StaticAttr.Model).
							Set(float64(gpu.UsedGlobalMemory))
						gpuFreeGlobalMemory.WithLabelValues(id2Str(gpu.StaticAttr.ID), gpu.StaticAttr.UUID, gpu.StaticAttr.Model).
							Set(float64(gpu.FreeGlobalMemory))
						gpuPower.WithLabelValues(id2Str(gpu.StaticAttr.ID), gpu.StaticAttr.UUID, gpu.StaticAttr.Model).
							Set(float64(gpu.Power))
						gpuEncoderUtilization.WithLabelValues(id2Str(gpu.StaticAttr.ID), gpu.StaticAttr.UUID, gpu.StaticAttr.Model).
							Set(float64(gpu.EncoderUtilization))
						gpuDecoderUtilization.WithLabelValues(id2Str(gpu.StaticAttr.ID), gpu.StaticAttr.UUID, gpu.StaticAttr.Model).
							Set(float64(gpu.DecoderUtilization))
						gpuMemoryUtilization.WithLabelValues(id2Str(gpu.StaticAttr.ID), gpu.StaticAttr.UUID, gpu.StaticAttr.Model).
							Set(float64(gpu.MemoryUtilization))
						gpuAttrMemorySize.WithLabelValues(id2Str(gpu.StaticAttr.ID), gpu.StaticAttr.UUID, gpu.StaticAttr.Model).
							Set(float64(gpu.StaticAttr.MemorySizeMB))
						gpuAttrMultiProcessorCount.WithLabelValues(id2Str(gpu.StaticAttr.ID), gpu.StaticAttr.UUID, gpu.StaticAttr.Model).
							Set(float64(gpu.StaticAttr.MultiprocessorCount))
						gpuAttrSharedDecoderCount.WithLabelValues(id2Str(gpu.StaticAttr.ID), gpu.StaticAttr.UUID, gpu.StaticAttr.Model).
							Set(float64(gpu.StaticAttr.SharedDecoderCount))
						gpuAttrSharedEncoderCount.WithLabelValues(id2Str(gpu.StaticAttr.ID), gpu.StaticAttr.UUID, gpu.StaticAttr.Model).
							Set(float64(gpu.StaticAttr.SharedEncoderCount))
					}
				case <-sigs:
					return
				}
			}
		}
	}()
}

func init() {
	prometheus.MustRegister(gpuUsedGlobalMemory)
	prometheus.MustRegister(gpuFreeGlobalMemory)
	prometheus.MustRegister(gpuPower)
	prometheus.MustRegister(gpuEncoderUtilization)
	prometheus.MustRegister(gpuDecoderUtilization)
	prometheus.MustRegister(gpuMemoryUtilization)
	prometheus.MustRegister(gpuAttrMemorySize)
	prometheus.MustRegister(gpuAttrMultiProcessorCount)
	prometheus.MustRegister(gpuAttrSharedDecoderCount)
	prometheus.MustRegister(gpuAttrSharedEncoderCount)
}

func id2Str(id uint) string {
	return strconv.Itoa(int(id))
}

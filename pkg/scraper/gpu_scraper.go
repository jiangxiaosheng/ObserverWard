package scraper

import (
	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"k8s.io/klog"
)

type GPUMetrics struct {
	GPUs []*MetricsSnapshotPerGPU
}

func (g *GPUMetrics) Clone() *GPUMetrics {
	res := &GPUMetrics{}
	for i := range g.GPUs {
		res.GPUs = append(res.GPUs, g.GPUs[i].Clone())
	}
	return res
}

type MetricsSnapshotPerGPU struct {
	device     *nvml.Device
	StaticAttr GPUStaticAttr

	FreeGlobalMemory   uint64
	UsedGlobalMemory   uint64
	Power              uint
	EncoderUtilization uint
	DecoderUtilization uint
	MemoryUtilization  uint
	Temperature        uint
}

type GPUStaticAttr struct {
	ID                  uint
	UUID                string
	MemorySizeMB        uint64
	MultiprocessorCount uint32
	SharedDecoderCount  uint32
	SharedEncoderCount  uint32
	Model               string
	Bandwidth           uint
}

func (s *GPUStaticAttr) clone() GPUStaticAttr {
	return GPUStaticAttr{
		ID:                  s.ID,
		UUID:                s.UUID,
		MemorySizeMB:        s.MemorySizeMB,
		MultiprocessorCount: s.MultiprocessorCount,
		SharedDecoderCount:  s.SharedDecoderCount,
		SharedEncoderCount:  s.SharedEncoderCount,
	}
}

func (m *MetricsSnapshotPerGPU) Clone() *MetricsSnapshotPerGPU {
	return &MetricsSnapshotPerGPU{
		device:             m.device,
		StaticAttr:         m.StaticAttr.clone(),
		UsedGlobalMemory:   m.UsedGlobalMemory,
		Power:              m.Power,
		EncoderUtilization: m.EncoderUtilization,
		DecoderUtilization: m.DecoderUtilization,
		MemoryUtilization:  m.MemoryUtilization,
	}
}

type GPUMetricScraper struct {
}

func NewGPUMetrics() (*GPUMetrics, error) {
	gpuMetrics := &GPUMetrics{}

	err := nvml.Init()
	if err != nil {
		klog.Exitf("nvml Library Initialized Error: %v", err)
	}

	devicesCount, err := nvml.GetDeviceCount()
	if err != nil {
		klog.Errorf("Getting Device Count Error: %v", err)
		return nil, err
	}

	gpuMetrics.GPUs = make([]*MetricsSnapshotPerGPU, devicesCount)
	for i := uint(0); i < devicesCount; i++ {
		device, err := nvml.NewDevice(i)
		if err != nil {
			klog.Errorf("Getting Device %d Error: %v", i, err)
			return nil, err
		}
		gpuMetrics.GPUs[i] = &MetricsSnapshotPerGPU{
			device: device,
			StaticAttr: GPUStaticAttr{
				ID:   i,
				UUID: device.UUID,
			},
		}
	}
	return gpuMetrics, nil
}

func NewGPUMetricsScraper() *GPUMetricScraper {
	return &GPUMetricScraper{}
}

func (s *GPUMetricScraper) ScrapeGPUMetrics(gpuMetrics *GPUMetrics) error {
	for i, gpuMetricsSnapshot := range gpuMetrics.GPUs {
		device := gpuMetricsSnapshot.device
		st, err := device.Status()
		if err != nil {
			klog.Errorf("Getting Device %d Status Error: %v", i, err)
			return err
		}
		//attr, err := device.GetAttributes()
		//if err != nil {
		//	klog.Errorf("Getting Device %d Attributes Error: %v", i, err)
		//	return err
		//}
		gpuMetricsSnapshot.FreeGlobalMemory = *st.Memory.Global.Free
		gpuMetricsSnapshot.UsedGlobalMemory = *st.Memory.Global.Used
		gpuMetricsSnapshot.Power = *st.Power
		gpuMetricsSnapshot.EncoderUtilization = *st.Utilization.Encoder
		gpuMetricsSnapshot.DecoderUtilization = *st.Utilization.Decoder
		gpuMetricsSnapshot.MemoryUtilization = *st.Utilization.Memory
		gpuMetricsSnapshot.Temperature = *st.Temperature
		//gpuMetricsSnapshot.StaticAttr.MultiprocessorCount = attr.MultiprocessorCount
		//gpuMetricsSnapshot.StaticAttr.SharedDecoderCount = attr.SharedDecoderCount
		//gpuMetricsSnapshot.StaticAttr.SharedEncoderCount = attr.SharedEncoderCount
		//gpuMetricsSnapshot.StaticAttr.MemorySizeMB = attr.MemorySizeMB
		gpuMetricsSnapshot.StaticAttr.Model = *device.Model
	}
	return nil
}

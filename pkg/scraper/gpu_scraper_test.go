package scraper

import (
	"fmt"
	"testing"
)

func TestGPUScraper(t *testing.T) {
	sc := NewGPUMetricsScraper()
	metrics, err := NewGPUMetrics()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	err = sc.ScrapeGPUMetrics(metrics)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	for _, g := range metrics.GPUs {
		fmt.Printf("uuid: %v\n", g.StaticAttr.UUID)
		fmt.Printf("power: %d\n", g.Power)
		fmt.Printf("free memory: %v\n", g.FreeGlobalMemory)
		fmt.Printf("used memory: %v\n", g.UsedGlobalMemory)
	}
}

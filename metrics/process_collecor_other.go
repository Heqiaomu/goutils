package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

func canCollectProcess() bool {
	_, err := procfs.NewDefaultFS()
	return err == nil
}

func (c *processCollector) processCollect(ch chan<- prometheus.Metric) {
	pid, err := c.pidFn()
	if err != nil {
		c.reportError(ch, nil, err)
		return
	}

	p, err := procfs.NewProc(pid)
	if err != nil {
		c.reportError(ch, nil, err)
		return
	}

	if stat, err := p.Stat(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.cpuTotal, prometheus.CounterValue, stat.CPUTime(), c.LabelsValue...)
		ch <- prometheus.MustNewConstMetric(c.vsize, prometheus.GaugeValue, float64(stat.VirtualMemory()), c.LabelsValue...)
		ch <- prometheus.MustNewConstMetric(c.rss, prometheus.GaugeValue, float64(stat.ResidentMemory()), c.LabelsValue...)
		if startTime, err := stat.StartTime(); err == nil {
			ch <- prometheus.MustNewConstMetric(c.startTime, prometheus.GaugeValue, startTime, c.LabelsValue...)
		} else {
			c.reportError(ch, c.startTime, err)
		}
	} else {
		c.reportError(ch, nil, err)
	}

	if fds, err := p.FileDescriptorsLen(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.openFDs, prometheus.GaugeValue, float64(fds), c.LabelsValue...)
	} else {
		c.reportError(ch, c.openFDs, err)
	}

	if limits, err := p.Limits(); err == nil {
		ch <- prometheus.MustNewConstMetric(c.maxFDs, prometheus.GaugeValue, float64(limits.OpenFiles), c.LabelsValue...)
		ch <- prometheus.MustNewConstMetric(c.maxVsize, prometheus.GaugeValue, float64(limits.AddressSpace), c.LabelsValue...)
	} else {
		c.reportError(ch, nil, err)
	}
}

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

type memStatsMetrics []struct {
	desc    *prometheus.Desc
	eval    func(*runtime.MemStats) float64
	valType prometheus.ValueType
}

type serviceGoCollector struct {
	goroutinesDesc *prometheus.Desc
	threadsDesc    *prometheus.Desc
	gcDesc         *prometheus.Desc
	goInfoDesc     *prometheus.Desc

	// ms... are memstats related.
	msLast          *runtime.MemStats // Previously collected memstats.
	msLastTimestamp time.Time
	msMtx           sync.Mutex // Protects msLast and msLastTimestamp.
	msMetrics       memStatsMetrics
	msRead          func(*runtime.MemStats) // For mocking in tests.
	msMaxWait       time.Duration           // Wait time for fresh memstats.
	msMaxAge        time.Duration           // Maximum allowed age of old memstats.
	LabelsValue     []string
}

type ServiceGoCollectorOpts struct {
	// PidFn returns the PID of the process the collector collects metrics
	// for. It is called upon each collection. By default, the PID of the
	// current process is used, as determined on construction time by
	// calling os.Getpid().
	LabelsKey   []string
	LabelsValue []string
	// If non-empty, each of the collected metrics is prefixed by the
	// provided string and an underscore ("_").
	Namespace string
	// If true, any error encountered during collection is reported as an
	// invalid metric (see NewInvalidMetric). Otherwise, errors are ignored
	// and the collected metrics will be incomplete. (Possibly, no metrics
	// will be collected at all.) While that's usually not desired, it is
	// appropriate for the common "mix-in" of process metrics, where process
	// metrics are nice to have, but failing to collect them should not
	// disrupt the collection of the remaining metrics.
	ReportErrors bool
}

func NewServiceGoCollector(opts ServiceGoCollectorOpts) prometheus.Collector {
	return &serviceGoCollector{
		goroutinesDesc: prometheus.NewDesc(
			"go_goroutines",
			"Number of goroutines that currently exist.",
			opts.LabelsKey, nil),
		threadsDesc: prometheus.NewDesc(
			"go_threads",
			"Number of OS threads created.",
			opts.LabelsKey, nil),
		gcDesc: prometheus.NewDesc(
			"go_gc_duration_seconds",
			"A summary of the pause duration of garbage collection cycles.",
			opts.LabelsKey, nil),
		goInfoDesc: prometheus.NewDesc(
			"go_info",
			"Information about the Go environment.",
			opts.LabelsKey, prometheus.Labels{"version": runtime.Version()}),
		msLast:    &runtime.MemStats{},
		msRead:    runtime.ReadMemStats,
		msMaxWait: time.Second,
		msMaxAge:  5 * time.Minute,
		msMetrics: memStatsMetrics{
			{
				desc: prometheus.NewDesc(
					memstatNamespace("alloc_bytes"),
					"Number of bytes allocated and still in use.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.Alloc) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("alloc_bytes_total"),
					"Total number of bytes allocated, even if freed.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.TotalAlloc) },
				valType: prometheus.CounterValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("sys_bytes"),
					"Number of bytes obtained from system.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.Sys) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("lookups_total"),
					"Total number of pointer lookups.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.Lookups) },
				valType: prometheus.CounterValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("mallocs_total"),
					"Total number of mallocs.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.Mallocs) },
				valType: prometheus.CounterValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("frees_total"),
					"Total number of frees.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.Frees) },
				valType: prometheus.CounterValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("heap_alloc_bytes"),
					"Number of heap bytes allocated and still in use.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.HeapAlloc) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("heap_sys_bytes"),
					"Number of heap bytes obtained from system.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.HeapSys) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("heap_idle_bytes"),
					"Number of heap bytes waiting to be used.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.HeapIdle) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("heap_inuse_bytes"),
					"Number of heap bytes that are in use.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.HeapInuse) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("heap_released_bytes"),
					"Number of heap bytes released to OS.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.HeapReleased) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("heap_objects"),
					"Number of allocated objects.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.HeapObjects) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("stack_inuse_bytes"),
					"Number of bytes in use by the stack allocator.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.StackInuse) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("stack_sys_bytes"),
					"Number of bytes obtained from system for stack allocator.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.StackSys) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("mspan_inuse_bytes"),
					"Number of bytes in use by mspan structures.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.MSpanInuse) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("mspan_sys_bytes"),
					"Number of bytes used for mspan structures obtained from system.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.MSpanSys) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("mcache_inuse_bytes"),
					"Number of bytes in use by mcache structures.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.MCacheInuse) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("mcache_sys_bytes"),
					"Number of bytes used for mcache structures obtained from system.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.MCacheSys) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("buck_hash_sys_bytes"),
					"Number of bytes used by the profiling bucket hash table.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.BuckHashSys) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("gc_sys_bytes"),
					"Number of bytes used for garbage collection system metadata.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.GCSys) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("other_sys_bytes"),
					"Number of bytes used for other system allocations.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.OtherSys) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("next_gc_bytes"),
					"Number of heap bytes when next garbage collection will take place.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.NextGC) },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("last_gc_time_seconds"),
					"Number of seconds since 1970 of last garbage collection.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return float64(ms.LastGC) / 1e9 },
				valType: prometheus.GaugeValue,
			}, {
				desc: prometheus.NewDesc(
					memstatNamespace("gc_cpu_fraction"),
					"The fraction of this program's available CPU time used by the GC since the program started.",
					opts.LabelsKey, nil,
				),
				eval:    func(ms *runtime.MemStats) float64 { return ms.GCCPUFraction },
				valType: prometheus.GaugeValue,
			},
		},
		LabelsValue: opts.LabelsValue,
	}
}

func (c *serviceGoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.goroutinesDesc
	ch <- c.threadsDesc
	ch <- c.gcDesc
	ch <- c.goInfoDesc
	for _, i := range c.msMetrics {
		ch <- i.desc
	}
}

func (c *serviceGoCollector) Collect(ch chan<- prometheus.Metric) {
	var (
		ms   = &runtime.MemStats{}
		done = make(chan struct{})
	)
	// Start reading memstats first as it might take a while.
	go func() {
		c.msRead(ms)
		c.msMtx.Lock()
		c.msLast = ms
		c.msLastTimestamp = time.Now()
		c.msMtx.Unlock()
		close(done)
	}()

	ch <- prometheus.MustNewConstMetric(c.goroutinesDesc, prometheus.GaugeValue, float64(runtime.NumGoroutine()), c.LabelsValue...)
	n, _ := runtime.ThreadCreateProfile(nil)
	ch <- prometheus.MustNewConstMetric(c.threadsDesc, prometheus.GaugeValue, float64(n), c.LabelsValue...)

	var stats debug.GCStats
	stats.PauseQuantiles = make([]time.Duration, 5)
	debug.ReadGCStats(&stats)

	quantiles := make(map[float64]float64)
	for idx, pq := range stats.PauseQuantiles[1:] {
		quantiles[float64(idx+1)/float64(len(stats.PauseQuantiles)-1)] = pq.Seconds()
	}
	quantiles[0.0] = stats.PauseQuantiles[0].Seconds()
	ch <- prometheus.MustNewConstSummary(c.gcDesc, uint64(stats.NumGC), stats.PauseTotal.Seconds(), quantiles, c.LabelsValue...)

	ch <- prometheus.MustNewConstMetric(c.goInfoDesc, prometheus.GaugeValue, 1, c.LabelsValue...)

	timer := time.NewTimer(c.msMaxWait)
	select {
	case <-done: // Our own ReadMemStats succeeded in time. Use it.
		timer.Stop() // Important for high collection frequencies to not pile up timers.
		c.msCollect(ch, ms)
		return
	case <-timer.C: // Time out, use last memstats if possible. Continue below.
	}
	c.msMtx.Lock()
	if time.Since(c.msLastTimestamp) < c.msMaxAge {
		// Last memstats are recent enough. Collect from them under the lock.
		c.msCollect(ch, c.msLast)
		c.msMtx.Unlock()
		return
	}
	// If we are here, the last memstats are too old or don't exist. We have
	// to wait until our own ReadMemStats finally completes. For that to
	// happen, we have to release the lock.
	c.msMtx.Unlock()
	<-done
	c.msCollect(ch, ms)
}
func (c *serviceGoCollector) msCollect(ch chan<- prometheus.Metric, ms *runtime.MemStats) {
	for _, i := range c.msMetrics {
		ch <- prometheus.MustNewConstMetric(i.desc, i.valType, i.eval(ms), c.LabelsValue...)
	}
}
func memstatNamespace(s string) string {
	return "go_memstats_" + s
}

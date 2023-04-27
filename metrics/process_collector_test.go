package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestNewServiceProcessCollector(t *testing.T) {
	pd := make(chan<- *prometheus.Desc, 100)
	// pm := make(chan<- prometheus.Metric, 100)
	// mm := make(prometheus.Labels)
	// mm["key"] = "value"
	// desc := prometheus.NewDesc("fqName", "help", []string{"v1", "v2"}, mm)
	// metric, _ := prometheus.NewConstMetric(desc, prometheus.ValueType(1), 1)

	type args struct {
		opts ServiceProcessCollectorOpts
	}
	tests := []struct {
		name string
		args args
		want prometheus.Collector
	}{
		{
			"TestNewServiceProcessCollector",
			args{
				ServiceProcessCollectorOpts{
					LabelsKey:    []string{"key"},
					LabelsValue:  []string{"value"},
					Namespace:    "ns",
					ReportErrors: false,
				},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := NewServiceProcessCollector(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("NewServiceProcessCollector() = %v, want %v", got, tt.want)
			// }
			c := NewServiceProcessCollector(tt.args.opts)
			c.Describe(pd)
			close(pd)
		})
	}
}

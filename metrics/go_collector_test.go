package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestNewServiceGoCollector(t *testing.T) {

	pd := make(chan<- *prometheus.Desc, 100)
	pm := make(chan<- prometheus.Metric, 100)
	mm := make(prometheus.Labels)
	mm["key"] = "value"
	desc := prometheus.NewDesc("fqName", "help", []string{"v1", "v2"}, mm)
	metric, _ := prometheus.NewConstMetric(desc, prometheus.ValueType(1), 1)

	type args struct {
		opts ServiceGoCollectorOpts
	}
	tests := []struct {
		name string
		args args
		want prometheus.Collector
	}{
		{
			"TestNewServiceGoCollector",
			args{
				ServiceGoCollectorOpts{
					LabelsKey:   []string{"Key"},
					LabelsValue: []string{"Value"},
				},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := NewServiceGoCollector(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("NewServiceGoCollector() = %v, want %v", got, tt.want)
			// }
			gc := NewServiceGoCollector(tt.args.opts)
			gc.Describe(pd)
			gc.Collect(pm)
			pd <- desc
			pm <- metric
			close(pd)
			close(pm)
		})
	}
}

package server

import (
	"fmt"
	"testing"
	"time"
)

func TestWorker_Start(t *testing.T) {

	type args struct {
		wa   WorkerAction
		task interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "123",
			args: args{
				wa: func(task interface{}, exit chan bool) {

					for {
						select {
						case <-exit:
							fmt.Println(789)
							return
						default:
							time.Sleep(3 * time.Second)
							fmt.Println(task)
						}
					}
				},
				task: "abc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFactory()
			w := f.NewWorker(nil)
			w1 := f.GetWorker(w.ID)
			w1.Start(tt.args.wa, tt.args.task)
			fmt.Println(w1.Doing)
			f.FireWorker(w1.ID)
			f.Destroy()
			fmt.Println(w1.Doing)
		})
	}
}

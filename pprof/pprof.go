package pprof

import (
	"net/http"
	"net/http/pprof"
)

type PProf interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

func StartPProf(p PProf) {
	p.HandleFunc("/debug/pprof/", pprof.Index)
	p.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	p.HandleFunc("/debug/pprof/profile", pprof.Profile)
	p.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	p.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

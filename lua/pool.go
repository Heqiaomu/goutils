package lua

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"sync"
)

// lStatePool is the pool of global Lua State
// call Shutdown when cloud close
type lStatePool struct {
	m     sync.Mutex
	saved []*lua.LState
	plm   []*PreLoadModule
}

type PreLoadModule struct {
	Name string
	Func lua.LGFunction
}

// LuaPool Global LState pool
var LuaPool = &lStatePool{
	saved: make([]*lua.LState, 0, 4),
}

func PreLoadModules(modules ...PreLoadModule) {
	var plm = make([]*PreLoadModule, 0)
	for _, module := range modules {
		nModule := module
		if module.Name != "" && module.Func != nil {
			plm = append(plm, &nModule)
			continue
		}
		fmt.Printf("module(%v) is not correct type\n", module)
	}
	LuaPool = &lStatePool{
		plm: plm,
	}
}

func (pl *lStatePool) Get() *lua.LState {
	pl.m.Lock()
	defer pl.m.Unlock()
	n := len(pl.saved)
	if n == 0 {
		return pl.new()
	}
	x := pl.saved[n-1]
	pl.saved = pl.saved[0 : n-1]
	return x
}

func (pl *lStatePool) new() *lua.LState {
	L := lua.NewState()
	for _, module := range pl.plm {
		L.PreloadModule(module.Name, module.Func)
	}
	return L
}

func (pl *lStatePool) Put(L *lua.LState) {
	pl.m.Lock()
	defer pl.m.Unlock()
	length := L.GetTop()
	if length != 0 {
		// clear luastate
		L.Pop(length)
	}
	pl.saved = append(pl.saved, L)
}

func (pl *lStatePool) Shutdown() {
	for _, L := range pl.saved {
		L.Close()
	}
}

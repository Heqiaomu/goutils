package lua

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func Test_lStatePool(t *testing.T) {
	// type fields struct {
	// 	m     sync.Mutex
	// 	saved []*lua.LState
	// 	plm   []*PreLoadModule
	// }
	// tests := []struct {
	// 	name   string
	// 	fields fields
	// 	want   *lua.LState
	// }{
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		pl := &lStatePool{
	// 			m:     tt.fields.m,
	// 			saved: tt.fields.saved,
	// 			plm:   tt.fields.plm,
	// 		}
	// 		if got := pl.Get(); !reflect.DeepEqual(got, tt.want) {
	// 			t.Errorf("lStatePool.Get() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }

	LuaPool.Put(lua.NewState())
	LuaPool.Get()
	LuaPool.Shutdown()
}

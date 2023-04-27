package lua

import (
	"fmt"
	"github.com/fatih/structs"
	lua "github.com/yuin/gopher-lua"
	"reflect"
	"strconv"
	"time"
)

// ErrorReturn trans err to return, the return must be 2 param,
// the first is data, and the second is error
func ErrorReturn(l *lua.LState, err error) int {
	l.Push(lua.LNil)
	l.Push(lua.LString(err.Error()))
	return 2
}

// DataReturn trans data to table and return 2 param
// the first is data, and the second is error
func DataReturn(l *lua.LState, data interface{}) int {
	// 判断，interface转为[]interface{}
	v := reflect.ValueOf(data)
	var resTable *lua.LTable
	if v.Kind() != reflect.Slice {
		resTable = SimpleStruct2LTable(data)
	} else {
		resTable = SimpleArrayStruct2LTable(data)
	}

	l.Push(resTable)
	l.Push(lua.LNil)
	return 2
}

func SimpleArrayStruct2LTable(ss interface{}) *lua.LTable {
	// 判断，interface转为[]interface{}
	v := reflect.ValueOf(ss)
	if v.Kind() != reflect.Slice {
		fmt.Printf("ss not slice type")
	}

	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}

	res := lua.LTable{}
	for index, is := range ret {
		m := structs.Map(is)
		table := Map2Table(m)
		res.RawSetInt(index+1, table)
	}
	return &res
}

// SimpleStruct2LTable simple trans struct to table
func SimpleStruct2LTable(s interface{}) *lua.LTable {
	res := lua.LTable{}
	m := structs.Map(s)
	res = *Map2Table(m)
	return &res
}

// Map2Table converts a Go map to a lua table
func Map2Table(m map[string]interface{}) *lua.LTable {
	// Main table pointer
	resultTable := &lua.LTable{}

	// Loop map
	for key, element := range m {

		switch element.(type) {
		case float64, float32:
			resultTable.RawSetString(key, lua.LNumber(element.(float64)))
		case int64, uint64, int32, uint32, int16, uint16, int8, uint8, int, uint:
			bits := reflect.TypeOf(element).Bits()
			i64, err := strconv.ParseInt(fmt.Sprintf("%d", element), 10, bits)
			if err != nil {
				panic(err)
			}
			resultTable.RawSetString(key, lua.LNumber(i64))
		case string:
			resultTable.RawSetString(key, lua.LString(element.(string)))
		case bool:
			resultTable.RawSetString(key, lua.LBool(element.(bool)))
		case []byte:
			resultTable.RawSetString(key, lua.LString(string(element.([]byte))))
		case map[string]interface{}:

			// Get table from map
			tble := Map2Table(element.(map[string]interface{}))

			resultTable.RawSetString(key, tble)
		case map[string]string:
			var tmp = make(map[string]interface{}, 0)
			for key, value := range element.(map[string]string) {
				tmp[key] = value
			}
			tble := Map2Table(tmp)
			resultTable.RawSetString(key, tble)

		case time.Time:
			resultTable.RawSetString(key, lua.LNumber(element.(time.Time).Unix()))

		case []map[string]interface{}:

			// Create slice table
			sliceTable := &lua.LTable{}

			// Loop element
			for _, s := range element.([]map[string]interface{}) {

				// Get table from map
				tble := Map2Table(s)

				sliceTable.Append(tble)
			}

			// Set slice table
			resultTable.RawSetString(key, sliceTable)

		case []interface{}:
			// Create slice table
			sliceTable := &lua.LTable{}

			// Loop interface slice
			for _, s := range element.([]interface{}) {

				// Switch interface type
				switch s.(type) {
				case map[string]interface{}:

					// Convert map to table
					t := Map2Table(s.(map[string]interface{}))

					// Append result
					sliceTable.Append(t)

				case int64, uint64, int32, uint32, int16, uint16, int8, uint8, int, uint:
					bits := reflect.TypeOf(s).Bits()
					i64, err := strconv.ParseInt(fmt.Sprintf("%d", s), 10, bits)
					if err != nil {
						panic(err)
					}
					sliceTable.Append(lua.LNumber(i64))

				case float64:

					// Append result as number
					sliceTable.Append(lua.LNumber(s.(float64)))

				case string:

					// Append result as string
					sliceTable.Append(lua.LString(s.(string)))

				case bool:

					// Append result as bool
					sliceTable.Append(lua.LBool(s.(bool)))
				}
			}

			// Append to main table
			resultTable.RawSetString(key, sliceTable)
		}
	}

	return resultTable
}

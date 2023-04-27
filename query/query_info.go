package query

import (
	"fmt"
	"github.com/Heqiaomu/goutil/auth"
	"github.com/pkg/errors"
	"reflect"
	"regexp"
	"strings"
)

type QueryInfo struct {
	start   int
	size    int
	order   []Order
	filter  []Filter
	mapping map[string]FilterMappingFunc
}

const (
	comma = ","
	space = " "
)

const (
	Desc = "desc"
	Asc  = "asc"
)

type Order struct {
	Colum string
	Sort  string
}

type Filter struct {
	Colum  string
	Action Action
	Value  interface{}
}

var filterRegex = regexp.MustCompile(`(.+[^\<])<(.*[^\<])>(.+)`)

type FilterMappingFunc func(v string) (mappingK, mappingV string)

func New() *QueryInfo {
	return &QueryInfo{
		mapping: make(map[string]FilterMappingFunc),
	}
}

// AddSortCol add sort mode on colum
func (q *QueryInfo) AddSortCol(colum, mode string) *QueryInfo {
	q.order = append(q.order, Order{
		Colum: colum,
		Sort:  mode,
	})
	return q
}

var NoValue = "<invalid Value>"

// ParseParams Please 只会处理filter的内容，path上的参数不会自动转化
func (q *QueryInfo) ParseParams(params interface{}) (*QueryInfo, error) {
	el := reflect.ValueOf(params).Elem()

	startField := el.FieldByName("Start")
	sizeField := el.FieldByName("Size")
	if !startField.IsValid() || startField.Int() == 0 {
		q.start = 1
	} else {
		q.start = int(startField.Int())
	}
	if !sizeField.IsValid() || sizeField.Int() == 0 {
		q.size = 10
	} else {
		q.size = int(sizeField.Int())
	}

	sort := el.FieldByName("Sort")
	if sort.IsValid() && len(sort.String()) != 0 {
		err := q.parseRawSort(sort.String())
		if err != nil {
			return nil, errors.Wrap(err, "分页解析错误")
		}
	}
	// 加上默认排序
	if len(q.order) == 0 {
		q.order = append(q.order, Order{
			Colum: "time_create",
			Sort:  "desc",
		})
	}
	filter := el.FieldByName("Filter")
	if filter.IsValid() && len(filter.String()) != 0 {
		err := q.parseRawFilter(filter.String())
		if err != nil {
			return nil, errors.Wrap(err, "分页解析错误")
		}
	}

	return q, nil
}

func (q *QueryInfo) Mapping(col string, mapping FilterMappingFunc) *QueryInfo {
	q.mapping[col] = mapping
	return q
}

func (q *QueryInfo) MappingColum(col string, mappingCol string) *QueryInfo {
	q.mapping[col] = func(v string) (mappingK, mappingV string) {
		return mappingCol, v
	}
	return q
}

func (q *QueryInfo) FilterColumns(cols ...string) *QueryInfo {
	for _, col := range cols {
		temp := col
		q.mapping[col] = func(v string) (mappingK, mappingV string) {
			return temp, v
		}
	}
	return q
}

func (q *QueryInfo) OriginalColumns(cols ...string) *QueryInfo {
	for _, col := range cols {
		temp := col
		q.mapping[col] = func(v string) (mappingK, mappingV string) {
			return temp, v
		}
	}
	return q
}

func (q *QueryInfo) parseRawSort(orderStr string) error {
	orders := strings.Split(orderStr, comma)
	var order []Order
	for _, sort := range orders {
		split := strings.Split(sort, space)
		if len(split) < 2 {
			return errors.New("排序规则错误")
		}
		mapFunc, ok := q.mapping[split[0]]
		column := split[0]
		if ok {
			_, v := mapFunc(split[0])
			column = v
		}
		order = append(order, Order{
			// 替换一次mapping 的结果
			Colum: snakeString(column),
			Sort:  split[1],
		})
	}
	q.order = append(q.order, order...)
	return nil
}

func (q *QueryInfo) parseRawFilter(filterStr string) error {
	var filter []Filter
	filterParams := strings.Split(filterStr, comma)
	for _, filterParam := range filterParams {
		match := filterRegex.FindStringSubmatch(filterParam)
		if len(match) < 4 {
			return errors.New("filter格式错误")
		}
		action, ok := ActionMap[match[2]]
		if !ok {
			return errors.Errorf("[%s]操作不存在", match[2])
		}
		col := match[1]
		mapFunc, ok := q.mapping[col]
		if !ok {
			continue
		}
		parseValue := action.parseValue(match[3])
		col, parseValue = mapFunc(parseValue)
		filter = append(filter, Filter{
			Colum:  col,
			Action: action,
			Value:  parseValue,
		})
	}
	q.filter = append(q.filter, filter...)
	return nil
}

func (q *QueryInfo) EqualFilter(colum string, value interface{}) *QueryInfo {
	q.filter = append(q.filter, Filter{
		Colum:  colum,
		Value:  value,
		Action: EQ,
	})
	return q
}

func (q *QueryInfo) Order(colum, sort string) *QueryInfo {
	q.order = append(q.order, Order{
		Colum: colum,
		Sort:  sort,
	})
	return q
}

// InFilter todo: 需要重新修改入参
func (q *QueryInfo) InFilter(colum string, value ...interface{}) *QueryInfo {
	if value == nil || len(value) == 0 {
		return q
	}
	if len(value) == 1 {
		i, ok := value[0].([]interface{})
		if ok {
			// 将内部的interface{}取出
			value = i
		}
	}
	// 数组转成string 不需要 括号
	var temp = make([]string, len(value))
	for i, v := range value {
		temp[i] = fmt.Sprintf("%v", v)
	}
	in := strings.Join(temp, "@")
	q.filter = append(q.filter, Filter{
		Colum:  colum,
		Value:  in,
		Action: IN,
	})
	return q
}

// Remove 移除指定 filter 只要column、action 都相同就会移除
func (q *QueryInfo) Remove(filter *Filter) *QueryInfo {
	var f []Filter
	for _, v := range q.filter {
		if v.Colum != filter.Colum && v.Action != filter.Action {
			f = append(f, v)
		}
	}
	q.filter = f
	return q
}

// RemoveSpecified 移除指定 filter column、action、value 均相同才删除
func (q *QueryInfo) RemoveSpecified(filter *Filter) *QueryInfo {
	var f []Filter
	var deleteIndex int
	for i, v := range q.filter {
		if v.Colum == filter.Colum && v.Action == filter.Action && v.Value == filter.Value {
			deleteIndex = i
		}
	}
	for i, v := range q.filter {
		if i != deleteIndex {
			f = append(f, v)
		}
	}

	q.filter = f
	return q
}

func (q *QueryInfo) AddFilter(colum string, value interface{}, action Action) *QueryInfo {
	q.filter = append(q.filter, Filter{
		Colum:  colum,
		Value:  value,
		Action: action,
	})
	return q
}

// FindFiler 只要column、action 都相同 就认为找到
func (q *QueryInfo) FindFiler(filter *Filter) []*Filter {
	var f []*Filter
	for _, v := range q.filter {
		if v.Colum == filter.Colum && v.Action == filter.Action {
			f = append(f, &v)
		}
	}
	return f
}

func (q *QueryInfo) NoEqualFilter(colum string, value interface{}) *QueryInfo {
	q.filter = append(q.filter, Filter{
		Colum:  colum,
		Value:  value,
		Action: NE,
	})
	return q
}

func (q *QueryInfo) Start(start int) *QueryInfo {
	q.start = start
	return q
}

func (q *QueryInfo) Size(size int) *QueryInfo {
	q.size = size
	return q
}

func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		// 判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	// ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}

func (q *QueryInfo) GetStart() int {
	return q.start
}

func (q *QueryInfo) GetSize() int {
	return q.size
}

func (q *QueryInfo) GetOrder() []Order {
	return q.order
}

func (q *QueryInfo) GetFilter() []Filter {
	return q.filter
}

func (q *QueryInfo) Auth(auth *auth.Auth) *QueryInfo {
	q.EqualFilter("viewer", auth.Viewer)
	if auth.Admin {
		return q
	}

	if auth.Viewer == "USER" {
		return q.EqualFilter("user_id", auth.UserID)
	} else {
		return q.EqualFilter("group_id", auth.GroupID)
	}

}

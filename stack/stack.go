package stack

type Stack interface {
	Pop() interface{}
	Push(ele interface{})
	Top() interface{}
	IsEmpty() bool
	Size() int
}

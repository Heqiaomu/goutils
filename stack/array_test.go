package stack

import "testing"

func TestArrayStack(t *testing.T) {
	stack := NewArrayStack()
	t.Log(stack.IsEmpty())
	t.Log(stack.Size())

	stack.Push(888)
	t.Log(stack.IsEmpty())
	t.Log(stack.Size())

	t.Log(stack.Top())
	t.Log(stack.IsEmpty())
	t.Log(stack.Size())

	t.Log(stack.Pop())
	t.Log(stack.IsEmpty())
	t.Log(stack.Size())

	stack.Push(987)
	stack.Push(841)
	t.Log(stack.IsEmpty())
	t.Log(stack.Size())

	t.Log(stack.Top())
	t.Log(stack.IsEmpty())
	t.Log(stack.Size())

	t.Log(stack.Pop())
	t.Log(stack.IsEmpty())
	t.Log(stack.Size())
}

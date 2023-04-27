package stack

type ArrayStack struct {
	top      int
	elements []interface{}
}

func NewArrayStack() *ArrayStack {
	return &ArrayStack{
		top:      -1,
		elements: make([]interface{}, 0),
	}
}

func (s *ArrayStack) Pop() interface{} {
	if s == nil {
		return nil
	}
	if s.IsEmpty() {
		return nil
	}
	ele := s.elements[s.top]
	s.elements = s.elements[0:s.top]
	s.top = s.top - 1
	return ele
}

func (s *ArrayStack) Push(ele interface{}) {
	if s == nil {
		return
	}
	s.elements = append(s.elements, ele)
	s.top = s.top + 1
}

func (s *ArrayStack) Top() interface{} {
	if s == nil {
		return nil
	}
	if s.IsEmpty() {
		return nil
	}
	return s.elements[s.top]
}

func (s *ArrayStack) IsEmpty() bool {
	if s == nil {
		return true
	}
	if len(s.elements) == 0 {
		return true
	}
	if s.top <= -1 {
		return true
	}
	return false
}

func (s *ArrayStack) Size() int {
	if s == nil {
		return 0
	}
	return len(s.elements)
}

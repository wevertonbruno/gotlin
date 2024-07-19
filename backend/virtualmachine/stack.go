package virtualmachine

import (
	"gotlin/backend/virtualmachine/chunk"
)

type stack struct {
	top    int
	values [StackMax]chunk.Value
}

func (s *stack) push(value chunk.Value) {
	s.values[s.top] = value
	s.top++
}

func (s *stack) pop() chunk.Value {
	s.top--
	return s.values[s.top]
}

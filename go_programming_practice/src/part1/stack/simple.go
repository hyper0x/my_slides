package stack

import (
	"errors"
)

type SimpleStack struct {
	capacity  uint
	cursor    uint
	container []interface{}
}

func (self *SimpleStack) build(rebuilding bool) {
	if self.capacity == 0 {
		err := errors.New("The capacity is ZERO!")
		panic(err)
	}
	if rebuilding || self.container == nil {
		self.container = make([]interface{}, self.capacity)
		self.cursor = 0
	}
}

func (self *SimpleStack) Clear() {
	self.build(true)
}

func (self *SimpleStack) Len() uint {
	self.build(false)
	return self.cursor
}

func (self *SimpleStack) Cap() uint {
	self.build(false)
	return uint(cap(self.container))
}

func (self *SimpleStack) Peek() interface{} {
	self.build(false)
	result := self.container[self.cursor-1]
	return result
}

func (self *SimpleStack) Pop() (interface{}, bool) {
	self.build(false)
	if self.cursor == 0 {
		return nil, false
	}
	result := self.container[self.cursor-1]
	self.container[self.cursor-1] = nil
	self.cursor -= 1
	return result, true
}

func (self *SimpleStack) Push(value interface{}) bool {
	self.build(false)
	if self.cursor == uint(cap(self.container)) {
		return false
	}
	self.container[self.cursor] = value
	self.cursor += 1
	return true
}

func NewSimpleStack(myCapacity uint) Stack {
	return &SimpleStack{capacity: myCapacity}
}

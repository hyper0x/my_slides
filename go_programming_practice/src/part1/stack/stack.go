package stack

type Stack interface {
	Clear()
	Len() uint
	Cap() uint
	Peek() interface{}
	Pop() (interface{}, bool)
	Push(value interface{}) bool
}

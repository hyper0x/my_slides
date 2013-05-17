package stack

import (
	"fmt"
	"runtime/debug"
	"testing"
)

func testWithPanic(f func(), t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			errorMsg := fmt.Sprintf("%s", err)
			if errorMsg == "The capacity is ZERO!" {
				t.Log("Ignore the cap error.\n")
			} else {
				t.Errorf("Fatal Error: %s\n", err)
				debug.PrintStack()
				t.FailNow()
			}
		}
	}()
	f()
}

func TestInterface(t *testing.T) {
	ss := &SimpleStack{}
	if _, ok := interface{}(ss).(Stack); !ok {
		t.Error("The instence of SimpleStack is NOT a Stack!\n")
	}
}

func TestClear(t *testing.T) {
	testWithPanic(func() {
		ss := &SimpleStack{}
		ss.Clear()
	}, t)
	ss := &SimpleStack{capacity: 5}
	ss.Clear()
}

func TestLen(t *testing.T) {
	testWithPanic(func() {
		ss := &SimpleStack{}
		ss.Len()
	}, t)
	ss := &SimpleStack{capacity: 5}
	expectedLen := uint(0)
	if ss.Len() != expectedLen {
		t.Errorf("The length is %d, but should be %d.\n", ss.Len(), expectedLen)
	}
}

func TestCap(t *testing.T) {
	testWithPanic(func() {
		ss := &SimpleStack{}
		ss.Len()
	}, t)
	expectedCap := uint(5)
	ss := &SimpleStack{capacity: expectedCap}
	if ss.Cap() != expectedCap {
		t.Errorf("The capacity is %d, but should be %d.\n", ss.Cap(), expectedCap)
	}
}

func TestOps(t *testing.T) {
	testWithPanic(func() {
		ss := &SimpleStack{}
		ss.Len()
	}, t)
	ss := &SimpleStack{capacity: 1}
	element0 := 1
	var result bool
	result = ss.Push(element0)
	if !result {
		t.Errorf("The result of 'Push(%d)' is %v, but should be %v.\n", element0, result, true)
	}
	element1 := 2
	result = ss.Push(element1)
	if result {
		t.Errorf("The result of 'Push(%d)' is %v, but should be %v. (because the stack is full)\n", element1, result, false)
	}
	expectedElement := element0
	objectResult := ss.Peek()
	if objectResult == nil {
		t.Error("The result of 'Peek()' is NIL.\n")
	}
	if v, ok := interface{}(objectResult).(int); ok {
		if v != expectedElement {
			t.Errorf("The result of 'Peek()' is %v, but should be %v.\n", v, expectedElement)
		}
	}
	objectResult, ok := ss.Pop()
	if !ok {
		t.Error("The result of 'Pop()' should be ok.\n")
	}
	if objectResult == nil {
		t.Error("The result of 'Pop()' is NIL.\n")
	}
	if v, ok := interface{}(objectResult).(int); ok {
		if v != expectedElement {
			t.Errorf("The result 'Pop()' is %v, but should be %v.\n", v, expectedElement)
		}
	}
	objectResult, ok = ss.Pop()
	if ok {
		t.Error("The result of 'Pop()' should be not ok.\n")
	}
	if objectResult != nil {
		t.Error("The result of 'Pop()' should be NIL.\n")
	}
}

package main

import (
	"errors"
	"math/rand"
	"time"
)

var (
	ErrNotSupported = errors.New("operation not supported")
	ErrBadRoll      = errors.New("roll can only support values greater than 0")
)

// Op is an arithmetic operator
type Op interface {
	Default() (int, error)
	Value(int) (int, error)
	Eval(int, int) (int, error)
}

// AddOp is an arithmetic add
type AddOp struct{}

func (op AddOp) Default() (int, error)      { return 0, ErrNotSupported }
func (op AddOp) Value(a int) (int, error)   { return 0, ErrNotSupported }
func (op AddOp) Eval(a, b int) (int, error) { return a + b, nil }

// SubOp is an aritmetic subtract
type SubOp struct{}

func (op SubOp) Default() (int, error)      { return 0, ErrNotSupported }
func (op SubOp) Value(a int) (int, error)   { return 0, ErrNotSupported }
func (op SubOp) Eval(a, b int) (int, error) { return a - b, nil }

// MultOp is an arithmetic multiply
type MultOp struct{}

func (op MultOp) Default() (int, error)      { return 0, ErrNotSupported }
func (op MultOp) Value(a int) (int, error)   { return 0, ErrNotSupported }
func (op MultOp) Eval(a, b int) (int, error) { return a * b, nil }

// RollOp rolls an n-sided die x times
type RollOp struct{}

func (op RollOp) Default() (int, error)    { return 1, nil }
func (op RollOp) Value(a int) (int, error) { return 0, ErrNotSupported }
func (op RollOp) Eval(a, b int) (int, error) {
	if a <= 0 || b <= 0 {
		return 0, ErrBadRoll
	}
	total := 0
	for i := 0; i < a; i++ {
		rand.Seed(time.Now().UnixNano())
		total += (rand.Int() % a) + 1
	}
	return total, nil
}

// NoOp is an arithmetic equality
type NoOp struct{}

func (op NoOp) Default() (int, error)      { return 0, ErrNotSupported }
func (op NoOp) Value(a int) (int, error)   { return a, nil }
func (op NoOp) Eval(a, b int) (int, error) { return 0, ErrNotSupported }

// Operator translates aritmetic characters into their respective operations
type Operator byte

// IsOperator returns true if a given byte can be translated into an op
func IsOperator(b byte) bool {
	return Operator(b).get() != nil
}

func (o Operator) get() Op {
	switch o {
	case '+':
		return AddOp{}
	case '-':
		return SubOp{}
	case '*':
		return MultOp{}
	case 'd', 'D':
		return RollOp{}
	case 0:
		return NoOp{}
	default:
		return nil
	}
}

// Order returns the hierarchical order of a given operation
func (o Operator) Order() int {
	switch o.get().(type) {
	case NoOp:
		return 0
	case AddOp, SubOp:
		return 1
	case MultOp:
		return 2
	case RollOp:
		return 3
	default:
		return -1
	}
}

func (o Operator) Default() (int, error) {
	if op := o.get(); op != nil {
		return op.Default()
	}
	return 0, ErrNotSupported
}

func (o Operator) Value(a int) (int, error) {
	if op := o.get(); op != nil {
		return op.Value(a)
	}
	return 0, ErrNotSupported
}

func (o Operator) Eval(a, b int) (int, error) {
	if op := o.get(); op != nil {
		return op.Eval(a, b)
	}
	return 0, ErrNotSupported
}

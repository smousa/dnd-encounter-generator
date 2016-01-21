package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"unicode"
)

type SyntaxError string

func (s SyntaxError) Error() string {
	return fmt.Sprintf("syntax error: %s", s)
}

func IsDigit(b byte) bool      { return unicode.IsDigit(rune(b)) }
func IsSpace(b byte) bool      { return unicode.IsSpace(rune(b)) }
func IsOpenParen(b byte) bool  { return b == '(' }
func IsCloseParen(b byte) bool { return b == ')' }

const (
	Add  Operator = '+'
	Sub  Operator = '-'
	Mult Operator = '*'
	Roll Operator = 'd'
)

type Operator byte

func IsOperator(b byte) bool {
	return Operator(b).Order() > 0
}

func (op Operator) Value(a, b int) (int, error) {
	switch op {
	case Add:
		return a + b, nil
	case Sub:
		return a - b, nil
	case Mult:
		return a * b, nil
	case Roll:
		if a > 0 && b > 0 {
			sum := 0
			for i := 0; i < a; i++ {
				rand.Seed(time.Now().UnixNano())
				sum += (rand.Int() % b) + 1
			}
			return sum, nil
		}
	}
	return 0, errors.New("cannot compute")
}

func (op Operator) Order() int {
	switch op {
	case Add, Sub:
		return 1
	case Mult:
		return 2
	case Roll:
		return 3
	}
	return 0
}

type Term struct {
	Value    int
	Operator Operator
}

func Reduce(terms ...Term) ([]Term, error) {
	switch len(terms) {
	case 0, 1:
		return terms, nil
	case 2:
		t, err := terms[0].Compute(terms[1])
		if err != nil {
			return nil, err
		}
		return []Term{t}, nil
	}
	for len(terms) > 2 {
		if terms[0].Order() >= terms[1].Order() {
			t, err := terms[0].Compute(terms[1])
			if err != nil {
				return nil, err
			}
			terms = append([]Term{t}, terms[2:]...)
		} else {
			t, err := terms[1].Compute(terms[2])
			if err != nil {
				return nil, err
			}
			terms = append([]Term{t}, terms[3:]...)
		}
	}
	return terms, nil
}

func (a Term) Compute(b Term) (t Term, err error) {
	v, err := a.Operator.Value(a.Value, b.Value)
	if err != nil {
		return
	}
	return Term{Value: v, Operator: b.Operator}, nil
}

func (a Term) Order() int {
	return a.Operator.Order()
}

type Phrase struct {
	Equation []byte
}

func (p Phrase) Compute(index int) (value int, idx int, err error) {
	var isOpenParen bool
	var terms []Term
	for idx = index; idx < len(p.Equation); idx++ {
		c := p.Equation[idx]
		if IsOpenParen(c) {
			if isOpenParen {
				value, idx, err = p.Compute(idx)
				if err != nil {
					return
				}
				if num := len(terms); num > 0 && terms[num-1].Operator.Order() == 0 {
					return 0, 0, SyntaxError(fmt.Sprintf("at char %d: missing operator", idx))
				}
				terms = append(terms, Term{Value: value})
			}
			isOpenParen = true
		} else if IsCloseParen(c) {
			if !isOpenParen {
				return 0, 0, SyntaxError(fmt.Sprintf("at char %d: missing '('", idx))
			}
			isOpenParen = false
			break
		} else if IsOperator(c) {
			num := len(terms)
			if num == 0 || terms[num-1].Operator.Order() > 0 {
				return 0, 0, SyntaxError(fmt.Sprintf("at char %d: missing operand", idx))
			}
			terms[num-1].Operator = Operator(c)
		} else if IsDigit(c) {
			if num := len(terms); num > 0 && terms[num-1].Operator.Order() == 0 {
				return 0, 0, SyntaxError(fmt.Sprintf("at char %d: missing operator", idx))
			}
			value, idx, err = p.getNumber(idx)
			if err != nil {
				return
			}
			terms = append(terms, Term{Value: value})
		} else if IsSpace(c) {
			continue
		} else {
			return 0, 0, SyntaxError(fmt.Sprintf("at char %d: unknown operator '%s'", idx, c))
		}
		if len(terms) == 3 {
			if terms, err = Reduce(terms...); err != nil {
				return 0, 0, err
			}
		}
	}
	if num := len(terms); num == 0 {
		return 0, 0, SyntaxError(fmt.Sprintf("at char %d: no data to evaluate", idx))
	} else if terms[num-1].Operator.Order() > 0 {
		return 0, 0, SyntaxError(fmt.Sprintf("at char %d: missing operand", idx))
	}
	for len(terms) > 1 {
		if terms, err = Reduce(terms...); err != nil {
			return 0, 0, err
		}
	}
	return terms[0].Value, idx, nil
}

func (p Phrase) getNumber(index int) (value int, idx int, err error) {
	var operand []byte
	for idx = index; idx < len(p.Equation); idx++ {
		c := p.Equation[idx]
		if IsDigit(c) {
			operand = append(operand, c)
		} else {
			break
		}
	}
	v, err := strconv.ParseInt(string(operand), 10, 32)
	if err != nil {
		return 0, 0, err
	}
	return int(v), idx, nil
}

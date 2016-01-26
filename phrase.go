package main

import (
	"fmt"
	"strconv"
	"unicode"
)

// SyntaxError is an error type
type SyntaxError string

func (s SyntaxError) Error() string {
	return fmt.Sprintf("syntax error: %s", s)
}

// IsDigit returns true if a byte is a digit
func IsDigit(b byte) bool { return unicode.IsDigit(rune(b)) }

// IsSpace returns true if a byte is whitespace
func IsSpace(b byte) bool { return unicode.IsSpace(rune(b)) }

// IsOpenParen returns true if a byte is an open parenthesis
func IsOpenParen(b byte) bool { return b == '(' }

// IsCloseParen returns true if a byte is a closed parenthesis
func IsCloseParen(b byte) bool { return b == ')' }

// Phrase describes a piece of the equation
type Phrase struct {
	Equation []byte
}

// Calculate computes the value of the equation
func Evaluate(equation []byte) (int, error) {
	v, err := Phrase{equation}.Eval()
	return v, err
}

// Eval computes the value of an equation
func (p Phrase) Eval() (int, error) {
	v, _, err := p.eval(-1)
	return v, err
}

// eval uses recursion to compute the value of the equation from certain
// indices. A -1 index indicates that this is the root calculation
func (p Phrase) eval(index int) (value int, idx int, err error) {
	var isOpen bool     // are we evaluating the equation from within a context?
	var isDefault bool  // did we default on the previous term?
	var isNegative bool // is it a negative number?

	var terms []Term // terms of the equation
	// set the idx start value
	if index < 0 {
		idx = 0
	} else {
		idx = index
	}
	for idx < len(p.Equation) {
		numTerms := len(terms)
		c := p.Equation[idx]
		if IsOpenParen(c) {
			if index < 0 || numTerms > 0 || isOpen {
				// evaluate within a new context
				value, idx, err = p.eval(idx)
				if err != nil {
					return
				}
				if numTerms > 0 && Order(terms[numTerms-1]) <= 0 {
					// satisfies use case of x(y+z) = x*(y+z)
					terms[numTerms-1].Operator = Operator('*')
				}
				terms = append(terms, Term{Operand: value})
				isDefault = false
			} else {
				// the context is in the scope of the parenthetical value
				isOpen = true
			}
		} else if IsCloseParen(c) {
			// end of context, exit loop
			if isOpen {
				isOpen = false
				break
			} else {
				return 0, 0, SyntaxError(fmt.Sprintf("at char %d: missing '('", idx))
			}
		} else if IsOperator(c) {
			op := Operator(c)
			if numTerms == 0 || Order(terms[numTerms-1]) > 0 {
				if isDefault || isNegative {
					return 0, 0, SyntaxError(fmt.Sprintf("at char %d: missing operator", idx))
				}
				if c == '-' {
					// set negative value
					isNegative = true
				} else {
					// get the default value for the term
					t, err := Default(op)
					if err != nil {
						return 0, 0, SyntaxError(fmt.Sprintf("at char %d: missing operator", idx))
					}
					terms = append(terms, t)
					isDefault = true
				}
			} else {
				terms[numTerms-1].Operator = op
			}
		} else if IsDigit(c) {
			// satisfies use case 4 5 (error)
			if numTerms > 0 && Order(terms[numTerms-1]) <= 0 {
				return 0, 0, SyntaxError(fmt.Sprintf("at char %d: missing operator", idx))
			}
			value, idx, err = p.getNumber(idx)
			if err != nil {
				return
			}
			if isNegative {
				value = -value
			}
			terms = append(terms, Term{Operand: value})
			isNegative = false
			isDefault = false
		} else if IsSpace(c) {
			// ignore spaces
			continue
		} else {
			// all other charaters are bogus
			return 0, 0, SyntaxError(fmt.Sprintf("at char %d: unknown operator '%s'", idx, c))
		}
		idx++
	}
	if numTerms := len(terms); numTerms == 0 {
		return 0, 0, SyntaxError(fmt.Sprintf("at char %d: no data to evaluate", idx))
	} else if Order(terms[numTerms-1]) > 0 {
		return 0, 0, SyntaxError(fmt.Sprintf("at char %d: missing operand", idx))
	}
	value, err = Reduce(terms...)
	if err != nil {
		return 0, 0, err
	}
	return
}

// getNumber returns the numerical value of the phrase at a given index
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

package main

import (
	"bytes"
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

const (
	OrderAdd  = 1
	OrderSub  = 1
	OrderMult = 2
	OrderRoll = 3
)

var OpMap = map[byte]Operation{
	'+': Add{},
	'-': Sub{},
	'*': Mult{},
	'd': Roll{},
}

type Operation interface {
	Do(a, b int) int
	Order() int
}

type Add struct{}

func (op Add) Do(a, b int) int { return a + b }
func (op Add) Order() int      { return OrderAdd }

type Sub struct{}

func (op Sub) Do(a, b int) int { return a - b }
func (op Sub) Order() int      { return OrderSub }

type Mult struct{}

func (op Mult) Do(a, b int) int { return a * b }
func (op Mult) Order() int      { return OrderMult }

type Roll struct{}

func (op Roll) Do(a, b int) (total int) {
	if a <= 0 || b <= 0 {
		return
	}
	for i := 0; i < a; i++ {
		rand.Seed(time.Now().UnixNano())
		total += (rand.Int() % b) + 1
	}
	return
}
func (op Roll) Order() int { return OrderRoll }

type Phrase struct {
	A  []byte
	B  []byte
	Op Operation
}

func IsDigit(b byte) bool      { return unicode.IsDigit(rune(b)) }
func IsSpace(b byte) bool      { return unicode.IsSpace(rune(b)) }
func IsOpenParen(b byte) bool  { return b == '(' }
func IsCloseParen(b byte) bool { return b == ')' }

func Number(data []byte) (int, bool) {
	if val, err := strconv.ParseInt(string(data), 10, 32); err != nil {
		return 0, false
	} else {
		return int(val), true
	}
}

func TrimParen(data []byte) []byte {
	return bytes.TrimLeft(bytes.TrimRight(data, ")"), "(")
}

func Compute(exec []byte) (int, error) {
	phrase, err := compile(exec)
	if err != nil {
		return 0, err
	}
	bNum, ok := Number(phrase.B)
	if !ok {
		if bNum, err = Compute(TrimParen(phrase.B)); err != nil {
			return 0, err
		}
	}
	if phrase.Op != nil {
		aNum, ok := Number(phrase.A)
		if !ok {
			if aNum, err = Compute(TrimParen(phrase.A)); err != nil {
				return 0, err
			}
		}
		return phrase.Op.Do(aNum, bNum), nil
	}
	return bNum, nil
}

func compile(exec []byte) (*Phrase, error) {
	exec = bytes.ToLower(exec)
	phrase := &Phrase{}
	paren := 0
	for i, c := range exec {
		if op, ok := OpMap[c]; ok {
			if paren > 0 && (phrase.Op == nil || phrase.Op.Order() > op.Order()) {
				if len(phrase.B) == 0 {
					return nil, SyntaxError(fmt.Sprintf("char %d: '%s' missing operand", i, c))
				}
				phrase.A = append(phrase.A, phrase.B...)
				phrase.B = []byte{}
				phrase.Op = op
			} else {
				phrase.B = append(phrase.B, c)
			}
		} else if IsDigit(c) {
			phrase.B = append(phrase.B, c)
		} else if IsOpenParen(c) {
			phrase.B = append(phrase.B, c)
			paren++
		} else if IsCloseParen(c) {
			if paren <= 0 {
				return nil, SyntaxError(fmt.Sprintf("char %d: '%s' missing '('", i, c))
			}
			phrase.B = append(phrase.B, c)
			paren--
		} else if IsSpace(c) {
			continue
		} else {
			return nil, SyntaxError(fmt.Sprintf("char %d: invalid character '%s'", i, c))
		}
	}
	if paren != 0 {
		return nil, SyntaxError("missing ')'")
	}
	if phrase.Op != nil && len(phrase.B) == 0 {
		return nil, SyntaxError("missing operand")
	}

	return phrase, nil
}

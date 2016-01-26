package main

import . "gopkg.in/check.v1"

var _ = Suite(&TermTestSuite{})

type TermTestSuite struct{}

func (s *TermTestSuite) TestPrecedes(c *C) {
	t := Term{Operator: Operator(0)}
	c.Check(t.Precedes(Term{Operator: Operator(0)}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('+')}), Equals, false)
	c.Check(t.Precedes(Term{Operator: Operator('-')}), Equals, false)
	c.Check(t.Precedes(Term{Operator: Operator('*')}), Equals, false)
	c.Check(t.Precedes(Term{Operator: Operator('d')}), Equals, false)
	t = Term{Operator: Operator('+')}
	c.Check(t.Precedes(Term{Operator: Operator(0)}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('+')}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('-')}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('*')}), Equals, false)
	c.Check(t.Precedes(Term{Operator: Operator('d')}), Equals, false)
	t = Term{Operator: Operator('-')}
	c.Check(t.Precedes(Term{Operator: Operator(0)}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('+')}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('-')}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('*')}), Equals, false)
	c.Check(t.Precedes(Term{Operator: Operator('d')}), Equals, false)
	t = Term{Operator: Operator('*')}
	c.Check(t.Precedes(Term{Operator: Operator(0)}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('+')}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('-')}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('*')}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('d')}), Equals, false)
	t = Term{Operator: Operator('d')}
	c.Check(t.Precedes(Term{Operator: Operator(0)}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('+')}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('-')}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('*')}), Equals, true)
	c.Check(t.Precedes(Term{Operator: Operator('d')}), Equals, true)
}

func (s *TermTestSuite) TestEval(c *C) {
	t1 := Term{Operand: 2, Operator: Operator('+')}
	t2 := Term{Operand: 3, Operator: Operator('-')}
	result, err := t1.Eval(t2)
	c.Check(err, IsNil)
	c.Check(result, DeepEquals, Term{Operand: 2 + 3, Operator: Operator('-')})
	result, err = t2.Eval(t1)
	c.Check(err, IsNil)
	c.Check(result, DeepEquals, Term{Operand: 3 - 2, Operator: Operator('+')})
	t2.Operator = Operator(0)
	result, err = t2.Eval(t1)
	c.Check(err, Equals, ErrNotSupported)
	c.Check(result, DeepEquals, Term{})
}

func (s *TermTestSuite) TestValue(c *C) {
	t1 := Term{Operand: 2, Operator: Operator('+')}
	result, err := t1.Value()
	c.Check(err, Equals, ErrNotSupported)
	c.Check(result, Equals, 0)
	t1 = Term{Operand: 3, Operator: Operator(0)}
	result, err = t1.Value()
	c.Check(err, IsNil)
	c.Check(result, Equals, 3)
}

func (s *TermTestSuite) TestDefault(c *C) {
	t, err := Default(Operator(0))
	c.Check(err, Equals, ErrNotSupported)
	c.Check(t, DeepEquals, Term{})
	t, err = Default(Operator('+'))
	c.Check(err, Equals, ErrNotSupported)
	c.Check(t, DeepEquals, Term{})
	t, err = Default(Operator('-'))
	c.Check(err, Equals, ErrNotSupported)
	c.Check(t, DeepEquals, Term{})
	t, err = Default(Operator('*'))
	c.Check(err, Equals, ErrNotSupported)
	c.Check(t, DeepEquals, Term{})
	t, err = Default(Operator('d'))
	c.Check(err, IsNil)
	c.Check(t, DeepEquals, Term{Operand: 1, Operator: Operator('d')})
}

func (s *TermTestSuite) TestOrder(c *C) {
	c.Check(Order(Term{Operator: Operator('!')}), Equals, -1)
	c.Check(Order(Term{Operator: Operator(0)}), Equals, 0)
	c.Check(Order(Term{Operator: Operator('+')}), Equals, 1)
	c.Check(Order(Term{Operator: Operator('-')}), Equals, 1)
	c.Check(Order(Term{Operator: Operator('*')}), Equals, 2)
	c.Check(Order(Term{Operator: Operator('d')}), Equals, 3)
}

func (s *TermTestSuite) TestReduce(c *C) {
	value, err := Reduce()
	c.Check(err, IsNil)
	c.Check(value, Equals, 0)

	value, err = Reduce(Term{Operand: 2})
	c.Check(err, IsNil)
	c.Check(value, Equals, 2)
	value, err = Reduce(Term{Operand: 2, Operator: Operator('+')})
	c.Check(err, Equals, ErrNotSupported)
	c.Check(value, Equals, 0)

	value, err = Reduce(Term{Operand: 2, Operator: Operator('+')}, Term{Operand: 3})
	c.Check(err, IsNil)
	c.Check(value, Equals, 2+3)

	value, err = Reduce(
		Term{Operand: 2, Operator: Operator('-')},
		Term{Operand: 3, Operator: Operator('*')},
		Term{Operand: 4},
	)
	c.Check(err, IsNil)
	c.Check(value, Equals, 2-3*4)
	value, err = Reduce(
		Term{Operand: 2, Operator: Operator('*')},
		Term{Operand: 3, Operator: Operator('-')},
		Term{Operand: 4},
	)
	c.Check(err, IsNil)
	c.Check(value, Equals, 2*3-4)
}

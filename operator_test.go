package main

import . "gopkg.in/check.v1"

var _ = Suite(&OperatorTestSuite{})

type OperatorTestSuite struct{}

func (s *OperatorTestSuite) doAdd(c *C, op Op) {
	value, err := op.Default()
	c.Check(err, Equals, ErrNotSupported)
	c.Check(value, Equals, 0)
	value, err = op.Value(1)
	c.Check(err, Equals, ErrNotSupported)
	c.Check(value, Equals, 0)
	value, err = op.Eval(2, 3)
	c.Check(err, IsNil)
	c.Check(value, Equals, 2+3)
}

func (s *OperatorTestSuite) TestAdd(c *C) {
	s.doAdd(c, AddOp{})
	s.doAdd(c, Operator('+'))
}

func (s *OperatorTestSuite) doSub(c *C, op Op) {
	value, err := op.Default()
	c.Check(err, IsNil)
	c.Check(value, Equals, 0)
	value, err = op.Value(1)
	c.Check(err, Equals, ErrNotSupported)
	c.Check(value, Equals, 0)
	value, err = op.Eval(2, 3)
	c.Check(err, IsNil)
	c.Check(value, Equals, 2-3)
}

func (s *OperatorTestSuite) TestSub(c *C) {
	s.doSub(c, SubOp{})
	s.doSub(c, Operator('-'))
}

func (s *OperatorTestSuite) doMult(c *C, op Op) {
	value, err := op.Default()
	c.Check(err, Equals, ErrNotSupported)
	c.Check(value, Equals, 0)
	value, err = op.Value(1)
	c.Check(err, Equals, ErrNotSupported)
	c.Check(value, Equals, 0)
	value, err = op.Eval(2, 3)
	c.Check(err, IsNil)
	c.Check(value, Equals, 2*3)
}

func (s *OperatorTestSuite) TestMult(c *C) {
	s.doMult(c, MultOp{})
	s.doMult(c, Operator('*'))
}

func (s *OperatorTestSuite) doRoll(c *C, op Op) {
	value, err := op.Default()
	c.Check(err, Equals, nil)
	c.Check(value, Equals, 1)
	value, err = op.Value(1)
	c.Check(err, Equals, ErrNotSupported)
	c.Check(value, Equals, 0)
	value, err = op.Eval(2, 4)
	c.Check(err, IsNil)
	c.Logf("Rolled 2d4 and got %d", value)
	c.Check(value >= 2, Equals, true)
	c.Check(value <= 4*2, Equals, true)
}

func (s *OperatorTestSuite) TestRoll(c *C) {
	s.doRoll(c, RollOp{})
	s.doRoll(c, Operator('d'))
	s.doRoll(c, Operator('D'))
}

func (s *OperatorTestSuite) doNoOp(c *C, op Op) {
	value, err := op.Default()
	c.Check(err, Equals, ErrNotSupported)
	c.Check(value, Equals, 0)
	value, err = op.Value(1)
	c.Check(err, IsNil)
	c.Check(value, Equals, 1)
	value, err = op.Eval(2, 3)
	c.Check(err, Equals, ErrNotSupported)
	c.Check(value, Equals, 0)
}

func (s *OperatorTestSuite) TestNoOp(c *C) {
	s.doNoOp(c, NoOp{})
	s.doNoOp(c, Operator(0))
}

func (s *OperatorTestSuite) doBad(c *C, op Op) {
	value, err := op.Default()
	c.Check(err, Equals, ErrNotSupported)
	c.Check(value, Equals, 0)
	value, err = op.Value(1)
	c.Check(err, Equals, ErrNotSupported)
	c.Check(value, Equals, 0)
	value, err = op.Eval(2, 3)
	c.Check(err, Equals, ErrNotSupported)
	c.Check(value, Equals, 0)
}

func (s *OperatorTestSuite) TestBad(c *C) {
	s.doBad(c, Operator('('))
	s.doBad(c, Operator(')'))
	s.doBad(c, Operator(' '))
	s.doBad(c, Operator('7'))
	s.doBad(c, Operator('a'))
	s.doBad(c, Operator('G'))
}

func (s *OperatorTestSuite) TestIsOperator(c *C) {
	c.Check(IsOperator('*'), Equals, true)
	c.Check(IsOperator('('), Equals, false)
	c.Check(IsOperator(')'), Equals, false)
	c.Check(IsOperator('-'), Equals, true)
	c.Check(IsOperator('+'), Equals, true)
	c.Check(IsOperator(' '), Equals, false)
	c.Check(IsOperator('7'), Equals, false)
	c.Check(IsOperator('a'), Equals, false)
	c.Check(IsOperator('G'), Equals, false)
	c.Check(IsOperator('d'), Equals, true)
	c.Check(IsOperator('D'), Equals, true)
	c.Check(IsOperator(0), Equals, true)
}

func (s *OperatorTestSuite) TestOrder(c *C) {
	c.Check(Operator('*').Order(), Equals, 2)
	c.Check(Operator('(').Order(), Equals, -1)
	c.Check(Operator(')').Order(), Equals, -1)
	c.Check(Operator('-').Order(), Equals, 1)
	c.Check(Operator('+').Order(), Equals, 1)
	c.Check(Operator(' ').Order(), Equals, -1)
	c.Check(Operator('7').Order(), Equals, -1)
	c.Check(Operator('a').Order(), Equals, -1)
	c.Check(Operator('G').Order(), Equals, -1)
	c.Check(Operator('d').Order(), Equals, 3)
	c.Check(Operator('D').Order(), Equals, 3)
	c.Check(Operator(0).Order(), Equals, 0)
}

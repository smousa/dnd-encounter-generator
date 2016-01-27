package main

import . "gopkg.in/check.v1"

var _ = Suite(&PhraseTestSuite{})

type PhraseTestSuite struct{}

func (s *PhraseTestSuite) TestEvaluate(c *C) {
	equation := "  \t\t 6 "
	value, err := Evaluate([]byte(equation))
	c.Check(err, IsNil)
	c.Check(value, Equals, 6)

	equation = "-       6"
	value, err = Evaluate([]byte(equation))
	c.Check(err, IsNil)
	c.Check(value, Equals, -6)

	equation = "6 - -4"
	value, err = Evaluate([]byte(equation))
	c.Check(err, IsNil)
	c.Check(value, Equals, 6 - -4)

	equation = "6---4"
	value, err = Evaluate([]byte(equation))
	c.Check(IsSyntaxError(err), Equals, true)
	c.Check(value, Equals, 0)

	equation = "6+4-"
	value, err = Evaluate([]byte(equation))
	c.Check(IsSyntaxError(err), Equals, true)
	c.Check(value, Equals, 0)

	equation = "6++4"
	value, err = Evaluate([]byte(equation))
	c.Check(IsSyntaxError(err), Equals, true)
	c.Check(value, Equals, 0)

	equation = "d6"
	value, err = Evaluate([]byte(equation))
	c.Check(err, IsNil)
	c.Check(value >= 1 && value <= 6, Equals, true)
	c.Logf("Rolled a d6 and got a %d", value)

	equation = "10d6"
	value, err = Evaluate([]byte(equation))
	c.Check(err, IsNil)
	c.Check(value >= 10 && value <= 60, Equals, true)
	c.Logf("Rolled a 10d6 and got a %d", value)

	equation = "(10 + 5)*4"
	value, err = Evaluate([]byte(equation))
	c.Check(err, IsNil)
	c.Check(value, Equals, (10+5)*4)

	equation = "10d5 + 7(1d4-1)"
	value, err = Evaluate([]byte(equation))
	c.Check(err, IsNil)
	c.Check(value >= 10 && value <= 71, Equals, true)
	c.Logf("Rolled a 10d5+7(1d4-1) and got a %d", value)
}

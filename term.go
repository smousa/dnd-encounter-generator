package main

// Term describes the initial operand and its operator.
type Term struct {
	Operand  int
	Operator Operator
}

// Precedes determines if the order of the current term is greater than what is
// passed.
func (a Term) Precedes(b Term) bool {
	return a.Operator.Order() >= b.Operator.Order()
}

// Eval evaluates two terms and returns a new Term of the calculated value
// and the 'b' term's operand.
func (a Term) Eval(b Term) (Term, error) {
	value, err := a.Operator.Eval(a.Operand, b.Operand)
	if err != nil {
		return Term{}, err
	}
	return Term{Operand: value, Operator: b.Operator}, nil
}

// Value returns the evaluated value of a provided term.
func (a Term) Value() (int, error) {
	return a.Operator.Value(a.Operand)
}

// Default returns a term given a its default operand.
func Default(op Operator) (Term, error) {
	val, err := op.Default()
	if err != nil {
		return Term{}, err
	}
	return Term{Operand: val, Operator: op}, nil
}

// Order returns the hierarchical order of the term
func Order(term Term) int {
	return term.Operator.Order()
}

// Reduce returns the evaluated value of the prescribed list of terms.
func Reduce(terms ...Term) (value int, err error) {
	numTerms := len(terms)
	if numTerms == 0 {
		return 0, nil
	}
	for numTerms > 1 {
		if terms, err = reduce(terms); err != nil {
			return 0, err
		}
		numTerms = len(terms)
	}
	return terms[0].Value()
}

// reduce evaluates the hierachical order of a list of terms, and reduces the
// number of remaining terms by one.
func reduce(terms []Term) ([]Term, error) {
	switch len(terms) {
	case 0, 1:
		return terms, nil
	case 2:
		t, err := terms[0].Eval(terms[1])
		if err != nil {
			return nil, err
		}
		return []Term{t}, nil
	default:
		if terms[0].Precedes(terms[1]) {
			newterms, err := reduce(terms[:2])
			if err != nil {
				return nil, err
			}
			return append(newterms, terms[2:]...), nil
		} else {
			newterms, err := reduce(terms[1:])
			if err != nil {
				return nil, err
			}
			return append(terms[:1], newterms...), nil
		}
	}
}

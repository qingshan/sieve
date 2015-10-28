package ast

func Equals(a, b Node) bool {
	switch av := a.(type) {
	case *File:
		switch bv := b.(type) {
		case *File:
			return av.Name == bv.Name && equalsCommands(av.List, bv.List)
		}
	case *ControlCommand:
		switch bv := b.(type) {
		case *ControlCommand:
			return av.Name == bv.Name && Equals(av.Test, bv.Test) && equalsCommands(av.Block, bv.Block)
		}
	case *StopCommand:
		if _, ok := b.(*StopCommand); ok {
			return true
		}
	case *GenericCommand:
		switch bv := b.(type) {
		case *GenericCommand:
			return av.Name == bv.Name && equalsArguments(av.Arguments, bv.Arguments)
		}
	case *CommentCommand:
		switch bv := b.(type) {
		case *CommentCommand:
			return av.Style == bv.Style && av.Text == bv.Text
		}
	case *TrueTest:
		if _, ok := b.(*TrueTest); ok {
			return true
		}
	case *FalseTest:
		if _, ok := b.(*FalseTest); ok {
			return true
		}
	case *NotTest:
		switch bv := b.(type) {
		case *NotTest:
			return Equals(av.Test, bv.Test)
		}
	case *AnyofTest:
		switch bv := b.(type) {
		case *AnyofTest:
			return equalsTests(av.Tests, bv.Tests)
		}
	case *AllofTest:
		switch bv := b.(type) {
		case *AllofTest:
			return equalsTests(av.Tests, bv.Tests)
		}
	case *GenericTest:
		switch bv := b.(type) {
		case *GenericTest:
			return av.Name == bv.Name && equalsArguments(av.Arguments, bv.Arguments)
		}
	case NumberArgument:
		switch bv := b.(type) {
		case NumberArgument:
			return string(av) == string(bv)
		}
	case TagArgument:
		switch bv := b.(type) {
		case TagArgument:
			return string(av) == string(bv)
		}
	case StringArgument:
		switch bv := b.(type) {
		case StringArgument:
			if len(av) != len(bv) {
				return false
			}
			for i := range av {
				if av[i] != bv[i] {
					return false
				}
			}
			return true
		}
	case nil:
		return b == nil
	}
	return false
}

func equalsCommands(av, bv []Command) bool {
	if len(av) != len(bv) {
		return false
	}
	for i := range av {
		if !Equals(av[i], bv[i]) {
			return false
		}
	}
	return true
}

func equalsTests(av, bv []Test) bool {
	if len(av) != len(bv) {
		return false
	}
	for i := range av {
		if !Equals(av[i], bv[i]) {
			return false
		}
	}
	return true
}

func equalsArguments(av, bv []Argument) bool {
	if len(av) != len(bv) {
		return false
	}
	for i := range av {
		if !Equals(av[i], bv[i]) {
			return false
		}
	}
	return true
}

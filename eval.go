//-----------------------------------------------------------------------------
// Copyright (c) 2022 Detlef Stern
//
// This file is part of sxpf.
//
// sxpf is licensed under the latest version of the EUPL // (European Union
// Public License). Please see file LICENSE.txt for your rights and obligations
// under this license.
//-----------------------------------------------------------------------------

package sxpf

// Environment provides methods to evaluate a s-expression.
type Environment interface {
	// LookupForm returns the form associated with the given symbol.
	LookupForm(*Symbol) (*Form, error)

	// Evaluate the string. In many cases, strings evaluate to itself.
	EvaluateString(*String) (Value, error)

	// Evaluate the symbol. In many cases this result in returning a value
	// found in some internal lookup tables.
	EvaluateSymbol(*Symbol) (Value, error)

	// Evaluate the given list. In many cases this means to evaluate the first
	// element to a form and then call the form with the remaning elements
	// (possibly evaluated) as parameters.
	EvaluateList(*List) (Value, error)
}

// Evaluate the given s-expression value in the given environment.
func Evaluate(env Environment, value Value) (Value, error) {
	switch val := value.(type) {
	case *Symbol:
		return env.EvaluateSymbol(val)
	case *String:
		return env.EvaluateString(val)
	case *List:
		return env.EvaluateList(val)
	default:
		// Other types evaluate to themself
		return value, nil
	}
}

// EvaluateCall by trying to evaluate the first slice element as a form.
// If the first slice element is a form, the last returned value is true.
func EvaluateCall(env Environment, vals []Value) (Value, error, bool) {
	if len(vals) == 0 {
		return nil, nil, false
	}
	if sym, ok := vals[0].(*Symbol); ok {
		form, err := env.LookupForm(sym)
		if err != nil {
			return nil, err, true
		}
		params := vals[1:]
		if !form.IsSpecial() {
			var err error
			params, err = EvaluateSlice(env, params)
			if err != nil {
				return nil, err, true
			}
		}
		res, err := form.Call(env, params)
		return res, err, true
	}
	return nil, nil, false
}

// EvaluateSlice by evaluating all slice elements, returning a slice of
// the same length with the result values.
func EvaluateSlice(env Environment, vals []Value) (res []Value, err error) {
	res = make([]Value, len(vals))
	for i, value := range vals {
		res[i], err = Evaluate(env, value)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
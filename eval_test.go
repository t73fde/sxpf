//-----------------------------------------------------------------------------
// Copyright (c) 2022 Detlef Stern
//
// This file is part of sxpf.
//
// sxpf is licensed under the latest version of the EUPL // (European Union
// Public License). Please see file LICENSE.txt for your rights and obligations
// under this license.
//-----------------------------------------------------------------------------

package sxpf_test

import (
	"bytes"
	"testing"

	"github.com/t73fde/sxpf"
)

func TestEvaluate(t *testing.T) {
	testcases := []struct {
		src string
		exp string
	}{
		{"a", "A"},
		{`"a"`, `"a"`},
		{"(CAT a b)", `"AB"`},
		{"(QUOTE (A b) c)", "((A B) C)"},
	}
	env := newTestEnv()
	for i, tc := range testcases {
		expr, err := sxpf.ReadString(tc.src)
		if err != nil {
			t.Error(err)
			continue
		}
		val, err := sxpf.Evaluate(env, expr)
		if err != nil {
			t.Error(err)
			continue
		}
		got := val.String()
		if got != tc.exp {
			t.Errorf("%d: %v should evaluate to %v, but got: %v", i, tc.src, tc.exp, got)
		}
	}
}

type testEnv struct {
	symMap *sxpf.SymbolMap
}

func newTestEnv() *testEnv {
	symMap := sxpf.NewSymbolMap(nil)
	for _, form := range testForms {
		symMap.Add(sxpf.NewSymbol(form.Name()), form)
	}
	return &testEnv{symMap: symMap}
}

var testForms = []*sxpf.Form{
	sxpf.NewPrimForm(
		"CAT",
		false,
		func(env sxpf.Environment, args []sxpf.Value) (sxpf.Value, error) {
			var buf bytes.Buffer
			for _, arg := range args {
				buf.WriteString(arg.String())
			}
			return sxpf.NewString(buf.String()), nil
		},
	),
	sxpf.NewPrimForm(
		"QUOTE",
		true,
		func(env sxpf.Environment, args []sxpf.Value) (sxpf.Value, error) {
			return sxpf.NewList(args...), nil
		},
	),
}

func (te *testEnv) LookupForm(sym *sxpf.Symbol) (*sxpf.Form, error) {
	if val, found := te.symMap.Lookup(sym); found {
		if form, ok := val.(*sxpf.Form); ok {
			return form, nil
		}
	}
	return nil, sxpf.ErrNotFormBound(sym)
}

func (*testEnv) EvaluateSymbol(sym *sxpf.Symbol) (sxpf.Value, error) {
	return sym, nil
}

func (*testEnv) EvaluateString(str *sxpf.String) (sxpf.Value, error) { return str, nil }
func (e *testEnv) EvaluateList(lst *sxpf.List) (sxpf.Value, error) {
	vals := lst.GetValue()
	res, err, done := sxpf.EvaluateCall(e, vals)
	if done {
		return res, err
	}
	result, err := sxpf.EvaluateSlice(e, vals)
	if err != nil {
		return nil, err
	}
	return sxpf.NewList(result...), nil
}

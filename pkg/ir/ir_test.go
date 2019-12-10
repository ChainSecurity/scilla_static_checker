package ir_test

import (
	"testing"

	"github.com/ChainSecurity/scilla_static_checker/pkg/ir"
)

func TestExample(t *testing.T) {
	void := ir.EnumType{}
	term := ir.DataVar{&void}
	got, ok := term.DataType.(*ir.EnumType)
	if !ok {
		t.Fatalf("DataVar.Type not a ConsType")
	}
	expected := &void
	if got != expected {
		t.Errorf("DataVar.Type not a ConsType\ngot:\t\t%v\nexpected:\t%v", got, expected)
	}
}

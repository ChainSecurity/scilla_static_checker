package ir

import (
	"errors"
	"fmt"
	//"strings"
)

type Visitor interface {
	Visit(x IRNode) (w Visitor)
}

func Walk(v Visitor, node IRNode, prev_node IRNode) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *Proc:
		for i, _ := range n.Vars {
			Walk(v, &n.Vars[i], n)
		}
		for i, _ := range n.Plan {
			Walk(v, n.Plan[i], n)
		}
		if n.Jump != nil {
			Walk(v, n.Jump, n)
		}
	case *Save:
		for i, _ := range n.Path {
			Walk(v, n.Path[i], n)
		}
		Walk(v, n.Data, n)
	case *Load:
		for i, _ := range n.Path {
			Walk(v, n.Path[i], n)
		}
	case *Accept:
		// do nothing
	case *Send:
		Walk(v, n.Data, n)
	case *DataVar:
		Walk(v, n.DataType, n)
	case *TypeVar:
		Walk(v, n.Kind, n)
	case *SetKind:
		// do nothing
	case *IntType, *NatType, *RawType, *StrType, *BnrType, *ExcType, *MsgType:
		// do nothing
	case *MapType:
		Walk(v, n.KeyType, n)
		Walk(v, n.ValType, n)
	case *Int:
		Walk(v, n.IntType, n)
	case *Nat:
		Walk(v, n.NatType, n)
	case *Raw:
		Walk(v, n.RawType, n)
	case *Str:
		Walk(v, n.StrType, n)
	case *Bnr:
		Walk(v, n.BnrType, n)
	case *Exc:
		Walk(v, n.ExcType, n)
	case *Msg:
		Walk(v, n.MsgType, n)
		for k, _ := range n.Data {
			Walk(v, n.Data[k], n)
		}
	case *Map:
		Walk(v, n.MapType, n)
	case *AllDD:
		for i, _ := range n.Vars {
			Walk(v, n.Vars[i], n)
		}
		Walk(v, n.Term, n)
	case *AllTD:
		for i, _ := range n.Vars {
			Walk(v, n.Vars[i], n)
		}
		Walk(v, n.Term, n)
	case *AllTT:
		for i, _ := range n.Vars {
			Walk(v, n.Vars[i], n)
		}
		Walk(v, n.Term, n)
	case *AppDD:
		for i, _ := range n.Args {
			Walk(v, n.Args[i], n)
		}
		Walk(v, n.To, n)
	case *AppTD:
		for i, _ := range n.Args {
			Walk(v, n.Args[i], n)
		}
		Walk(v, n.To, n)
	case *AppTT:
		for i, _ := range n.Args {
			Walk(v, n.Args[i], n)
		}
		Walk(v, n.To, n)
	case *AbsDD:
		for i, _ := range n.Vars {
			Walk(v, &n.Vars[i], n)
		}
		Walk(v, n.Term, n)
	case *AbsTD:
		for i, _ := range n.Vars {
			Walk(v, &n.Vars[i], n)
		}
		Walk(v, n.Term, n)
	case *AbsTT:
		for i, _ := range n.Vars {
			Walk(v, &n.Vars[i], n)
		}
		Walk(v, n.Term, n)
	case *EnumType:
		for _, ts := range *n {
			for i, _ := range ts {
				Walk(v, ts[i], n)
			}
		}
	case *Enum:
		Walk(v, n.EnumType, n)
		for i, _ := range n.Data {
			Walk(v, n.Data[i], n)
		}
	case nil:
		fmt.Printf("NIL %T\n", prev_node)
	default:
		//fmt.Printf("- %T\n", n)
		panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
}

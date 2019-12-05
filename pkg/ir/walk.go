package ir

import (
	"errors"
	"fmt"
	//"strings"
)

type Visitor interface {
	Visit(x Node, prevX Node) (w Visitor)
}

func Walk(v Visitor, node Node, prev_node Node) {
	if v = v.Visit(node, prev_node); v == nil {
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
	case *Event:
		Walk(v, n.Data, n)
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
			Walk(v, &n.Vars[i], n)
		}
		Walk(v, n.Term, n)
	case *AllTD:
		for i, _ := range n.Vars {
			Walk(v, &n.Vars[i], n)
		}
		Walk(v, n.Term, n)
	case *AllTT:
		for i, _ := range n.Vars {
			Walk(v, &n.Vars[i], n)
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
		for _, ts := range n.Constructors {
			for i, _ := range ts {
				Walk(v, ts[i], n)
			}
		}
	case *PickProc:
		Walk(v, n.From, n)
		for i, _ := range n.With {
			Walk(v, &n.With[i], n)
		}
	case *ProcCase:
		Walk(v, &n.Bind, n)
		Walk(v, &n.Body, n)
	case *Bind:
		Walk(v, n.BindType, n)
		if n.Cond != nil {
			Walk(v, n.Cond, n)
		}
	case *Cond:
		for i, _ := range n.Data {
			Walk(v, &n.Data[i], n)
		}
	case *Builtin:
		//do Nothing
	case *Enum:
		Walk(v, n.EnumType, n)
		for i, _ := range n.Data {
			Walk(v, n.Data[i], n)
		}
	case *CallProc:
		for i, _ := range n.Args {
			Walk(v, n.Args[i], n)
		}
		Walk(v, n.To, n)
	case *PickData:
		Walk(v, n.From, n)
		for i, _ := range n.With {
			Walk(v, &n.With[i], n)
		}
	case *DataCase:
		Walk(v, &n.Bind, n)
		Walk(v, n.Body, n)
	case *ProcType:
		for i, _ := range n.Vars {
			Walk(v, &n.Vars[i], n)
		}
	default:
		//fmt.Printf("- %T\n", n)
		panic(errors.New(fmt.Sprintf("Unhandled type: %T", n)))
	}
}

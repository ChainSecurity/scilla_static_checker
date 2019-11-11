package ir

import (
	"errors"
	"fmt"
)

func setDefaultType(h map[string]Type, k string, t Type) (r Type) {
	var set bool
	if r, set = h[k]; !set {
		h[k] = t
		r = t
		set = true
	}
	return
}

func stackMapCopy(s map[string][]Data) map[string][]Data {
	sCopy := map[string][]Data{}
	for k, vals := range s {
		valsCopy := make([]Data, len(vals))
		copy(valsCopy, vals)
		sCopy[k] = valsCopy
	}
	return sCopy
}

func stackMapPush(s map[string][]Data, k string, v Data) {
	s[k] = append(s[k], v)
}

func stackMapPop(s map[string][]Data, k string) {
	n := len(s[k]) - 1
	s[k] = s[k][:n]
}

func stackMapPeek(s map[string][]Data, k string) (Data, bool) {
	l := len(s[k])
	if l == 0 {
		return nil, false
	}
	return s[k][l-1], true
}

func (builder *CFGBuilder) TypeOf(d Data) Type {
	switch x := d.(type) {
	case *Nat:
		return x.NatType
	case *DataVar:
		return x.DataType
	case *Load:
		if len(x.Path) != 0 {
			panic(errors.New("Load with non empty path is not implemented"))
		}
		fmt.Println(builder.fieldTypeMap, x.Slot, builder.fieldTypeMap["b"])
		return builder.fieldTypeMap[x.Slot]
	case *AbsTD:
		return builder.TypeOf(x.Term)
	case *AbsDD:
		return builder.TypeOf(x.Term)
	case *AppDD:
		return builder.TypeOf(x.To)
	case *AppTD:
		return builder.TypeOf(x.To)
	case *Builtin:
		return x.BuiltinType
	case *Enum:
		return x.EnumType
	default:
		panic(errors.New(fmt.Sprintf("builder.TypeOf not implemented %T\n", d)))
	}
	//return nil
}

func KindOf(t Type) Kind {
	fmt.Printf("builder.KindOf not implemented %T\n", t)
	return nil
}

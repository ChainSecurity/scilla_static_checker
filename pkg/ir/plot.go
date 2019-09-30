package ir

import (
	"errors"
	"fmt"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"strings"
)

type dotNode struct {
	id    int64
	ports []string
	name  string
}

func (n dotNode) ID() int64 { return n.id }
func (n dotNode) Attributes() []encoding.Attribute {
	var ports = []string{}
	for _, p := range n.ports {
		ports = append(ports, fmt.Sprintf("<%s> %s", p, p))
	}
	label := fmt.Sprintf("{%s | { %s }}", n.name, strings.Join(ports, " | "))
	attrs := []encoding.Attribute{
		{Key: "shape", Value: "record"},
		{Key: "label", Value: label},
	}
	return attrs
}

type dotPortedEdge struct {
	from, to    graph.Node
	fromPort    string
	fromCompass string
	toPort      string
	toCompass   string
}

func (e dotPortedEdge) From() graph.Node { return e.from }
func (e dotPortedEdge) To() graph.Node   { return e.to }
func (e dotPortedEdge) ReversedEdge() graph.Edge {
	e.from, e.to = e.to, e.from
	e.fromPort, e.toPort = e.toPort, e.fromPort
	e.fromCompass, e.toCompass = e.toCompass, e.fromCompass
	return e
}

func (e dotPortedEdge) Weight() float64 { return 0 }

func (e dotPortedEdge) FromPort() (port, compass string) {
	return e.fromPort, e.fromCompass
}
func (e dotPortedEdge) ToPort() (port, compass string) {
	return e.toPort, e.toCompass
}

type dotBuilder struct {
	nodeCounter int64
	nodes       []*dotNode
	edges       []*dotPortedEdge
}

func dotWalkType(b *dotBuilder, t Type) graph.Node {
	n := dotNode{
		b.getId(),
		[]string{},
		"Type",
	}
	return &n
}

func dotWalkWhen(b *dotBuilder, w *When) graph.Node {
	n := dotNode{
		b.getId(),
		[]string{"Data", fmt.Sprintf("Case: %s", w.Case)},
		"When",
	}
	b.nodes = append(b.nodes, &n)
	for _, v := range w.Data {
		vNode := dotWalkBind(b, v)
		e := dotPortedEdge{
			from:     n,
			to:       vNode,
			fromPort: "Data"}
		b.edges = append(b.edges, &e)

	}
	return &n
}

func dotWalkBind(b *dotBuilder, d *Bind) graph.Node {
	n := dotNode{
		b.getId(),
		[]string{"BindType", "When"},
		"Bind",
	}
	var m graph.Node
	m = dotWalkType(b, d.BindType)
	var e *dotPortedEdge
	e = &dotPortedEdge{
		from:     n,
		to:       m,
		fromPort: "BindType"}
	b.edges = append(b.edges, e)
	if d.When != nil {
		m = dotWalkWhen(b, d.When)
		e = &dotPortedEdge{
			from:     n,
			to:       m,
			fromPort: "BindType"}
		b.edges = append(b.edges, e)
	}
	return &n
}

func dotWalkDataCase(b *dotBuilder, d *DataCase) graph.Node {
	n := dotNode{
		b.getId(),
		[]string{"Bind", "Body"},
		"DataCase",
	}
	var e *dotPortedEdge
	e = &dotPortedEdge{
		from:     n,
		to:       dotWalkData(b, d.Body),
		fromPort: "Body"}
	b.edges = append(b.edges, e)
	e = &dotPortedEdge{
		from:     n,
		to:       dotWalkBind(b, d.Bind),
		fromPort: "Bind"}
	b.edges = append(b.edges, e)
	return &n
}

func dotWalkData(b *dotBuilder, d Data) graph.Node {
	switch x := d.(type) {
	case *AbsDD:
		n := dotNode{
			b.getId(),
			[]string{"Vars", "Term"},
			"AbsDD",
		}
		b.nodes = append(b.nodes, &n)
		for _, v := range x.Vars {
			vNode := dotWalkData(b, v)
			e := dotPortedEdge{
				from:     n,
				to:       vNode,
				fromPort: "Vars"}
			b.edges = append(b.edges, &e)

		}
		tNode := dotWalkData(b, x.Term)
		e := dotPortedEdge{
			from:     n,
			to:       tNode,
			fromPort: "Term"}
		b.edges = append(b.edges, &e)
		return &n
	case *DataVar:
		n := dotNode{
			b.getId(),
			[]string{"DataType"},
			"DataVar",
		}
		b.nodes = append(b.nodes, &n)
		tNode := dotWalkType(b, x.DataType)
		e := dotPortedEdge{
			from:     n,
			to:       tNode,
			fromPort: "DataType"}
		b.edges = append(b.edges, &e)
		return &n
	case *PickData:
		n := dotNode{
			b.getId(),
			[]string{"From", "With"},
			"PickData",
		}
		b.nodes = append(b.nodes, &n)
		for _, dc := range x.With {
			e := dotPortedEdge{
				from:     n,
				to:       dotWalkDataCase(b, &dc),
				fromPort: "With"}
			b.edges = append(b.edges, &e)

		}
		e := dotPortedEdge{
			from:     n,
			to:       dotWalkData(b, x.From),
			fromPort: "From"}
		b.edges = append(b.edges, &e)
		return &n
	case *Nat:
		n := dotNode{
			b.getId(),
			[]string{"NatType", fmt.Sprintf("Data: %s", x.Data)},
			"Nat",
		}
		b.nodes = append(b.nodes, &n)
		tNode := dotWalkType(b, x.NatType)
		e := dotPortedEdge{
			from:     n,
			to:       tNode,
			fromPort: "NatType"}
		b.edges = append(b.edges, &e)
		return &n
	default:
		panic(errors.New(fmt.Sprintf("unhandeled type: %T", x)))
	}
}

func (b *dotBuilder) getId() int64 {
	id := b.nodeCounter
	b.nodeCounter++
	return id
}
func directedPortedAttrGraphFrom(b *dotBuilder) graph.Directed {
	dg := simple.NewDirectedGraph()
	for _, e := range b.edges {
		dg.SetEdge(e)
	}
	return dg
}

func Plot(b *CFGBuilder) {
	keys := make([]string, 0, len(b.GlobalVarMap))
	for key := range b.GlobalVarMap {
		keys = append(keys, key)
	}
	d := dotBuilder{0, []*dotNode{}, []*dotPortedEdge{}}
	dotWalkData(&d, b.GlobalVarMap[keys[0]])
	g := directedPortedAttrGraphFrom(&d)
	got, err := dot.Marshal(g, "asd", "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Printf(string(got))

}

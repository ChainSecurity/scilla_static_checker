package ast

import (
	"encoding/json"
	"errors"
)

func getNodeType(rawMsg *json.RawMessage) (string, error) {
	type ntype struct {
		NodeType string `json:"node_type"`
	}
	var n ntype
	err := json.Unmarshal(*rawMsg, &n)
	if err != nil {
		return "", err
	}
	return n.NodeType, nil
}

func unmarshalExpression(rawMsg *json.RawMessage) (Expression, error) {
	ntype, err := getNodeType(rawMsg)
	if err != nil {
		return nil, err
	}

	switch ntype {
	case "LiteralExpression":
		var m LiteralExpression
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "VarExpression":
		var m VarExpression
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "LetExpression":
		var m LetExpression
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "MessageExpression":
		var m MessageExpression
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "FunExpression":
		var m FunExpression
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "AppExpression":
		var m AppExpression
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "ConstrExpression":
		var m ConstrExpression
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "MatchExpression":
		var m MatchExpression
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "BuiltinExpression":
		var m BuiltinExpression
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "TFunExpression":
		var m TFunExpression
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "TAppExpression":
		var m TAppExpression
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "FixpointExpression":
		var m FixpointExpression
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	default:
		return nil, errors.New("Unsupported type found!")
	}
}

func unmarshalLiteral(rawMsg *json.RawMessage) (Literal, error) {
	ntype, err := getNodeType(rawMsg)
	if err != nil {
		return nil, err
	}
	switch ntype {
	case "StringLiteral":
		var m StringLiteral
		err := json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "BNumLiteral":
		var m BNumLiteral
		err := json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "ByStrLiteral":
		var m ByStrLiteral
		err := json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "ByStrXLiteral":
		var m ByStrXLiteral
		err := json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "IntLiteral":
		var m IntLiteral
		err := json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "UintLiteral":
		var m UintLiteral
		err := json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "MapLiteral":
		var m MapLiteral
		err := json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "ADTValLiteral":
		var m ADTValLiteral
		err := json.Unmarshal(*rawMsg, &m)
		return &m, err
	default:
		return nil, errors.New("Unsupported type found!")
	}
}

func unmarshalStatement(rawMsg *json.RawMessage) (Statement, error) {
	ntype, err := getNodeType(rawMsg)
	if err != nil {
		return nil, err
	}
	switch ntype {
	case "LoadStatement":
		var m LoadStatement
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "StoreStatement":
		var m StoreStatement
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "BindStatement":
		var m BindStatement
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "MapUpdateStatement":
		var m MapUpdateStatement
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "MapGetStatement":
		var m MapGetStatement
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "MatchStatement":
		var m MatchStatement
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "ReadFromBCStatement":
		var m ReadFromBCStatement
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "AcceptPaymentStatement":
		var m AcceptPaymentStatement
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "SendMsgsStatement":
		var m SendMsgsStatement
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "CreateEvntStatement":
		var m CreateEvntStatement
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "CallProcStatement":
		var m CallProcStatement
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "ThrowStatement":
		var m ThrowStatement
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	default:
		return nil, errors.New("Unsupported type found!")
	}
}

func unmarshalPattern(rawMsg *json.RawMessage) (Pattern, error) {
	ntype, err := getNodeType(rawMsg)
	if err != nil {
		return nil, err
	}
	switch ntype {
	case "WildcardPattern":
		var m WildcardPattern
		err := json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "BinderPattern":
		var m BinderPattern
		err := json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "ConstructorPattern":
		var m ConstructorPattern
		err := json.Unmarshal(*rawMsg, &m)
		return &m, err
	default:
		return nil, errors.New("Unsupported type found!")
	}
}

func unmarshalPayload(rawMsg *json.RawMessage) (Payload, error) {
	ntype, err := getNodeType(rawMsg)
	if err != nil {
		return nil, err
	}
	switch ntype {
	case "PayloadLitral":
		var m PayloadLitral
		err := json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "PayloadVariable":
		var m PayloadVariable
		err := json.Unmarshal(*rawMsg, &m)
		return &m, err
	default:
		return nil, errors.New("Unsupported type found!")
	}
}

func unmarshalLibEntry(rawMsg *json.RawMessage) (LibEntry, error) {
	ntype, err := getNodeType(rawMsg)
	if err != nil {
		return nil, err
	}
	switch ntype {
	case "LibraryVariable":
		var m LibraryVariable
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	case "LibraryType":
		var m LibraryType
		err = json.Unmarshal(*rawMsg, &m)
		return &m, err
	default:
		return nil, errors.New("Unsupported type found!")
	}
}

func (le *LetExpression) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMsg *json.RawMessage
	rawMsg = objMap["expression"]
	e, err := unmarshalExpression(rawMsg)
	if err != nil {
		return err
	}

	rawMsg = objMap["body"]
	bd, err := unmarshalExpression(rawMsg)
	if err != nil {
		return err
	}

	type core struct {
		AnnotatedNode
		Var     *Identifier `json:"variable"`
		VarType string      `json:"variable_type"` //Optional
	}

	var c core
	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	le.Expr = &e
	le.Body = &bd
	le.AnnotatedNode = c.AnnotatedNode
	le.Var = c.Var
	le.VarType = c.VarType
	return nil
}

func (le *LiteralExpression) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMsg *json.RawMessage
	rawMsg = objMap["value"]
	v, err := unmarshalLiteral(rawMsg)
	if err != nil {
		return err
	}

	type core struct {
		AnnotatedNode
	}

	var c core
	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	le.AnnotatedNode = c.AnnotatedNode
	le.Val = &v
	return nil
}

func (f *Field) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMsg *json.RawMessage
	rawMsg = objMap["expression"]
	e, err := unmarshalExpression(rawMsg)
	if err != nil {
		return err
	}

	type core struct {
		Name *Identifier `json:"field_name"`
		Type string      `json:"field_type"`
	}

	var c core
	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	f.Name = c.Name
	f.Type = c.Type
	f.Expr = &e
	return nil
}

func (mec *MatchExpressionCase) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMsg *json.RawMessage
	rawMsg = objMap["pattern"]
	p, err := unmarshalPattern(rawMsg)
	if err != nil {
		return err
	}

	rawMsg = objMap["expression"]
	e, err := unmarshalExpression(rawMsg)
	if err != nil {
		return err
	}

	mec.Pat = &p
	mec.Expr = &e
	return nil
}

func (bs *BindStatement) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMsg *json.RawMessage
	rawMsg = objMap["rhs_expr"]
	e, err := unmarshalExpression(rawMsg)
	if err != nil {
		return err
	}

	type core struct {
		AnnotatedNode
		Lhs *Identifier `json:"lhs"`
	}

	var c core
	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	bs.RhsExpr = &e
	bs.AnnotatedNode = c.AnnotatedNode
	bs.Lhs = c.Lhs
	return nil
}

func (fe *FunExpression) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMsg *json.RawMessage
	rawMsg = objMap["rhs_expr"]
	e, err := unmarshalExpression(rawMsg)
	if err != nil {
		return err
	}

	type core struct {
		AnnotatedNode
		FunType string      `json:"fun_type"`
		Lhs     *Identifier `json:"lhs"`
	}

	var c core
	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	fe.RhsExpr = &e
	fe.AnnotatedNode = c.AnnotatedNode
	fe.Lhs = c.Lhs
	fe.FunType = c.FunType
	return nil
}

func (ma *MessageArgument) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMsg *json.RawMessage
	rawMsg = objMap["payload"]
	p, err := unmarshalPayload(rawMsg)
	if err != nil {
		return err
	}

	ma.Pl = &p

	type core struct {
		Var string `json:"variable"`
	}

	var c core
	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	ma.Var = c.Var
	return nil
}

func (pll *PayloadLitral) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMsg *json.RawMessage
	rawMsg = objMap["literal"]
	l, err := unmarshalLiteral(rawMsg)
	if err != nil {
		return err
	}

	pll.Lit = &l

	return nil
}

func (l *LibraryVariable) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMsg *json.RawMessage
	rawMsg = objMap["expression"]
	e, err := unmarshalExpression(rawMsg)
	if err != nil {
		return err
	}

	l.Expr = &e

	type core struct {
		VarType string `json:"variable_type"`
	}

	var c core
	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	l.VarType = c.VarType
	return nil
}

func (comp *Component) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMsgs []*json.RawMessage
	err = json.Unmarshal(*objMap["body"], &rawMsgs)
	if err != nil {
		return err
	}

	comp.Body = make([]*Statement, len(rawMsgs))
	for index, rawMsg := range rawMsgs {
		s, err := unmarshalStatement(rawMsg)
		if err != nil {
			return err
		}
		comp.Body[index] = &s
	}

	type core struct {
		Name   *Identifier  `json:"name"`
		Params []*Parameter `json:"params"`
	}

	var c core
	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	comp.Name = c.Name
	return nil
}

func (msc *MatchStatementCase) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMsgs []*json.RawMessage
	err = json.Unmarshal(*objMap["pattern_body"], &rawMsgs)
	if err != nil {
		return err
	}

	msc.Body = make([]*Statement, len(rawMsgs))
	for index, rawMsg := range rawMsgs {
		e, err := unmarshalStatement(rawMsg)
		if err != nil {
			return err
		}
		msc.Body[index] = &e
	}

	var rawMsg *json.RawMessage
	rawMsg = objMap["pattern"]
	p, err := unmarshalPattern(rawMsg)
	if err != nil {
		return err
	}

	msc.Pat = &p

	return nil
}

func (cp *ConstructorPattern) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMsgs []*json.RawMessage
	err = json.Unmarshal(*objMap["patterns"], &rawMsgs)
	if err != nil {
		return err
	}

	cp.Pats = make([]*Pattern, len(rawMsgs))
	for index, rawMsg := range rawMsgs {
		p, err := unmarshalPattern(rawMsg)
		if err != nil {
			return err
		}
		cp.Pats[index] = &p
	}

	type core struct {
		ConstrName string `json:"constructor_name"`
	}

	var c core
	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	cp.ConstrName = c.ConstrName
	return nil
}

func (l *Library) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMsgs []*json.RawMessage
	err = json.Unmarshal(*objMap["library_entries"], &rawMsgs)
	if err != nil {
		return err
	}

	l.Entries = make([]*LibEntry, len(rawMsgs))
	for index, rawMsg := range rawMsgs {
		e, err := unmarshalLibEntry(rawMsg)
		if err != nil {
			return err
		}
		l.Entries[index] = &e
	}

	type core struct {
		Name *Identifier `json:"library_name"`
	}

	var c core
	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	l.Name = c.Name
	return nil
}

func Parse_cmod(b []byte) *ContractModule {
	var c ContractModule
	if err := json.Unmarshal(b, &c); err != nil {
		//fmt.Println(err)
		panic(err)
	}
	//fmt.Println(c)
	return &c
}

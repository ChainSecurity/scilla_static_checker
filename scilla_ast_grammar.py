from __future__ import annotations

from typing import List, Optional

from grammar import production, abstract_production, ProductionOps


@abstract_production
class AstNode(ProductionOps):
    pass

@abstract_production
class Statement(AstNode):
    loc: Location

@abstract_production
class Expression(AstNode):
    loc: Location

@abstract_production
class Literal(AstNode):
    pass

@abstract_production
class Payload(AstNode):
    pass

@abstract_production
class Pattern(AstNode):
    pass

@abstract_production
class LibraryEntry(AstNode):
    name: Identifier

@production
class Location(AstNode):
    source_file: str
    line: int
    column: int

@production
class Identifier(AstNode):
    loc: Location
    identifier: str

@production
class CtrDef(AstNode):
    crt_def_name: Identifier
    c_arg_types: List[str]

@production
class StringLiteral(Literal):
    value: str

@production
class BNumLiteral(Literal):
    value: str

@production
class ByStrLiteral(Literal):
    value: str

@production
class ByStrXLiteral(Literal):
    value: str

@production
class IntLiteral(Literal):
    value: str

@production
class UintLiteral(Literal):
    value: str

@production
class MapValue(AstNode):
    key: Literal
    value: Literal

@production
class MapLiteral(Literal):
    key_type: str
    value_type: str
    node_type: str
    value: List[MapValue]

@production
class ADTValueLiteral(Literal):
    pass

@production
class LiteralExpression(Expression):
    value: Literal

@production
class VarExpression(Expression):
    variable: Identifier

@production
class LetExpression(Expression):
    variable: Identifier
    variable_type: Optional[str]
    expression: Expression
    body: Expression

@production
class PayloadLitral(Payload):
    literal: Literal

@production
class PayloadVariable(Payload):
    value: Identifier

@production
class MessageArgument(AstNode):
    variable: str
    payload: Payload

@production
class MessageExpression(Expression):
    arguments: List[MessageArgument]

@production
class FunExpression(Expression):
    lhs: Identifier
    rhs: Expression
    fun_type: str

@production
class AppExpression(Expression):
    lhs: Identifier
    rhs: List[Identifier]

@production
class ConstrExpression(Expression):
    types: List[str]
    constructor_name: str
    arguments: List[Identifier]

@production
class WildcardPattern(Pattern):
    pass

@production
class BinderPattern(Pattern):
    variable: Identifier

@production
class ConstructorPattern(Pattern):
    constructor_name: str
    patterns: List[Pattern]

@production
class MatchExpressionCase(AstNode):
    pattern: Pattern
    expression: Expression

@production
class MatchStatementCase(AstNode):
    pattern: Pattern
    pattern_body: List[Statement]

@production
class MatchExpression(Expression):
    lhs: Identifier
    rhs: List[MatchExpressionCase]

@production
class Builtin(AstNode):
    loc: Location
    builtin_type: str

@production
class BuiltinExpression(Expression):
    arguments: List[Identifier]
    builtin_function: Builtin

@production
class TFunExpression(Expression):
    pass

@production
class TAppExpression(Expression):
    pass

@production
class FixpointExpression(Expression):
    pass

@production
class LoadStatement(Statement):
    lhs: Identifier
    rhs: Identifier

@production
class StoreStatement(Statement):
    lhs: Identifier
    rhs: Identifier

@production
class BindStatement(Statement):
    lhs: Identifier
    rhs: Expression

@production
class MapUpdateStatement(Statement):
    map_name: Identifier
    rhs: Optional[Identifier]
    keys: List[Identifier]

@production
class MapGetStatement(Statement):
    map_name: Identifier
    lhs: Identifier
    keys: List[Identifier]
    is_value_retrieve: bool

@production
class MatchStatement(Statement):
    arg: Identifier
    cases: List[MatchStatementCase]

@production
class ReadFromBCStatement(Statement):
    lhs: Identifier
    rhs: str

@production
class AcceptPaymentStatement(Statement):
    pass

@production
class SendMsgsStatement(Statement):
    arg: Identifier

@production
class CreateEvntStatement(Statement):
    arg: Identifier

@production
class CallProcStatement(Statement):
    arg: Identifier
    messages: List[Identifier]

@production
class ThrowStatement(Statement):
    arg: Optional[Identifier]

@production
class LibraryVariable(LibraryEntry):
    variable_type: Optional[str]
    expression: Expression

@production
class LibraryType(LibraryEntry):
    ctr_defs: List[CtrDef]

@production
class Library(AstNode):
    library_name: Identifier
    library_entries: List[LibraryEntry]

@production
class ExternalLibrary(AstNode):
    name: Identifier
    alias: Optional[Identifier]

@production
class ContractModule(AstNode):
    scilla_major_version: int
    name: Identifier
    library: Optional[Library]
    external_libraries: List[ExternalLibrary]
    contract: Contract

@production
class Field(AstNode):
    field_name: Identifier
    field_type: str
    expression: Expression

@production
class Parameter(AstNode):
    parameter_name: Identifier
    parameter_type: str

@production
class Component(AstNode):
    component_name: Identifier
    params: List[Parameter]
    body: List[Statement]

@production
class Contract(AstNode):
    name: Identifier
    params: List[Parameter]
    fields: List[Field]
    components: List[Component]

package main

const (
	NodeTypeNumber = iota
	NodeTypeString
	NodeTypeExpression
	NodeTypeStatement
	NodeTypeStatements
	NodeTypeNumberValue
	NodeTypeStringValue
)

type Node interface {
	eval() Node
	asNumber() float64
	asString() string
	nodeType() int
	isTruthy() bool

	// for debugging
	String() string
}

package main

import "fmt"

type Primitive interface {
	Node
	isPrimitive() bool
}

type NumberValue struct {
	number float64
}

func NewNumberValue(f float64) *NumberValue {
	n := &NumberValue{number: f}
	return n
}

func (n *NumberValue) eval() Node {
	return n
}

func (n *NumberValue) asNumber() float64 {
	return n.number
}

func (n *NumberValue) asString() string {
	return fmt.Sprintf("%g", n.number)
}

func (n *NumberValue) nodeType() int {
	return NodeTypeNumberValue
}

func (n *NumberValue) String() string {
	return fmt.Sprintf("[Type: NumberValue]: value is '%g'", n.number)
}

func (n *NumberValue) isPrimitive() bool {
	return true
}

type StringValue struct {
	str string
}

func NewStringValue(str string) *StringValue {
	s := &StringValue{str: str}
	return s
}

func (s *StringValue) eval() Node {
	return s
}

func (s *StringValue) asNumber() float64 {
	return 0
}

func (s *StringValue) asString() string {
	return s.str
}

func (s *StringValue) nodeType() int {
	return NodeTypeStringValue
}

func (s *StringValue) String() string {
	return fmt.Sprintf("[Type: StringValue]: value is '%s'", s.str)
}

func (s *StringValue) isPrimitive() bool {
	return true
}

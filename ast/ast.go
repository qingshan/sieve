package ast

import (
	"bytes"
	"fmt"
)

type Node interface {
	String() string
}

type Test interface {
	Node
	testNode()
}

type Command interface {
	Node
	commandNode()
}

type Argument interface {
	Node
	argumentNode()
}

type CommentCommand struct {
	Style string
	Text  string
}

func (c *CommentCommand) String() string {
	return stringDepth(c, 0)
}

func (c *CommentCommand) commandNode() {
}

type StopCommand struct {
}

func (c *StopCommand) String() string {
	return stringDepth(c, 0)
}

func (c *StopCommand) commandNode() {
}

type GenericCommand struct {
	Name      string
	Arguments []Argument
}

func (c *GenericCommand) String() string {
	return stringDepth(c, 0)
}

func (c *GenericCommand) commandNode() {
}

type ControlCommand struct {
	Name  string
	Test  Test
	Block []Command
}

func (c *ControlCommand) String() string {
	return stringDepth(c, 0)
}

func (c *ControlCommand) commandNode() {
}

type TrueTest struct {
}

func (t *TrueTest) String() string {
	return "true";
}

func (c *TrueTest) testNode() {
}

type FalseTest struct {
}

func (t *FalseTest) String() string {
	return "false";
}

func (c *FalseTest) testNode() {
}

type NotTest struct {
	Test Test
}

func (t *NotTest) String() string {
	return "not " + t.Test.String()
}

func (c *NotTest) testNode() {
}

type AllofTest struct {
	Tests []Test
}

func (t *AllofTest) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("allof")
	buffer.WriteString(" (")
	for i, v := range t.Tests {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(v.String())
	}
	buffer.WriteString(")")
	return buffer.String()
}


func (c *AllofTest) testNode() {
}

type AnyofTest struct {
	Tests []Test
}

func (c *AnyofTest) testNode() {
}

func (t *AnyofTest) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("anyof")
	buffer.WriteString(" (")
	for i, v := range t.Tests {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(v.String())
	}
	buffer.WriteString(")")
	return buffer.String()
}


type GenericTest struct {
	Name      string
	Arguments []Argument
}

func (t *GenericTest) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(t.Name)
	for _, a := range t.Arguments {
		buffer.WriteString(" ")
		buffer.WriteString(a.String())
	}
	return buffer.String()
}

func (c *GenericTest) testNode() {
}

type NumberArgument struct {
	Value string
}

func (a *NumberArgument) String() string {
	return string(a.Value)
}

func (c *NumberArgument) argumentNode() {
}

type TagArgument struct {
	Value string
}

func (a *TagArgument) String() string {
	return string(a.Value)
}

func (c *TagArgument) argumentNode() {
}

type StringArgument struct {
	Value []string
}

func (a *StringArgument) String() string {
	switch len(a) {
	case 1:
		return a[0]
	default:
		var buffer bytes.Buffer
		buffer.WriteString("[")
		for i, v := range a.Value {
			if i > 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(v)
		}
		buffer.WriteString("]")
		return buffer.String();
	}
}

func (c *StringArgument) argumentNode() {
}

type File struct {
	Name string
	List []Command
}

func (f *File) String() string {
	b := new(bytes.Buffer)
	for _, c := range f.List {
		fmt.Fprintln(b, c)
	}
	return b.String()
}

func stringDepth(c Command, d int) string {
	var buffer bytes.Buffer
	for i := 0; i < d; i++ {
		buffer.WriteString("\t")
	}
	switch v := c.(type) {
	case *ControlCommand:
		buffer.WriteString(v.Name)
		if v.Test != nil {
			buffer.WriteString(" ")
			buffer.WriteString(v.Test.String())
		}
		buffer.WriteString(" ")
		buffer.WriteString("{\n")
		for _, bv := range v.Block {
			buffer.WriteString(stringDepth(bv, d + 1))
			buffer.WriteString("\n")
		}
		for i := 0; i < d; i++ {
			buffer.WriteString("\t")
		}
		buffer.WriteString("}")
	case *StopCommand:
		buffer.WriteString("stop;")
	case *CommentCommand:
		switch v.Style {
		case "line":
			buffer.WriteString(v.Style)
			buffer.WriteString(" ")
			buffer.WriteString(v.Text)
		default:
			buffer.WriteString("/*")
			buffer.WriteString(v.Text)
			buffer.WriteString("*/")
		}
	case *GenericCommand:
		buffer.WriteString(v.Name)
		for _, a := range v.Arguments {
			buffer.WriteString(" ")
			buffer.WriteString(a.String())
		}
		buffer.WriteString(";")
	}
	return buffer.String();
}

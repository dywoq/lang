package ast

import "github.com/dywoq/lang/pkg/token"

type Node interface {
	node()
}

type Tree struct {
	File       string   `json:"file"`
	Statements []Node   `json:"statements"`
	Global     []string `json:"global"`
}

type Declaration struct {
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
	Value      Node   `json:"value"`
}

type Function struct {
	Args []FunctionArgument `json:"args"`
	Body []Node             `json:"body"`
}

type FunctionArgument struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Variadic bool   `json:"variadic"`
}

type Instruction struct {
	Name string                `json:"name"`
	Args []InstructionArgument `json:"args"`
}

type InstructionArgument struct {
	Kind  token.Kind `json:"kind"`
	Value Node       `json:"value"`
}

type Value struct {
	Value string     `json:"value"`
	Kind  token.Kind `json:"kind"`
}

type ModifierConversion struct {
	Name  string `json:"name"`
	Value Node   `json:"value"`
}

func (Tree) node()                {}
func (Declaration) node()         {}
func (Function) node()            {}
func (FunctionArgument) node()    {}
func (Instruction) node()         {}
func (InstructionArgument) node() {}
func (Value) node()               {}
func (ModifierConversion) node()  {}

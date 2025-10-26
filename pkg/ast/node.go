package ast

type Node interface {
	node()
}

type Tree struct {
	File       string   `json:"file"`
	Statements Node     `json:"statements"`
	Global     []string `json:"global"`
}

type Variable struct {
	Identifier string `json:"identifier"`
	Exported   bool   `json:"exported"`
	Const      bool   `json:"const"`
	Consteval  bool   `json:"consteval"`
	Copied     bool   `json:"copied"`
	Type       string `json:"type"`
	Value      Node   `json:"value"`
}

type Function struct {
	ReturnType string             `json:"return_type"`
	Args       []FunctionArgument `json:"args"`
	Body       []Node             `json:"body"`
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
	Consteval bool `json:"consteval"`
	Copy      bool `json:"copy"`
	Value     Node `json:"value"`
}

func (Tree) node()                {}
func (Variable) node()            {}
func (Function) node()            {}
func (FunctionArgument) node()    {}
func (Instruction) node()         {}
func (InstructionArgument) node() {}

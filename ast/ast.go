package ast

// A Tree is defined by a single Root Node
type Tree Node

// A Node
type Node interface {
	// Type is used to identify nodes in a fast way
	Type() NodeType
}

// A BranchNode may hold 1:n child nodes
type BranchNode interface {
	Node
	AddChild(child Node)
	GetChildNodes() []Node
}

// NodeType defines the type of a specific node
type NodeType int

const (
	Root NodeType = iota + 1
	Doctype
	HTML
	Element
	Comment
	Text
	Component
	DynComponent
	Slot
	Head
	Body
	Style
	Script
	Link
	For
	If
	Template
	Content
)

// The RootNode is the top of an AST Tree (or Subtree)
type RootNode struct {
	BranchNodeImpl
}

func (n RootNode) Type() NodeType {
	return Root
}

// BranchNode is the base for all Nodes that may hold children
type BranchNodeImpl struct {
	ChildNodes []Node
}

func (bn *BranchNodeImpl) AddChild(child Node) {
	bn.ChildNodes = append(bn.ChildNodes, child)
}

func (bn *BranchNodeImpl) GetChildNodes() []Node {
	return bn.ChildNodes
}

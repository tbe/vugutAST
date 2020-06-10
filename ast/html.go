package ast

type AttributeHolder struct {
	StaticAttrs     map[string]string
	ExprAttrs       map[string]string
	AttributeLister string
}

type DoctypeNode struct {
	Data string
}

func (dn DoctypeNode) Type() NodeType {
	return Doctype
}

type HTMLNode struct {
	BranchNodeImpl
	AttributeHolder
}

func (n HTMLNode) Type() NodeType {
	return HTML
}

type ElementNode struct {
	BranchNodeImpl
	AttributeHolder
}

func (en ElementNode) Type() NodeType {
	return Element
}

type CommentNode struct {
	Comment string
}

func (cn CommentNode) Type() NodeType {
	return Comment
}

type TextNode struct {
	Text string
}

func (tn TextNode) Type() NodeType {
	return Text
}

type HeadNode struct {
	BranchNodeImpl
}

func (hn HeadNode) Type() NodeType {
	return Head
}

type BodyNode struct {
	BranchNodeImpl
	AttributeHolder
}

func (bn BodyNode) Type() NodeType {
	return Body
}

type StyleNode struct {
	BranchNodeImpl
	AttributeHolder
}

func (sn StyleNode) Type() NodeType {
	return Style
}

type ScriptNode struct {
	BranchNodeImpl
	AttributeHolder
}

func (sn ScriptNode) Type() NodeType {
	return Script
}

type LinkNode struct {
	BranchNodeImpl
	AttributeHolder
}

func (ln LinkNode) Type() NodeType {
	return Link
}

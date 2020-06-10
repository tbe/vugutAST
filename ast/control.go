package ast

type TemplateNode struct {
	BranchNodeImpl
}

func (tn TemplateNode) Type() NodeType {
	return Template
}

type ForNode struct {
	BranchNodeImpl
	Condition string
}

func (fn ForNode) Type() NodeType {
	return For
}

type IfNode struct {
	BranchNodeImpl
	Condition string
}

func (ifn IfNode) Type() NodeType {
	return If
}

type ContentNode struct {
	Content string
}

func (cn ContentNode) Type() NodeType {
	return Content
}
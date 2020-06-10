package ast

type ComponentNode struct {
	BranchNodeImpl
	Package               string
	Name                  string
	StaticAttributes      map[string]string
	ExprAttributes        map[string]string
	DynamicAttributes     map[string]string
	DynamicExprAttributes map[string]string
}

func (cn ComponentNode) Type() NodeType {
	return Component
}

type DynComponentNode struct {
	Expression string
}

func (dcn DynComponentNode) Type() NodeType {
	return DynComponent
}

type SlotNode struct {
	BranchNodeImpl
	Name string
}

func (sn SlotNode) Type() NodeType {
	return Slot
}

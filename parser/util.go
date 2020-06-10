package parser

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/vugu/html"

	"github.com/tbe/vuguAST/ast"
)

func hasUpperFirst(s string) bool {
	frune, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(frune)
}

type attributes struct {
	lowerCase map[string]string
	upperCase map[string]string
	expr      map[string]string
	vgIf      string
	vgFor     string
	vgContent string
	vgAttr    string
}

func parseAttributes(token html.Token) *attributes {
	attrs := &attributes{
		lowerCase: make(map[string]string),
		upperCase: make(map[string]string),
		expr:      make(map[string]string),
	}
	for _, a := range token.Attr {
		switch a.Key {
		case "vg-attr", ":":
			attrs.vgAttr = a.Val
		case "vg-if":
			attrs.vgIf = a.Val
		case "vg-for":
			attrs.vgFor = a.Val
		case "vg-content", "vg-html":
			attrs.vgContent = a.Val
		default:
			if strings.HasPrefix(a.OrigKey, ":") {
				attrs.expr[strings.TrimPrefix(a.OrigKey, ":")] = a.Val
			} else if hasUpperFirst(a.OrigKey) {
				attrs.upperCase[a.OrigKey] = a.Val
			} else {
				attrs.lowerCase[a.Key] = a.Val
			}
		}
	}

	return attrs
}

func createControlNodes(attrs *attributes) ast.BranchNode {
	var n ast.BranchNode
	if attrs.vgFor != "" {
		n = &ast.ForNode{
			Condition: attrs.vgFor,
		}
	}

	if attrs.vgIf != "" {
		ifn := &ast.ForNode{
			Condition: attrs.vgIf,
		}
		if n != nil {
			n.AddChild(ifn)
		} else {
			n = ifn
		}
	}
	return n
}

func normalizeAttrs(attrs *attributes) ast.AttributeHolder {
	// create an attributeHolder struct
	attrHolder := ast.AttributeHolder{
		StaticAttrs:     attrs.lowerCase,
		ExprAttrs:       attrs.expr,
		AttributeLister: attrs.vgAttr,
	}
	// append the uppercase attributes
	for k, v := range attrs.upperCase {
		// TODO: this could override some attributes, we should catch this and warn about it
		attrHolder.StaticAttrs[strings.ToLower(k)] = v
	}

	return attrHolder
}

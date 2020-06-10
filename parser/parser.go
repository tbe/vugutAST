package parser

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/vugu/html"

	"github.com/tbe/vuguAST/ast"
)

type Parser struct {
}

func (p *Parser) Parse(r io.Reader, fname string) (ast.Tree, error) {
	tree := &ast.RootNode{}
	tokenizer := html.NewTokenizer(r)

	return tree, p.traverse(tokenizer, tree)

}

func (p *Parser) traverse(tokenizer *html.Tokenizer, parent ast.BranchNode) error {
Loop:
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				break Loop
			}
			return err
		case html.DoctypeToken:
			parent.AddChild(&ast.DoctypeNode{Data: string(tokenizer.Text())})
		case html.TextToken:
			parent.AddChild(&ast.TextNode{Text: string(tokenizer.Text())})
		case html.CommentToken:
			parent.AddChild(&ast.CommentNode{Comment: string(tokenizer.Text())})
		case html.SelfClosingTagToken, html.StartTagToken:
			child, err := p.parseTokenByTag(tokenizer)
			if err != nil {
				return err
			}
			parent.AddChild(child)
		case html.EndTagToken:
			break Loop
		}
	}
	return nil
}

func (p *Parser) parseTokenByTag(tokenizer *html.Tokenizer) (ast.Node, error) {
	_, origTag, _ := tokenizer.TagNameAndOrig()

	// check if we have a component
	if strings.Contains(string(origTag), ":") {
		// check if this is an component
		compParts := strings.SplitN(string(origTag), ":", 2)
		if hasUpperFirst(compParts[1]) {
			return p.parseComponent(tokenizer, compParts[0], compParts[1])
		}
	}

	token := tokenizer.Token()
	isTraversable := token.Type == html.StartTagToken
	attrs := parseAttributes(token)

	// catch our special cases
	switch token.Data {
	case "vg-template":
		return p.parseBranchNode(&token, tokenizer, &ast.TemplateNode{}, attrs)
	case "vg-comp":
		// verify we have the required expr attribute
		expr, ok := attrs.lowerCase["expr"]
		if !ok {
			return nil, errors.New("missing expr attribute for vg-comp")
		}
		compnode := &ast.DynComponentNode{Expression: expr}
		// vg-comp is not a branch node, but we have to make sure that we reach the closing tag
		if isTraversable {
			// we create a new root node as "dummy" to satisfy the requirements of traverse
			if err := p.traverse(tokenizer, &ast.RootNode{}); err != nil {
				return nil, err
			}
		}
		cntrlNode := createControlNodes(attrs)
		if cntrlNode != nil {
			cntrlNode.AddChild(compnode)
			return cntrlNode, nil
		}
		return cntrlNode, nil

	case "vg-slot":
		// a vg-slot needs a name
		name, ok := attrs.lowerCase["name"]
		if !ok {
			return nil, errors.New("missing name attribute for vg-slot")
		}
		return p.parseBranchNode(&token, tokenizer, &ast.SlotNode{Name: name}, attrs)
	case "html":
		return p.parseBranchNode(&token, tokenizer, &ast.HTMLNode{AttributeHolder: normalizeAttrs(attrs)}, attrs)
	case "head":
		return p.parseBranchNode(&token, tokenizer, &ast.HeadNode{}, attrs)
	case "body":
		return p.parseBranchNode(&token, tokenizer, &ast.BodyNode{AttributeHolder: normalizeAttrs(attrs)}, attrs)
	case "script":
		return p.parseBranchNode(&token, tokenizer, &ast.ScriptNode{AttributeHolder: normalizeAttrs(attrs)}, attrs)
	case "style":
		return p.parseBranchNode(&token, tokenizer, &ast.ScriptNode{AttributeHolder: normalizeAttrs(attrs)}, attrs)
	case "link":
		return p.parseBranchNode(&token, tokenizer, &ast.LinkNode{AttributeHolder: normalizeAttrs(attrs)}, attrs)
	default:
		return p.parseBranchNode(&token, tokenizer, &ast.ElementNode{AttributeHolder: normalizeAttrs(attrs)}, attrs)

	}
}

func (p *Parser) parseComponent(tokenizer *html.Tokenizer, pkg, name string) (ast.Node, error) {
	// get the attributes first
	token := tokenizer.Token()
	attrs := parseAttributes(token)

	comp := &ast.ComponentNode{
		Package:               pkg,
		Name:                  name,
		StaticAttributes:      attrs.upperCase,
		ExprAttributes:        make(map[string]string),
		DynamicAttributes:     attrs.lowerCase,
		DynamicExprAttributes: make(map[string]string),
	}

	for k, v := range attrs.expr {
		if hasUpperFirst(k) {
			comp.ExprAttributes[k] = v
		} else {
			comp.DynamicExprAttributes[k] = v
		}
	}

	if token.Type != html.SelfClosingTagToken {
		// traverse here, if this is not self closing
		if err := p.traverse(tokenizer, comp); err != nil {
			return nil, err
		}
	}

	// TODO: what should we do with vg-content here? is this the default slot?

	// check if we have control blocks in before
	if cnode := createControlNodes(attrs); cnode != nil {
		cnode.AddChild(comp)
		return cnode, nil
	}

	return comp, nil
}

func (p *Parser) parseBranchNode(token *html.Token, tokenizer *html.Tokenizer, node ast.BranchNode, attrs *attributes) (ast.Node, error) {
	isTraversable := token.Type != html.SelfClosingTagToken

	// if we have vg-content set, we use this and ignore all the other stuff
	if attrs.vgContent != "" {
		node.AddChild(&ast.ContentNode{Content: attrs.vgContent})
		if isTraversable {
			dummyNode := &ast.RootNode{}
			if err := p.traverse(tokenizer, dummyNode); err != nil {
				return nil, err
			}
			if len(dummyNode.GetChildNodes()) > 0 {
				return nil, fmt.Errorf("`%s` contains content and has a `vg-content` attribute", token.OrigData)
			}
		}
	} else {
		if isTraversable {
			if err := p.traverse(tokenizer, node); err != nil {
				return nil, err
			}
		}
	}

	// check for control functions
	cntrlNode := createControlNodes(attrs)
	if cntrlNode != nil {
		cntrlNode.AddChild(node)
		return cntrlNode, nil
	}

	return node, nil
}


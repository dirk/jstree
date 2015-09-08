package jstree

import "fmt"
import "reflect"
import "strings"

func indentFor(indent int) string {
	s := strings.Repeat(" ", indent * 2)
	return s[:len(s) - 1]
}

func (p *Program) Print() {
	println("Program")
	indent := 0
	for _, node := range p.body {
		printNode(node, indent + 1)
	}
}

func (e *ExportDefaultDeclaration) Print(indent int) {
	println(indentFor(indent), "ExportDefault")
	printNode(e.declaration, indent + 1)
}

func (e *ExportNamedDeclaration) Print(indent int) {
	println(indentFor(indent), "ExportNamedDeclaration")
	for _, node := range e.specifiers {
		printNode(node, indent + 1)
	}
}

func (e *ExportSpecifier) Print(indent int) {
	exported := e.exported.name
	local    := e.local.name
	fmt.Printf("%s ExportSpecifier local:%s exported:%s\n", indentFor(indent), local, exported)
}

func (i *ImportDeclaration) Print(indent int) {
	fmt.Printf("%s ImportDeclaration '%s'\n", indentFor(indent), i.source.value)
	for _, node := range i.specifiers {
		printNode(node, indent + 1)
	}
}

func (i *ImportDefaultSpecifier) Print(indent int) {
	println(indentFor(indent), "ImportDefaultSpecifier")
	printNode(i.local, indent + 1)
}

func (i *ImportSpecifier) Print(indent int) {
	imported := i.imported.name
	local    := i.local.name
	fmt.Printf("%s ImportSpecifier local:%s imported:%s\n", indentFor(indent), local, imported)
}

func (f *FunctionDeclaration) Print(indent int) {
	println(indentFor(indent), "FunctionDeclaration")
	printNode(f.id, indent + 1)
	printNode(f.body, indent + 1)
}

func (i *Identifier) Print(indent int) {
	fmt.Printf("%s Identifier '%s'\n", indentFor(indent), i.name)
}

func (b *BlockStatement) Print(indent int) {
	println(indentFor(indent), "BlockStatement")
	for _, node := range b.Body {
		printNode(node, indent + 1)
	}
}

func (r *ReturnStatement) Print(indent int) {
	println(indentFor(indent), "ReturnStatement")
	if r.argument != nil {
		printNode(r.argument, indent + 1)
	}
}

func (b *BinaryExpression) Print(indent int) {
	fmt.Printf("%s BinaryExpression '%s'\n", indentFor(indent), b.operator)
	printNode(b.left,  indent + 1)
	printNode(b.right, indent + 1)
}

func (n *Literal) Print(indent int) {
	fmt.Printf("%s Literal '%s'\n", indentFor(indent), n.value)
}

func (d *VariableDeclaration) Print(indent int) {
	fmt.Printf("%s VariableDeclaration %s\n", indentFor(indent), d.kind)
	for _, declarator := range d.declarations {
		printNode(declarator, indent + 1)
	}
}

func (d *VariableDeclarator) Print(indent int) {
	println(indentFor(indent), "VariableDeclarator")
	printNode(d.id,   indent + 1)
	printNode(d.init, indent + 1)
}

func printNode(node interface{}, indent int) {
	switch n := node.(type) {
	case *BinaryExpression: n.Print(indent)
	case *BlockStatement: n.Print(indent)
	case *ExportDefaultDeclaration: n.Print(indent)
	case *ExportNamedDeclaration:   n.Print(indent)
	case *ExportSpecifier:          n.Print(indent)
	case *ImportDeclaration: n.Print(indent)
	case *ImportDefaultSpecifier: n.Print(indent)
	case *ImportSpecifier: n.Print(indent)
	case *FunctionDeclaration: n.Print(indent)
	case *Identifier: n.Print(indent)
	case *Literal: n.Print(indent)
	case *ReturnStatement: n.Print(indent)
	case *VariableDeclaration: n.Print(indent)
	case *VariableDeclarator: n.Print(indent)
	default:
		fmt.Println(indentFor(indent), "Unknown type:", reflect.TypeOf(node))
	}
}


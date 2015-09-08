package jstree

import "fmt"
import "log"
import "reflect"

import "github.com/bitly/go-simplejson"

type Position struct {
	start int
	end int
}

type Program struct {
	Position
	body []interface{}
}

type FunctionDeclaration struct {
	Position
	id interface{}
	generator bool
	expression bool
	params []interface{}
	body interface{}
}

type Identifier struct {
	Position
	name string
}

type ExportDefaultDeclaration struct {
	Position
	declaration interface{}
}

type ExportNamedDeclaration struct {
	Position
	specifiers []*ExportSpecifier
}

type ExportSpecifier struct {
	Position
	exported *Identifier
	local *Identifier
}

type ImportDefaultSpecifier struct {
	Position
	local *Identifier
}

type ImportSpecifier struct {
	Position
	imported *Identifier
	local *Identifier
}

type ImportDeclaration struct {
	Position
	specifiers []interface{}
	source *Literal
}

type ReturnStatement struct {
	Position
	argument interface{}
}

type Literal struct {
	Position
	value string
}

type BinaryExpression struct {
	Position
	left interface{}
	operator string
	right interface{}
}

type BlockStatement struct {
	Position
	Body []interface{}
}

type ExpressionStatement struct {
	Position
	Expression interface{}
}

type VariableDeclaration struct {
	Position
	kind string
	declarations []*VariableDeclarator
}

type VariableDeclarator struct {
	Position
	id interface{}
	init interface{}
}

type UpdateExpression struct {
	Position
	Operator string
	Prefix   bool
	Argument interface{}
}

type ForStatement struct {
	Position
	Init interface{}
	Test interface{}
	Update interface{}
	Body interface{}
}


func parseBody(sourceBody *simplejson.Json) ([]interface{}, error) {
	body := make([]interface{}, 0)
	for index := 0; ; index++ {
		sourceNode := sourceBody.GetIndex(index)
		// Break if we hit a nil node
		_, present := sourceNode.Map()
		if present != nil { break }

		node, err := parseNode(sourceNode)
		if err != nil { return nil, err }

		body = append(body, node)
	}
	return body, nil
}

func isEmptyJson(j *simplejson.Json) bool {
	// Create a pointer to the zero value of the simplejson.Json type
	// and get the interface{} pointer to it.
	zeroInterface := reflect.New(reflect.TypeOf(j).Elem()).Interface()

	// Convert it to a properly typed pointer for the DeepEqual comparison.
	zero := (zeroInterface).(*simplejson.Json)

	return reflect.DeepEqual(zero, j)
}

func parseNodeIfPresent(n *simplejson.Json) (interface{}, error) {
	if isEmptyJson(n) {
		return nil, nil
	}
	return parseNode(n)
}

func ParseForStatement(f *simplejson.Json) (*ForStatement, error) {
	init, err := parseNodeIfPresent(f.Get("init"))
	if err != nil { return nil, err }

	test, err := parseNodeIfPresent(f.Get("test"))
	if err != nil { return nil, err }

	update, err := parseNodeIfPresent(f.Get("update"))
	if err != nil { return nil, err }

	body, err := parseNodeIfPresent(f.Get("body"))
	if err != nil { return nil, err }

	return &ForStatement{
		Position: newPosition(f),
		Init:     init,
		Test:     test,
		Update:   update,
		Body:     body,
	}, nil
}

func ParseFunctionDeclaration(f *simplejson.Json) (*FunctionDeclaration, error) {
	generator, err := f.Get("generator").Bool()
	if err != nil { return nil, err }

	expression, err := f.Get("expression").Bool()
	if err != nil { return nil, err }

	id, err := ParseIdentifier(f.Get("id"))
	if err != nil { return nil, err }

	body, err := ParseBlockStatement(f.Get("body"))
	if err != nil { return nil, err }

	return &FunctionDeclaration{
		Position:   newPosition(f),
		id:         id,
		generator:  generator,
		expression: expression,
		params:     []interface{}{},
		body:       body,
	}, nil
}

func ParseBlockStatement(b *simplejson.Json) (*BlockStatement, error) {
	typ := b.Get("type").MustString()
	if typ != "BlockStatement" {
		log.Fatal("Expected BlockStatement, got %s", typ)
	}

	body, err := parseBody(b.Get("body"))
	if err != nil { return nil, err }

	return &BlockStatement{
		Position: newPosition(b),
		Body:     body,
	}, nil
}

func ParseExpressionStatement(e *simplejson.Json) (*ExpressionStatement, error) {
	expr, err := parseNode(e.Get("expression"))
	if err != nil { return nil, err }

	return &ExpressionStatement{
		Position:   newPosition(e),
		Expression: expr,
	}, nil


}

func ParseExportDefaultDeclaration(e *simplejson.Json) (*ExportDefaultDeclaration, error) {
	declaration, err := parseNode(e.Get("declaration"))
	if err != nil { return nil, err }

	return &ExportDefaultDeclaration{
		Position:    newPosition(e),
		declaration: declaration,
	}, nil
}

func ParseExportNamedDeclaration(e *simplejson.Json) (*ExportNamedDeclaration, error) {
	// fmt.Printf("%#v\n", e)

	exportSpecifiers := e.Get("specifiers")
	specifiers := make([]*ExportSpecifier, 0)

	for i := 0; ; i++ {
		exportSpecifier := exportSpecifiers.GetIndex(i)
		// Try to convert to a map and break if it's not there
		_, present := exportSpecifier.Map()
		if present != nil { break }

		maybeSpecifier, err := parseNode(exportSpecifier)
		if err != nil { return nil, err }

		specifier, isSpecifier := maybeSpecifier.(*ExportSpecifier)
		if !isSpecifier {
			return nil, fmt.Errorf("Expected ImportSpecifier, got %#v", maybeSpecifier)
		}
		specifiers = append(specifiers, specifier)
	}

	return &ExportNamedDeclaration{
		Position:   newPosition(e),
		specifiers: specifiers,
	}, nil
}

func ParseExportSpecifier(s *simplejson.Json) (*ExportSpecifier, error) {
	// fmt.Printf("%#v\n", s)

	exported, err := ParseIdentifier(s.Get("exported"))
	if err != nil { return nil, err }

	local, err := ParseIdentifier(s.Get("local"))
	if err != nil { return nil, err }

	return &ExportSpecifier{
		Position: newPosition(s),
		exported: exported,
		local:    local,
	}, nil
}

func ParseImportDeclaration(i *simplejson.Json) (*ImportDeclaration, error) {
	source, err := ParseLiteral(i.Get("source"))
	if err != nil { return nil, err }

	importSpecifiers := i.Get("specifiers")
	specifiers := make([]interface{}, 0)

	for i := 0; ; i++ {
		importSpecifier := importSpecifiers.GetIndex(i)
		// Try to convert to a map and break if it's not there
		_, present := importSpecifier.Map()
		if present != nil { break }

		specifier, err := parseNode(importSpecifier)
		if err != nil { return nil, err }

		_, isDefaultSpecifier := specifier.(*ImportDefaultSpecifier)
		_, isSpecifier        := specifier.(*ImportSpecifier)

		if !isDefaultSpecifier && !isSpecifier {
			return nil, fmt.Errorf("Expected ImportDefaultSpecifier or ImportSpecifier, got %#v", specifier)
		}
		specifiers = append(specifiers, specifier)
	}

	return &ImportDeclaration{
		Position:   newPosition(i),
		specifiers: specifiers,
		source:     source,
	}, nil
}

func ParseImportDefaultSpecifier(s *simplejson.Json) (*ImportDefaultSpecifier, error) {
	local, err := ParseIdentifier(s.Get("local"))
	if err != nil { return nil, err }

	return &ImportDefaultSpecifier{
		Position: newPosition(s),
		local:    local,
	}, nil
}

func ParseImportSpecifier(s *simplejson.Json) (*ImportSpecifier, error) {
	imported, err := ParseIdentifier(s.Get("imported"))
	if err != nil { return nil, err }

	local, err := ParseIdentifier(s.Get("local"))
	if err != nil { return nil, err }

	return &ImportSpecifier{
		Position: newPosition(s),
		imported: imported,
		local:    local,
	}, nil
}

func ParseIdentifier(i *simplejson.Json) (*Identifier, error) {
	return &Identifier{
		Position: newPosition(i),
		name:     i.Get("name").MustString(),
	}, nil
}

func ParseReturnStatement(r *simplejson.Json) (*ReturnStatement, error) {
	argument, err := parseNode(r.Get("argument"))
	if err != nil { return nil, err }

	return &ReturnStatement{
		Position: newPosition(r),
		argument: argument,
	}, nil
}

func ParseLiteral(n *simplejson.Json) (*Literal, error) {
	return &Literal{
		Position: newPosition(n),
		value:    n.Get("value").MustString(),
	}, nil
}

func ParseBinaryExpression(b *simplejson.Json) (*BinaryExpression, error) {
	left, err := parseNode(b.Get("left"))
	if err != nil { return nil, err }

	right, err := parseNode(b.Get("right"))
	if err != nil { return nil, err }

	return &BinaryExpression{
		Position: newPosition(b),
		left:     left,
		operator: b.Get("operator").MustString(),
		right:    right,
	}, nil
}

func ParseVariableDeclaration(d *simplejson.Json) (interface{}, error) {
	var declarator *VariableDeclarator
	var ok bool

	decls := make([]*VariableDeclarator, 0)
	sourceDeclarations := d.Get("declarations")
	for i := 0; ; i++ {
		sourceDeclarator := sourceDeclarations.GetIndex(i)
		// Try to convert to a map and break if it's not there
		_, present := sourceDeclarator.Map()
		if present != nil { break }

		maybeDeclarator, err := parseNode(sourceDeclarator)
		if err != nil { return nil, err }

		declarator, ok = maybeDeclarator.(*VariableDeclarator)
		if !ok {
			return nil, fmt.Errorf("Expected VariableDeclarator, got %#v", maybeDeclarator)
		}
		decls = append(decls, declarator)
	}

	return &VariableDeclaration{
		Position:     newPosition(d),
		kind:         d.Get("kind").MustString(),
		declarations: decls,
	}, nil
}

func ParseVariableDeclarator(d *simplejson.Json) (*VariableDeclarator, error) {
	id, err := parseNode(d.Get("id"))
	if err != nil { return nil, err }

	init, err := parseNode(d.Get("init"))
	if err != nil { return nil, err }

	return &VariableDeclarator{
		Position: newPosition(d),
		id:       id,
		init:     init,
	}, nil
}

func ParseUpdateExpression(u *simplejson.Json) (*UpdateExpression, error) {
	operator := u.Get("operator").MustString()

	prefix, err := u.Get("prefix").Bool()
	if err != nil { return nil, err }

	argument, err := parseNode(u.Get("argument"))
	if err != nil { return nil, err }

	return &UpdateExpression{
		Position: newPosition(u),
		Operator: operator,
		Prefix:   prefix,
		Argument: argument,
	}, nil
}

// Single-threaded program parsing -------------------------------------------

func ParseProgram(p *simplejson.Json) (*Program, error) {
	checkProgramType(p)

	body, err := parseBody(p.Get("body"))
	if err != nil { return nil, err }

	return &Program{
		Position: newPosition(p),
		body:     body,
	}, nil
}

func checkProgramType(p *simplejson.Json) {
	typ := p.Get("type").MustString()
	if typ != "Program" {
		log.Fatal("Expected Program, got %s", typ)
	}

}

// Parallel program parsing --------------------------------------------------

type ParallelParseResponse struct {
	index int
	node  interface{}
}

func ParallelParseProgram(p *simplejson.Json) *Program {
	checkProgramType(p)

	body := p.Get("body")
	bodySlice := body.MustArray()
	total := len(bodySlice)
	// Result array and channel
	bodyResult  := make([]interface{}, total)
	bodyChannel := make(chan ParallelParseResponse)

	for i, _ := range bodySlice {
		node := body.GetIndex(i)
		go parallelParseNode(bodyChannel, i, node)
	}

	received := 0
	for response := range bodyChannel {
		bodyResult[response.index] = response.node
		// Note that we've received a parse response and check if finished
		received++
		if received == total { break }
	}
	return &Program{
		Position: newPosition(p),
		body:     bodyResult,
	}
}

func parallelParseNode(ch chan<- ParallelParseResponse, index int, json *simplejson.Json) {
	node, err := parseNode(json)
	if err != nil { log.Fatal(err) }

	ch <- ParallelParseResponse{index, node}
}


// Utilities -----------------------------------------------------------------

func parseNode(n *simplejson.Json) (interface{}, error) {
	typ := n.Get("type").MustString()
	switch typ {
	case "BinaryExpression":         return ParseBinaryExpression(n)
	case "BlockStatement":           return ParseBlockStatement(n)
	case "ExpressionStatement":      return ParseExpressionStatement(n)
	case "ForStatement":             return ParseForStatement(n)
	case "FunctionDeclaration":      return ParseFunctionDeclaration(n)
	case "ExportDefaultDeclaration": return ParseExportDefaultDeclaration(n)
	case "ExportNamedDeclaration":   return ParseExportNamedDeclaration(n)
	case "ExportSpecifier":          return ParseExportSpecifier(n)
	case "ImportDeclaration":        return ParseImportDeclaration(n)
	case "ImportDefaultSpecifier":   return ParseImportDefaultSpecifier(n)
	case "ImportSpecifier":          return ParseImportSpecifier(n)
	case "Identifier":               return ParseIdentifier(n)
	case "ReturnStatement":          return ParseReturnStatement(n)
	case "Literal":                  return ParseLiteral(n)
	case "UpdateExpression":         return ParseUpdateExpression(n)
	case "VariableDeclaration":      return ParseVariableDeclaration(n)
	case "VariableDeclarator":       return ParseVariableDeclarator(n)
	default:                         return nil, fmt.Errorf("Can't parse type: %s\n", typ)
	}
}

func newPosition(p *simplejson.Json) Position {
	return Position{
		start: p.Get("start").MustInt(),
		end:   p.Get("end").MustInt(),
	}
}



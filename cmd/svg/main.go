package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type tokens struct {
	toks []string
	pos  int
}

func newTokens(input string) *tokens {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "{", "{ ")
	input = strings.ReplaceAll(input, "}", " }")
	return &tokens{toks: strings.Split(input, " ")}
}

func main() {
	var input string
	if len(os.Args) > 1 {
		input = os.Args[1]
	} else {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input = string(b)
	}
	toks := newTokens(input)
	expr := parse(toks)
	f, err := os.CreateTemp("/tmp", "expr-*.svg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := printSVG(f, expr); err != nil {
		panic(err)
	}
	cmd := exec.Command("open", f.Name())
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	fmt.Print(input)
}

func parse(t *tokens) any {
	left := parseOperand(t)
	if t.done() {
		return left
	}
	op := t.next()
	right := parseOperand(t)
	return expr{left: left, op: op, right: right}
}

func parseOperand(t *tokens) any {
	tok := t.next()
	if tok != "{" {
		return tok
	}
	operand := parse(t)
	if !t.done() && t.peek() == "}" {
		t.next()
	}
	return operand
}

func (t *tokens) done() bool {
	return t.pos >= len(t.toks)
}

func (t *tokens) next() string {
	token := t.peek()
	t.pos++
	return token
}

func (t *tokens) peek() string {
	if t.done() {
		panic("no more tokens")
	}
	return t.toks[t.pos]
}

// expr is a node in the AST
type expr struct {
	left  any
	op    string
	right any
}

type Node struct {
	S      string
	X, Y   float64
	Lx, Ly float64 // lineTo coordinates
}

func scaleX(x float64) float64 { return x*25 + 50 }
func scaleY(y float64) float64 { return y*15 + 10 }

func newNodes(v any, x, y, dx, lx, ly float64) []Node {
	node := Node{
		X:  scaleX(x),
		Y:  scaleY(y),
		Lx: scaleX(lx),
		Ly: scaleY(ly),
	}
	e, ok := v.(expr)
	if !ok {
		node.S = fmt.Sprintf("%v", v)
		return []Node{node}
	}
	node.S = e.op
	nodes := []Node{node}
	nodes = append(nodes, newNodes(e.left, x-dx, y+1, dx/2, x, y)...)
	nodes = append(nodes, newNodes(e.right, x+dx, y+1, dx/2, x, y)...)
	return nodes
}

//go:embed expr.svg.tmpl
var templ string

func printSVG(w io.Writer, expr any) error {
	nodes := newNodes(expr, 0, 0, 1, 0, 0)
	data := struct {
		Expr           string
		MaxX, MaxY     float64
		LabelX, LabelY float64
		Nodes          []Node
	}{
		Expr:  fmt.Sprintf("%v", expr),
		Nodes: nodes,
	}
	for _, n := range nodes {
		data.MaxX = max(data.MaxX, n.X+10)
		data.MaxY = max(data.MaxY, n.Y+30)
	}
	data.LabelX = data.MaxX / 2
	data.LabelY = data.MaxY - 10
	textOffset := func(f float64) float64 { return f + 2 }
	funcs := template.FuncMap{
		"textOffset": textOffset,
	}
	t := template.Must(template.New("svg").Funcs(funcs).Parse(templ))
	return t.Execute(w, data)
}

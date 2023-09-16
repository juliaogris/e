# e

Is is a cut down version of [evy].

This repository demonstrates how to build a minimal lexer and parser from
scratch without using any external libraries.

Try the `e` command with:

```sh
go run main.go test.e
```

[evy]: https://github.com/foxygoat/evy

## Syntax grammar

```
prog = { stmt } .
stmt = decl | assign .

decl  = ident ":" type .
ident = LETTER { LETTER | DIGIT } .
type  = "num" | "string" | "bool" .

assign = ident "=" expr .

expr    = operand | unary_expr |
          binary_expr .
operand = literal | ident | group .
literal = /* e.g. "abc", 1, 2.34, true, false */ .
group   = "(" expr ")" .

unary_expr = UNARY_OP expr .
UNARY_OP   = "-" | "!" .

binary_expr = expr BINARY_OP expr .
BINARY_OP = "*" | "/" | "%" |
            "+" | "-" |
            "<" | "<=" | ">" | ">=" |
            "==" | "!=" |
            "and" |
            "or" .
```

## `pratt` command

The `pratt` command is a stripped-down expression parser based on the
[Pratt parser]. To help users understand Pratt parsing, the command also
provides a naive, recursive, right-associative expression parser and a naive,
iterative, left-associative expression parser.

Try it with:

```
go run ./cmd/pratt/main.go '1 + 2 * 3'
```

[Pratt parser]: https://en.wikipedia.org/wiki/Pratt_parser

## `svg` command

The `svg` command is meant to be used with the `pratt` command. It generates
and opens an SVG image of the expression tree in the default SVG viewer.

Try it with:

```
go install ./cmd/...
pratt '1 * 2 + 3 * 4' | svg
```

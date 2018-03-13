package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	env := NewEnv() // for variable table
	for {
		err := readEvalPrint(in, env)
		if err != nil {
			break
		}
	}
}

func readEvalPrint(in *bufio.Reader, env *Env) error {
	fmt.Printf("> ")
	line, err := in.ReadBytes('\n')
	if err != nil {
		return err
	}
	ast, nbuf := Read(line)
	var _ = nbuf
	v := Eval(ast, env)
	fmt.Println(v)
	for index, element := range env.Var {
		fmt.Println("Index:", index, "Element:", element)
	}
	return nil
}

//Eval ...
func Eval(ast Ast, env *Env) Ast {
	return ast
}

//Read ...
func Read(b []byte) (ast Ast, n []byte) {
	return parseStatement(b)
}

func getNum(buf []byte) (n int, nbuf []byte) {
	n = digitVal(buf[0])
	nbuf = buf[1:]
	for len(nbuf) > 0 {
		if d := digitVal(nbuf[0]); d >= 0 {
			n = n*10 + d
		} else {
			break
		}
		nbuf = nbuf[1:]
	}
	return n, nbuf
}

func digitVal(b byte) int {
	return int(b - '0') // convert int
}

func getSymbol(buf []byte) (sym string, nbuf []byte) {
	var i int
	for i = 0; i < len(buf); i++ {
		if !isAlpha(buf[i]) && !isDigit(buf[i]) {
			break
		}
	}
	return string(buf[0:i]), buf[i:]
}
func isAlpha(buf byte) bool {
	return regexp.MustCompile(`[A-Za-z]`).Match([]byte{buf})
}
func isDigit(buf byte) bool {
	return regexp.MustCompile(`[0-9]`).Match([]byte{buf})
}

//Ast ...
type Ast interface {
	String() string
	Eval(env *Env) Ast
}

//Num ...
type Num int

func (n Num) String() string {
	return fmt.Sprintf("%d", int(n))
}

//Eval ...
func (n Num) Eval(env *Env) Ast {
	return n
}

//Symbol ...
type Symbol string

//Symbol.String() ...
func (s Symbol) String() string {
	return string(s)
}

//Eval ...
func (s Symbol) Eval(env *Env) Ast {
	return s
}

//Env ...
type Env struct{ Var map[string]Ast }

//NewEnv ...
func NewEnv() *Env {
	env := new(Env)
	env.Var = make(map[string]Ast)
	return env
}

//Set ...
func Set(env *Env, key string, n int) {
	env.Var[key] = Num(n)
}

func envValue(env *Env, key string) (n int, ok bool) {
	if v, found := env.Var[key]; found {
		n, v := getNum([]byte(v.String()))
		var _ = v
		return n, true
	}
	return 0, false
}

//AssignOp ...
type AssignOp struct {
	Var  Symbol
	Expr Ast
}

//AssignOp.String ...
func (b AssignOp) String() string {
	return string(b.Var)
}

//Eval ...
func (b AssignOp) Eval(env *Env) Ast {
	return b
}

func parseStatement(b []byte) (stmt Ast, n []byte) {
	stmt, n = parseExpression(b)
	if n[0] == '=' {
		if sym, ok := stmt.(Symbol); ok {
			var expr Ast
			expr, n = parseExpression(n[1:])
			stmt = AssignOp{Var: sym, Expr: expr}
		} else {
			panic("value is not symbol:" + stmt.String())
		}
	}
	return stmt, n
}

//BinOp ...
type BinOp struct {
	Op byte
	Left,
	Right Ast
}

//BinOp.String ...
func (b BinOp) String() string {
	return string(b.Op)
}

//Eval ...
func (b BinOp) Eval(env *Env) Ast {
	return b
}

func parseExpression(b []byte) (expr Ast, n []byte) {
	expr, n = parseTerm(b)
	for n[0] == '+' || n[0] == '-' {
		op := n[0]
		var term Ast
		term, n = parseTerm(n[1:])
		expr = BinOp{Op: op, Left: expr, Right: term}
	}
	return expr, n
}

func parseTerm(b []byte) (term Ast, n []byte) {
	term, n = parseFactor(b)
	for n[0] == '*' || n[0] == '/' {
		op := n[0]
		var factor Ast
		factor, n = parseFactor(n[1:])
		term = BinOp{Op: op, Left: term, Right: factor}
	}
	return term, n
}

func parseFactor(b []byte) (factor Ast, n []byte) {
	if isAlpha(b[0]) {
		factor, n := getSymbol(b[:1])
		var _ = n
		return Symbol(factor), b[1:]
	} else if isDigit(b[0]) {
		factor, n := getNum(b[:1])
		var _ = n
		return Num(factor), b[1:]
	} else {
		panic("invalid factor:" + string(b))
	}
}

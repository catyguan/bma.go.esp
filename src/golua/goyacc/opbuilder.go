package goyacc

import (
	"fmt"
	"strconv"
	"strings"
)

type OP uint8

const (
	OP_NONE    = OP(0)
	OP_VALUE   = OP(1)
	OP_ADD     = OP(2)
	OP_SUB     = OP(3)
	OP_MUL     = OP(4)
	OP_DIV     = OP(5)
	OP_PMUL    = OP(6)
	OP_MOD     = OP(7)
	OP_LT      = OP(8)
	OP_GT      = OP(9)
	OP_LTEQ    = OP(10)
	OP_GTEQ    = OP(11)
	OP_EQ      = OP(12)
	OP_NOTEQ   = OP(13)
	OP_STRADD  = OP(14)
	OP_RETURN  = OP(15)
	OP_BLOCK   = OP(16)
	OP_EXPLIST = OP(17)
	OP_LOCAL   = OP(18)
	OP_VAR     = OP(19)
	OP_AND     = OP(20)
	OP_OR      = OP(21)
	OP_ASSIGN  = OP(22)
	OP_IF      = OP(23)
	OP_UNTIL   = OP(24)
	OP_WHILE   = OP(25)
	OP_FOR     = OP(26)
	OP_FORIN   = OP(27)
	OP_NOT     = OP(28)
	OP_LEN     = OP(29)
	OP_NSIGN   = OP(30)
	OP_MEMBER  = OP(31)
	OP_FIELD   = OP(32)
	OP_TABLE   = OP(33)
	OP_ARRAY   = OP(34)
)

// type NILVALUE bool

// func IsNilValue(v interface{}) bool {
// 	if v == nil {
// 		return true
// 	}
// 	_, ok := v.(NILVALUE)
// 	return ok
// }

var OPNames = []string{
	"NONE", "VALUE",
	"+", "-", "*", "/", "^", "%",
	"<", ">", "<=", ">=", "==", "~=", "..",
	"return", "block", "explist", "local",
	"var", "and", "or", "=",
	"if", "until", "while", "for", "for-in",
	"not", "#", "-sign",
	"member", "field", "table", "array",
}

func toNode(yylex yyLexer, val *yySymType) (Node, error) {
	if val == nil {
		return nil, nil
	}
	if val.op == OP_VALUE {
		r := new(Node0)
		r.op = OP_VALUE
		r.Value = val.value
		return r, nil
	}
	if val.value != nil {
		if n, ok := val.value.(Node); ok {
			return n, nil
		}
	}
	if val.op == OP_NONE {
		if val.token.kind == NAME {
			r := new(Node0)
			r.op = OP_VAR
			r.Value = val.token.image
			return r, nil
		}
		return nil, nil
	}
	return nil, nil
	// return nil, fmt.Errorf("unknow node(%d)", val.token.kind)
}

func op1(yylex yyLexer, lval *yySymType, op OP, v1 *yySymType) {
	if yyDebug >= 2 {
		fmt.Println("op1 >> ", op, v1)
	}
	n1, err1 := toNode(yylex, v1)
	if err1 != nil {
		yylex.Error(err1.Error())
		return
	}
	r := new(Node1)
	r.op = op
	r.Child = n1
	lval.op = OP_NONE
	lval.value = r
	if yyDebug >= 2 {
		fmt.Println("op1 end: ", r)
	}
}

func op2(yylex yyLexer, lval *yySymType, op OP, v1 *yySymType, v2 *yySymType) {
	if yyDebug >= 2 {
		fmt.Println("op2 >> ", op, v1, v2)
	}
	if v1.op == OP_VALUE && v2.op == OP_VALUE {
		ok, v, err := ExecOp2(op, v1.value, v2.value)
		if err != nil {
			yylex.Error(err.Error())
			return
		}
		if ok {
			lval.op = OP_VALUE
			lval.value = v
			if yyDebug >= 2 {
				fmt.Println("op2 merge value >> ", op, v1.value, v2.value, v)
			}
			return
		}
	}
	n1, err1 := toNode(yylex, v1)
	if err1 != nil {
		yylex.Error(err1.Error())
		return
	}
	n2, err2 := toNode(yylex, v2)
	if err2 != nil {
		yylex.Error(err2.Error())
		return
	}
	r := new(Node2)
	r.op = op
	r.Child1 = n1
	r.Child2 = n2
	lval.op = OP_NONE
	lval.value = r
	if yyDebug >= 2 {
		fmt.Println("op2 end: ", r, n1, n2)
	}
}

func opAppend(yylex yyLexer, lval *yySymType, v1 *yySymType, v2 *yySymType) {
	doOpAppend(yylex, lval, v1, v2, OP_BLOCK)
}

func opExpList(yylex yyLexer, lval *yySymType, v1 *yySymType, v2 *yySymType) {
	doOpAppend(yylex, lval, v1, v2, OP_EXPLIST)
}

func doOpAppend(yylex yyLexer, lval *yySymType, v1 *yySymType, v2 *yySymType, op OP) {
	if yyDebug >= 2 {
		fmt.Println("opAppend >> ", v1, v2)
	}
	n1, err1 := toNode(yylex, v1)
	if err1 != nil {
		yylex.Error(err1.Error())
		return
	}
	n2, err2 := toNode(yylex, v2)
	if err2 != nil {
		yylex.Error(err2.Error())
		return
	}
	if n1 == nil {
		lval.op = OP_NONE
		lval.value = n2
		if yyDebug >= 2 {
			fmt.Println("opAppend end: nil, n2", n2)
		}
		return
	}
	if n2 == nil {
		lval.op = OP_NONE
		lval.value = n1
		if yyDebug >= 2 {
			fmt.Println("opAppend end: n1, nil", n1)
		}
		return
	}
	nn1, ok1 := n1.(*NodeN)
	if !ok1 {
		tmp := new(NodeN)
		tmp.op = op
		tmp.Childs = []Node{n1}
		nn1 = tmp
		if yyDebug >= 2 {
			fmt.Println("opAppend new block")
		}
	}
	lval.op = OP_NONE
	lval.value = nn1
	if nn2, ok2 := n2.(*NodeN); ok2 {
		if op == OP_BLOCK && nn2.op == OP_BLOCK {
			for _, cn := range nn2.Childs {
				nn1.AddChild(cn)
			}
			if yyDebug >= 2 {
				fmt.Println("opAppend end: merge block", len(nn1.Childs))
			}
			return
		}
	}
	nn1.AddChild(n2)
	if yyDebug >= 2 {
		fmt.Println("opAppend end: ", len(nn1.Childs))
	}
}

func opLocal(yylex yyLexer, lval *yySymType, nsval *yySymType, expl *yySymType) {
	if yyDebug >= 2 {
		fmt.Println("opLocal", nsval, expl)
	}
	n, err := toNode(yylex, expl)
	if err != nil {
		yylex.Error(err.Error())
		return
	}
	var ns []string
	if nsval.value == nil {
		ns = []string{nsval.token.image}
	} else {
		ns = nsval.value.([]string)
	}
	r := new(NodeLocal)
	r.op = OP_LOCAL
	r.Names = ns
	r.ExpList = n
	lval.op = OP_NONE
	lval.value = r
	if yyDebug >= 2 {
		fmt.Println("opLocal end: ", r)
	}
}

func opFlag(lval *yySymType, op OP) {
	lval.op = op
}

func opVar(lval *yySymType, val1 *yySymType) {
	r := new(Node0)
	r.op = OP_VAR
	r.Value = val1.token.image
	lval.op = OP_NONE
	lval.value = r
}

func opValueExt(lval *yySymType, v interface{}) {
	r := new(Node0)
	r.op = OP_VALUE
	r.Value = v

	lval.op = OP_NONE
	lval.value = r
}

func opValue(yylex yyLexer, lval *yySymType) {
	err := func() error {
		lval.op = OP_VALUE
		switch lval.token.kind {
		case NIL:
			lval.value = nil
		case TRUE:
			lval.value = true
		case FALSE:
			lval.value = false
		case STRING:
			lval.value = lval.token.image
		case NUMBER:
			s := lval.token.image
			if !strings.Contains(s, ".") {
				v32, err1 := strconv.ParseInt(s, 10, 32)
				if err1 == nil {
					lval.value = int32(v32)
					break
				}
				nerr := err1.(*strconv.NumError)
				if nerr.Err == strconv.ErrRange {
					v64, err2 := strconv.ParseInt(s, 10, 64)
					if err2 != nil {
						return err2
					}
					lval.value = v64
					break
				}
				if nerr.Err != strconv.ErrSyntax {
					return err1
				}
			}
			f64, err3 := strconv.ParseFloat(s, 64)
			if err3 != nil {
				return err3
			}
			lval.value = f64
			break
		}
		return nil
	}()
	if err != nil {
		yylex.Error(err.Error())
	}
	if yyDebug >= 2 {
		fmt.Println("opValue end: ", lval)
	}
}

func nameAppend(yylex yyLexer, lval *yySymType, val1 *yySymType, val2 *yySymType) {
	if yyDebug >= 2 {
		fmt.Println("nameAppend", val1, val2)
	}
	var ns []string
	if val1.value == nil {
		ns = []string{val1.token.image}
	} else {
		ns, _ = val1.value.([]string)
	}
	if val2 != nil {
		if val2.value == nil {
			ns = append(ns, val2.token.image)
		} else {
			lt, _ := val2.value.([]string)
			for _, s := range lt {
				ns = append(ns, s)
			}
		}
	}
	lval.op = OP_NONE
	lval.value = ns
	if yyDebug >= 2 {
		fmt.Println("nameAppend end: ", ns)
	}
}

func endChunk(yylex yyLexer, lval *yySymType) {
	if yyDebug >= 2 {
		fmt.Println("endChunk", lval)
	}
	n, err := toNode(yylex, lval)
	if err != nil {
		yylex.Error(err.Error())
		return
	}
	p := yylex.(*Parser)
	p.chunk = n
}

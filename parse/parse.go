package parse

import (
    "fmt"
    "os"
    . "g9cc_v2/common"
    . "g9cc_v2/utils"
)

// Recursive-descent parser

var pos int
var tokens []interface{}

func expect(ty int) {
    t := tokens[pos].(*Token)
    if t.Ty != ty {
        Error(fmt.Sprintf("%c (%d) expected, but got %c (%d)", ty, ty, t.Ty, t.Ty))
    }
    pos++
}

func consume(ty int) bool {
    t := tokens[pos].(*Token)
    if t.Ty != ty {
        return false
    }
    pos++
    return true
}

func new_node(op int, lhs *Node, rhs *Node) *Node {
    return &Node{Ty: op, Lhs: lhs, Rhs: rhs}
}

func fail(i int) {
    t := tokens[pos].(*Token)
    fmt.Printf("unexpected token: %s\n", t.Input)
    os.Exit(1)
}

func term() *Node {
    t := tokens[pos].(*Token)
    pos++

    if t.Ty == TK_NUM {
        return &Node{Ty: ND_NUM, Val: t.Val}
    }

    if t.Ty == TK_IDENT {
        return &Node{Ty: ND_IDENT, Name: t.Name}
    }

    Error(fmt.Sprintf("number or identifier expected, but got %s", t.Input))
    return nil
}

func mul() *Node {
    lhs := term()

    for {
        t := tokens[pos].(*Token)
        if t.Ty == '*' || t.Ty == '/' {
            pos++
            lhs = new_node(t.Ty, lhs, term())
        } else {
            break
        }
    }
    return lhs
}

func expr() *Node {
    lhs := mul()

    for {
        t := tokens[pos].(*Token)
        if t.Ty == '+' || t.Ty == '-' {
            pos++
            lhs = new_node(t.Ty, lhs, mul())
        } else {
            break
        }
    }

    return lhs
}

func assign() *Node {
    lhs := expr()
    if consume('=') {
        return new_node('=', lhs, expr())
    }
    return lhs
}

func stmt() *Node {

    node := &Node{Ty: ND_COMP_STMT}

    for {
        t := tokens[pos].(*Token)
        if t.Ty == TK_EOF {
            return node
        }

        e := new(Node)

        if t.Ty == TK_RETURN {
            // `return` tokenを読み飛ばす
            pos++
            e.Ty = ND_RETURN
            e.Expr = assign()
        } else {
            e.Ty = ND_EXPR_STMT
            e.Expr = assign()
        }

        node.Stmts = append(node.Stmts, e)

        expect(';')
    }

    // error
    return nil
}

func Parse(t []interface{}) *Node {
    tokens = t

    return stmt()
}

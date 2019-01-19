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

func new_node(op int, lhs *Node, rhs *Node) *Node {
    return &Node{Ty: op, Lhs: lhs, Rhs: rhs}
}

func fail(i int) {
    t := tokens[pos].(*Token)
    fmt.Printf("unexpected token: %s\n", t.Input)
    os.Exit(1)
}

func num() *Node {
    t := tokens[pos].(*Token)
    if t.Ty == TK_NUM {
        pos++
        return &Node{Ty: ND_NUM, Val: t.Val}
    }
    Error(fmt.Sprintf("number expected, but got %s", t.Input))
    return nil
}

func mul() *Node {
    lhs := num()

    for {
        t := tokens[pos].(*Token)
        if t.Ty == '*' || t.Ty == '/' {
            pos++
            lhs = new_node(t.Ty, lhs, num())
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

func stmt() *Node {

    node := &Node{Ty: ND_COMP_STMT}

    for {
        t := tokens[pos].(*Token)
        if t.Ty == TK_EOF {
            return node
        }

        e := new(Node)

        if t.Ty == TK_RETURN {
            pos++
            e.Ty = ND_RETURN
            e.Expr = expr()
        } else {
            Error("unexpected token")
        }

        node.Stmts = append(node.Stmts, e)

        expect(';')
    }

}

func Parse(t []interface{}) *Node {
    tokens = t
    pos = 0

    return stmt()
}
